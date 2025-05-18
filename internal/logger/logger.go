package logger

import (
	"log"
	"os"
)

var Verbose bool = false
var l *log.Logger

func init() {
	l = log.New(os.Stdout, "", 0)
}

func Print(v ...any) {
	l.SetFlags(0)
	l.SetPrefix(reset.String())
	l.Print(v...)
}

func Println(v ...any) {
	l.SetFlags(0)
	l.SetPrefix(reset.String())
	l.Println(v...)
}

func Printf(format string, v ...any) {
	l.SetFlags(0)
	l.SetPrefix(reset.String())
	l.Printf(format, v...)
}

func Info(v ...any) {
	if !Verbose {
		return
	}
	l.SetFlags(log.Ltime)
	l.SetPrefix(fgBlueBright.String() + "[INFO] ")
	l.Print(v...)
}

func Infoln(v ...any) {
	if !Verbose {
		return
	}
	l.SetFlags(log.Ltime)
	l.SetPrefix(fgBlueBright.String() + "[INFO] ")
	l.Println(v...)
}

func Infof(format string, v ...any) {
	if !Verbose {
		return
	}
	l.SetFlags(log.Ltime)
	l.SetPrefix(fgBlueBright.String() + "[INFO] ")
	l.Printf(format, v...)
}

func Error(v ...any) {
	if Verbose {
		l.SetFlags(log.Ltime | log.Llongfile)
	} else {
		l.SetFlags(0)
	}
	l.SetPrefix(fgRed.String() + "[ERRORE] ")
	l.Print(v...)
}

func Errorln(v ...any) {
	if Verbose {
		l.SetFlags(log.Ltime | log.Llongfile)
	} else {
		l.SetFlags(0)
	}
	l.SetPrefix(fgRed.String() + "[ERRORE] ")
	l.Println(v...)
}

func Errorf(format string, v ...any) {
	if Verbose {
		l.SetFlags(log.Ltime | log.Llongfile)
	} else {
		l.SetFlags(0)
	}
	l.SetPrefix(fgRed.String() + "[ERRORE] ")
	l.Printf(format, v...)
}

func Fatal(v ...any) {
	if Verbose {
		l.SetFlags(log.Ltime | log.Llongfile)
	} else {
		l.SetFlags(0)
	}
	l.SetPrefix(fgRedBright.String() + "[CRITICO] ")
	l.Fatal(v...)
}

func Fatalln(v ...any) {
	if Verbose {
		l.SetFlags(log.Ltime | log.Llongfile)
	} else {
		l.SetFlags(0)
	}
	l.SetPrefix(fgRedBright.String() + "[CRITICO] ")
	l.Fatalln(v...)
}

func Fatalf(format string, v ...any) {
	if Verbose {
		l.SetFlags(log.Ltime | log.Llongfile)
	} else {
		l.SetFlags(0)
	}
	l.SetPrefix(fgRedBright.String() + "[CRITICO] ")
	l.Fatalf(format, v...)
}
