package logger

import (
	"io"
	"log"
	"os"
)

type ScraperLogger struct {
	Info  *log.Logger
	Error *log.Logger
}

func NewScraperLogger() (*ScraperLogger, error) {
	if err := os.MkdirAll("logs/scraper", 0o755); err != nil {
		return nil, err
	}

	infoFile, err := os.OpenFile(
		"logs/scraper/info.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile(
		"logs/scraper/error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	return &ScraperLogger{
		Info: log.New(
			io.MultiWriter(os.Stdout, infoFile),
			Blue+"[SCRAPER]"+Green+"[INFO] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
		Error: log.New(
			io.MultiWriter(os.Stderr, errorFile),
			Blue+"[SCRAPER]"+Red+"[ERROR] "+Reset,
			log.Ldate|log.Ltime|log.Lshortfile,
		),
	}, nil
}
