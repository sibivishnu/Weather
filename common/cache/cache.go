package cache

//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"path"
	"time"
)

const (
	CACHE_BUST__GLOBAL      = ""
	CACHE_BUST__DEVICE      = ""
	CACHE_BUST__FILE        = ""
	CACHE_BUST__DEVICE_INFO = ":di-v4:"
)

type (
	RedisInstance struct {
		RedisSession *redis.Client
	}
)

var (
	localCache map[string][]byte
)

// ----------------------------------------------
// Exports
// ----------------------------------------------
func (redisInstance RedisInstance) GetCachedFile(folder string, file string, expiration time.Duration) ([]byte, error) {

	if localCache == nil {
		localCache = make(map[string][]byte)
	}

	var v []byte
	var ok bool
	if v, ok = localCache[file]; ok {
		return v, nil
	} else {
		body, err := ioutil.ReadFile(path.Join(folder, file))
		if err == nil {
			//redisInstance.SaveRedisData(body, key, expiration)
			log.Printf("[Redis] Saved file to in memory cache %s:%s", folder, file)
			localCache[file] = body
		} else {
			log.Printf("[Redis] Unable to read file %s:%s| %v", folder, file, err)
		}
		return body, err
	}

	/*
		hasher := md5.New()
		hasher.Write([]byte(file))
		hash := hex.EncodeToString(hasher.Sum(nil))
		key := "file.cache:" + hash

		r, err := redisInstance.GetCachedData(key)
		if (err != nil) {
			// Loading the template file, which depends on the category
			body, err := ioutil.ReadFile(path.Join(folder, file))
			if err == nil {
				redisInstance.SaveRedisData(body, key, expiration)
			} else {
				log.Printf("[Cache] Unable to read file %s:%s| %v", folder, file, err)
			}
			return body, err
		}
		return r, err
	*/
}

// Write the data to the cache
func (redisInstance RedisInstance) SaveRedisData(data []byte, key string, expiration time.Duration) error {
	key = key + CACHE_BUST__GLOBAL
	err := redisInstance.RedisSession.Set(key, data, expiration).Err()
	if err != nil {
		log.Printf(err.Error())
		log.Printf("[Redis] Unable to save key: %s to the cache| %v", key, err)
		return err
	}
	return nil
}

func (redisInstance RedisInstance) GetCachedData(key string) ([]byte, error) {
	key = key + CACHE_BUST__GLOBAL
	data, err := redisInstance.RedisSession.Get(key).Bytes()
	if err != nil {
		log.Printf("[Redis] Cache data not found for cache key: %s| %v", key, err)
		return nil, err
	}
	return data, nil
}

func (redisInstance RedisInstance) QueryCache(filter string) (map[string]string, error) {
	var cursor uint64
	var n int
	var err error
	returnList := make(map[string]string)

	for {
		var keys []string
		keys, cursor, err = redisInstance.RedisSession.Scan(cursor, filter, 100).Result()
		if err != nil {
			log.Println(err.Error())

		}

		n += len(keys)
		for _, key := range keys {
			value, err := redisInstance.RedisSession.Get(key).Result()
			if err != nil {
				log.Println(err.Error())
			} else {
				returnList[key] = value
			}
		}

		if cursor == 0 {
			break
		}
	}

	return returnList, err

}

func (redisInstance RedisInstance) RemoveKeyFromCache(key string) error {
	key = key + CACHE_BUST__GLOBAL
	err := redisInstance.RedisSession.Del(key).Err()
	if err != nil {
		log.Printf("[Redis] Unable to delete key : %s from the cache| %v", key, err)
		return err
	}
	return nil
}

func SetupRedis(redisHost *string, password *string, db int) *redis.Client {
	log.Printf("[Redis] Setting up Connection: %s", *redisHost)
	client := redis.NewClient(&redis.Options{
		Addr:     *redisHost,
		Password: *password,
		DB:       db, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("[Redis] Unable To Connect: %s", *redisHost)
		panic(err)
	}

	return client
}
