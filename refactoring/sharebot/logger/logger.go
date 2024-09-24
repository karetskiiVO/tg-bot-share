package logger

import "log"

// Logger is type witch can make changes for log easier 
type Logger struct {
	log.Logger
}

// Warning create warning output same way as fmt.Print
func (l *Logger) Warning(v ...any) {
	l.Println(append([]any{"[warn]"}, v)...)
}
// Error create error output same way as fmt.Print
func (l *Logger) Error(v ...any)   {
	l.Println(append([]any{"[err]"}, v)...)
}
// Info create information output same way as fmt.Print
func (l *Logger) Info(v ...any) {
	l.Println(append([]any{"[info]"}, v)...)
}

// Warningf create warning output same way as fmt.Printf
func (l *Logger) Warningf(format string, v ...any) {
	l.Printf(format, v...)
}
// Errorf create error output same way as fmt.Printf
func (l *Logger) Errorf(format string, v ...any)   {
	l.Printf(format, v...)
}
// Infof create information output same way as fmt.Printf
func (l *Logger) Infof(format string, v ...any) {
	l.Printf(format, v...)
}
