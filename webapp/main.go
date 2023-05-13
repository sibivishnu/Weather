package main

//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Packages
//----------------------------------------------
import (
	"log"
	"net/http"
	"os"

	"../common"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"google.golang.org/api/option"
)

// ----------------------------------------------
// Environment Variable Names
// ----------------------------------------------
const (
	ENV_ACCU_API_KEY          = "ACCU_API_KEY"
	ENV_REDIS_HOST            = "REDIS_HOST"
	FLAG_HTTP_PORT            = "HTTP_PORT"
	FLAG_HTTP_HOST            = "HTTP_HOST"
	FLAG_HTTP_SCHEME          = "HTTP_SCHEME"
	ENV_FIREBASE_SERVICE_FILE = "FIREBASE_APPLICATION_CREDENTIALS"
)

// ----------------------------------------------
// Global Variables
// ----------------------------------------------
var (
	options        map[string]interface{}
	httpHost       string
	httpScheme     string
	firebaseClient *auth.Client
)

//==============================================
// Functions
//==============================================

// ----------------------------------------------
// setupHTTP - prepare http routes.
// ----------------------------------------------
func setupHTTP(port string) {
	log.Println("[WebApp] Starting the http server on port : " + port)
	router := mux.NewRouter()

	// Forecast Calls
	router.HandleFunc("/api/v1.1/forecast/id/{id}", actionGetForecastData).Methods("GET")
	router.HandleFunc("/api/v2.0/forecast/id/{id}", actionGetForecastDataVer2).Methods("GET")
	router.HandleFunc("/api/v2.2/forecast/id/{id}", actionGetForecastDataJson).Methods("GET")

	router.HandleFunc("/api/v2.3/forecast/id/{id}/hourly", actionGetHourlyForecastDataJson).Methods("GET")
	router.HandleFunc("/api/v2.3/forecast/id/{id}/daily", actionGetDailyForecastDataJson).Methods("GET")

	// Test Data Calls
	router.HandleFunc("/api/v2.0/forecast/test/id/{id}", actionGetTestForecastData).Methods("GET")

	// Data Stream Calls
	router.HandleFunc("/api/v1.1/forecast/data-streams/id/{id}", actionGetForecastDataStreams).Methods("GET")

	// Device Location Calls
	router.HandleFunc("/api/v1.1/forecast/client/pc/{postal_code}/cc/{country_code}", actionGetLocationByPostalCode).Methods("GET")
	router.HandleFunc("/api/v1.1/forecast/client/cityorpc/{pc_or_city}/cc/{country_code}", actionGetLocationByCityOrPostalCode).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1.1/forecast/client/location/device/{device_id}", actionSetDeviceLocation).Methods("PUT", "POST")

	// Admin Calls
	router.HandleFunc("/api/v1.1/forecast/admin/id/{id}", actionAdminGetForecastData).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v2.0/forecast/admin/id/{id}", actionAdminGetForecastDataVer2).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v2.2/forecast/admin/id/{id}", actionAdminGetForecastDataJson).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1.1/forecast/admin/location/device/{device_id}", actionAdminUpdateDeviceLocation).Methods("PUT", "POST", "OPTIONS")
	router.HandleFunc("/api/v1.1/forecast/admin/getRanges/WeatherService/{cat_type}", actionAdminGetCategoryRanges).Methods("GET", "OPTIONS")

	// Root
	router.HandleFunc("/", actionDisplayCheckPage).Methods("GET")

	http.ListenAndServe(":"+port, router)
}

// ----------------------------------------------
// main - main program entry point
// ----------------------------------------------
func main() {
	app := cli.NewApp()
	app.Name = "Weather Service"
	app.Usage = "Weather Cache Service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  FLAG_HTTP_PORT,
			Value: "5000",
			Usage: "Http Port",
		},
	}

	app.Action = runIt
	app.Run(os.Args)
}

// ----------------------------------------------
// runIt - Application Action
// ----------------------------------------------
func runIt(runtimeContext *cli.Context) {
	log.Println("[WebApp] begin")

	//-------------------------------------------------
	// Init Globals
	//-------------------------------------------------
	options = make(map[string]interface{})
	options["redis.host"] = os.Getenv(ENV_REDIS_HOST)
	options["accuweather.key"] = os.Getenv(ENV_ACCU_API_KEY)
	options["datastore.project"] = "lax-gateway" // os.Getenv(ENV_PROJECT_ID)
	options["config.categories"] = "/conf/categories.json"
	common_init.LoadCommonEnvironment(options)

	// Locals
	firebaseServiceFile := os.Getenv(ENV_FIREBASE_SERVICE_FILE)
	opt := option.WithCredentialsFile(firebaseServiceFile)

	// Globals
	httpHost = os.Getenv(FLAG_HTTP_HOST)
	httpScheme = os.Getenv(FLAG_HTTP_SCHEME)

	// Configure FireBase App
	app, err := firebase.NewApp(common.CTX, nil, opt)
	if err != nil {
		log.Println("[WebApp] Firebase NewApp Error")
		panic(err)
	}

	// Configure Firebase Client
	firebaseClient, err = app.Auth(common.CTX)
	if err != nil {
		log.Println("[WebApp] Firebase Auth Error")
		panic(err)
	}

	// Prepare Http Request Handlers
	setupHTTP(runtimeContext.String(FLAG_HTTP_PORT))
}
