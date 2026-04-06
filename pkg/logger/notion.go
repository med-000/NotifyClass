package logger

import (
	"io"
	"log"
	"os"
)

type NotionLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewNotionLogger() (*NotionLogger, error) {
	if err := os.MkdirAll("logs/notion", 0o755); err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile(
		"logs/notion/info.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile(
		"logs/notion/error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	warnFile, err := os.OpenFile(
		"logs/notion/warn.log",
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

	return &NotionLogger{
		Info: log.New(
			infoWriter,
			BrightMagenta+"[NOTION]"+Green+"[INFO] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Error: log.New(
			errorWriter,
			BrightMagenta+"[NOTION]"+Red+"[ERROR] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Warn: log.New(
			warnWriter,
			BrightMagenta+"[NOTION]"+Yellow+"[WARN] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}, nil
}
