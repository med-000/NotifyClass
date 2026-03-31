package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/med-000/notifyclass/pkg/service"

	"github.com/robfig/cron/v3"
)

var running int32 = 0

func main() {
	// 秒ありcron（ログ確認しやすい）
	c := cron.New(cron.WithSeconds())

	// 毎時0分0秒
	_, err := c.AddFunc("0 0 * * * *", func() {
		// 多重実行防止
		if !atomic.CompareAndSwapInt32(&running, 0, 1) {
			log.Println("skip: already running")
			return
		}
		defer atomic.StoreInt32(&running, 0)

		start := time.Now()
		log.Println("cron start:", start)

		if err := service.RunNotify(); err != nil {
			log.Println("error:", err)
			return
		}

		log.Println("cron end:", time.Since(start))
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	log.Println("cron started")

	// プロセスを止めない
	select {}
}
