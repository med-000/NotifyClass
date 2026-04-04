package logger

import (
	"io"
	"log"
	"os"
)

type RepositoryLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewRepositoryLogger() (*RepositoryLogger, error) {
	if err := os.MkdirAll("logs/repository", 0o755); err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile(
		"logs/repository/info.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile(
		"logs/repository/error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	warnFile, err := os.OpenFile(
		"logs/warn/info.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	appFile, err := os.OpenFile(
		"logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}
	infoWriter := io.MultiWriter(
		os.Stdout,
		infoFile,
		appFile,
	)

	errorWriter := io.MultiWriter(
		os.Stderr,
		errorFile,
		appFile,
	)

	warnWriter := io.MultiWriter(
		os.Stderr,
		warnFile,
		appFile,
	)

	return &RepositoryLogger{
		Info: log.New(
			infoWriter,
			Magenta+"[REPOSITORY]"+Green+"[INFO] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Error: log.New(
			errorWriter,
			Magenta+"[REPOSITORY]"+Red+"[ERROR] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Warn: log.New(
			warnWriter,
			Blue+"[REPOSITORY]"+Yellow+"[WARN] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}, nil
}
