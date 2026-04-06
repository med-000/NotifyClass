package logger

import (
	"io"
	"log"
	"os"
)

type DiscordLogger struct {
	Info  *log.Logger
	Error *log.Logger
	Warn  *log.Logger
}

func NewDiscordLogger() (*DiscordLogger, error) {
	if err := os.MkdirAll("logs/discord", 0o755); err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile(
		"logs/discord/info.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile(
		"logs/discord/error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	warnFile, err := os.OpenFile(
		"logs/discord/warn.log",
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

	return &DiscordLogger{
		Info: log.New(
			infoWriter,
			BrightCyan+"[DISCORD]"+Green+"[INFO] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Error: log.New(
			errorWriter,
			BrightCyan+"[DISCORD]"+Red+"[ERROR] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Warn: log.New(
			warnWriter,
			BrightCyan+"[DISCORD]"+Yellow+"[WARN] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}, nil
}
