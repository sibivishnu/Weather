package main

//----------------------------------------------
// Packages
//----------------------------------------------
import (
	"../common"
	"../common/cache"
	"../common/const/device"
	"../common/providers/weather_api"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gopkg.in/guregu/null.v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//----------------------------------------------
// Constants
//----------------------------------------------

//==============================================
// Type Definitions
//==============================================

//----------------------------------------------
// @ReplyMessage
//----------------------------------------------
type ReplyMessage struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

//----------------------------------------------
// @AdminResponse
//----------------------------------------------
type AdminResponse struct {
	Zip         string `json:"zip"`
	City        string `json:"city"`
	Location    string `json:"location"`
	LastUpdated string `json:"last_updated"`
	Forecast    string `json:"forecast"`
}

//----------------------------------------------
// @AdminResponseV2
//----------------------------------------------
type AdminResponseV2 struct {
	Geo         device.Geo
	Location    weather_api.PostalCodeResponse
	LastUpdated LastUpdated
	Forecast    weather_api.ApiResponseInterface
}

type AdminResponseFormatV1p2 struct {
	Geo         device.Geo
	Location    weather_api.PostalCodeResponse
	LastUpdated LastUpdated
	Forecast    weather_api.ApiResponseInterface
}



type AdminResponseFormatV1p3 struct {
	Geo         device.Geo
	Location    weather_api.PostalCodeResponse
	LastUpdated LastUpdated
	Forecast    weather_api.ApiResponseInterface
}

type LastUpdated struct {
	OneDayForecast null.String
	CurrentForecast null.String
	TenDayForecast null.String
	TwentyFourHourForecast null.String
}

//----------------------------------------------
// @DeviceLocation
//----------------------------------------------
type DeviceLocation struct {
	CountryCode      string `json:"country_code"`
	PostalCode       string `json:"postal_code"`
	CityOrPostalCode string `json:"city_or_postal_code"`
	ACWKey           string `json:"acw_key"`
}


//==============================================
// Protocols
//==============================================
func (s AdminResponseV2) JsonResponse(version string) (string, error) {
	switch version {
	case "1.1e":
		fallthrough
	case "1.1":
		var j, err = json.Marshal(s)
		return string(j), err
		//var out bytes.Buffer
		//json.Indent(&out, j, "", "  ")
		//return string(out.Bytes()), err
	case "1.2e":
		fallthrough
	case "1.2":
		var response AdminResponseFormatV1p2
		response.Geo = s.Geo
		response.Location = s.Location
		response.LastUpdated = s.LastUpdated

		if s.Forecast != nil {
			f, _ := s.Forecast.ResponseFormat(version)
			response.Forecast = f
		} else {
			response.Forecast = nil
		}

		var j, err = json.Marshal(response)
		return string(j), err
	case "1.3e":
		fallthrough
	case "1.3":
		var response AdminResponseFormatV1p2
		response.Geo = s.Geo
		response.Location = s.Location
		response.LastUpdated = s.LastUpdated

		if s.Forecast != nil {
			f, _ := s.Forecast.ResponseFormat(version)
			response.Forecast = f
		} else {
			response.Forecast = nil
		}

		var j, err = json.Marshal(response)
		return string(j), err

	case "1.4e":
		fallthrough
	case "1.4":
		var response AdminResponseFormatV1p3
		response.Geo = s.Geo
		response.Location = s.Location
		response.LastUpdated = s.LastUpdated

		if s.Forecast != nil {
			f, _ := s.Forecast.ResponseFormat(version)
			response.Forecast = f
		} else {
			response.Forecast = nil
		}

		var j, err = json.Marshal(response)
		return string(j), err


	case "1.5e":
		fallthrough
	case "1.5":
		var response AdminResponseFormatV1p3
		response.Geo = s.Geo
		response.Location = s.Location
		response.LastUpdated = s.LastUpdated

		if s.Forecast != nil {
			f, _ := s.Forecast.ResponseFormat(version)
			response.Forecast = f
		} else {
			response.Forecast = nil
		}

		var j, err = json.Marshal(response)
		return string(j), err


	//var out bytes.Buffer
	//json.Indent(&out, j, "", "  ")
	//return string(out.Bytes()), err
	default:
		return "{\"outcome\": \"unknown version\"}", nil
	}
}


//==============================================
// Functions - Api Actions
//==============================================

//----------------------------------------------
// @actionGetForecastData
// [GET] /api/v1.1/forecast/id/{id}
//----------------------------------------------
func actionGetForecastData(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))

	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}
	
	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, err.Error())
		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "<anonymous:true>")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		//res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	// Get Forecast
	if testOverride(deviceID) {
		res = location.GetWeatherForecastTest(display.Category, display.ID)
	} else {
		res = location.GetWeatherForecast(display.Category, display.ID, "BASIC", firmwareVersion)
	}

	// Update Device Request Details
	updateDeviceRequestEntry(deviceID)

	// syncElixirBackend
	syncElixirBackend(display)

	// Return response
	io.WriteString(rw, res)
}

//----------------------------------------------
// @actionGetForecastDataVer2
// [GET] /api/v2.0/forecast/id/{id}
//----------------------------------------------
func actionGetForecastDataVer2(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]
	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))

	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}

	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, err.Error())






		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "<anonymous:true>")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		//res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	// Get Forecast
	if testOverride(deviceID) {
		res = location.GetWeatherForecastTest(display.Category, display.ID)
	} else {
		res = location.GetWeatherForecastV2(display.Category, display.ID, firmwareVersion)
	}

	// Update Device Request Details
	updateDeviceRequestEntry(deviceID)

	// syncElixirBackend
	syncElixirBackend(display)

	// Return response
	io.WriteString(rw, res)
}

//----------------------------------------------
// @actionGetForecastDataJson - nullable support with json payload.
// [GET] /api/v2.2/forecast/id/{id}
//----------------------------------------------
func actionGetForecastDataJson(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]
	(rw).Header().Set("Content-Type", "application/json")

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))
	callSubVersion := strings.TrimSpace(r.FormValue("v"))
	i8nSet := strings.TrimSpace(r.FormValue("i8nV"))


	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}

	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, "{\"anonymous\": true, \"error\": true, \"msg\": \"" + err.Error() + "\"}")
		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "{\"anonymous\": true}")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		//res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	// Get Forecast
	version := "1.4"
	if callSubVersion == "3" {
		version = "1.5"
	}
	if callSubVersion == "4" {
		version = "1.3"
	}

	if i8nSet == "2" {
		version = version + "e"
	}

	if testOverride(deviceID) {
		// NYI, need json formatter. res = location.GetWeatherForecastTest(display.Category, display.ID)
		res, _ = location.NullableGetWeatherForecastJson(display.Category, display.ID, firmwareVersion, callSubVersion).JsonResponse(version)
	} else {
		res, _ = location.NullableGetWeatherForecastJson(display.Category, display.ID, firmwareVersion, callSubVersion).JsonResponse(version)
	}

	// Update Device Request Details
	updateDeviceRequestEntry(deviceID)

	// syncElixirBackend
	syncElixirBackend(display)

	// Return response
	io.WriteString(rw, res)
}


//----------------------------------------------
// @actionGetHourlyForecastDataJson - nullable support with json payload.
// [GET] /api/v2.3/forecast/id/{id}/hourly
//----------------------------------------------
func actionGetHourlyForecastDataJson(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]
	(rw).Header().Set("Content-Type", "application/json")

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))
	callSubVersion := strings.TrimSpace(r.FormValue("v"))
	i8nSet := strings.TrimSpace(r.FormValue("i8nV"))


	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}

	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, "{\"anonymous\": true, \"error\": true, \"msg\": \"" + err.Error() + "\"}")
		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "{\"anonymous\": true}")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		//res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	// Get Forecast
	version := "1.4"
	if callSubVersion == "3" {
		version = "1.5"
	}
	if callSubVersion == "4" {
		version = "1.3"
	}

	if i8nSet == "2" {
		version = version + "e"
	}

	if testOverride(deviceID) {
		res, _ = location.NullableGetWeatherForecastJsonExtended(display.Category, display.ID, firmwareVersion, callSubVersion, false, false, true, false).JsonResponse(version)
	} else {
		res, _ = location.NullableGetWeatherForecastJsonExtended(display.Category, display.ID, firmwareVersion, callSubVersion, false,false, true, false).JsonResponse(version)
	}

	// Update Device Request Details
	// updateDeviceRequestEntry(deviceID)

	// syncElixirBackend
	// syncElixirBackend(display)

	// Return response
	io.WriteString(rw, res)
}


//----------------------------------------------
// @actionGetDailyForecastDataJson - nullable support with json payload.
// [GET] /api/v2.3/forecast/id/{id}/hourly
//----------------------------------------------
func actionGetDailyForecastDataJson(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]
	(rw).Header().Set("Content-Type", "application/json")

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))
	callSubVersion := strings.TrimSpace(r.FormValue("v"))
	i8nSet := strings.TrimSpace(r.FormValue("i8nV"))


	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}

	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, "{\"anonymous\": true, \"error\": true, \"msg\": \"" + err.Error() + "\"}")
		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "{\"anonymous\": true}")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		//res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	// Get Forecast
	version := "1.4"
	if callSubVersion == "3" {
		version = "1.5"
	}
	if callSubVersion == "4" {
		version = "1.3"
	}

	if i8nSet == "2" {
		version = version + "e"
	}

	if testOverride(deviceID) {
		res, _ = location.NullableGetWeatherForecastJsonExtended(display.Category, display.ID, firmwareVersion, callSubVersion, false,true, false, true).JsonResponse(version)
	} else {
		res, _ = location.NullableGetWeatherForecastJsonExtended(display.Category, display.ID, firmwareVersion, callSubVersion, false,true, false, true).JsonResponse(version)
	}

	// Update Device Request Details
	updateDeviceRequestEntry(deviceID)

	// syncElixirBackend
	syncElixirBackend(display)

	// Return response
	io.WriteString(rw, res)
}

//----------------------------------------------
// @actionGetTestForecastData
// [GET] /api/v2.0/forecast/test/id/{id}
//----------------------------------------------
// API for test device forecast data see ticket https://github.com/lacrossetech/weather-service/issues/16
func actionGetTestForecastData(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]

	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}

	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, err.Error())
		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "<anonymous:true>")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		//res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	// Get Forecast
	res = location.GetWeatherForecastTest(display.Category, display.ID)

	// Return response
	io.WriteString(rw, res)
}

//----------------------------------------------
// @actionGetForecastDataStreams
// [GET] /api/v1.1/forecast/data-streams/id/{id}
//----------------------------------------------
func actionGetForecastDataStreams(rw http.ResponseWriter, r *http.Request) {
	var res string
	vars := mux.Vars(r)
	deviceID := vars["id"]

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))

	// 1. Check for unsupported devices
	if !supportedDevice(deviceID) {
		io.WriteString(rw, "device blocked")
		return
	}

	// Load device from Redis
	display, anonymous, err := getDevice(deviceID)
	if err != nil {
		// No data found from the cache
		io.WriteString(rw, "device not found")
		return
	}

	// HMAC digest check
	err = hmacCheck(r, display.PSK)
	if err != nil {
		io.WriteString(rw, err.Error())
		return
	}

	// Check for Anonymous
	if anonymous {
		io.WriteString(rw, "<anonymous:true>")
		return
	}

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		// res = "Location " + display.Geo.ACWKey + "|" + display.Geo.Zip + " : " + display.Geo.CountryCode + " Not found"
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	res = location.GetWeatherForecast(display.Category, display.ID, "DATASTREAMS", firmwareVersion)
	io.WriteString(rw, res)
}


//----------------------------------------------
// @actionGetLocationByPostalCode
// [GET] /api/v1.1/forecast/client/pc/{postal_code}/cc/{country_code}
//----------------------------------------------
func actionGetLocationByPostalCode(rw http.ResponseWriter, r *http.Request) {

	if isTokenValid(r) == false {
		sendApiOutcomeResponse(rw, http.StatusUnauthorized, errors.New("Wrong bearer token"))
		return
	}

	vars := mux.Vars(r)
	countryCode := strings.TrimSpace(vars["country_code"])
	postalCode := strings.TrimSpace(vars["postal_code"])

	resBytes, err := weather_api.SearchAllLocationsPerCountry(postalCode, countryCode)
	if err != nil {
		log.Printf(err.Error())
		io.WriteString(rw, "{\"locations\":[]}")
		return
	}

	rw.Write(resBytes)

}

//----------------------------------------------
// @actionGetLocationByCityOrPostalCode
// [GET] /api/v1.1/forecast/client/cityorpc/{pc_or_city}/cc/{country_code}
//----------------------------------------------
func actionGetLocationByCityOrPostalCode(rw http.ResponseWriter, r *http.Request) {

	setResponseHeaders(&rw, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isTokenValid(r) == false {
		sendApiOutcomeResponse(rw, http.StatusUnauthorized, errors.New("Wrong bearer token"))
		return
	}

	vars := mux.Vars(r)
	countryCode := strings.TrimSpace(vars["country_code"])
	postalCode := strings.TrimSpace(vars["pc_or_city"])

	resBytes, err := weather_api.SearchAllLocationsPerPCorCity(postalCode, countryCode)
	if err != nil {
		log.Printf(err.Error())
		io.WriteString(rw, "{\"locations\":[]}")
		return
	}

	rw.Write(resBytes)

}

//----------------------------------------------
// @actionSetDeviceLocation
// [PUT,POST] /api/v1.1/forecast/client/location/device/{device_id}
//----------------------------------------------
func actionSetDeviceLocation(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	deviceID := strings.TrimSpace(vars["device_id"])



	if isTokenValid(r) == false {
		sendApiOutcomeResponse(rw, http.StatusUnauthorized, errors.New("Wrong bearer token"))
		log.Printf("handleDeviceLocationUpdate error wrong barreer token for device %s", deviceID)
		return
	}

	// Get the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendApiOutcomeResponse(rw, http.StatusInternalServerError, err)
		log.Printf("handleDeviceLocationUpdate error Could not read body for device %s", deviceID)
		return
	}

	// Unmarshal the body
	var dl DeviceLocation
	err = json.Unmarshal(body, &dl)
	if err != nil {
		sendApiOutcomeResponse(rw, http.StatusInternalServerError, err)
		log.Printf("handleDeviceLocationUpdate error Could not unmarshal body for device %s", deviceID)
		return
	}
	log.Printf("handleDeviceLocationUpdate Device received from the request: %+v\n device: %s", dl, deviceID)

	// Get the device from cache
	key := "device:" + deviceID
	redisInstance := &cache.RedisInstance{RedisSession: common.RedisClient}
	data, err := redisInstance.GetCachedData(key)
	var d device.Device
	if err == nil {
		json.Unmarshal(data, &d)
	} else {
		sendApiOutcomeResponse(rw, http.StatusInternalServerError, errors.New("Wrong device ID"))
		log.Printf("handleDeviceLocationUpdate error device %s not found", deviceID)
		return
	}

	// Set new location values
	if strings.TrimSpace(dl.CountryCode) == "" && strings.TrimSpace(dl.PostalCode) == "" && strings.TrimSpace(dl.ACWKey) == "" && strings.TrimSpace(dl.CityOrPostalCode) == "" {
		d.Geo.Anonymous = true
	} else if strings.TrimSpace(dl.CountryCode) == "" || (strings.TrimSpace(dl.PostalCode) == "" && strings.TrimSpace(dl.CityOrPostalCode) == "") || strings.TrimSpace(dl.ACWKey) == "" {
		sendApiOutcomeResponse(rw, http.StatusInternalServerError, errors.New("Wrong values provided, must me either all null(anonymous) or all not null"))
		return
	} else {
		d.Geo.Anonymous = false
	}

	d.Geo.CountryCode = dl.CountryCode
	if strings.TrimSpace(dl.CityOrPostalCode) != "" {
		location, _ := weather_api.GetLocationFromPC(dl.ACWKey)
		if location.Type == weather_api.LocationTypeCity {
			d.Geo.City = dl.CityOrPostalCode
			d.Geo.Zip = ""
		} else {
			d.Geo.Zip = dl.CityOrPostalCode
			d.Geo.City = ""
		}
	} else {
		d.Geo.Zip = dl.PostalCode
		d.Geo.City = ""
	}

	// Reset the geo refresh count because the user location has changed
	d.GeoRefreshCount = 0

	d.Geo.ACWKey = dl.ACWKey
	// Save the device information into the cache

	dataBytes, _ := json.Marshal(d)
	err = redisInstance.SaveRedisData(dataBytes, key, 0)
	if err != nil {
		sendApiOutcomeResponse(rw, http.StatusInternalServerError, err)
		return
	}

	sendApiOutcomeResponse(rw, http.StatusOK, nil)

}


//----------------------------------------------
// @actionDisplayCheckPage
// [GET] /
//----------------------------------------------
func actionDisplayCheckPage(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, "OK")
}

//==============================================
// Functions - Admin Api Actions
//==============================================

//----------------------------------------------
// @actionAdminGetForecastData
// [GET] /api/v1.1/forecast/admin/id/{id}
//----------------------------------------------
func actionAdminGetForecastData(rw http.ResponseWriter, r *http.Request) {
	setResponseHeaders(&rw, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	res := AdminResponse{}

	// Api Args
	vars := mux.Vars(r)
	deviceID := vars["id"]
	details := strings.TrimSpace(r.FormValue("details"))

	// Load device from Redis
	display, _, err := getDevice(deviceID)
	if err != nil {
		io.WriteString(rw, "device not found")
		return
	}

	res.Zip = display.Geo.Zip
	res.City = display.Geo.City

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}
	res.Location = location.Key


	// Update LastUpdate Field
	res.LastUpdated, _ = forecastLastUpdatedString(display, location)

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))

	// Grab Forecast
	if details == "true" {
		res.Forecast = location.GetWeatherForecast(display.Category, display.ID, "BASIC", firmwareVersion)
	}

	json.NewEncoder(rw).Encode(res)
}

//----------------------------------------------
// @actionAdminGetForecastDataVer2
// [GET] /api/v2.0/forecast/admin/id/{id}
//----------------------------------------------
func actionAdminGetForecastDataVer2(rw http.ResponseWriter, r *http.Request) {
	setResponseHeaders(&rw, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	res := AdminResponse{}

	// Api Args
	vars := mux.Vars(r)
	deviceID := vars["id"]
	details := strings.TrimSpace(r.FormValue("details"))

	// Load device from Redis
	display, _, err := getDevice(deviceID)
	if err != nil {
		io.WriteString(rw, "device not found")
		return
	}

	res.Zip = display.Geo.Zip
	res.City = display.Geo.City

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}
	res.Location = location.Key


	// Update LastUpdate Field
	res.LastUpdated, _ = forecastLastUpdatedString(display, location)

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))

	// Grab Forecast
	if details == "true" {
		// Get Extended Info
		res.Forecast = location.GetWeatherForecastV2(display.Category, display.ID, firmwareVersion)
	}

	json.NewEncoder(rw).Encode(res)
}

//----------------------------------------------
// @actionAdminGetForecastDataJson
// [GET] /api/v2.2/forecast/admin/id/{id}
//----------------------------------------------
func actionAdminGetForecastDataJson(rw http.ResponseWriter, r *http.Request) {
	setResponseHeaders(&rw, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	(rw).Header().Set("Content-Type", "application/json")

	res := AdminResponseV2{}

	// Pending
	// scriptOverride := strings.TrimSpace(r.FormValue("s"))
	// timeStamp := strings.TrimSpace(r.FormValue("t"))
	firmwareVersion := strings.TrimSpace(r.FormValue("fw"))
	callSubVersion := strings.TrimSpace(r.FormValue("v"))
	if callSubVersion == "" {
		callSubVersion = "3"
	}
	i8nSet := strings.TrimSpace(r.FormValue("i8nV"))

	// Api Args
	vars := mux.Vars(r)
	deviceID := vars["id"]
	details := strings.TrimSpace(r.FormValue("details"))
	version := strings.TrimSpace(r.FormValue("version"))

	// Load device from Redis
	display, _, err := getDevice(deviceID)
	if err != nil {
		io.WriteString(rw, "device not found")
		return
	}

	res.Geo = display.Geo

	// Load Location Information
	var location weather_api.PostalCodeResponse
	location, err = getDeviceLocation(display)
	if err!= nil {
		log.Printf("Location %s,%s,%s not found, device : %s", display.Geo.ACWKey, display.Geo.Zip, display.Geo.CountryCode, deviceID)
		log.Printf(err.Error())
		log.Printf("Device: %+v\n", display)
		return
	}

	res.Location = location
	// Update LastUpdate Field
	_, res.LastUpdated = forecastLastUpdatedString(display, location)

	// Grab Forecast
	if details == "true" {
		//json, _ := location.NullableGetWeatherForecastJson(display.Category, display.ID, "BASIC").JsonResponse("1.2")
		res.Forecast = location.NullableGetWeatherForecastJson(display.Category, display.ID, firmwareVersion, callSubVersion)
	}

	if version == "" {
		version = "1.1"
		if i8nSet == "2" {
			version = version + "e"
		}
		j, _ := res.JsonResponse(version)
		io.WriteString(rw, j)
	} else {
		if i8nSet == "2" {
			version = version + "e"
		}
		j, _ := res.JsonResponse(version)
		io.WriteString(rw, j)
	}
}

//----------------------------------------------
// @actionAdminUpdateDeviceLocation
// [PUT,POST] /api/v1.1/forecast/admin/location/device/{device_id}
//----------------------------------------------
func actionAdminUpdateDeviceLocation(rw http.ResponseWriter, r *http.Request) {
	setResponseHeaders(&rw, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(r)
	deviceID := strings.TrimSpace(vars["device_id"])

	if isTokenValid(r) == false {
		sendApiOutcomeResponse(rw, http.StatusUnauthorized, errors.New("Wrong bearer token"))
		log.Printf("handleDeviceLocationUpdate error wrong barreer token for device %s", deviceID)
		return
	}

	err := deviceLocationUpdate(r, deviceID)

	if err != nil {
		sendApiOutcomeResponse(rw, http.StatusInternalServerError, err)
	} else {
		sendApiOutcomeResponse(rw, http.StatusOK, nil)
	}
}

//----------------------------------------------
// @actionAdminGetCategoryRanges
// [GET] /api/v1.1/forecast/admin/getRanges/WeatherService/{cat_type}
//----------------------------------------------
func actionAdminGetCategoryRanges(rw http.ResponseWriter, r *http.Request) {
	setResponseHeaders(&rw, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if isTokenValid(r) == false {
		sendApiOutcomeResponse(rw, http.StatusUnauthorized, errors.New("Wrong bearer token"))
		return
	}

	vars := mux.Vars(r)
	cat_type := strings.TrimSpace(vars["cat_type"])

	catList := device.LoadCategoryConf("/conf/categories.json")
	ranges := catList.GetDeviceCatRanges(cat_type)

	json.NewEncoder(rw).Encode(ranges)

}

//==============================================
// Functions - Support
//==============================================

//----------------------------------------------
// @tapGatewayGeo
//----------------------------------------------
func tapGatewayGeo(deviceID string) {

	geoRefreshLink := "https://ingv2.lacrossetechnology.com/api/v1.1/gateways/" + deviceID + "/displays/" + deviceID + "/geo"

	resp, err := http.Get(geoRefreshLink)
	if err != nil {
		log.Printf("Geo API Refresh Error: %+v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Geo API Refresh Error: %+v\n", err)
	}

	log.Printf("Geo API Refresh Response: %s\n", string(body))

}

//----------------------------------------------
// @sendApiOutcomeResponse
//----------------------------------------------
func sendApiOutcomeResponse(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	res := ReplyMessage{}
	res.Code = code
	if err != nil {
		res.Message = err.Error()
		log.Printf("Error: %+v\n", err)
	} else {
		res.Message = "Success"
	}
	json.NewEncoder(w).Encode(res)
}

//----------------------------------------------
// @deviceLocationUpdate
//----------------------------------------------
// This function should also be used by the handleDeviceLocationUpdate call
func deviceLocationUpdate(r *http.Request, deviceID string) error {

	// Get the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("handleDeviceLocationUpdate error Could not read body for device %s", deviceID)
		return err
	}

	// Unmarshal the body
	var dl DeviceLocation
	err = json.Unmarshal(body, &dl)
	if err != nil {
		log.Printf("handleDeviceLocationUpdate error Could not unmarshal body for device %s", deviceID)
		return err
	}
	log.Printf("handleDeviceLocationUpdate Device received from the request: %+v\n device: %s", dl, deviceID)

	// Get the device from cache
	key := "device:" + deviceID
	redisInstance := &cache.RedisInstance{RedisSession: common.RedisClient}
	data, err := redisInstance.GetCachedData(key)
	var d device.Device
	if err == nil {
		json.Unmarshal(data, &d)
	} else {
		log.Printf("handleDeviceLocationUpdate error device %s not found", deviceID)
		return err
	}

	// Set new location values
	if strings.TrimSpace(dl.CountryCode) == "" && strings.TrimSpace(dl.PostalCode) == "" && strings.TrimSpace(dl.ACWKey) == "" && strings.TrimSpace(dl.CityOrPostalCode) == "" {
		d.Geo.Anonymous = true
	} else if strings.TrimSpace(dl.CountryCode) == "" || (strings.TrimSpace(dl.PostalCode) == "" && strings.TrimSpace(dl.CityOrPostalCode) == "") || strings.TrimSpace(dl.ACWKey) == "" {
		return errors.New("Wrong values provided, must me either all null(anonymous) or all not null")
	} else {
		d.Geo.Anonymous = false
	}

	d.Geo.CountryCode = dl.CountryCode
	if strings.TrimSpace(dl.CityOrPostalCode) != "" {
		location, _ := weather_api.GetLocationFromPC(dl.ACWKey)
		if location.Type == weather_api.LocationTypeCity {
			d.Geo.City = dl.CityOrPostalCode
			d.Geo.Zip = ""
		} else {
			d.Geo.Zip = dl.CityOrPostalCode
			d.Geo.City = ""
		}
	} else {
		d.Geo.Zip = dl.PostalCode
		d.Geo.City = ""
	}

	// Reset the geo refresh count because the user location has changed
	d.GeoRefreshCount = 0

	d.Geo.ACWKey = dl.ACWKey
	// Save the device information into the cache

	dataBytes, _ := json.Marshal(d)
	err = redisInstance.SaveRedisData(dataBytes, key, 0)
	if err != nil {
		return err
	}

	return nil
}

//----------------------------------------------
// @hmacCheck
//----------------------------------------------
func hmacCheck(r *http.Request, psk string) error {

	// Build full request URL
	u := r.URL
	uri := u.RequestURI()
	url := httpScheme + "://" + httpHost + uri

	token := r.Header.Get("x-hmac-token")

	message := "[GET] " + url + "\n----- body -----\n"

	hmac := computeHmac(message, psk)
	//log.Printf("computed %s", hmac)

	// Checking the hmac agains't the provided token
	if hmac != token {
		log.Printf("Hmac token wrong : %s %s", token, url)
		err := errors.New("Wrong hmac token")
		return err
	}

	return nil
}

//----------------------------------------------
// @isTokenValid
// Firebase token validation
// We are expecting the token to be passed to Authorization header, on the format : "Barrear <Token>"
//----------------------------------------------
func isTokenValid(r *http.Request) bool {

	valid := false

	tokenStr := r.Header.Get("Authorization")
	userAgent := r.Header.Get("User-Agent")

	tokenArr := strings.Split(tokenStr, " ")
	if len(tokenArr) <= 1 {
		log.Printf("Malformed authorization header")
		return false
	}
	token := strings.TrimSpace(tokenArr[1])

	_, err := firebaseClient.VerifyIDToken(common.CTX, token)
	if err != nil {
		log.Printf("Token ERROR %s not valid, error : %s user-agent : %s", token, err.Error(), userAgent)
	} else {
		log.Printf("Token Validation succeeded, user-agent : %s", userAgent)
		valid = true
	}

	return valid

}

//----------------------------------------------
// @computeHmac
//----------------------------------------------
func computeHmac(message string, secret string) string {

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//----------------------------------------------
// @setResponseHeaders
// CORS validation
//----------------------------------------------
func setResponseHeaders(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

//----------------------------------------------
//
//----------------------------------------------
func supportedDevice(deviceId string) bool {
	if strings.ToUpper(deviceId) == "2CF272" {
		//Management requested exception case
		return false
	}
	return true
}

//----------------------------------------------
//
//----------------------------------------------
func syncElixirBackend(display device.Device) {
	// Tap Elixir Backend to ensure synchronization
	if display.GeoRefreshCount >= 0 && display.GeoRefreshCount < weather_api.GEO_REFRESH_API_HIT_AMOUNT {
		// Run the query

		go tapGatewayGeo(display.ID)
		display.GeoRefreshCount++
		updateDevice(display)
	}
}

//----------------------------------------------
//
//----------------------------------------------
func updateDevice(display device.Device) error {
	key := "device:" + display.ID
	dataBytes, _ := json.Marshal(display)
	redisInstance := &cache.RedisInstance{RedisSession: common.RedisClient}
	return redisInstance.SaveRedisData(dataBytes, key, 0)
}

//----------------------------------------------
//
//----------------------------------------------
func testOverride(deviceId string) bool {
	if deviceId == "2CFB3D" || deviceId == "6334DE" || deviceId == "63A252" || deviceId == "2D937D" || deviceId == "2DB9F8" || deviceId == "2D480E" || deviceId == "602750" {
		return true
	} else {
		return false
	}
}

//----------------------------------------------
//
//----------------------------------------------
func getDevice(deviceId string) (device.Device, bool, error) {
	var d device.Device
	key := "device:" + deviceId
	redisInstance := &cache.RedisInstance{RedisSession: common.RedisClient}
	data, err := redisInstance.GetCachedData(key)

	// No data found from the cache
	if err != nil {
		return d, false, err
	} else {
		json.Unmarshal(data, &d)
		anonymous := false

		// Fix for Elixir/Appengine use of 'USA' isntead of 'US'
		if strings.Contains(d.Geo.CountryCode, "USA") {
			d.Geo.CountryCode = "US"
		}

		if (strings.TrimSpace(d.Geo.Zip) == "" && strings.TrimSpace(d.Geo.City) == "") || d.Geo.Anonymous == true {
			anonymous = true
		}

		return d, anonymous, nil
	}
}

//----------------------------------------------
//
//----------------------------------------------
func getDeviceLocation(display device.Device) (weather_api.PostalCodeResponse, error) {
	if display.Geo.ACWKey != "" {
		return weather_api.GetLocationFromPC(display.Geo.ACWKey)
	} else {
		return weather_api.GetLocation(strings.TrimSpace(display.Geo.Zip), strings.TrimSpace(display.Geo.CountryCode))
	}
}

//----------------------------------------------
//
//----------------------------------------------
func updateDeviceRequestEntry(deviceId string) {
	redisInstance := &cache.RedisInstance{RedisSession: common.RedisClient}
	nowStr := time.Now().Format("02:01:2006 15:04:05")
	redisInstance.SaveRedisData([]byte(nowStr), "devicerequested:"+deviceId, 0)
}

func forecastLastUpdatedString(display device.Device, location weather_api.PostalCodeResponse) (string, LastUpdated) {
	redisInstance := &cache.RedisInstance{RedisSession: common.RedisClient}
	var res = ""
	var lu LastUpdated

	switch display.Category {
	case device.CAT1:
		updateDateKey := "forecastupdate:1day:" + location.Key
		data, _ := redisInstance.GetCachedData(updateDateKey)
		res = string(data)
		lu.OneDayForecast = null.NewString(string(data), true)
	case device.CAT2:
		updateDateKey := "forecastupdate:current:" + location.Key
		data, _ := redisInstance.GetCachedData(updateDateKey)
		res = "currentcondition api : " + string(data)
		lu.CurrentForecast = null.NewString(string(data), true)

		updateDateKey = "forecastupdate:1day:" + location.Key
		data, _ = redisInstance.GetCachedData(updateDateKey)
		lu.OneDayForecast = null.NewString(string(data), true)
		res = res + " ;1day api : " + string(data)
	case device.CAT3:
		updateDateKey := "forecastupdate:10day:" + location.Key
		data, _ := redisInstance.GetCachedData(updateDateKey)
		res = "10day api : " + string(data)
		lu.TenDayForecast = null.NewString(string(data), true)

		updateDateKey = "forecastupdate:24hour:" + location.Key
		data, _ = redisInstance.GetCachedData(updateDateKey)
		lu.TwentyFourHourForecast = null.NewString(string(data), true)
		res = res + " ;24hour api : " + string(data)
	}
	return res, lu
}
