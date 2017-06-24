package logs

import "sync"

type Logs interface {
	Log(name string) Log
}

type Config struct {
	Formats []FormatConfig
	Loggers []LoggerConfig
}

func NewConfig() Config {
	return Config{
		Loggers: []LoggerConfig{
			LoggerConfig{
				Type:  LoggerConsole,
				Level: LevelInfo,
			},
		},
	}
}

func New(config Config) Logs {
	formats := map[string]Format{}
	for _, fc := range config.Formats {
		if _, ok := formats[fc.Name]; ok {
			panic("logs: Duplicate format \"" + fc.Name + "\"")
		}
		formats[fc.Name] = newFormat(fc)
	}
	if _, ok := formats[""]; !ok {
		formats[""] = newDefaultFormat()
	}

	loggers := []Logger{}
	for _, lc := range config.Loggers {
		f := formats[lc.Format]
		if f == nil {
			panic("logs: Undefined format \"" + lc.Format + "\"")
		}

		l := newLogger(lc, f)
		loggers = append(loggers, l)
	}

	return &logs{
		logs:    make(map[string]Log),
		formats: formats,
		loggers: loggers,
	}
}

type logs struct {
	mu      sync.Mutex
	logs    map[string]Log
	loggers []Logger
	formats map[string]Format
}

func (logs *logs) Log(name string) Log {
	logs.mu.Lock()
	defer logs.mu.Unlock()

	log, ok := logs.logs[name]
	if ok {
		return log
	}

	log = newLog(logs, name)
	logs.logs[name] = log
	return log
}
