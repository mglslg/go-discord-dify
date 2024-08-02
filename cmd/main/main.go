package main

import (
	"flag"
	"github.com/mglslg/go-discord-dify/cmd/g"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger *log.Logger

func main() {
	//设置日期格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var configFilePath string

	flag.StringVar(&configFilePath, "config", "", "path to config file")
	flag.Parse()

	if configFilePath == "" {
		log.Panic("config file must be specified !")
	}

	g.InitConfig(configFilePath)

	logFile := g.InitLogger()
	logger = g.Logger

	//todo check配置完整性，自定义的那几项配置配置

	g.InitSecretConfig()
	g.InitUserSession()

	session, err := initDiscordSession()

	if err != nil {
		logger.Fatal("Error g discord session:", err)
		return
	} else {
		logger.Println("Session init successfully")
	}

	err = session.Open()
	if err != nil {
		logger.Fatal("Error opening connection:", err)
		return
	}

	g.AppContext.BotId = session.State.User.ID

	logger.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = session.Close()
	if err != nil {
		log.Fatal("session closing failed")
	}

	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Fatal("logFile closing failed")
		}
	}(logFile)
}
