package g

import (
	"fmt"
	"github.com/mglslg/go-discord-dify/cmd/g/ds"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

var (
	Logger         *log.Logger
	AppContext     ds.AppContext
	SecToken       ds.Token
	UserSessionMap map[string]*ds.UserSession
)

// InitConfig readConfig reads the config file and unmarshals it into the config variable
func InitConfig(configPath string) {
	//解析默认配置文件
	log.Println("Reading default config file...")
	defaultConfigFile, err := os.ReadFile("config/default_config.yaml")
	if err != nil {
		log.Panic("Read default config failed...", err)
	}

	err = yaml.Unmarshal(defaultConfigFile, &AppContext)
	if err != nil {
		log.Panic("Resolve default config file failed!", err)
	}
	log.Println("Default Config file read successfully!")

	//解析自定义配置文件
	log.Println("Reading custom config file...")
	file, err := os.ReadFile(configPath)
	if err != nil {
		//log.Fatal("[ERROR]Read custom config failed...", err)
		log.Panic("Read custom config failed...", err)
		//panic(err)
	}

	err = yaml.Unmarshal(file, &AppContext)
	if err != nil {
		log.Panic("Resolve custom config file failed!", err)
	}

	// 将当前配置文件路径保存到环境变量中
	AppContext.ConfigFilePath = configPath
	log.Println("Custom config file read successfully!")
}

func InitLogger() *os.File {
	currentDate := time.Now().Format("2006-01-02")

	logPath := AppContext.LogFilePath

	log.Println("logPath:", logPath)

	// Check if the logs directory exists, create it if it does not exist
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(logPath, 0755); mkErr != nil {
			log.Panicf("Unable to create log directory: %v", mkErr)
		}
	}

	logFileName := fmt.Sprintf("%s/%s-%s.log", logPath, currentDate, AppContext.BotName)

	// Create a log file
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicf("Unable to open log file: %v", err)
	}

	// Create a logger
	Logger = log.New(io.MultiWriter(os.Stderr, f), "", log.LstdFlags|log.Lshortfile)

	return f
}

func InitSecretConfig() {
	log.Println("Init secret config...")

	if AppContext.BotToken != "" {
		SecToken.Discord = AppContext.BotToken
	} else {
		SecToken.Discord = os.Getenv("DISCORD_BOT_TOKEN")
	}

	if AppContext.DifyToken != "" {
		SecToken.Dify = AppContext.DifyToken
	} else {
		SecToken.Dify = os.Getenv("DIFY_TOKEN")
	}

	if SecToken.Discord == "" {
		log.Panic("DISCORD_BOT_TOKEN is not set, either set it in the environment or in the config file")
	}
	if SecToken.Dify == "" {
		log.Panic("DIFY_TOKEN is not set, either set it in the environment or in the config file")
	}

	Logger.Println("Secret Config file read successfully! DiscordToken:", SecToken.Discord, ",DifyToken:", SecToken.Dify)
}

func InitUserSession() {
	UserSessionMap = make(map[string]*ds.UserSession)
}

// GetUserSession Get the current user session, create it if it does not exist
func GetUserSession(authorId string, channelId string, authorName string) *ds.UserSession {
	key := getUserChannelId(authorId, channelId)
	_, exists := UserSessionMap[key]
	if !exists {
		UserSessionMap[key] = newUserSession(authorId, channelId, authorName)
	}
	return UserSessionMap[key]
}

func getLogFilePath(filePath string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("获取日志文件路径出错")
	}

	// 获取main文件所在目录
	baseDir := path.Dir(filename)

	if filepath.IsAbs(filePath) {
		return filePath
	} else {
		fullPath, err := filepath.Abs(path.Join(baseDir, filePath))
		if err != nil {
			log.Fatal("获取日志文件路径出错:", err)
		}
		return fullPath
	}
}

func newUserSession(authorId string, channelId string, authorName string) *ds.UserSession {
	userChannelId := getUserChannelId(authorId, channelId)
	return &ds.UserSession{
		UserId:        authorId,
		UserName:      authorName,
		UserChannelID: userChannelId,
		ChannelID:     channelId,
	}
}

func getUserChannelId(authorId string, channelId string) string {
	return authorId + "_" + channelId
}
