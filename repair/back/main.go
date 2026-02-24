package main

import (
	"fmt"
	"os"
	"repair/global/config"
	"repair/global/log"
	"repair/models"
	"repair/services"
	"runtime"

	"math/rand"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	log.Init()
	log.Info().Str("Version", config.Version).Str("Mode", config.Mode).Msg("Start")

	models.InitCache()

	os.MkdirAll(fmt.Sprintf("%v/detail", config.UploadPath), 0755)

	go services.Cron()
	//go services.Fcm()	cd /Users/seongjaehyeon/anb_origin2/repair/back
	scp .env.yml root@141.164.54.133:~/repair/back/
	go services.Notify()
	services.Http()
}
