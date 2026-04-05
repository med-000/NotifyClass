package logger

import (
	"io"
	"log"
	"os"
)

type ServiceLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewServiceLogger() (*ServiceLogger, error) {
	if err := os.MkdirAll("logs/Service", 0o755); err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile(
		"logs/Service/info.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile(
		"logs/Service/error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	warnFile, err := os.OpenFile(
		"logs/Service/warn.log",
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

	return &ServiceLogger{
		Info: log.New(
			infoWriter,
			BrightBlue+"[SERVICE]"+Green+"[INFO] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Error: log.New(
			errorWriter,
			BrightBlue+"[SERVICE]"+Red+"[ERROR] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Warn: log.New(
			warnWriter,
			BrightBlue+"[SERVICE]"+Yellow+"[WARN] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}, nil
}
