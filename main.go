package main

import (
	"bytes"
	"database/sql"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache"
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

	cache := ttlcache.NewCache()
	cache.SetTTL(time.Duration(GetEnvAsInt("SESSION_TTL_MIN", 60)) * time.Minute)

	app := &application{
		tgClient:       &TbotWrapper{bot.Client()},
		cache:          cache,
		attachmentsDir: GetEnv("ATTACHMENTS_DIR", "attachments"),
		token:          token,
		logicScript:    loadScript(),
		vmFactory:      VmFactoryImpl{},
		dbClient:       setupDB(),
	}

	setupTimer(app)

	//bind handlers
	bot.HandleMessage("", app.messageHandler)
	bot.HandleCallback(app.callbackHandler)

	//start bot
	log.Fatal(bot.Start())
}

func setupDB() *sql.DB {
	if GetEnv("DB_DRIVER", "") == "" || GetEnv("DB_CONN_STR", "") == "" {
		return nil
	}

	dbClient, err := sql.Open(GetEnv("DB_DRIVER", ""), GetEnv("DB_CONN_STR", ""))
	if err != nil {
		log.Fatal("Failed to connect to db ", err)
	}
	if err = dbClient.Ping(); err != nil {
		if err != nil {
			log.Fatal("Failed to ping db ", err)
		}
	}

	return dbClient
}

//timer function
func setupTimer(app *application) {

	if GetEnv("timer", "") == "" {
		return
	}

	duration, err := time.ParseDuration(GetEnv("timer", ""))
	if err != nil {
		log.Fatal("Error parsing time duration for timer ", err)
	} else {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				app.handleMessage(nil, nil)
			}
		}()
	}
}

func loadScript() string {

	if GetEnv("SCRIPTS", "") == "" {
		log.Fatal("No script specified ")
	}

	scripts := strings.Split(GetEnv("SCRIPTS", ""), ",")

	var b bytes.Buffer
	for _, scriptPath := range scripts {
		script, err := ReadFile(scriptPath)
		if err != nil {
			log.Fatal("Failed to load script ", err)
		}
		b.WriteString(script)
		b.WriteString("\n")
	}

	return b.String()
}
