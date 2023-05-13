package device
//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
	"../../../common"
	"../../cache"
	"cloud.google.com/go/datastore"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"time"
)

//----------------------------------------------
// Globals
//----------------------------------------------
var (
	Categories *CategoryConfList
)

//----------------------------------------------
// Types
//----------------------------------------------
type (
	CategoryConf struct {
		Category string          `bson:"category" json:"category"`
		Ranges   []CategoryRange `bson:"ranges" json:"ranges"`
	}

	CategoryRange struct {
		Min      string `bson:"min" json:"min"`
		Max      string `bson:"max" json:"max"`
		MinFloat float32
		MaxFloat float32
	}

	CategoryConfList struct {
		CategoryConf []CategoryConf
	}

	Device struct {
		ID              string
		PSK             string
		Category        string
		GeoRefreshCount int // This variable control the amount of times we called Keith Geo refresh api
		Geo             Geo
	}

	TimeLoopSettings struct {
		Enabled bool
		Mode int
		LoopOffset int64
		LoopStart int64
		LoopEnd int64
	}

	TimeCompressionSettings struct {
		Enabled bool
		AccelerationRate float64 // e.g. 1 to 1, 60 to 1, .25 to 1
		StartTime int64 // initial time we use to calculate the current accelerated value.
		TimeOffset int64 // add/subtract x seconds to the clock
	}

	TimeZoneOverrideSettings struct {
		Enabled bool
		Sign int
		HourOffset int
		MinuteOffset int
	}

	ForecastScriptingSettings struct {
		Enabled bool
		Mode int
	}

	ExtendedDeviceInfo struct {
		ID string
		DataScript int
		TimeZoneOverride TimeZoneOverrideSettings
		TimeLoop TimeLoopSettings
		TimeCompression TimeCompressionSettings
		ForecastScripting ForecastScriptingSettings
		HasDateTimeBug bool
		Attributes map[string]int64
	}



	RawSensorEntity struct {
		Serial                string            `datastore:"serial"`
		Manufacturer          int               `datastore:"manufacturer"`
		Geo                   Geo               `datastore:"geo"`
		Handle                string            `datastore:"handle"`
		Weight                int               `datastore:"weight,omitempty,flatten"`
		VerificationCode      string            `datastore:"verificationCode"`
		Series                string            `datastore:"series"`
		SensorTypeEntityKey   *datastore.Key    `datastore:"sensorTypeEntityKey,omitempty,flatten"`
		SchemaVersion         int               `datastore:"schemaVersion"`
		OnDisplayCheckCache   string            `datastore:"onDisplayCheckCache"`
		ModifiedOn            time.Time         `datastore:"modifiedOn"`
		Longitude             int               `datastore:"longitude,omitempty,flatten"`
		LinkedSensors         *datastore.Entity `datastore:"linkedSensors,omitempty,flatten"`
		Latitude              int               `datastore:"latitude,omitempty,flatten"`
		LastSynchedVNext      time.Time         `datastore:"lastSynchedVNext"`
		LastSynched           time.Time         `datastore:"lastSynched"`
		Kind                  int               `datastore:"kind,omitempty,flatten"`
		InternalAttributes    int               `datastore:"internalAttributes,omitempty,flatten"`
		FlaggedForSynchVNext  bool              `datastore:"flaggedForSynchVNext"`
		FlaggedForSynch       bool              `datastore:"flaggedForSynch"`
		FlaggedForDeleteVNext bool              `datastore:"flaggedForDeleteVNext"`
		FlaggedForDelete      bool              `datastore:"flaggedForDelete"`
		Fields                interface{}       `datastore:"fields,omitempty,flatten"`
		CreatedOn             time.Time         `datastore:"createdOn"`
		ClaimToken            int               `datastore:"claimToken,omitempty,flatten"`
		ClaimCacheKey         int               `datastore:"claimCacheKey,omitempty,flatten"`
		ClaimCacheId          int               `datastore:"claimCacheId"`
		Batch                 int               `datastore:"batch"`
		Attributes            datastore.PropertyList `datastore:"attributes"`
		//Attributes            map[string]string `datastore:"attributes,flatten"`
	}

	Geo struct {
		Zip         string `datastore:"zip" bson:"zip" json:"zip,postal_code"`
		CountryCode string `datastore:"countryCode" bson:"countryCode" json:"countryCode"`
		Timezone    string `datastore:"timezone" bson:"timezone" json:"timezone"`
		Latitude    string `datastore:"latitude" bson:"latitude" json:"latitude"`
		Anonymous   bool   `datastore:"anonymous" bson:"anonymous" json:"anonymous"`
		Longitude   string `datastore:"longitude" bson:"longitude" json:"longitude"`
		ACWKey      string `datastore:"acw_key" bson:"acw_key" json:"acw_key"`
		City        string `datastore:"city" bson:"city" json:"city"`
	}

	DevicePubSub struct {
		Serial      string `bson:"serial" json:"serial"`
		Anonymous   bool   `bson:"anonymous" json:"anonymous"`
		Zip         string `bson:"zip" json:"zip"`
		Timezone    string `bson:"timezone" json:"timezone"`
		CountryCode string `bson:"countryCode" json:"countryCode"`
		Longitude   string `bson:"longitude" json:"longitude"`
		Latitude    string `bson:"latitude" json:"latitude"`
	}
)

//----------------------------------------------
// Exports
//----------------------------------------------
func LoadCategoryConf(filename string) *CategoryConfList {

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	var cc []CategoryConf
	json.Unmarshal(raw, &cc)

	for i := range cc {
		for j := range cc[i].Ranges {
			cc[i].Ranges[j].MinFloat = hexToFloat(cc[i].Ranges[j].Min)
			cc[i].Ranges[j].MaxFloat = hexToFloat(cc[i].Ranges[j].Max)
		}
	}

	return &CategoryConfList{CategoryConf: cc}
}

func (catConf *CategoryConfList) GetDeviceCat(deviceID string) string {

	for _, catStruct := range catConf.CategoryConf {
		for _, cRange := range catStruct.Ranges {
			if hexToFloat(deviceID) >= cRange.MinFloat && hexToFloat(deviceID) <= cRange.MaxFloat {
				return catStruct.Category
			}
		}
	}

	return CAT1

}

func (catConf *CategoryConfList) GetDeviceCatRanges(cat string) string {

	ranges := ""

	for _, catStruct := range catConf.CategoryConf {

		if cat != catStruct.Category {
			continue
		} else {
			for _, cRange := range catStruct.Ranges {
				if ranges != "" {
					ranges = ranges + ";"
				}
				ranges = ranges + cRange.Min + "::" + cRange.Max
			}
		}
	}

	return ranges

}

func GetDeviceCat(deviceID string) string {

	f := hexToFloat(deviceID)
	if f >= hexToFloat("280001") && f <= hexToFloat("2B3DE2") {
		return CAT1
	} else if f >= hexToFloat("2CF269") && f <= hexToFloat("2CF2CC") {
		return CAT2
	} else if f >= hexToFloat("2CF2CD") && f <= hexToFloat("2CF330") {
		return CAT3
	} else {
		return ""
	}

}

func GetExtendedDeviceInfo(ID string) (ExtendedDeviceInfo, error) {
	key := "device.attributes:" + cache.CACHE_BUST__GLOBAL + cache.CACHE_BUST__DEVICE_INFO + ID
	raw, err := common.RedisInstance.GetCachedData(key)
	if (err != nil) {
		return RefreshExtendedInfo(ID, nil)
	}
	var extendedInfo ExtendedDeviceInfo
	json.Unmarshal(raw, &extendedInfo)
	return extendedInfo, nil
}

func RefreshExtendedInfo(ID string, raw *RawSensorEntity) (ExtendedDeviceInfo, error) {
	log.Printf("[Device] Attributes Refresh Processing ID=%s", ID)
	extendedInfo := ExtendedDeviceInfo{}
	extendedInfo.ID = ID
	extendedInfo.Attributes = map[string]int64{}

	if raw == nil {
		var x RawSensorEntity
		q := datastore.NewQuery("SensorEntity").Filter("serial =", ID).Limit(1)
		it := common.DataStoreClient.Run(common.CTX, q)
		_, _ = it.Next(&x)
		raw = &x
	}

	if (raw.Serial != ID) {
		log.Printf("Unexpected Raw Data: %v", raw)
	}
	//extendedInfo.ID = raw.Serial

	if (raw.Attributes != nil) {
		for k := range raw.Attributes {
			av, _ := raw.Attributes[k].Value.(int64)
			extendedInfo.Attributes[raw.Attributes[k].Name] = av
		}
	}

	var v int64
	v = 0
	var ok bool

	//----------
	// DateTimeBug Settings
	//----------
	extendedInfo.HasDateTimeBug = (extendedInfo.Attributes["date-time-bug"] == 1)

	//----------
	// TimeCompression Settings
	//----------
	extendedInfo.TimeCompression.Enabled = false
	extendedInfo.TimeCompression.AccelerationRate = 1.0
	extendedInfo.TimeCompression.StartTime = time.Date(2019, 1,1,0,0,0,0, time.UTC).Unix()
	extendedInfo.TimeCompression.TimeOffset = 0

	if v, ok = extendedInfo.Attributes["time-compress"]; ok {
		if (v != 0) {
			extendedInfo.TimeCompression.Enabled = true
		}
		if (v > 0) {
			extendedInfo.TimeCompression.AccelerationRate = float64(v)
		} else if (v < 0) {
			extendedInfo.TimeCompression.AccelerationRate = (1.0 / float64(v))
		}
	}
	if v, ok = extendedInfo.Attributes["time-compress.start"]; ok {
		extendedInfo.TimeCompression.StartTime = v
	}
	if v, ok = extendedInfo.Attributes["time-compress.offset"]; ok {
		if (v != 0) {
			extendedInfo.TimeCompression.Enabled = true
			extendedInfo.TimeCompression.TimeOffset = v
		}
	}

	//----------
	// TimeZoneOverride Settings
	//----------
	extendedInfo.TimeZoneOverride.Enabled = false
	extendedInfo.TimeZoneOverride.Sign = 1
	extendedInfo.TimeZoneOverride.HourOffset = 0
	extendedInfo.TimeZoneOverride.MinuteOffset = 0
	if v, ok = extendedInfo.Attributes["time-zone-override"]; ok {
		if (v != 0) {
			extendedInfo.TimeZoneOverride.Enabled = true
			if (v < 0) {
				extendedInfo.TimeZoneOverride.Sign = -1
				v = -v
			}
			extendedInfo.TimeZoneOverride.MinuteOffset = int(v % 100)
			extendedInfo.TimeZoneOverride.HourOffset = int((v - (v % 100)) / 100)
		}
	}

	//----------
	// TimeLoop Settings
	//----------
	extendedInfo.TimeLoop.Enabled = false
	extendedInfo.TimeLoop.LoopOffset = 0
	extendedInfo.TimeLoop.LoopStart = 0
	extendedInfo.TimeLoop.LoopEnd = 0

	if v, ok = extendedInfo.Attributes["time-loop"]; ok {
		vm := int(v)
		if (vm > 0 && vm <= TimeLoopBackendDriver) {
			extendedInfo.TimeLoop.Mode = vm
			extendedInfo.TimeLoop.Enabled = true
			if (vm == TimeLoopCustom) {
				if v, ok = extendedInfo.Attributes["time-loop.start"]; ok {
					extendedInfo.TimeLoop.LoopStart = v
				}

				if v, ok = extendedInfo.Attributes["time-loop.end"]; ok {
					extendedInfo.TimeLoop.LoopEnd = v
				}

				if (extendedInfo.TimeLoop.LoopEnd  < extendedInfo.TimeLoop.LoopStart) {
					s := extendedInfo.TimeLoop.LoopEnd
					extendedInfo.TimeLoop.LoopEnd  = extendedInfo.TimeLoop.LoopStart
					extendedInfo.TimeLoop.LoopStart = s
				}
			}
		} else if vm != 0 {
			log.Printf("[Device] Invalid TimeLoop Mode %d", v)
		}

		if (extendedInfo.TimeLoop.Enabled) {
			if v, ok = extendedInfo.Attributes["time-loop.offset"]; ok {
				extendedInfo.TimeLoop.LoopOffset = v
			}
		}
	}


	//----------
	// ForecastScripting Settings
	//----------
	extendedInfo.ForecastScripting.Enabled = false
	extendedInfo.ForecastScripting.Mode = ForecastScriptDisabled
	if v, ok = extendedInfo.Attributes["forecast-script"]; ok {
		if (v >= ForecastScriptStaticA && v <= ForecastScriptBackendDriver) {
			extendedInfo.ForecastScripting.Mode = int(v)
			extendedInfo.ForecastScripting.Enabled = true
		} else if (v != 0) {
			log.Printf("[Device] Invalid Forecast Script Mode %d", int(v))
		}
	}

	//----------
	// Persist
	//---------

	key := "device.attributes:" + cache.CACHE_BUST__GLOBAL + cache.CACHE_BUST__DEVICE_INFO + ID
	dataBytes, err := json.Marshal(extendedInfo)
	if err != nil {
		log.Printf("[Device] Marshall ExtendedInfo Failed: %s, %o", key, err)
	} else {
		log.Printf("[Device] Attributes Persist (%s) %v", key, extendedInfo)
		common.RedisInstance.SaveRedisData(dataBytes, key, 0)
	}
	return extendedInfo, err
}

//----------------------------------------------
// Local Funcs
//----------------------------------------------
func hexToFloat(h string) float32 {

	un, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}

	un32 := uint32(un)

	f := math.Float32frombits(un32)
	return f

}
