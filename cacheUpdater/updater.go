package main
//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"../common"
	"../common/const/device"
	"../common/providers/weather_api"
	"bufio"
	"bytes"
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

//----------------------------------------------
// Exports
//----------------------------------------------

//----------------------------------------------
// Local Funcs
//----------------------------------------------
func copyFile() {
	if scpServerHost == "" {
		log.Printf("NO SCP Settings Provided, Bypassing SCP Remote File Step.")
	} else {
		log.Printf("COPYING DEVICE FILE FROM REMOTE SERVER : scpServerRSA: %s, scpServerUser: %s, scpServerHost: %s, remoteFilePath: %s, targetFolder: %s", scpServerRSA, scpServerUser, scpServerHost, remoteFilePath, targetFolder)

		var out_cmd, err_cmd bytes.Buffer
		cmd := exec.Command("scp", "-o", "StrictHostKeyChecking=no", "-i", scpServerRSA, scpServerUser+"@"+scpServerHost+":"+remoteFilePath, targetFolder)
		cmd.Stdout = &out_cmd
		cmd.Stderr = &err_cmd

		err := cmd.Run()

		if err != nil {
			fmt.Println("Scp Command error :", err_cmd.String())
		}
	}
}

func runForecastUpdater() {
	log.Printf("forecast runing")
	locationListMap, _ := common.RedisInstance.QueryCache("activelocations*")


	for locKey, locVal := range locationListMap {
		keyArr := strings.Split(locKey, ":")
		if len(keyArr) < 2 {
			log.Printf("Wrong location key %s", locKey)
			continue
		}

		valArr := strings.Split(locVal, ":")
		if len(valArr) < 2 {
			log.Printf("Wrong location value %s", locVal)
			continue
		}

		rawLocKey := keyArr[1]
		timeZone := valArr[0]
		devCategory := valArr[1]

		//accuWeather.SetLocalDateAndHour(timeZone)
		weatherTime, err := weather_api.GetLocalDateAndHour(timeZone)
		if err != nil {
			log.Printf("Could not get time for the timeZone %s", timeZone)
			continue
		}

		log.Printf("Forecast refreshing for Category : %s, Accuweather key : %s, Timezone : %s, Local Date : %s, Local Time : %s, Hour range : %s", devCategory, rawLocKey, timeZone, weatherTime.LocalDate, weatherTime.LocalTime, weatherTime.HourRange)
		switch devCategory {
		case device.CAT2, device.CAT1:
			weather_api.QueryAccuDayForecastAPI(rawLocKey, timeZone, "1day", weatherTime)
		case device.CAT3:
			weather_api.QueryAccuHourForecastAPI(rawLocKey, "24hour", weatherTime)
			weather_api.QueryAccuDayForecastAPI(rawLocKey, timeZone, "10day", weatherTime)
		}

	}

}

func dstUpdateProcess() {

	// Zip code
	runDSTUpdater("zip:*")

	// Postal Code
	runDSTUpdater("postalcode:*")

}

func runDSTUpdater(filter string) {

	log.Printf("forecast runing")
	locationListMap, _ := common.RedisInstance.QueryCache(filter)
	now := time.Now()

	for locKey, locVal := range locationListMap {
		var pc weather_api.PostalCodeResponse
		err := json.Unmarshal([]byte(locVal), &pc)

		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, pc.TimeZone.NextOffsetChange)
		if err != nil {
			//log.Printf("runDSTUpdater : Could not parse time %s for key %s", pc.TimeZone.NextOffsetChange, locKey)
			continue
		}

		if t.Before(now) {
			log.Printf("DST reached for key : %s", locKey)
			common.RedisInstance.RemoveKeyFromCache(locKey)
		}

	}

}

func runDeviceGeoRefreshUpdater() {

	log.Printf("Updating geo refresh count")
	deviceListMap, _ := common.RedisInstance.QueryCache("devicerequested:*")

	for deviceKey, _ := range deviceListMap {
		keyArr := strings.Split(deviceKey, ":")
		if len(keyArr) < 2 {
			log.Printf("Wrong device requested record key %s", deviceKey)
			continue
		}

		deviceID := keyArr[1]
		key := "device:" + deviceID

		data, err := common.RedisInstance.GetCachedData(key)

		// No data found from the cache
		if err != nil {
			log.Printf("Wrong device requested record key %s", deviceKey)
			continue
		}

		var d device.Device
		json.Unmarshal(data, &d)

		d.GeoRefreshCount = 0
		dataBytes, _ := json.Marshal(d)
		err = common.RedisInstance.SaveRedisData(dataBytes, key, 0)
		if err != nil {
			log.Printf("Wrong device requested record key %s", deviceKey)
			continue
		}

	}

}

func listenGeo() {
	log.Println("[Listen] Geo: " + subscriptionName + " @ " + topicName)
	client, err := pubsub.NewClient(common.CTX, projectID)
	if err != nil {
		panic(err)
	}

	/*
	log.Println(projectID)
	log.Println(subscriptionName)
	log.Println(topicName)
	*/

	sub := client.Subscription(subscriptionName)

	// We check if the subscription exists
	ok, err := sub.Exists(common.CTX)
	if err != nil {
		log.Printf("[Listen] Geo: Exception %v", err)
		panic(err)
	}

	topic := client.Topic(topicName)
	//attribute := client.Topic(attributeTopic)

	// If the subscription doesn't exist, then we need to create one
	if !ok {
		log.Println("[Listen] Geo: Subscription does not exist. Creating one")
		sub, err = client.CreateSubscription(common.CTX, subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
		if (err != nil) {
			log.Printf("[Listen] Geo: Exception %v", err)
			panic(err)
		}

	}

	// Receive messages
	log.Println("[Listen] Geo: Begining listening on the queue")
	err = sub.Receive(common.CTX, func(ctx context.Context, m *pubsub.Message) {

		log.Println("[Listen] Geo: New pub/sub request")
		var dev device.Device
		var dps device.DevicePubSub

		// We need to extract the device UUID from the payload
		json.Unmarshal(m.Data, &dps)
		key := "device:" + dps.Serial
		log.Printf("[Listen] Geo: Serial : %s Zip : %s", dps.Serial, dps.Zip)

		deviceByte, _ := common.RedisInstance.GetCachedData(key)
		json.Unmarshal(deviceByte, &dev)

		if dev.Geo.ACWKey != "" {
			log.Printf("[Listen] Geo:We already have an ACW Key %s in the cache for device %s, The pub/sub received will be ignored", dev.Geo.ACWKey, dev.ID)
		} else {

			dev.Geo.Zip = dps.Zip
			dev.Geo.CountryCode = dps.CountryCode
			dev.Geo.Timezone = dps.Timezone
			dev.Geo.Anonymous = dps.Anonymous
			dev.Geo.Latitude = dps.Latitude

			newDeviceByte, _ := json.Marshal(dev)
			common.RedisInstance.SaveRedisData(newDeviceByte, key, 0)
			log.Println("[Listen] Geo: Device Data updated successfully")
		}

		m.Ack()
	})

	// If context don't get canceled and program exit then something is wrong
	if err != context.Canceled {
		panic(err)
	}
}

func listenAttr() {
	log.Println("[Listen] Attribute: " + attributeSubscription + " @ " + attributeTopic)
	client, err := pubsub.NewClient(common.CTX, projectID)
	if err != nil {
		panic(err)
	}
	sub := client.Subscription(attributeSubscription)

	// We check if the subscription exists
	ok, err := sub.Exists(common.CTX)
	if err != nil {
		log.Printf("[Listen] Attribute Exception %v", err)
		panic(err)
	}

	topic := client.Topic(attributeTopic)

	// If the subscription doesn't exist, then we need to create one
	if !ok {
		log.Println("[Listen] Attribute: Subscription does not exist. Creating one")
		sub, err = client.CreateSubscription(common.CTX, attributeSubscription, pubsub.SubscriptionConfig{Topic: topic})
		if (err != nil) {
			log.Printf("[Listen] Attribute: failed to create sub %v", err)
			panic(err)
		}
	}

	// Receive messages
	log.Println("[Listen] Attribute: Receiving")
	err = sub.Receive(common.CTX, func(ctx context.Context, m *pubsub.Message) {
		log.Println("[Listen] Attribute: New pub/sub request")
		var dps device.DevicePubSub

		// We need to extract the device UUID from the payload
		json.Unmarshal(m.Data, &dps)
		device.RefreshExtendedInfo(dps.Serial, nil)
		m.Ack()
	})

	// If context don't get canceled and program exit then something is wrong
	if err != context.Canceled {
		panic(err)
	}
}

func listenDevicePubSub() {
	listenAttr()
	listenGeo()
}

func runCacheIDUpdater() {

	// Copy the file from SCP Server
	copyFile()

	file, err := os.Open(devicesFile)
	if err != nil {
		fmt.Println("Error opening the file :", err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//var wg sync.WaitGroup

	log.Printf("Cache update process started")

	for scanner.Scan() {
		job := &Job{}
		job.Line = scanner.Text()
		JobQueue <- *job
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Cache update process completed")
	//wg.Wait()

	// Once all goroutines have finished processing, empty the file
	// TODO : Delete the file from the scp server using ssh commands
}

func cacheID(line string) {
	lineArr := strings.Split(line, ",")

	if len(lineArr) < 1 {
		log.Printf("Can't extract ID and PSK for %s", line)
		return
	}

	dev := device.Device{}
	dev.ID = lineArr[0]
	dev.PSK = lineArr[1]
	log.Printf("Device: %s", lineArr[0])

	// Query device Geo
	//log.Printf("Datastore query for %s started", dev.ID)
	q := datastore.NewQuery("SensorEntity").Filter("serial =", dev.ID).Limit(1)
	it := common.DataStoreClient.Run(common.CTX, q)
	var x device.RawSensorEntity
	_, _ = it.Next(&x)

	/*
		if err != nil {
			log.Printf("Error fetching next: %v", err)
		}
	*/

	// Set the Geo we got from datastore
	dev.Geo = x.Geo

	// Set the right category
	dev.Category = device.Categories.GetDeviceCat(dev.ID)
	//log.Printf(dev.ID)
	//log.Printf(dev.Category)
	if (dev.Category == device.CAT2 || dev.Category == device.CAT3) && strings.TrimSpace(dev.Geo.Zip) == "" {
		dev.Geo.Zip = "17036"
	}

	dataBytes, _ := json.Marshal(dev)
	key := "device:" + dev.ID

	data, err := common.RedisInstance.GetCachedData(key)
	var d device.Device
	// No data found from the cache
	if err == nil {
		json.Unmarshal(data, &d)
	}

	if d.Geo.Zip != dev.Geo.Zip {
		nowStr := time.Now().Format("02:01:2006")
		mismatch_key := "geomismatch:" + dev.ID + ":" + nowStr
		val := "Old : " + d.Geo.Zip + " New :" + dev.Geo.Zip
		common.RedisInstance.SaveRedisData([]byte(val), mismatch_key, 0)
	}
	common.RedisInstance.SaveRedisData(dataBytes, key, 0)

	// Load Device Attribute Information
	device.RefreshExtendedInfo(dev.ID, &x)
	return
}
