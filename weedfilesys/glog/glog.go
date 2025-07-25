package glog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type severity int32

// 获取严重程度
func (s *severity) get() severity {
	return severity(atomic.LoadInt32((*int32)(s)))
}

// 严重程度不就是级别吗，感觉能删
type Level int32

func (l *Level) get() Level {
	return Level(atomic.LoadInt32((*int32)(l)))
}

func (l *Level) set(level Level) {
	atomic.StoreInt32((*int32)(l), int32(level))
}

const (
	infoLog severity = iota
	warningLog
	errorLog
	fatalLog
	numSeverity //表示有几个严重程度
)
const severityChar = "IWEF"

// 我勒个，原来使用枚举代表下标来索引string
var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

type buffer struct {
	bytes.Buffer
	tmp  [64]byte // 临时字节数组是为了创建日志头
	next *buffer  //我真搞不懂，你创建一个新的buffer也是new一个，那为啥还要这个next干嘛
}

// 日志输出目标，但凡是要被用作输出目的地的对象比如控制台，文件……，都要实现这个接口
type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

// 表示用于vmodule 标志的一个过滤器，包括一个日志详细等级和一个用于匹配源文件的模式
// 🧠 应用场景：配合 -vmodule 参数
// 在命令行中，-vmodule 允许你按文件名设置不同的日志等级，比如：
// ./myserver -vmodule=net=3,main=2
// 表示：
//
// net.* 文件下的日志打印到 V(3) 的级别
//
// main.* 文件下的打印到 V(2) 的级别
//
// 每个规则在内部就被解析为一个 modulePat 对象
// 目前看来没啥用
type modulePat struct {
	pattern string // 要匹配的文件名或模式，比如main.go / *.pb.go
	literal bool   //如果为true，表示pattern是字面字符串 不带通配符
	level   Level  //日志详细级别，比如 v=2 表示Level 2
}

func (m *modulePat) match(file string) bool {
	if m.literal {
		return file == m.pattern
	}
	match, _ := filepath.Match(m.pattern, file)
	return match
}

// 目前看来没啥用
type moduleSpec struct {
	filter []modulePat
}

// 日志核心状态容器
type loggingT struct {
	toStderr     bool
	alsoToStderr bool

	// Level flag 表示什么等级以上的日志会输出到stderr，这是原子处理的是并发安全的
	stderrThreshold severity

	// 维护了一组日志输出用的临时内存缓冲区，减少内存分配
	freeList *buffer
	// 只保护了freeList，不与主日志锁mu共享，目的是让日志写入和缓冲池复用并发执行
	freeListMu sync.Mutex

	// 除了日志输出缓冲区，保障其他日志服务的锁
	mu sync.Mutex

	// 为每个level的日志提供对应的输出地点
	file [numSeverity]flushSyncWriter

	pcs [1]uintptr //难说

	vmap map[uintptr]Level // 缓存每个调用点所对应的日志等级

	filterLength int32 // 模块级别日志控制 链条 的长度，大于0表示启用的了模块界别控制

	traceLocation traceLocation // 在特定位置记录堆栈

	vmodule   moduleSpec // 模块级别日志控制配置，例如按文件名或包名设置不同的日志界别
	verbosity Level      // -v=2 决定默认打印哪些V日志

	exited bool // 标记日志系统是否已经退出，防止重复写日志
}

var logging loggingT

// 往buffer中填数字的工具函数，从右往左，从个位开始填
const digits = "0123456789"

func (buf *buffer) twoDigits(i, d int) {
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}
func (buf *buffer) nDigits(n, i, d int, pad byte) {
	j := n - 1
	for ; j >= 0 && d > 0; j-- {
		buf.tmp[i+j] = digits[d%10]
		d /= 10
	}
	for ; j >= 0; j-- {
		buf.tmp[i+j] = pad
	}
}
func (buf *buffer) someDigits(i, d int) int {
	// Print into the top, then copy down. We know there's space for at least
	// a 10-digit number.
	j := len(buf.tmp)
	for {
		j--
		buf.tmp[j] = digits[d%10]
		d /= 10
		if d == 0 {
			break
		}
	}
	return copy(buf.tmp[i:], buf.tmp[j:])
}

// 获取一个新的，准备投入使用的buffer
// freeList = b1 -> b2 -> b3 -> nil
// b = b1->, freeList:b2->b3->nil
// 断开b1到b2之间的链接,重置b1
func (l *loggingT) getNewBufferReadyToUse() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		b.next = nil
		b.Reset()
	}

	return b
}

// 格式化日志头部:	I0724 18:36:15.123456 myfile.go:42 ]
// 传日志等级，源文件名，文件的代码行号
func (l *loggingT) formatHeader(s severity, file string, line int) *buffer {
	now := time.Now()

	// 规正行号和级别
	if line < 0 {
		line = 0
	}
	if s > fatalLog {
		s = infoLog
	}

	// 从链表中拿一个buffer出来
	buf := l.getNewBufferReadyToUse()

	// now = 2025-07-24 18:36:15.654321
	_, mouth, day := now.Date()
	hour, minute, second := now.Clock()
	buf.tmp[0] = severityChar[s]
	buf.twoDigits(1, int(mouth))
	buf.twoDigits(3, int(day))
	buf.tmp[5] = ' '
	buf.twoDigits(6, hour)
	buf.tmp[8] = ':'
	buf.twoDigits(9, minute)
	buf.tmp[11] = ':'
	buf.twoDigits(12, second)
	buf.tmp[14] = '.'
	buf.nDigits(6, 15, now.Nanosecond()/1000, '0')
	buf.tmp[21] = ' '
	// I0724 18:36:15.654321

	buf.Write(buf.tmp[:22])
	buf.WriteString(file) // 写入 "main.go"
	buf.tmp[0] = ':'
	n := buf.someDigits(1, line) // 将 123 写入 tmp[1..3]
	buf.tmp[n+1] = ' '
	buf.Write(buf.tmp[:n+2])
	// I0724 18:36:15.654321 main.go:123
	return buf
}

// header 返回一个缓冲区，其中包含格式化后的日志头，以及用户代码所在的文件名和行号信息
// 参数depth制定了当前调用栈网上追溯多少层，以及定位应记录的源代码行（即日志调用发生的位置）
// 生成的日志行格式如下：
//
//	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...
//
// 其中各字段含义如下：
//
//	L                单个字符，表示日志级别（例如 'I' 表示 INFO）
//	mm               月份（补零，例如 5 月表示为 "05"）
//	dd               日期（补零）
//	hh:mm:ss.uuuuuu  时间（小时:分钟:秒.微秒）
//	threadid         当前线程 ID（使用空格补齐）
//	file             源文件名
//	line             行号
//	msg              用户传入的日志内容
func (l *loggingT) header(s severity, depth int) (*buffer, string, int) {
	// 获取调用日志函数的代码位置（文件名+行号）
	_, file, line, ok := runtime.Caller(depth + 3)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	return l.formatHeader(s, file, line), file, line
}

// 追溯源文件日志调用的综合信息
type traceLocation struct {
	file string
	line int
}

func (t *traceLocation) isSet() bool {
	return t.line > 0
}

func (t *traceLocation) match(file string, line int) bool {
	if t.line != line {
		return false
	}
	if i := strings.LastIndex(file, "/"); i >= 0 {
		file = file[i+1:]
	}
	return t.file == file
}

// stacks 获取当前goroutinue堆栈信息，all代表是否获取所有goroutinue，否则就是当前
func stacks(all bool) []byte {
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	// 最多尝试五次分配字节缓冲区
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		// 将当前或者全部goroutinue的堆栈信息写入trace中，并返回实际写入的字节数
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}

	return trace
}

// flushAll() 是日志系统中非常关键的一个步骤，负责将所有日志缓冲区中的内容刷新到磁盘，
// 并确保真正写入操作系统文件（fsync）。它常用于日志程序退出、崩溃、Fatal 日志等场景中。
func (l *loggingT) flushAll() {
	// Flush from fatal down, in case there's trouble flushing.
	for s := fatalLog; s >= infoLog; s-- {
		file := l.file[s]
		if file != nil {
			file.Flush() // 清空写入缓冲区
			file.Sync()  // 写入磁盘（fsync）
		}
	}
}

var logExitFunc func(error)

func (l *loggingT) exit(err error) {
	fmt.Fprintf(os.Stderr, "glog: exiting because of error: %s\n", err)
	// If logExitFunc is set, we do that instead of exiting.
	if logExitFunc != nil {
		logExitFunc(err)
		return
	}
	l.flushAll()
	l.exited = true // os.Exit(2)
}

// syncBuffer 将 bufio.Writer 与其底层文件关联起来，
// 既可以使用底层文件的 Sync 方法，
// 又可以包装 Write 方法以支持日志文件的滚动（log rotation）。
// 由于某些方法名冲突，底层文件不能直接内嵌（embed）。
// syncBuffer 的所有方法调用时都需要持有日志系统的 l.mu 锁。
type syncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file   *os.File
	sev    severity
	nbytes uint64 // 写入文件的字节数
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	if sb.logger.exited {
		return
	}
	if sb.nbytes+uint64(len(p)) >= MaxSize {
		if err := sb.rotateFile(time.Now()); err != nil {
			sb.logger.exit(err)
		}
	}
	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	if err != nil {
		sb.logger.exit(err)
	}
	return
}

// Log file created at: 2025/07/24 16:12:45
// Running on machine: node01
// Binary: Built with gc go1.22.1 for linux/amd64
// Log line format: [IWEF]mmdd hh:mm:ss threadid file:line] msg
func (sb *syncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}
	var err error
	sb.file, _, err = create(severityName[sb.sev], now)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Running on machine: %s\n", host)
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [IWEF]mmdd hh:mm:ss threadid file:line] msg\n")
	n, err := sb.file.Write(buf.Bytes())
	sb.nbytes += uint64(n)
	return err
}

func (l *loggingT) createFiles(sev severity) error {
	now := time.Now()
	for s := sev; s >= infoLog && l.file[s] == nil; s-- {
		sb := &syncBuffer{
			logger: l,
			sev:    s,
		}
		if err := sb.rotateFile(now); err != nil {
			return err
		}
		l.file[s] = sb
	}
	return nil
}

// 如果 fatalNoStacks 非 0，则在退出时不打印所有 goroutine 的堆栈信息。
// 它允许 Exit 及相关函数使用 Fatal 日志（但不输出堆栈）。
var fatalNoStacks uint32

func (l *loggingT) lockAndFlushAll() {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()
}

func Flush() {
	logging.lockAndFlushAll()
}

// 这段函数 timeoutFlush 是一个带超时控制的日志刷新函数，
// 用于在调用 glog.Fatal() 或程序退出前安全地 flush（刷新）日志缓冲区，同时防止死锁风险。
func timeoutFlush(timeout time.Duration) {
	done := make(chan bool, 1)
	go func() {
		Flush() // calls logging.lockAndFlushAll()
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintln(os.Stderr, "glog: Flush took longer than", timeout)
	}
}

func (l *loggingT) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// Let big buffers die a natural death.
		return
	}
	l.freeListMu.Lock()
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

type OutputStats struct {
	lines int64
	bytes int64
}

var Stats struct {
	Info, Warning, Error OutputStats
}
var severityStats = [numSeverity]*OutputStats{
	infoLog:    &Stats.Info,
	warningLog: &Stats.Warning,
	errorLog:   &Stats.Error,
}

func (l *loggingT) output(s severity, buf *buffer, file string, line int, alsoToStderr bool) {
	l.mu.Lock()
	if l.traceLocation.isSet() {
		if l.traceLocation.match(file, line) {
			buf.Write(stacks(false)) // 只获取当前goroutinue
		}
	}

	data := buf.Bytes()
	if l.toStderr {
		os.Stderr.Write(data)
	} else {
		if alsoToStderr || l.alsoToStderr || s >= l.stderrThreshold.get() {
			os.Stderr.Write(data)
		}
		if l.file[s] == nil {
			if err := l.createFiles(s); err != nil {
				os.Stderr.Write(data)
				l.exit(err)
			}
		}
		switch s {
		case fatalLog:
			l.file[fatalLog].Write(data)
			fallthrough
		case errorLog:
			l.file[errorLog].Write(data)
			fallthrough
		case warningLog:
			l.file[warningLog].Write(data)
			fallthrough
		case infoLog:
			l.file[infoLog].Write(data)
		}
	}
	if s == fatalLog {
		// If we got here via Exit rather than Fatal, print no stacks.
		if atomic.LoadUint32(&fatalNoStacks) > 0 {
			l.mu.Unlock()
			timeoutFlush(10 * time.Second)
			os.Exit(1)
		}
		// Dump all goroutine stacks before exiting.
		// First, make sure we see the trace for the current goroutine on standard error.
		// If -logtostderr has been specified, the loop below will do that anyway
		// as the first stack in the full dump.
		if !l.toStderr {
			os.Stderr.Write(stacks(false))
		}
		// Write the stack trace for all goroutines to the files.
		trace := stacks(true)
		logExitFunc = func(error) {} // If we get a write error, we'll still exit below.
		for log := fatalLog; log >= infoLog; log-- {
			if f := l.file[log]; f != nil { // Can be nil if -logtostderr is set.
				f.Write(trace)
			}
		}
		l.mu.Unlock()
		timeoutFlush(10 * time.Second)
		os.Exit(255) // C++ uses -1, which is silly because it's anded with 255 anyway.
	}
	l.putBuffer(buf)
	l.mu.Unlock()
	if stats := severityStats[s]; stats != nil {
		atomic.AddInt64(&stats.lines, 1)
		atomic.AddInt64(&stats.bytes, int64(len(data)))
	}
}

func (l *loggingT) printf(s severity, format string, args ...interface{}) {
	buf, file, line := l.header(s, 0)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}
func (l *loggingT) print(s severity, args ...interface{}) {
	l.printDepth(s, 1, args...)
}

func (l *loggingT) printDepth(s severity, depth int, args ...interface{}) {
	buf, file, line := l.header(s, depth)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf, file, line, false)
}
func (l *loggingT) println(s severity, args ...interface{}) {
	buf, file, line := l.header(s, 0)
	fmt.Fprintln(buf, args...)
	l.output(s, buf, file, line, false)
}

func (l *loggingT) setV(pc uintptr) Level {
	// 获取对应的函数名
	fn := runtime.FuncForPC(pc)

	// 获取该函数所在的源文件名和行号
	file, _ := fn.FileLine(pc)

	// 去掉文件名末尾的 .go 后缀（保留文件名主体）
	if strings.HasPrefix(file, ".go") {
		file = file[:len(file)-3]
	}

	// 提取纯文件名: a/b/c/d.go => a/b/c/d => d
	if slash := strings.LastIndex(file, "/"); slash >= 0 {
		file = file[slash+1:]
	}

	for _, filter := range l.vmodule.filter {
		if filter.match(file) {
			l.vmap[pc] = filter.level
			return filter.level
		}
	}
	l.vmap[pc] = 0
	return 0
}

type Verbose bool

func V(level Level) Verbose {
	if logging.verbosity.get() >= level {
		return Verbose(true)
	}

	if atomic.LoadInt32(&logging.filterLength) > 0 {
		logging.mu.Lock()
		defer logging.mu.Unlock()
		// 如果写入栈帧的个数是0，就返回一个Verbose(false)
		if runtime.Callers(2, logging.pcs[:]) == 0 {
			return Verbose(false)
		}
		v, ok := logging.vmap[logging.pcs[0]]
		if !ok {
			v = logging.setV(logging.pcs[0])
		}
		return Verbose(v >= level)
	}
	return Verbose(false)
}

// Info is equivalent to the global Info function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Info(args ...interface{}) {
	if v {
		logging.print(infoLog, args...)
	}
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infoln(args ...interface{}) {
	if v {
		logging.println(infoLog, args...)
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infof(format string, args ...interface{}) {
	if v {
		logging.printf(infoLog, format, args...)
	}
}

func Fatal(args ...interface{}) {
	logging.print(fatalLog, args...)
}
func FatalDepth(depth int, args ...interface{}) {
	logging.printDepth(fatalLog, depth, args...)
}
func Fatalf(format string, v ...interface{}) {
	logging.printf(fatalLog, format, v...)
}

func Error(args ...interface{}) {
	logging.print(errorLog, args...)
}

func Errorf(format string, v ...interface{}) {
	logging.printf(errorLog, format, v...)
}

func Warning(args ...interface{}) {
	logging.print(warningLog, args...)
}

func Warningf(format string, v ...interface{}) {
	logging.printf(warningLog, format, v...)
}
