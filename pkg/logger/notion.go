package logger

import "log"

type NotionLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewNotionLogger() (*NotionLogger, error) {
	infoWriter, errorWriter, warnWriter, err := buildLogWriters("notion")
	if err != nil {
		return nil, err
	}

	return &NotionLogger{
		Info: newStdLogger(
			infoWriter,
			BrightMagenta+"[NOTION]"+Green+"[INFO] "+Reset,
		),
		Error: newStdLogger(
			errorWriter,
			BrightMagenta+"[NOTION]"+Red+"[ERROR] "+Reset,
		),
		Warn: newStdLogger(
			warnWriter,
			BrightMagenta+"[NOTION]"+Yellow+"[WARN] "+Reset,
		),
	}, nil
}
