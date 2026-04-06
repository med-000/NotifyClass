package logger

import "log"

type DiscordLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewDiscordLogger() (*DiscordLogger, error) {
	infoWriter, errorWriter, warnWriter, err := buildLogWriters("discord")
	if err != nil {
		return nil, err
	}

	return &DiscordLogger{
		Info: newStdLogger(
			infoWriter,
			BrightCyan+"[DISCORD]"+Green+"[INFO] "+Reset,
		),
		Error: newStdLogger(
			errorWriter,
			BrightCyan+"[DISCORD]"+Red+"[ERROR] "+Reset,
		),
		Warn: newStdLogger(
			warnWriter,
			BrightCyan+"[DISCORD]"+Yellow+"[WARN] "+Reset,
		),
	}, nil
}
