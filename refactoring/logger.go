package sharebot

import "log"

type Logger struct {
	log.Logger
}

func (l Logger) Warning(v ...any) {
	l.Println(append([]any{"[warn]"}, v)...)
}
func (l Logger) Error(v ...any)   {
	l.Println(append([]any{"[err]"}, v)...)
}
func (l Logger) Info(v ...any) {
	l.Println(append([]any{"[info]"}, v)...)
}

func (l Logger) Warningf(format string, v ...any) {
	l.Printf(format, v...)
}
func (l Logger) Errorf(format string, v ...any)   {
	l.Printf(format, v...)
}
func (l Logger) Infof(format string, v ...any) {
	l.Printf(format, v...)
}
