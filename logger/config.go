package logger

type Config struct {
	ConsoleEnabled bool
	ConsoleLevel   string
	ConsoleJson    bool

	FileEnabled bool
	Filename    string
	FileLevel   string
	FileJson    bool

	FileDirectory string // Directory to log to to when filelogging is enabled

	FileMaxSize int // MaxSize the max size in MB of the logfile before it's rolled

	FileMaxBackups int // MaxBackups the max number of rolled files to keep

	FileMaxAge int // MaxAge the max age in days to keep a logfile

	FileCompress bool // Compress files

	Caller bool
}
