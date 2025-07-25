package glog

import (
	"context"
	"fmt"
	reqid "my_backend/weedfilesys/util/request_id"
	"sync/atomic"
)

const requestIDField = "request_id"

func formatMetaTag(ctx context.Context) string {
	if requestID := reqid.Get(ctx); requestID != "" {
		return fmt.Sprintf("%s:%s", requestIDField, requestID)
	}
	return ""
}

// InfoCtx 是 Verbose.Info 的一个支持 context 的替代版本。
// 它会在日志等级满足（由 v 控制）的前提下，将信息写入 INFO 日志，并在 context 中存在请求 ID 时，自动将该请求 ID 添加到日志前面。
// 所有参数的处理方式与 fmt.Print 一致（直接拼接打印，无格式化占位符）。
func (v Verbose) InfoCtx(ctx context.Context, args ...interface{}) {
	if !v {
		return
	}
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.print(infoLog, args...)
}

func (v Verbose) InfolnCtx(ctx context.Context, args ...interface{}) {
	if !v {
		return
	}
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.println(infoLog, args...)
}

func (v Verbose) InfofCtx(ctx context.Context, format string, args ...interface{}) {
	if !v {
		return
	}
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		format = metaTag + " " + format
	}
	logging.printf(infoLog, format, args...)
}

// InfofCtx logs a formatted message at info level, prepending a request ID from
// the context if it exists. This is a context-aware alternative to Infof.
func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		format = metaTag + " " + format
	}
	logging.printf(infoLog, format, args...)
}

// InfoCtx logs a message at info level, prepending a request ID from the context
// if it exists. This is a context-aware alternative to Info.
func InfoCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.print(infoLog, args...)
}

// WarningCtx logs to the WARNING and INFO logs.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to Warning.
func WarningCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.print(warningLog, args...)
}

// WarningDepthCtx logs to the WARNING and INFO logs with a custom call depth.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to WarningDepth.
func WarningDepthCtx(ctx context.Context, depth int, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.printDepth(warningLog, depth, args...)
}

// WarninglnCtx logs to the WARNING and INFO logs.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Println.
// This is a context-aware alternative to Warningln.
func WarninglnCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.println(warningLog, args...)
}

// WarningfCtx logs to the WARNING and INFO logs.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Printf.
// This is a context-aware alternative to Warningf.
func WarningfCtx(ctx context.Context, format string, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		format = metaTag + " " + format
	}
	logging.printf(warningLog, format, args...)
}

// ErrorCtx logs to the ERROR, WARNING, and INFO logs.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to Error.
func ErrorCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.print(errorLog, args...)
}

// ErrorDepthCtx logs to the ERROR, WARNING, and INFO logs with a custom call depth.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to ErrorDepth.
func ErrorDepthCtx(ctx context.Context, depth int, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.printDepth(errorLog, depth, args...)
}

// ErrorlnCtx logs to the ERROR, WARNING, and INFO logs.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Println.
// This is a context-aware alternative to Errorln.
func ErrorlnCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.println(errorLog, args...)
}

// ErrorfCtx logs to the ERROR, WARNING, and INFO logs.
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Printf.
// This is a context-aware alternative to Errorf.
func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		format = metaTag + " " + format
	}
	logging.printf(errorLog, format, args...)
}

// FatalCtx logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to Fatal.
func FatalCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.print(fatalLog, args...)
}

// FatalDepthCtx logs to the FATAL, ERROR, WARNING, and INFO logs with a custom call depth,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to FatalDepth.
func FatalDepthCtx(ctx context.Context, depth int, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.printDepth(fatalLog, depth, args...)
}

// FatallnCtx logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Println.
// This is a context-aware alternative to Fatalln.
func FatallnCtx(ctx context.Context, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.println(fatalLog, args...)
}

// FatalfCtx logs to the FATAL, ERROR, WARNING, and INFO logs,
// including a stack trace of all running goroutines, then calls os.Exit(255).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Printf.
// This is a context-aware alternative to Fatalf.
func FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		format = metaTag + " " + format
	}
	logging.printf(fatalLog, format, args...)
}

// ExitCtx logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to ExitCtx
func ExitCtx(ctx context.Context, args ...interface{}) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.print(fatalLog, args...)
}

// ExitDepthCtx logs to the FATAL, ERROR, WARNING, and INFO logs with a custom call depth,
// then calls os.Exit(1). Prepends a request ID from the context if it exists.
// Arguments are handled in the manner of fmt.Print.
// This is a context-aware alternative to ExitDepth.
func ExitDepthCtx(ctx context.Context, depth int, args ...interface{}) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.printDepth(fatalLog, depth, args...)
}

// ExitlnCtx logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Println.
// This is a context-aware alternative to Exitln.
func ExitlnCtx(ctx context.Context, args ...interface{}) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		args = append([]interface{}{metaTag}, args...)
	}
	logging.println(fatalLog, args...)
}

// ExitfCtx logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Prepends a request ID from the context if it exists. Arguments are handled in the manner of fmt.Printf.
// This is a context-aware alternative to Exitf.
func ExitfCtx(ctx context.Context, format string, args ...interface{}) {
	atomic.StoreUint32(&fatalNoStacks, 1)
	if metaTag := formatMetaTag(ctx); metaTag != "" {
		format = metaTag + " " + format
	}
	logging.printf(fatalLog, format, args...)
}
