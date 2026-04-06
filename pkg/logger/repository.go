package logger

import "log"

type RepositoryLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewRepositoryLogger() (*RepositoryLogger, error) {
	infoWriter, errorWriter, warnWriter, err := buildLogWriters("repository")
	if err != nil {
		return nil, err
	}

	return &RepositoryLogger{
		Info: newStdLogger(
			infoWriter,
			Magenta+"[REPOSITORY]"+Green+"[INFO] "+Reset,
		),
		Error: newStdLogger(
			errorWriter,
			Magenta+"[REPOSITORY]"+Red+"[ERROR] "+Reset,
		),
		Warn: newStdLogger(
			warnWriter,
			Blue+"[REPOSITORY]"+Yellow+"[WARN] "+Reset,
		),
	}, nil
}
