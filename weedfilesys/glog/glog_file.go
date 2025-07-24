package glog

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// log file 的最大字节数
var MaxSize uint64 = 1024 * 1024 * 1800
var MaxFileCount = 5

// 可以成为log file的候选目录
var logDirs []string

// If non-empty, overrides the choice of directory in which to write logs.
// See createLogDirs for the full list of possible destinations.
var logDir = flag.String("logdir", "", "If non-empty, write log files in this directory")

func createLogDirs() {
	if *logDir != "" {
		logDirs = append(logDirs, *logDir)
	} else {
		logDirs = append(logDirs, os.TempDir())
	}
}

var (
	// 记录当前运行环境的信息
	pid      = os.Getpid()               // 当前进程的ID
	program  = filepath.Base(os.Args[0]) // 程序名
	host     = "unknownhost"             // 主机名
	username = "unknownuser"             // 用户名 去除域名前缀
)

// pid      = 32875
// program  = "myapp"
// host     = "server01"
// username = "CORPDOMAIN_john.doe"
func init() {
	h, err := os.Hostname()
	if err == nil {
		host = shortHostname(h)
	}

	currnt, err := user.Current()
	if err == nil {
		username = currnt.Username
	}
	// "CORPDOMAIN\\john.doe" => "CORPDOMAIN_john.doe"
	username = strings.Replace(username, `\`, "_", -1)
}

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

// 日志文件命名生成函数 name包含tag和起始时间，以及一个对应tag的符号链接名link
// 假设以下变量：
//
// go
// 复制
// 编辑
// program  = "myapp"
// host     = "node01"
// userName = "john"
// tag      = "INFO"
// pid      = 12345
// t        = 2025-07-24 14:30:05
// 调用：
// name, link := logName("INFO", t)
// 则得到：
// name = "myapp.node01.john.log.INFO.20250724-143005.12345"
// link = "myapp.INFO"
// 日志文件：
// myapp.node01.john.log.INFO.20250724-143005.12345  ← 实际的日志文件
// 软链接（可选）：
// myapp.INFO → myapp.node01.john.log.INFO.20250724-143005.12345
// 用于始终指向最新的 INFO 日志文件
func logName(tag string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s.%s.log.%s.%04d%02d%02d-%02d%02d%02d.%d",
		program,
		host,
		username,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
	return name, program + "." + tag
}

func prefix(tag string) string {
	return fmt.Sprintf("%s.%s.%s.log.%s", program, host, username, tag)
}

var onceLogDirs sync.Once

// create 创建一个新的日志文件，返回文件对象和文件名。
// 文件名包含日志类型 tag（如 "INFO", "FATAL" 等）和时间 t。
// 创建成功后，它还会尝试更新该类型对应的软链接（忽略错误）
func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	// 1.先保证日志目录存在在
	onceLogDirs.Do(createLogDirs)
	if len(logDirs) == 0 {
		return nil, "", errors.New("no log dirs")
	}

	// 2.生成文件名和软连接名
	name, link := logName(tag, t)
	logPrefix := prefix(tag)
	var lastErr error

	// 3.遍历所有日志目录尝试写入
	for _, dir := range logDirs {
		entries, _ := os.ReadDir(dir)
		var previousLogs []string
		// 找出当前目录下所有匹配该类型的旧日志文件（用 logPrefix 作为前缀匹配）
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), logPrefix) {
				previousLogs = append(previousLogs, entry.Name())
			}
		}
		if len(previousLogs) >= MaxFileCount {
			sort.Strings(previousLogs)
			for i, entry := range previousLogs {
				if i > len(previousLogs)-MaxFileCount {
					break
				}
				os.Remove(filepath.Join(dir, entry))
			}
		}

		// 创建新的日志文件
		fname := filepath.Join(dir, name)
		f, err := os.Create(fname)
		if err == nil {
			symlink := filepath.Join(dir, link)
			os.Remove(symlink)
			os.Symlink(name, symlink)
			return f, fname, nil
		}

		lastErr = err
	}

	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}
