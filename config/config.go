package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	LoggerId int64
	OwnerId  []int64

	MongoUri      string
	StartImage    string
	StartImageGif string
	Token         string
	StartTime     time.Time
)

func init() {
	godotenv.Load()
	StartTime = time.Now()

	parseToInt64 := func(val string) int64 {
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Error parsing int64: %v", err))
		}
		return parsed
	}

	Token = Getenv[string]("TOKEN", "8050656956:AAGsJ8EniqZ1Bhe6F5xSelX08C43kzqboQI", nil)
	StartImage = Getenv[string](
		"START_IMG_URL",
		"https://telegra.ph/file/ba238ec5e542d8754cea7-dc1786aa23ae1224f2.jpg",
		nil,
	)
	StartImageGif = Getenv[string](
		"START_IMG_GIF",
		"https://raw.githubusercontent.com/Vivekkumar-IN/assets/refs/heads/master/ezgif-408f355da640ed.gif",
		nil,, 
	)
	LoggerId = Getenv("LOGGER_ID", "-1002647107199", parseToInt64)
	MongoUri = Getenv[string]("MONGO_DB_URI", "mongodb+srv://marin:marin69@cluster0.zxaf7uc.mongodb.net/?retryWrites=true&w=majority", nil)

 OwnerId = Getenv("OWNER_ID", "7706682472", func(key string) []int64 {

id := strings.Split(key, " ")
var ids int64

for _, k := range(id) {
ids = append(ids, parseToInt64(k))

}


})

	OwnerId = Getenv("OWNER_ID", "7706682472", parseToInt64)

	if Token == "" {
		log.Panic("TOKEN environment variable is empty")
	}
	if MongoUri == "" {
		log.Panic("MONGO_DB_URI environment variable is empty")
	}
}

func Getenv[T any](key, defaultValue string, convert func(string) T) T {
	value := defaultValue
	if envValue, exists := os.LookupEnv(key); exists {
		value = envValue
	}

	if convert != nil {
		return convert(value)
	}

	return any(value).(T)
}
