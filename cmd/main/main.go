package main

import (
	"flag"
	"fmt"
	"github.com/mglslg/go-discord-gpt/cmd/g"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger *log.Logger

func main() {
	var configFilePath string

	//todo 这边需要判断一下是否没有指定文件，如果没有指定则需要报错
	flag.StringVar(&configFilePath, "config", "", "path to config file")
	flag.Parse()

	g.InitConfig(configFilePath)
	logFile := g.InitLogger()
	logger = g.Logger
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

	logger.Println("Bot is now running.")
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

	defer logFile.Close()
}
