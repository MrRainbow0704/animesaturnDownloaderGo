package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Verbose = false
	logger  = log.New(os.Stderr, "", 0)
)

func Logger() *log.Logger {
	return logger
}

// Info logs an informational message with a cyan "[INFO]" prefix.
// Arguments are handled in the manner of [fmt.Print].
func Info(v ...any) {
	if !Verbose {
		return
	}

	logger.SetPrefix(ColorizeNoReset("[INFO] ", ColorFgCyanBright))
	logger.Output(2, Colorize(fmt.Sprint(v...), ColorFgCyanBright))
	logger.SetPrefix("")
}

// Infof logs a formatted informational message with a cyan "[INFO]" prefix.
// Arguments are handled in the manner of [fmt.Printf].
func Infof(format string, v ...any) {
	if !Verbose {
		return
	}

	logger.SetPrefix(ColorizeNoReset("[INFO] ", ColorFgCyanBright))
	logger.Output(2, Colorize(fmt.Sprintf(format, v...), ColorFgCyanBright))
	logger.SetPrefix("")
}

// Infoln logs an informational message with a cyan "[INFO]" prefix, followed by a newline.
// Arguments are handled in the manner of [fmt.Println].
func Infoln(v ...any) {
	if !Verbose {
		return
	}

	logger.SetPrefix(ColorizeNoReset("[INFO] ", ColorFgCyanBright))
	logger.Output(2, Colorize(fmt.Sprintln(v...), ColorFgCyanBright))
	logger.SetPrefix("")
}

// Warn logs a warning message with a yellow "[WARN]" prefix.
// Arguments are handled in the manner of [fmt.Print].
func Warn(v ...any) {
	logger.SetPrefix(ColorizeNoReset("[WARN] ", ColorFgYellow))
	logger.Output(2, Colorize(fmt.Sprint(v...), ColorFgYellow))
	logger.SetPrefix("")
}

// Warnf logs a formatted warning message with a yellow "[WARN]" prefix.
// Arguments are handled in the manner of [fmt.Printf].
func Warnf(format string, v ...any) {
	logger.SetPrefix(ColorizeNoReset("[WARN] ", ColorFgYellow))
	logger.Output(2, Colorize(fmt.Sprintf(format, v...), ColorFgYellow))
	logger.SetPrefix("")
}

// Warnln logs a warning message with a yellow "[WARN]" prefix, followed by a newline.
// Arguments are handled in the manner of [fmt.Println].
func Warnln(v ...any) {
	logger.SetPrefix(ColorizeNoReset("[WARN] ", ColorFgYellow))
	logger.Output(2, Colorize(fmt.Sprintln(v...), ColorFgYellow))
	logger.SetPrefix("")
}

// Error logs an error message with a red "[ERROR]" prefix.
// Arguments are handled in the manner of [fmt.Print].
func Error(v ...any) {
	logger.SetPrefix(ColorizeNoReset("[ERROR] ", ColorFgRed))
	logger.Output(2, Colorize(fmt.Sprint(v...), ColorFgRed))
	logger.SetPrefix("")
}

// Errorf logs a formatted error message with a red "[ERROR]" prefix.
// Arguments are handled in the manner of [fmt.Printf].
func Errorf(format string, v ...any) {
	logger.SetPrefix(ColorizeNoReset("[ERROR] ", ColorFgRed))
	logger.Output(2, Colorize(fmt.Sprintf(format, v...), ColorFgRed))
	logger.SetPrefix("")
}

// Errorln logs an error message with a red "[ERROR]" prefix, followed by a newline.
// Arguments are handled in the manner of [fmt.Println].
func Errorln(v ...any) {
	logger.SetPrefix(ColorizeNoReset("[ERROR] ", ColorFgRed))
	logger.Output(2, Colorize(fmt.Sprintln(v...), ColorFgRed))
	logger.SetPrefix("")
}

// Fatal logs a fatal message with a magenta "[FATAL]" prefix and then calls [os.Exit](1).
// Arguments are handled in the manner of [fmt.Print].
func Fatal(v ...any) {
	s := fmt.Sprint(v...)
	logger.SetPrefix(ColorizeNoReset("[FATAL] ", ColorFgMagenta))
	logger.Output(2, Colorize(s, ColorFgMagenta))
	logger.SetPrefix("")
	os.Exit(1)
}

// Fatalf logs a formatted fatal message with a magenta "[FATAL]" prefix and then calls [os.Exit](1).
// Arguments are handled in the manner of [fmt.Printf].
func Fatalf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	logger.SetPrefix(ColorizeNoReset("[FATAL] ", ColorFgMagenta))
	logger.Output(2, Colorize(s, ColorFgMagenta))
	logger.SetPrefix("")
	os.Exit(1)
}

// Fatalln logs a fatal message with a magenta "[FATAL]" prefix, followed by a newline, and then calls os.Exit(1).
// Arguments are handled in the manner of [fmt.Println].
func Fatalln(v ...any) {
	s := fmt.Sprintln(v...)
	logger.SetPrefix(ColorizeNoReset("[FATAL] ", ColorFgMagenta))
	logger.Output(2, Colorize(s, ColorFgMagenta))
	logger.SetPrefix("")
	os.Exit(1)
}
