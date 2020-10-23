package main 

import (
	"time"
	"log"
	
	"github.com/robfig/cron"

	"ops-backend/view/asset/host"
	"ops-backend/pkg"
)

func main () {
	log.Println("Starting...")

	pkg.Setup()

    c := cron.New()
    c.AddFunc("*/10 * * * * *", func() {
        log.Println("Run asset.SyncHost...")
        host.SyncHost()
	})
	
	c.Start()

    t1 := time.NewTimer(time.Second * 10)
    for {
        select {
        case <-t1.C:
            t1.Reset(time.Second * 10)
        }
    }
}