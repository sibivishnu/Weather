package common

//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"cloud.google.com/go/datastore"
	"github.com/go-redis/redis"
	cache "github.com/sibivishnu/Weather/cacheUpdater"
	"golang.org/x/net/context"
)

// ----------------------------------------------
// Globals
// ----------------------------------------------
var (
	RedisClient     *redis.Client
	RedisInstance   *cache.RedisInstance
	CTX             context.Context
	DataStoreClient *datastore.Client
)
