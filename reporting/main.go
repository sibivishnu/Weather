package main

import (
	"github.com/go-redis/redis"
	"github.com/sibivishnu/Weather/common/cache"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

const (
	ENV_ACCU_API_KEY      = "ACCU_API_KEY"
	ENV_REDIS_HOST        = "REDIS_HOST"
	FLAG_HTTP_PORT        = "HTTP_PORT"
	ENV_SMTP_HOST         = "SMTP_HOST"
	ENV_SMTP_PORT         = "SMTP_PORT"
	ENV_SMTP_USERNAME     = "SMTP_USERNAME"
	ENV_SMTP_PASSWORD     = "SMTP_PASSWORD"
	ENV_SMTP_FROM_NAME    = "SMTP_FROM_NAME"
	ENV_SMTP_FROM_ADDRESS = "SMTP_FROM_ADDRESS"
	ENV_SMTP_RECIPIENTS   = "SMTP_RECIPIENTS"
)

var (
	redisClient   *redis.Client
	redisInstance *cache.RedisInstance

	smtpHost        string
	smtpPort        string
	smtpUsername    string
	smtpPassword    string
	smtpFromName    string
	smtpFromAddress string
	smtpRecipients  string
)

func main() {

	app := cli.NewApp()

	app.Name = "Weather Cache updater Service"
	app.Usage = "Weather Cache updater Service"
	app.Action = runIt
	app.Run(os.Args)

}

func runIt(c *cli.Context) {

	log.Println("main start")

	redisHost := os.Getenv(ENV_REDIS_HOST)
	redisClient = setupRedis(&redisHost)
	redisInstance = &cache.RedisInstance{RedisSession: redisClient}

	smtpHost = os.Getenv(ENV_SMTP_HOST)
	smtpPort = os.Getenv(ENV_SMTP_PORT)
	smtpUsername = os.Getenv(ENV_SMTP_USERNAME)
	smtpPassword = os.Getenv(ENV_SMTP_PASSWORD)
	smtpFromName = os.Getenv(ENV_SMTP_FROM_NAME)
	smtpFromAddress = os.Getenv(ENV_SMTP_FROM_ADDRESS)
	smtpRecipients = os.Getenv(ENV_SMTP_RECIPIENTS)

	sendDailyMail()

	mailSenderProcess := time.NewTicker(24 * time.Hour)
	done := make(chan bool)

	for {
		select {
		case <-mailSenderProcess.C:
			sendDailyMail()
		case <-done:
			return
		}
	}
}

// establish a connection to redis
func setupRedis(redisHost *string) *redis.Client {
	log.Println(*redisHost)
	client := redis.NewClient(&redis.Options{
		Addr:     *redisHost,
		Password: "",
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
