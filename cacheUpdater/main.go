package main
//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"../common/init"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//----------------------------------------------
// Constants
//----------------------------------------------
const (
	ENV_ACCU_API_KEY               = "ACCU_API_KEY"
	ENV_REDIS_HOST                 = "REDIS_HOST"
	FLAG_HTTP_PORT                 = "HTTP_PORT"
	ENV_MAX_QUEUE                  = "MAX_QUEUE"
	ENV_MAX_WORKER                 = "MAX_WORKER"
	ENV_SCP_SERVER_HOST            = "SCP_SERVER_HOST"
	ENV_SCP_SERVER_USER            = "SCP_SERVER_USER"
	ENV_SCP_SERVER_RSA             = "SCP_SERVER_RSA"
	ENV_DEVICE_REMOTE_FILE_PATH    = "DEVICE_REMOTE_FILE_PATH"
	ENV_DEVICE_LOCAL_TARGET_FOLDER = "DEVICE_LOCAL_TARGET_FOLDER"
	ENV_PROJECT_ID                 = "PROJECT_ID"
	ENV_SUBSCRIPTION_NAME          = "SUBSCRIPTION_NAME"
	ENV_TOPIC_NAME                 = "TOPIC_NAME"
	ENV_ATTRIBUTE_TOPIC_NAME       = "ATTRIBUTE_TOPIC_NAME"
)


//----------------------------------------------
// Globals
//----------------------------------------------
var (
	options map[string]interface{}

	scpServerHost  string
	scpServerUser  string
	scpServerRSA   string
	remoteFilePath string
	targetFolder   string
	devicesFile    string

	projectID        string
	subscriptionName string
	topicName        string

	attributeSubscription string
	attributeTopic string
)

//----------------------------------------------
// Local Funcs
//----------------------------------------------
func main() {
	app := cli.NewApp()
	app.Name = "Weather Cache updater Service"
	app.Usage = "Weather Cache updater Service"
	app.Action = runIt
	app.Run(os.Args)
}

func runIt(c *cli.Context) {
	log.Println("[CacheUpdater] Begin")

	// Load Paths from Environment
	scpServerHost = os.Getenv(ENV_SCP_SERVER_HOST)
	scpServerUser = os.Getenv(ENV_SCP_SERVER_USER)
	scpServerRSA = os.Getenv(ENV_SCP_SERVER_RSA)
	remoteFilePath = os.Getenv(ENV_DEVICE_REMOTE_FILE_PATH)
	targetFolder = os.Getenv(ENV_DEVICE_LOCAL_TARGET_FOLDER)
	projectID = os.Getenv(ENV_PROJECT_ID)
	subscriptionName = os.Getenv(ENV_SUBSCRIPTION_NAME)
	topicName = os.Getenv(ENV_TOPIC_NAME)
	maxQueue, _ := strconv.Atoi(os.Getenv(ENV_MAX_QUEUE))
	maxWorker, _ := strconv.Atoi(os.Getenv(ENV_MAX_WORKER))

	// Derive Attribute PubSub Items
	attributeSubscription = subscriptionName + "_Attr"
	attributeTopic = "AttrSync"

	// Prepare Full File Path for Upload
	devicesFile = filepath.Join(targetFolder, filepath.Base(remoteFilePath))

	//-------------------------------------------------
	// Init Globals
	//-------------------------------------------------
	options = make(map[string]interface{})
	options["redis.host"] = os.Getenv(ENV_REDIS_HOST)
	options["accuweather.key"] = os.Getenv(ENV_ACCU_API_KEY)
	options["datastore.project"] = "lax-gateway" // os.Getenv(ENV_PROJECT_ID)
	options["config.categories"] = "/conf/categories.json"
	common_init.LoadCommonEnvironment(options)


	//-----------------------------------------
	// Launch Services
	//-----------------------------------------
	JobQueue = make(chan Job, maxQueue)
	dispatcher := NewDispatcher(maxWorker)
	dispatcher.Run()
	go listenGeo()
	go listenAttr()
	log.Println("[cacheUpdater] runCacheIdUpdater()")

	runCacheIDUpdater()
	deviceCacheUpdater := time.NewTicker(120 * time.Minute)
	dstUpdater := time.NewTicker(6 * time.Hour)
	done := make(chan bool)
	for {
		select {
		case <-deviceCacheUpdater.C:
			runCacheIDUpdater()
		case <-dstUpdater.C:
			dstUpdateProcess()
		case <-done:
			return
		}
	}
}
