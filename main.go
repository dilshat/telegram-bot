package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot/v2"
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

	app := &application{
		tgClient:       &TbotWrapper{bot.Client()},
		attachmentsDir: GetEnv("ATTACHMENTS_DIR", "attachments"),
		token:          token,
		vmFactory:      VmFactoryImpl{},
	}

	if e := app.initialize(); e != nil {
		log.Fatal("Error initializing app ", e)
	}

	//bind handlers
	bot.HandleMessage("", app.messageHandler)
	bot.HandleCallback(app.callbackHandler)

	go func() {
		//let bot connect
		time.Sleep(time.Second * 3)
		app.onInit()
	}()

	//start bot
	log.Fatal(bot.Start())
}
