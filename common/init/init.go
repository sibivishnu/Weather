package init

//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"cloud.google.com/go/datastore"
	"github.com/go-redis/redis"

	//	cache "github.com/sibivishnu/Weather/cacheUpdater"
	"log"
	"time"

	"github.com/sibivishnu/Weather/common"
	common "github.com/sibivishnu/Weather/common"
	cache "github.com/sibivishnu/Weather/common/cache"
	"github.com/sibivishnu/Weather/common/const/device"
	"github.com/sibivishnu/Weather/common/providers/weather_api"
	"golang.org/x/net/context"
)

// ----------------------------------------------
// Global Variables
// ----------------------------------------------
var (
	redisClient *redis.Client
)

// ----------------------------------------------
// Exports
// ----------------------------------------------
func LoadCommonEnvironment(options map[string]interface{}) error {
	var v interface{}
	var ok bool

	//=============================================
	// Initialize Common
	//=============================================

	// 1. Setup Redis
	redisHost := "127.0.0.1"
	redisPassword := ""
	redisDB := 0

	if v, ok = options["redis.host"]; ok {
		redisHost = v.(string)
	}

	if v, ok = options["redis.pass"]; ok {
		redisPassword = v.(string)
	}

	if v, ok = options["redis.db"]; ok {
		redisDB = v.(int)
	}

	common.RedisClient = cache.SetupRedis(&redisHost, &redisPassword, redisDB)
	common.RedisInstance = &cache.RedisInstance{RedisSession: common.RedisClient}

	// 2. Setup DataStore Client
	common.CTX = context.Background()
	if v, ok = options["datastore.project"]; ok {
		common.DataStoreClient, _ = datastore.NewClient(common.CTX, v.(string))
	} else {
		log.Printf("[Common] Options should include datastore.project entry")
	}

	//=============================================
	// Initialize accuWeather
	//=============================================

	// 1. Setup LocationMap
	weather_api.LocationMap = make(map[string]*time.Location)

	// 2. Setup  AccuWeather Key
	if v, ok = options["accuweather.key"]; ok {
		weather_api.AccuApiKey = v.(string)
	} else {
		log.Printf("[Common] Options should include accuweather.key entry")
	}

	//=============================================
	// Initialize device
	//=============================================
	if v, ok = options["config.categories"]; ok {
		device.Categories = device.LoadCategoryConf(v.(string))
	} else {
		log.Printf("[Common] Options should include config.categories entry")
	}

	return nil
}
