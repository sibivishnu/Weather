package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	//	"fmt"
	//	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	QueryResponse struct {
		Duration  float64
		DeviceID  string
		Iteration int
	}
)

var (
	deviceList    []string
	responseMap   = make(map[string][]QueryResponse)
	lock          = sync.RWMutex{}
	queryInterval = 120
)

const (
	ITERATION_COUNT = 2
)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	// Initialize Jobs with one device per Job
	file, err := os.Open("/tmp/devices.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		deviceList = append(deviceList, scanner.Text())
	}

	for i := 0; i < ITERATION_COUNT; i++ {
		//log.Printf("Iteration : %d", i)
		var wg sync.WaitGroup
		for _, device := range deviceList {
			wg.Add(1)
			dArr := strings.Split(device, ",")
			if len(dArr) < 2 {
				log.Printf("Can't extract ID and PSK for %s", device)
				return
			}

			go processDevice(&wg, dArr[0], dArr[1], i)
		}
		wg.Wait()
		if i+1 == ITERATION_COUNT {
			break
		}
		time.Sleep(time.Duration(queryInterval) * time.Second)
	}

	// Display final report
	//for _, reportArr := range responseMap {
	//	for _, report := range reportArr {
	//		fmt.Printf("Device : %s Query duration : %f Iteration : %d \n", report.DeviceID, report.Duration, report.Iteration)
	//	}
	//}
}

func processDevice(wg *sync.WaitGroup, deviceID string, psk string, it int) {
	defer wg.Done()
	// Sleep for a random time
	waitTime := random(1, 60)
	//log.Printf("Starting the execution with Device %s Iteration: %d", deviceID, it)
	//log.Printf("Wait time for Device %s : %d", deviceID, waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)

	// Run the request
	queryResponse := queryWeatherService(deviceID, psk, it)
	queryResponse.Iteration = it
	//log.Printf("%s", deviceID)
	//fmt.Printf("%s,%f,%d", queryResponse.DeviceID, queryResponse.Duration, queryResponse.Iteration)
	//writeToMap(deviceID, queryResponse)
}

func writeToMap(device string, response QueryResponse) {
	lock.Lock()
	lock.RLock()
	defer lock.Unlock()
	defer lock.RUnlock()
	if _, ok := responseMap[device]; ok {
		responseMap[device] = append(responseMap[device], response)
	} else {
		var respArr []QueryResponse
		respArr = append(respArr, response)
		responseMap[device] = respArr
	}
}

func queryWeatherService(deviceID string, psk string, it int) QueryResponse {

	//log.Printf("Exe starting for %s", deviceID)

	client := &http.Client{}
	url := "https://ingv2.lacrossetechnology.com/api/v1.1/forecast/id/" + deviceID
	//url := "https://ingv2.lacrossetechnology.com/api/v2.0/forecast/id/" + deviceID
	//url := "https://ingv2.lacrossetechnology.com/api/v2.0/forecast/test/id/" + deviceID
	req, _ := http.NewRequest("GET", url, nil)

	message := "[GET] " + url + "\n----- body -----\n"
	hm := computeHmac(message, psk)
	log.Printf("%s, %s", deviceID, hm)
	req.Header.Add("x-hmac-token", hm)

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return QueryResponse{}
	}
	defer resp.Body.Close()
	//_, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return QueryResponse{}
	//}
	elapsed := time.Since(start).Seconds()

	//log.Printf("%s", deviceID)
	//qr := QueryResponse{}
	//qr.Duration = elapsed
	//qr.DeviceID = deviceID

	//fmt.Printf("%s,%f \n", deviceID, elapsed)
	log.Printf("Device: %s Duration: %f Iteration: %d", deviceID, elapsed, it)
	return QueryResponse{}

}

func computeHmac(message string, secret string) string {

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
