package common

//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	cache "./cache"
	"cloud.google.com/go/datastore"
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
)

//----------------------------------------------
// Globals
//----------------------------------------------
var (
	RedisClient *redis.Client
	RedisInstance *cache.RedisInstance
	CTX            context.Context
	DataStoreClient         *datastore.Client
)