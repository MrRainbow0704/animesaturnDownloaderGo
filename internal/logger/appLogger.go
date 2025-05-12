package logger

import (
	"context"
	"fmt"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppLogger struct {
	ctx *context.Context
}

func (a *AppLogger) SetContext(c *context.Context) {
	a.ctx = c
}

func (a *AppLogger) Print(message string) {
	if a.ctx == nil {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("console.log(\"%s\");", message))
	Print(message)
}

func (a *AppLogger) Trace(message string) {
	if a.ctx == nil {
		return
	}
	if !Verbose {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("window.notifications.info(\"%s\", 3000);", message))
	Info(message)
}

func (a *AppLogger) Debug(message string) {
	if a.ctx == nil {
		return
	}
	if !Verbose {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("window.notifications.info(\"%s\", 3000);", message))
	Info(message)
}

func (a *AppLogger) Info(message string) {
	if a.ctx == nil {
		return
	}
	if !Verbose {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("window.notifications.info(\"%s\", 3000);", message))
	Info(message)
}

func (a *AppLogger) Warning(message string) {
	if a.ctx == nil {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("window.notifications.error(\"%s\", 3000);", message))
	Error(message)
}

func (a *AppLogger) Error(message string) {
	if a.ctx == nil {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("window.notifications.error(\"%s\", 3000);", message))
	Error(message)
}

func (a *AppLogger) Fatal(message string) {
	if a.ctx == nil {
		return
	}
	wails.WindowExecJS(*a.ctx, fmt.Sprintf("window.notifications.fatal(\"%s\", 3000);", message))
	Fatal(message)
}
