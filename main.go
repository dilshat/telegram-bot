package main

import (
	"path/filepath"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/labstack/gommon/log"

	"github.com/joho/godotenv"
	"github.com/yanzay/tbot/v2"
)

const (
	logicScript   = "logic.js"
	libraryScript = "lib.js"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env ", err)
	}
}

func main() {

	token := GetEnv("TELEGRAM_TOKEN", "")
	bot := tbot.New(token)

	cache := ttlcache.NewCache()
	cache.SetTTL(time.Duration(GetEnvAsInt("SESSION_TTL_MIN", 60)) * time.Minute)

	logicScript, err := loadScript()
	if err != nil {
		log.Fatal("Failed to load logic script ", err)
	}

	app := &application{
		client:         &TbotWrapper{bot.Client()},
		cache:          cache,
		attachmentsDir: GetEnv("ATTACHMENTS_DIR", "attachments"),
		token:          token,
		logicScript:    logicScript,
		vmFactory:      VmFactoryImpl{},
	}

	//bind handlers
	bot.HandleMessage("", app.messageHandler)
	bot.HandleCallback(app.callbackHandler)

	//start bot
	log.Fatal(bot.Start())
}

func loadScript() (string, error) {
	libScriptPath := filepath.Join(GetEnv("SCRIPTS_PATH", "scripts"), libraryScript)
	libScript, err := ReadFile(libScriptPath)
	if err != nil {
		log.Error("Failed to load library script ", err)
	}

	logicScriptPath := filepath.Join(GetEnv("SCRIPTS_PATH", "scripts"), logicScript)
	logicScript, err := ReadFile(logicScriptPath)
	if err != nil {
		return "", err
	}

	return logicScript + "\n" + libScript, nil
}
