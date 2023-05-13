package nws
//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//----------------------------------------------
// Contants
//----------------------------------------------
const (
	NWS_BASE_URL      = "https://graphical.weather.gov"
	NWS_BASE_URL_PATH = "/xml/sample_products/browser_interface/ndfdXMLclient.php"
)

//----------------------------------------------
// Types
//----------------------------------------------
type MainStruct struct {
	Data Data `xml:"data"`
}

type Data struct {
	Location   []Location  `xml:"location"`
	Parameters []Parameter `xml:"parameters"`
}

type Location struct {
	LocationKey string `xml:"location-key"`
	Coords      Coords `xml:"point"`
}

type Parameter struct {
	ConvectiveHazards []ConvectiveHazard `xml:"convective-hazard"`
}

type ConvectiveHazard struct {
	SevereComponent SevereComponent `xml:"severe-component"`
}

type SevereComponent struct {
	Type  string `xml:"type,attr"`
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type Coords struct {
	Lattitude string `xml:"latitude,attr"`
	Longitude string `xml:"longitude,attr"`
}

//----------------------------------------------
// Exports
//----------------------------------------------
func GetSevereComponentMap(zip string) map[string]string {

	// We create the map and initialize it with null values in case nothing is returned
	severeMap := make(map[string]string)
	severeMap["tornadoes"] = "0"
	severeMap["hail"] = "0"

	log.Printf("GetSevereComponentMap for zip: %s", zip)

	now := time.Now()
	hourAfter := now.Add(1 * time.Hour)
	beginTime := now.Format("2006-01-02T15")
	beginTime = beginTime + ":00:00"
	endTime := hourAfter.Format("2006-01-02T15")
	endTime = endTime + ":00:00"

	var Url *url.URL
	Url, _ = url.Parse(NWS_BASE_URL)

	Url.Path += NWS_BASE_URL_PATH
	parameters := url.Values{}
	parameters.Add("zipCodeList", zip)
	parameters.Add("product", "time-series")
	parameters.Add("begin", beginTime)
	parameters.Add("end", endTime)
	parameters.Add("ptornado", "ptornado")
	parameters.Add("phail", "phail")
	Url.RawQuery = parameters.Encode()

	resp, err := http.Get(Url.String())
	if err != nil {
		log.Printf("GetSevereComponentMap Error: %+v\n", err)
		return severeMap
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GetSevereComponentMap Error: %+v\n", err)
		return severeMap
	}

	var ms MainStruct
	xml.Unmarshal(body, &ms)

	if len(ms.Data.Parameters) > 0 {
		for i := 0; i < len(ms.Data.Parameters[0].ConvectiveHazards); i++ {
			log.Printf("GetSevereComponentMap for zip: %s type : %s", zip, ms.Data.Parameters[0].ConvectiveHazards[i].SevereComponent.Type)
			log.Printf("GetSevereComponentMap for zip: %s type : %s", zip, ms.Data.Parameters[0].ConvectiveHazards[i].SevereComponent.Value)
			severeMap[ms.Data.Parameters[0].ConvectiveHazards[i].SevereComponent.Type] = ms.Data.Parameters[0].ConvectiveHazards[i].SevereComponent.Value
		}
	} else {
		return severeMap
	}

	return severeMap
}
