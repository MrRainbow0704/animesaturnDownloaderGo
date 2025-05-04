package logger

type AppLogger struct{}

func (a *AppLogger) Print(message string)   { Print(message) }
func (a *AppLogger) Trace(message string)   { Info(message) }
func (a *AppLogger) Debug(message string)   { Info(message) }
func (a *AppLogger) Info(message string)    { Info(message) }
func (a *AppLogger) Warning(message string) { Error(message) }
func (a *AppLogger) Error(message string)   { Error(message) }
func (a *AppLogger) Fatal(message string)   { Fatal(message) }
