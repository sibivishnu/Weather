package weather_api

//==============================================
// CopyRight 2019 La Crosse Technology, LTD.
//==============================================

//==============================================
// Imports
//==============================================
import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sibivishnu/Weather/common/const/device"
	"github.com/sibivishnu/Weather/common/nws"
	"gopkg.in/guregu/null.v3"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//================================================================
//================================================================
// Globals
//================================================================
//================================================================

//==============================================
// Globals - Tables
//==============================================

//==============================================
// Types
//==============================================

type (

	//==============================================
	// API
	//==============================================

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	ApiString string

	//----------------------------------------------
	//
	//----------------------------------------------

	//==============================================
	// Common
	//==============================================

	//----------------------------------------------
	// @PostalCodeResponse
	//----------------------------------------------
	/**
	 * @brief
	 */
	PostalCodeResponse struct {
		Version            int
		Key                string
		Type               string `bson:"Type" json:"Type"`
		Rank               int    `bson:"Rank" json:"Rank"`
		LocalizedName      string `bson:"LocalizedName" json:"LocalizedName"`
		EnglishName        string `bson:"EnglishName" json:"EnglishName"`
		PrimaryPostalCode  string `bson:"PrimaryPostalCode" json:"PrimaryPostalCode"`
		TimeZone           AccuTimeZone
		Region             AccuCountry
		AdministrativeArea AccuCountry
		Country            AccuCountry
		Code               string
	}

	//----------------------------------------------
	//
	//----------------------------------------------

	//----------------------------------------------
	// @AccuCountry
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuCountry struct {
		ID            string
		LocalizedName string
		EnglishName   string
	}

	//----------------------------------------------
	// @AccuTimeZone
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuTimeZone struct {
		Code             string
		Name             string
		GmtOffset        float64
		IsDaylightSaving bool
		NextOffsetChange string
	}

	//==============================================
	// Legacy Forecast
	//==============================================

	//----------------------------------------------
	// @WeatherTime
	//----------------------------------------------
	/**
	 * @brief
	 */
	WeatherTime struct {
		LocalTime string
		LocalDate string
		DayInfo   string
		HourRange string
		Iso8601   string
		DateTime  time.Time
	}

	//----------------------------------------------
	// @AccuHourlyForecastResponse
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuHourlyForecastResponse struct {
		DateTime                 string
		DateStr                  string
		TimeStr                  string
		MinTemp                  float64
		MaxTemp                  float64
		EpochDateTime            int
		WeatherIcon              int
		IconPhrase               string
		IsDaylight               bool
		Temperature              Temperature
		RealFeelTemperature      Temperature
		WetBulbTemperature       Temperature
		DewPoint                 Temperature
		Wind                     Wind
		WindGust                 Wind
		RelativeHumidity         int
		Visibility               Temperature
		Ceiling                  Temperature
		UVIndex                  int
		UVIndexText              string
		PrecipitationProbability int
		RainProbability          int
		SnowProbability          int
		IceProbability           int
		TotalLiquid              Temperature
		Rain                     Temperature
		Snow                     Temperature
		Ice                      Temperature
		CloudCover               int
	}

	//----------------------------------------------
	// @AccuCurrentForecastResponse
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuCurrentForecastResponse struct {
		LocalObservationDateTime string
		EpochTime                int
		WeatherText              string
		WeatherIcon              int
		IsDayTime                bool
		Temperature              CurrentTemp
	}

	//----------------------------------------------
	// @CurrentTemp
	//----------------------------------------------
	/**
	 * @brief
	 */
	CurrentTemp struct {
		Metric   Temperature
		Imperial Temperature
	}

	//----------------------------------------------
	// @Wind
	//----------------------------------------------
	/**
	 * @brief
	 */
	Wind struct {
		Speed     Temperature
		Direction Direction
	}

	//----------------------------------------------
	// @Direction
	//----------------------------------------------
	/**
	 * @brief
	 */
	Direction struct {
		Degrees   int
		Localized string
		English   string
	}

	//----------------------------------------------
	// @AccuHeadline
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuHeadline struct {
		EffectiveDate      string
		EffectiveEpochDate int
		Severity           int
		Text               string
		Category           string
		EndDate            string
		EndEpochDate       int
	}

	//----------------------------------------------
	// @Sun
	//----------------------------------------------
	/**
	 * @brief
	 */
	Sun struct {
		Rise      string
		EpochRise int64
		Set       string
		EpochSet  int64
		Phase     string
		Age       int
	}

	//----------------------------------------------
	// @AirAndPollen
	//----------------------------------------------
	/**
	 * @brief
	 */
	AirAndPollen struct {
		Name          string
		Value         int
		Category      string
		CategoryValue int
		Type          string
	}

	//----------------------------------------------
	// @DailyForecast
	//----------------------------------------------
	/**
	 * @brief
	 */
	DailyForecast struct {
		Headline       AccuHeadline
		DailyForecasts []AccuDailyForecast
	}

	//----------------------------------------------
	// @AccuTemplateCat1Struct
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuTemplateCat1Struct struct {
		DateStr       string
		TimeStr       string
		GmtOffset     float64
		Accu1d        *DailyForecast
		Headline      *AccuHeadline
		DailyForecast *AccuDailyForecast
		FlowControl   int
	}

	//----------------------------------------------
	// @AccuTemplateCat2Struct
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuTemplateCat2Struct struct {
		DateStr         string
		TimeStr         string
		GmtOffset       float64
		Accu1d          *DailyForecast
		Headline        *AccuHeadline
		DailyForecast   *AccuDailyForecast
		CurrentForecast *AccuCurrentForecastResponse
		FlowControl     int
	}

	//----------------------------------------------
	// @AccuTemplateCat3Struct
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuTemplateCat3Struct struct {
		DateStr string
		TimeStr string
		// Time Forecast was obtained and cached.
		ForecastTime  string
		GmtOffset     float64
		DailyForecast *AccuDailyForecast
		Accu7d        *[]AccuDailyForecast
		Accu24h       *[]AccuHourlyForecastResponse
		FlowControl   int
	}

	//----------------------------------------------
	// @MinMaxTemperature
	//----------------------------------------------
	/**
	 * @brief
	 */
	MinMaxTemperature struct {
		Minimum Temperature
		Maximum Temperature
	}

	//----------------------------------------------
	// @SummaryTemperature
	//----------------------------------------------
	/**
	 * @brief
	 */
	SummaryTemperature struct {
		Heating Temperature
		Cooling Temperature
	}

	//----------------------------------------------
	// @DayNightData
	//----------------------------------------------
	/**
	 * @brief
	 */
	DayNightData struct {
		Icon                     int
		IconPhrase               string
		ShortPhrase              string
		LongPhrase               string
		PrecipitationProbability int
		ThunderstormProbability  int
		RainProbability          int
		SnowProbability          int
		IceProbability           int
		HoursOfPrecipitation     int
		HoursOfRain              int
		HoursOfSnow              int
		HoursOfIce               int
		CloudCover               int
		Wind                     Wind
		WindGust                 Wind
		TotalLiquid              Temperature
		Rain                     Temperature
		Snow                     Temperature
		Ice                      Temperature
	}

	//----------------------------------------------
	// @LocationLookupDict
	//----------------------------------------------
	/**
	 * @brief
	 */
	LocationLookupDict struct {
		Locations []LocationLookupResponse `json:"locations"`
	}

	//----------------------------------------------
	// @LocationLookupResponse
	//----------------------------------------------
	/**
	 * @brief
	 */
	LocationLookupResponse struct {
		CityName           string `json:"city_name"`
		CountryCode        string `json:"country_code"`
		AcwKey             string `json:"acw_key"`
		AdministrativeArea string `json:"administrative_area"`
		TzName             string `json:"tz_name"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	Temperature struct {
		Value      float64
		ValueRound int
		Unit       string
		UnitType   int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuDailyForecast struct {
		Date                     string
		EpochDate                int64
		Sun                      Sun
		Moon                     Sun
		Temperature              MinMaxTemperature
		RealFeelTemperature      MinMaxTemperature
		RealFeelTemperatureShade MinMaxTemperature
		HoursOfSun               float32
		DegreeDaySummary         SummaryTemperature
		AirAndPollen             []AirAndPollen
		AirAndPollenMap          map[string]int
		AirAndPollenCategoryMap  map[string]string
		Day                      DayNightData
		Night                    DayNightData
		Actual                   *DayNightData
		NWSevereComponentMap     map[string]string
	}

	//==============================================
	// Nullable Versions
	//==============================================

	//----------------------------------------------
	// @NullableWeatherTime
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableWeatherTime struct {
		LocalTime null.String
		LocalDate null.String
		DayInfo   null.String
		HourRange null.String
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuHourlyForecast struct {
		DateTime null.String
		//DateStr                  null.String
		//TimeStr                  null.String
		//MinTemp                  null.Float
		//MaxTemp                  null.Float
		EpochDateTime            null.Int
		WeatherIcon              null.Int
		IconPhrase               null.String
		IsDaylight               null.Bool
		Temperature              NullableTemperature
		RealFeelTemperature      NullableTemperature
		WetBulbTemperature       NullableTemperature
		DewPoint                 NullableTemperature
		Wind                     NullableWind
		WindGust                 NullableWind
		RelativeHumidity         null.Int
		Visibility               NullableReading
		Ceiling                  NullableReading
		UVIndex                  null.Int
		UVIndexText              null.String
		PrecipitationProbability null.Int
		RainProbability          null.Int
		SnowProbability          null.Int
		IceProbability           null.Int
		TotalLiquid              NullableReading
		Rain                     NullableReading
		Snow                     NullableReading
		Ice                      NullableReading
		CloudCover               null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuCurrentForecastResponse struct {
		LocalObservationDateTime null.String
		EpochTime                null.Int
		WeatherText              null.String
		WeatherIcon              null.Int
		IsDayTime                null.Bool
		Temperature              NullableCurrentTemp
		TornadoProbability       null.Int
		HailProbability          null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableCurrentTemp struct {
		Metric   NullableTemperature
		Imperial NullableTemperature
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableTemperature struct {
		Value      null.Float
		ValueRound null.Int
		Unit       null.String
		UnitType   null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableReading struct {
		Value      null.Float
		ValueRound null.Int
		Unit       null.String
		UnitType   null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableSpeed struct {
		Value      null.Float
		ValueRound null.Int
		Unit       null.String
		UnitType   null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableWind struct {
		Speed     NullableSpeed
		Direction NullableDirection
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableDirection struct {
		Degrees   null.Int
		Localized null.String
		English   null.String
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuHeadline struct {
		EffectiveDate      null.String
		EffectiveEpochDate null.Int
		Severity           null.Int
		Text               null.String
		Category           null.String
		EndDate            null.String
		EndEpochDate       null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableSun struct {
		Rise      null.String
		EpochRise null.Int
		Set       null.String
		EpochSet  null.Int
		Phase     null.String
		Age       null.Int
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAirAndPollen struct {
		Name          null.String
		Value         null.Int
		Category      null.String
		CategoryValue null.Int
		Type          null.String
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableDailyForecast struct {
		Headline       NullableAccuHeadline
		Today          *NullableDayNightData
		DailyForecasts []NullableAccuDailyForecast
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableUniversalForecast struct {
		Date               null.String
		Time               null.String
		Category           null.Int
		GmtOffset          null.Float
		ForecastTime       null.String
		Headline           *NullableAccuHeadline
		Today              *NullableDayNightData
		Current            *NullableAccuCurrentForecastResponse
		Daily              *[]NullableAccuDailyForecast
		Hourly             *[]NullableAccuHourlyForecast
		NWSForecast        map[string]string
		FlowControl        null.Int
		ExtendedDeviceInfo device.ExtendedDeviceInfo
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableMinMaxTemperature struct {
		Minimum NullableTemperature
		Maximum NullableTemperature
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableSummaryTemperature struct {
		Heating NullableTemperature
		Cooling NullableTemperature
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableUniversalForecastFormatV1p3 struct {
		Date         null.String // ymd in gateway locale
		Time         null.String // hms in gateway locale
		GmtOffset    null.Float  // offset in gateway locale
		ForecastTime null.String // ymd hms in gateway locale time (time forecast was fetched from accuweather)
		Category     null.Int
		FlowControl  null.Int
		//Headline       *NullableAccuHeadline
		Today       ApiResponseInterface
		Current     ApiResponseInterface
		Daily       []ApiResponseInterface
		Hourly      []ApiResponseInterface
		NWSForecast map[string]string
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableUniversalForecastFormatV1p2 struct {
		Date         null.String
		Time         null.String
		GmtOffset    null.Float
		ForecastTime null.String
		Category     null.Int
		FlowControl  null.Int
		//Headline       *NullableAccuHeadline
		Today   ApiResponseInterface
		Current ApiResponseInterface
		Daily   []ApiResponseInterface
		Hourly  []ApiResponseInterface
		//Daily *[]NullableAccuDailyForecast
		//Hourly *[]NullableAccuHourlyForecast
		NWSForecast map[string]string
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuDailyForecast struct {
		Date                     null.String
		EpochDate                null.Int
		Sun                      NullableSun
		Moon                     NullableSun
		Temperature              NullableMinMaxTemperature
		RealFeelTemperature      NullableMinMaxTemperature
		RealFeelTemperatureShade NullableMinMaxTemperature
		HoursOfSun               null.Float
		DegreeDaySummary         NullableSummaryTemperature
		AirAndPollen             []NullableAirAndPollen
		AirAndPollenMap          map[string]int
		AirAndPollenCategoryMap  map[string]string
		Day                      NullableDayNightData
		Night                    NullableDayNightData
		//Actual                   *NullableDayNightData // Todo deprecated
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableDayNightData struct {
		Icon                     null.Int
		IconPhrase               null.String
		ShortPhrase              null.String
		LongPhrase               null.String
		PrecipitationProbability null.Int
		ThunderstormProbability  null.Int
		RainProbability          null.Int
		SnowProbability          null.Int
		IceProbability           null.Int
		HoursOfPrecipitation     null.Float
		HoursOfRain              null.Float
		HoursOfSnow              null.Float
		HoursOfIce               null.Float
		CloudCover               null.Int
		Wind                     NullableWind
		WindGust                 NullableWind
		TotalLiquid              NullableReading
		Rain                     NullableReading
		Snow                     NullableReading
		Ice                      NullableReading
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuDailyForecastFormatV1p3 struct {
		Date      null.String `json:"D"`
		EpochDate null.Int    `json:"U"`
		MoonPhase null.Int    `json:"MP"`
		SunRise   null.String `json:"Sr"`
		SunSet    null.String `json:"Ss"`
		MoonRise  null.String `json:"Mr"`
		MoonSet   null.String `json:"Ms"`

		MinTemperature null.Float `json:"Tl"`
		MaxTemperature null.Float `json:"Th"`
		//MinRealFeelTemperature   null.Float `json:"Fl"`
		//MaxRealFeelTemperature   null.Float `json:"Fh"`
		//MinRealFeelTemperatureShade null.Float `json:"FSl"`
		//MaxRealFeelTemperatureShade null.Float `json:"FSh"`
		HoursOfSun null.Float `json:"HoS"`
		//DaySummaryHeating        null.Float `json:"DsH"`
		//DaySummaryCooling        null.Float `json:"DsC"`

		// AirAndPollen
		UVIndex            null.Int `json:"UVi"`
		UVCategory         null.Int `json:"UVc"`
		AirQualityCategory null.Int `json:"AQc"`
		GrassCategory      null.Int `json:"Gc"`
		MoldCategory       null.Int `json:"Mc"`
		RagweedCategory    null.Int `json:"Rc"`
		TreeCategory       null.Int `json:"Tc"`

		Day   ApiResponseInterface
		Night ApiResponseInterface
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuHourlyForecastFormatV1p3 struct {
		DateTime null.String `json:"DT"`
		//DateStr                  null.String
		//TimeStr                  null.String
		//MinTemp                  null.Float
		//MaxTemp                  null.Float
		EpochDateTime null.Int   `json:"U"`
		WeatherIcon   null.Int   `json:"WX"`
		IsDaylight    null.Bool  `json:"isDL"`
		Temperature   null.Float `json:"T"`
		//RealFeelTemperature      null.Float `json:"F"`
		//WetBulbTemperature       null.Float `json:"WB"`
		//DewPoint                 null.Float `json:"DP"`
		WindSpeed   null.Float `json:"WS"`
		WindHeading null.Int   `json:"WH"`
		WindGust    null.Float `json:"GS"`
		//WindGustHeading          null.Int `json:"GH"`
		//RelativeHumidity         null.Int `json:"RHu"`
		//Visibility               null.Float `json:"V"`
		//Ceiling                  null.Float `json:"C"`
		//UVIndex                  null.Int `json:"UVi"`
		//UVCategory               null.Int `json:"UVc"`
		PrecipitationProbability null.Int `json:"Pp"`
		//RainProbability          null.Int `json:"Rp"`
		//SnowProbability          null.Int `json:"Sp"`
		//IceProbability           null.Int `json:"Ip"`
		//TotalLiquid              null.Float `json:"TLiq"`
		//Rain                     null.Float `json:"R"`
		//Snow                     null.Float `json:"S"`
		//Ice                      null.Float `json:"I"`
		//CloudCover               null.Int `json:"CC"`

	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuCurrentForecastResponseFormatV1p3 struct {
		Epoch              null.Int   `json:"U"`
		WeatherIcon        null.Int   `json:"WX"`
		IsDayTime          null.Bool  `json:"isDT"`
		Temperature        null.Float `json:"T"`
		TornadoProbability null.Int   `json:"TNp"`
		HailProbability    null.Int   `json:"Hp"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableDayNightDataFormatV1p3 struct {
		Icon                     null.Int `json:"WX"`
		PrecipitationProbability null.Int `json:"Pp"`
		ThunderstormProbability  null.Int `json:"Tp"`
		RainProbability          null.Int `json:"Rp"`
		SnowProbability          null.Int `json:"Sp"`
		IceProbability           null.Int `json:"Ip"`
		//HoursOfPrecipitation     null.Float `json:"Ph"`
		//HoursOfRain              null.Float `json:"Rh"`
		//HoursOfSnow              null.Float `json:"Sh"`
		//HoursOfIce               null.Float `json:"Ih"`
		CloudCover  null.Int   `json:"CC"`
		WindSpeed   null.Float `json:"WS"`
		WindHeading null.Int   `json:"WH"`
		WindGust    null.Float `json:"GS"`
		//TotalLiquid              null.Float `json:"TLiq"`
		Rain null.Float `json:"R"`
		Snow null.Float `json:"S"`
		Ice  null.Float `json:"I"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuDailyForecastFormatV1p2 struct {
		Date      null.String `json:"D"`
		EpochDate null.Int    `json:"U"`
		MoonPhase null.Int    `json:"MP"`
		SunRise   null.String `json:"Sr"`
		SunSet    null.String `json:"Ss"`
		MoonRise  null.String `json:"Mr"`
		MoonSet   null.String `json:"Ms"`

		MinTemperature              null.Float `json:"Tl"`
		MaxTemperature              null.Float `json:"Th"`
		MinRealFeelTemperature      null.Float `json:"Fl"`
		MaxRealFeelTemperature      null.Float `json:"Fh"`
		MinRealFeelTemperatureShade null.Float `json:"FSl"`
		MaxRealFeelTemperatureShade null.Float `json:"FSh"`
		HoursOfSun                  null.Float `json:"HoS"`
		DaySummaryHeating           null.Float `json:"DsH"`
		DaySummaryCooling           null.Float `json:"DsC"`

		// AirAndPollen
		UVIndex            null.Int `json:"UVi"`
		UVCategory         null.Int `json:"UVc"`
		AirQualityCategory null.Int `json:"AQc"`
		GrassCategory      null.Int `json:"Gc"`
		MoldCategory       null.Int `json:"Mc"`
		RagweedCategory    null.Int `json:"Rc"`
		TreeCategory       null.Int `json:"Tc"`

		Day   ApiResponseInterface
		Night ApiResponseInterface
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuHourlyForecastFormatV1p2 struct {
		DateTime null.String `json:"DT"`
		//DateStr                  null.String
		//TimeStr                  null.String
		//MinTemp                  null.Float
		//MaxTemp                  null.Float
		EpochDateTime            null.Int   `json:"U"`
		WeatherIcon              null.Int   `json:"WX"`
		IsDaylight               null.Bool  `json:"isDL"`
		Temperature              null.Float `json:"T"`
		RealFeelTemperature      null.Float `json:"F"`
		WetBulbTemperature       null.Float `json:"WB"`
		DewPoint                 null.Float `json:"DP"`
		WindSpeed                null.Float `json:"WS"`
		WindHeading              null.Int   `json:"WH"`
		WindGust                 null.Float `json:"GS"`
		WindGustHeading          null.Int   `json:"GH"`
		RelativeHumidity         null.Int   `json:"RHu"`
		Visibility               null.Float `json:"V"`
		Ceiling                  null.Float `json:"C"`
		UVIndex                  null.Int   `json:"UVi"`
		UVCategory               null.Int   `json:"UVc"`
		PrecipitationProbability null.Int   `json:"Pp"`
		RainProbability          null.Int   `json:"Rp"`
		SnowProbability          null.Int   `json:"Sp"`
		IceProbability           null.Int   `json:"Ip"`
		TotalLiquid              null.Float `json:"TLiq"`
		Rain                     null.Float `json:"R"`
		Snow                     null.Float `json:"S"`
		Ice                      null.Float `json:"I"`
		CloudCover               null.Int   `json:"CC"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableAccuCurrentForecastResponseFormatV1p2 struct {
		Epoch              null.Int   `json:"U"`
		WeatherIcon        null.Int   `json:"WX"`
		IsDayTime          null.Bool  `json:"isDT"`
		Temperature        null.Float `json:"T"`
		TornadoProbability null.Int   `json:"TNp"`
		HailProbability    null.Int   `json:"Hp"`

		// Deprecated
		TornadoProbabilityDeprecated null.Int `json:"tornadoes"`
		// Deprecated
		HailProbabilityDeprecated null.Int `json:"hail"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableDayNightDataFormatV1p2 struct {
		Icon                     null.Int   `json:"WX"`
		PrecipitationProbability null.Int   `json:"Pp"`
		ThunderstormProbability  null.Int   `json:"Tp"`
		RainProbability          null.Int   `json:"Rp"`
		SnowProbability          null.Int   `json:"Sp"`
		IceProbability           null.Int   `json:"Ip"`
		HoursOfPrecipitation     null.Float `json:"Ph"`
		HoursOfRain              null.Float `json:"Rh"`
		HoursOfSnow              null.Float `json:"Sh"`
		HoursOfIce               null.Float `json:"Ih"`
		CloudCover               null.Int   `json:"CC"`
		WindSpeed                null.Float `json:"WS"`
		WindHeading              null.Int   `json:"WH"`
		WindGust                 null.Float `json:"GS"`
		TotalLiquid              null.Float `json:"TLiq"`
		Rain                     null.Float `json:"R"`
		Snow                     null.Float `json:"S"`
		Ice                      null.Float `json:"I"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableLocationLookupDict struct {
		Locations []NullableLocationLookupResponse `json:"locations"`
	}

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	NullableLocationLookupResponse struct {
		CityName           null.String `json:"city_name"`
		CountryCode        null.String `json:"country_code"`
		AcwKey             null.String `json:"acw_key"`
		AdministrativeArea null.String `json:"administrative_area"`
		TzName             null.String `json:"tz_name"`
	}

	//==============================================
	// JSON Facades
	//==============================================

	//----------------------------------------------
	//
	//----------------------------------------------
	/**
	 * @brief
	 */
	ForecastResponseVersionOne struct {
		Date            null.String
		Time            null.String
		GmtOffset       null.Float
		FlowControl     null.Int
		CurrentForecast NullableAccuCurrentForecastResponse
		DailyForecasts  []NullableAccuDailyForecast
		HourlyForecasts []NullableAccuHourlyForecast
	}
)

//==============================================
// Protocols - MarshalJSON
//==============================================

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func AccuIconDetails(icon null.Int) AccuIcon {
	if icon.Valid {
		if val, ok := AccuweatherIcons[int(icon.Int64)]; ok {
			return val
		}
	}
	return AccuweatherIcons[0]
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func DisplayIconByAccuIcon(icon null.Int) null.Int {
	a := AccuIconDetails(icon)
	if a.IconNumber != 0 {
		return null.NewInt(int64(a.DisplayIcon), true)
	} else {
		return null.NewInt(0, false)
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (v NullableDirection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Degrees)
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (v NullableTemperature) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (v NullableReading) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Value)
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (v NullableSpeed) MarshalJSON() ([]byte, error) {
	// Note previous implementation converted Metric KM/h to MPH
	//var t null.Int
	//if v.Value.Valid {
	//	t.SetValid(int64(Round(v.Value.Float64 / 3.6)))
	//}
	return json.Marshal(v.Value)
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (v NullableCurrentTemp) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Metric.Value)
}

//==============================================
// Protocols - JsonResponse
//==============================================
/**
 * @brief
 */
func VersionedJson(s ApiResponseInterface, version string) (string, error) {
	r, err := s.ResponseFormat(version)
	if err != nil {
		return "null", err
	} else {
		var j, err = json.Marshal(r)
		if err != nil {
			return "null", err
		} else {
			return string(j), err
			//var out bytes.Buffer
			//json.Indent(&out, j, "", "  ")
			//return string(out.Bytes()), err
		}
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s ApiString) JsonResponse(version string) (string, error) {
	return string(s), nil
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableAccuHourlyForecast) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuHourlyForecastFormatV1p2) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuHourlyForecastFormatV1p3) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuDailyForecast) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuDailyForecastFormatV1p2) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuDailyForecastFormatV1p3) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableDayNightData) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableDayNightDataFormatV1p2) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableDayNightDataFormatV1p3) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuCurrentForecastResponse) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuCurrentForecastResponseFormatV1p2) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableAccuCurrentForecastResponseFormatV1p3) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableUniversalForecast) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableUniversalForecastFormatV1p2) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

/**
 * @brief
 */
func (s NullableUniversalForecastFormatV1p3) JsonResponse(version string) (string, error) {
	return VersionedJson(s, version)
}

//==============================================
// Protocols - ResponseFormat
//==============================================
/**
 * @brief
 */
func getFormattedDateFromEpoch(v null.String) null.String {
	if v.Valid {
		layout := "2006-01-02T15:04:05-07:00"
		t, err := time.Parse(layout, v.String)
		if err != nil {
			log.Printf("GetWeatherForecastTest - Error parsing the time : %s %s", v.String, err.Error())
			var r null.String
			return r
		} else {
			return null.NewString(t.Format("15:04"), true)
		}
	} else {
		return v
	}
}

/**
 * @brief
 */
func (s ApiString) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] ApiString")
	return s, nil
}

/**
 * @brief
 */
func (s NullableAccuHourlyForecast) ResponseFormat(version string) (ApiResponseInterface, error) {

	var weatherCategoryToI8NCode = WeatherCategoryToFirmwareMap

	switch version {
	case "1.1e":
		fallthrough
	case "1.1":
		return s, nil
	case "1.2e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.2":
		var r NullableAccuHourlyForecastFormatV1p2
		r.DateTime = s.DateTime
		//r.DateStr = s.DateStr
		//r.TimeStr = s.TimeStr
		//r.MinTemp = s.MinTemp
		//r.MaxTemp = s.MaxTemp
		r.EpochDateTime = s.EpochDateTime

		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)

		r.IsDaylight = s.IsDaylight
		r.Temperature = s.Temperature.Value
		r.RealFeelTemperature = s.RealFeelTemperature.Value
		r.WetBulbTemperature = s.WetBulbTemperature.Value
		r.DewPoint = s.DewPoint.Value
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		r.WindGustHeading = s.WindGust.Direction.Degrees
		r.RelativeHumidity = s.RelativeHumidity
		r.Visibility = s.Visibility.Value
		r.Ceiling = s.Ceiling.Value
		r.UVIndex = s.UVIndex
		if s.UVIndexText.Valid {
			r.UVCategory = null.NewInt(int64(ValueToEnumMap[s.UVIndexText.String]), true)
		}

		r.PrecipitationProbability = s.PrecipitationProbability
		r.RainProbability = s.RainProbability
		r.SnowProbability = s.SnowProbability
		r.IceProbability = s.IceProbability
		r.TotalLiquid = s.TotalLiquid.Value
		r.Rain = s.Rain.Value
		r.Snow = s.Snow.Value
		r.Ice = s.Ice.Value
		r.CloudCover = s.CloudCover
		return r, nil

	case "1.3e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.3":
		var r NullableAccuHourlyForecastFormatV1p2
		r.DateTime = s.DateTime
		//r.DateStr = s.DateStr
		//r.TimeStr = s.TimeStr
		//r.MinTemp = s.MinTemp
		//r.MaxTemp = s.MaxTemp
		r.EpochDateTime = s.EpochDateTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)

		r.IsDaylight = s.IsDaylight
		r.Temperature = s.Temperature.Value
		r.RealFeelTemperature = s.RealFeelTemperature.Value
		r.WetBulbTemperature = s.WetBulbTemperature.Value
		r.DewPoint = s.DewPoint.Value
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		r.WindGustHeading = s.WindGust.Direction.Degrees
		r.RelativeHumidity = s.RelativeHumidity
		r.Visibility = s.Visibility.Value
		r.Ceiling = s.Ceiling.Value
		r.UVIndex = s.UVIndex
		if s.UVIndexText.Valid {
			r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[s.UVIndexText.String]), true)
		}

		r.PrecipitationProbability = s.PrecipitationProbability
		r.RainProbability = s.RainProbability
		r.SnowProbability = s.SnowProbability
		r.IceProbability = s.IceProbability
		r.TotalLiquid = s.TotalLiquid.Value
		r.Rain = s.Rain.Value
		r.Snow = s.Snow.Value
		r.Ice = s.Ice.Value
		r.CloudCover = s.CloudCover
		return r, nil
	case "1.4e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.4":
		var r NullableAccuHourlyForecastFormatV1p3
		r.DateTime = s.DateTime
		r.EpochDateTime = s.EpochDateTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)

		r.IsDaylight = s.IsDaylight
		r.Temperature = s.Temperature.Value
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		//r.WindGustHeading = s.WindGust.Direction.Degrees

		// r.RelativeHumidity = s.RelativeHumidity
		// r.Visibility = s.Visibility.Value
		// r.UVIndex = s.UVIndex
		// r.RealFeelTemperature = s.RealFeelTemperature.Value
		// if s.UVIndexText.Valid {
		// 	r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[[s.UVIndexText.String]), true)
		// }

		r.PrecipitationProbability = s.PrecipitationProbability
		// r.RainProbability = s.RainProbability
		// r.SnowProbability = s.SnowProbability
		// r.IceProbability = s.IceProbability
		// r.TotalLiquid = s.TotalLiquid.Value
		// r.Rain = s.Rain.Value
		// r.Snow = s.Snow.Value
		// r.Ice = s.Ice.Value
		// r.CloudCover = s.CloudCover
		return r, nil

	case "1.5e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.5":
		var r NullableAccuHourlyForecastFormatV1p3
		r.DateTime = s.DateTime
		r.EpochDateTime = s.EpochDateTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)

		r.IsDaylight = s.IsDaylight
		r.Temperature = s.Temperature.Value
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		//r.WindGustHeading = s.WindGust.Direction.Degrees

		// r.RelativeHumidity = s.RelativeHumidity
		// r.Visibility = s.Visibility.Value
		// r.UVIndex = s.UVIndex
		// r.RealFeelTemperature = s.RealFeelTemperature.Value

		// if s.UVIndexText.Valid {
		//	r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[Extended[s.UVIndexText.String]), true)
		// }

		r.PrecipitationProbability = s.PrecipitationProbability
		// r.RainProbability = s.RainProbability
		// r.SnowProbability = s.SnowProbability
		// r.IceProbability = s.IceProbability
		// r.TotalLiquid = s.TotalLiquid.Value
		// r.Rain = s.Rain.Value
		// r.Snow = s.Snow.Value
		// r.Ice = s.Ice.Value
		// r.CloudCover = s.CloudCover
		return r, nil

	default:
		return s, errors.New("unsupported version")
	}
}

/**
 * @brief
 */
func (s NullableAccuHourlyForecastFormatV1p3) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuHourlyForecastFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

/**
 * @brief
 */
func (s NullableAccuHourlyForecastFormatV1p2) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuHourlyForecastFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

/**
 * @brief
 */
func (s NullableAccuDailyForecast) ResponseFormat(version string) (ApiResponseInterface, error) {
	var weatherCategoryToI8NCode = WeatherCategoryToFirmwareMap

	switch version {
	case "1.1e":
		fallthrough
	case "1.1":
		return s, nil

	case "1.2e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.2":
		var r NullableAccuDailyForecastFormatV1p2
		r.Date = s.Date
		r.EpochDate = s.EpochDate

		if s.Moon.Phase.Valid {
			if val, ok := MoonPhaseToFirmwareMap[s.Moon.Phase.String]; ok {
				if ok {
					r.MoonPhase = null.NewInt(int64(val), true)
				}
			}
		}

		r.SunRise = getFormattedDateFromEpoch(s.Sun.Rise)
		r.SunSet = getFormattedDateFromEpoch(s.Sun.Set)
		r.MoonRise = getFormattedDateFromEpoch(s.Moon.Rise)
		r.MoonSet = getFormattedDateFromEpoch(s.Moon.Set)

		r.MinTemperature = s.Temperature.Minimum.Value
		r.MaxTemperature = s.Temperature.Maximum.Value
		r.MinRealFeelTemperature = s.RealFeelTemperature.Minimum.Value
		r.MaxRealFeelTemperature = s.RealFeelTemperature.Maximum.Value
		r.MinRealFeelTemperatureShade = s.RealFeelTemperatureShade.Minimum.Value
		r.MaxRealFeelTemperatureShade = s.RealFeelTemperatureShade.Maximum.Value
		r.HoursOfSun = s.HoursOfSun
		r.DaySummaryCooling = s.DegreeDaySummary.Cooling.Value
		r.DaySummaryHeating = s.DegreeDaySummary.Heating.Value

		if val, ok := s.AirAndPollenMap["UVIndex"]; ok {
			if ok {
				r.UVIndex = null.NewInt(int64(val), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["UVIndex"]; ok {
			if ok {
				r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["AirQuality"]; ok {
			if ok {
				r.AirQualityCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Grass"]; ok {
			if ok {
				r.GrassCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Mold"]; ok {
			if ok {
				r.MoldCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Ragweed"]; ok {
			if ok {
				r.RagweedCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Tree"]; ok {
			if ok {
				r.TreeCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}

		r.Day, _ = s.Day.ResponseFormat(version)
		r.Night, _ = s.Night.ResponseFormat(version)

		return r, nil

	case "1.3e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.3":
		var r NullableAccuDailyForecastFormatV1p2
		r.Date = s.Date
		r.EpochDate = s.EpochDate

		if s.Moon.Phase.Valid {
			if val, ok := MoonPhaseToFirmwareMap[s.Moon.Phase.String]; ok {
				if ok {
					r.MoonPhase = null.NewInt(int64(val), true)
				}
			}
		}

		r.SunRise = getFormattedDateFromEpoch(s.Sun.Rise)
		r.SunSet = getFormattedDateFromEpoch(s.Sun.Set)
		r.MoonRise = getFormattedDateFromEpoch(s.Moon.Rise)
		r.MoonSet = getFormattedDateFromEpoch(s.Moon.Set)

		r.MinTemperature = s.Temperature.Minimum.Value
		r.MaxTemperature = s.Temperature.Maximum.Value
		//r.MinRealFeelTemperature = s.RealFeelTemperature.Minimum.Value
		//r.MaxRealFeelTemperature = s.RealFeelTemperature.Maximum.Value
		//r.MinRealFeelTemperatureShade = s.RealFeelTemperatureShade.Minimum.Value
		//r.MaxRealFeelTemperatureShade = s.RealFeelTemperatureShade.Maximum.Value
		r.HoursOfSun = s.HoursOfSun
		r.DaySummaryCooling = s.DegreeDaySummary.Cooling.Value
		r.DaySummaryHeating = s.DegreeDaySummary.Heating.Value

		if val, ok := s.AirAndPollenMap["UVIndex"]; ok {
			if ok {
				r.UVIndex = null.NewInt(int64(val), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["UVIndex"]; ok {
			if ok {
				r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["AirQuality"]; ok {
			if ok {
				r.AirQualityCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Grass"]; ok {
			if ok {
				r.GrassCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Mold"]; ok {
			if ok {
				r.MoldCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Ragweed"]; ok {
			if ok {
				r.RagweedCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Tree"]; ok {
			if ok {
				r.TreeCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}

		r.Day, _ = s.Day.ResponseFormat(version)
		r.Night, _ = s.Night.ResponseFormat(version)

		return r, nil

	case "1.4e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.4":
		var r NullableAccuDailyForecastFormatV1p3
		r.Date = s.Date
		r.EpochDate = s.EpochDate

		if s.Moon.Phase.Valid {
			if val, ok := MoonPhaseToFirmwareMap[s.Moon.Phase.String]; ok {
				if ok {
					r.MoonPhase = null.NewInt(int64(val), true)
				}
			}
		}

		r.SunRise = getFormattedDateFromEpoch(s.Sun.Rise)
		r.SunSet = getFormattedDateFromEpoch(s.Sun.Set)
		r.MoonRise = getFormattedDateFromEpoch(s.Moon.Rise)
		r.MoonSet = getFormattedDateFromEpoch(s.Moon.Set)

		r.MinTemperature = s.Temperature.Minimum.Value
		r.MaxTemperature = s.Temperature.Maximum.Value
		r.HoursOfSun = s.HoursOfSun

		if val, ok := s.AirAndPollenMap["UVIndex"]; ok {
			if ok {
				r.UVIndex = null.NewInt(int64(val), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["UVIndex"]; ok {
			if ok {
				r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["AirQuality"]; ok {
			if ok {
				r.AirQualityCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Grass"]; ok {
			if ok {
				r.GrassCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Mold"]; ok {
			if ok {
				r.MoldCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Ragweed"]; ok {
			if ok {
				r.RagweedCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Tree"]; ok {
			if ok {
				r.TreeCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}

		r.Day, _ = s.Day.ResponseFormat(version)
		r.Night, _ = s.Night.ResponseFormat(version)

		return r, nil

	case "1.5e":
		weatherCategoryToI8NCode = WeatherCategoryToFirmwareMapExtended
		fallthrough
	case "1.5":
		var r NullableAccuDailyForecastFormatV1p3
		r.Date = s.Date
		r.EpochDate = s.EpochDate

		if s.Moon.Phase.Valid {
			if val, ok := MoonPhaseToFirmwareMap[s.Moon.Phase.String]; ok {
				if ok {
					r.MoonPhase = null.NewInt(int64(val), true)
				}
			}
		}

		r.SunRise = getFormattedDateFromEpoch(s.Sun.Rise)
		r.SunSet = getFormattedDateFromEpoch(s.Sun.Set)
		r.MoonRise = getFormattedDateFromEpoch(s.Moon.Rise)
		r.MoonSet = getFormattedDateFromEpoch(s.Moon.Set)

		r.MinTemperature = s.Temperature.Minimum.Value
		r.MaxTemperature = s.Temperature.Maximum.Value
		r.HoursOfSun = s.HoursOfSun

		if val, ok := s.AirAndPollenMap["UVIndex"]; ok {
			if ok {
				r.UVIndex = null.NewInt(int64(val), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["UVIndex"]; ok {
			if ok {
				r.UVCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["AirQuality"]; ok {
			if ok {
				r.AirQualityCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Grass"]; ok {
			if ok {
				r.GrassCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Mold"]; ok {
			if ok {
				r.MoldCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Ragweed"]; ok {
			if ok {
				r.RagweedCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}
		if val, ok := s.AirAndPollenCategoryMap["Tree"]; ok {
			if ok {
				r.TreeCategory = null.NewInt(int64(weatherCategoryToI8NCode[val]), true)
			}
		}

		r.Day, _ = s.Day.ResponseFormat(version)
		r.Night, _ = s.Night.ResponseFormat(version)

		return r, nil

	default:
		return s, errors.New("unsupported version")
	}

}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableAccuDailyForecastFormatV1p3) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuDailyForecastFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableAccuDailyForecastFormatV1p2) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuDailyForecastFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableDayNightData) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableDayNightData")
	switch version {
	case "1.1e":
		fallthrough
	case "1.1":
		return s, nil
	case "1.2e":
		fallthrough
	case "1.2":
		var r NullableDayNightDataFormatV1p2
		r.Icon = DisplayIconByAccuIcon(s.Icon)
		r.PrecipitationProbability = s.PrecipitationProbability
		r.ThunderstormProbability = s.ThunderstormProbability
		r.RainProbability = s.RainProbability
		r.SnowProbability = s.SnowProbability
		r.IceProbability = s.IceProbability
		r.HoursOfPrecipitation = s.HoursOfPrecipitation
		r.HoursOfRain = s.HoursOfRain
		r.HoursOfSnow = s.HoursOfSnow
		r.HoursOfIce = s.HoursOfIce
		r.CloudCover = s.CloudCover
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		r.TotalLiquid = s.TotalLiquid.Value
		r.Rain = s.Rain.Value
		r.Snow = s.Snow.Value
		r.Ice = s.Ice.Value
		return r, nil
	case "1.3e":
		fallthrough
	case "1.3":
		var r NullableDayNightDataFormatV1p2
		r.Icon = DisplayIconByAccuIcon(s.Icon)
		r.PrecipitationProbability = s.PrecipitationProbability
		r.ThunderstormProbability = s.ThunderstormProbability
		r.RainProbability = s.RainProbability
		r.SnowProbability = s.SnowProbability
		r.IceProbability = s.IceProbability
		r.HoursOfPrecipitation = s.HoursOfPrecipitation
		r.HoursOfRain = s.HoursOfRain
		r.HoursOfSnow = s.HoursOfSnow
		r.HoursOfIce = s.HoursOfIce
		r.CloudCover = s.CloudCover
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		r.TotalLiquid = s.TotalLiquid.Value
		r.Rain = s.Rain.Value
		r.Snow = s.Snow.Value
		r.Ice = s.Ice.Value
		return r, nil
	case "1.4e":
		fallthrough
	case "1.4":
		var r NullableDayNightDataFormatV1p3
		r.Icon = DisplayIconByAccuIcon(s.Icon)
		r.PrecipitationProbability = s.PrecipitationProbability
		r.ThunderstormProbability = s.ThunderstormProbability
		r.RainProbability = s.RainProbability
		r.SnowProbability = s.SnowProbability
		r.IceProbability = s.IceProbability
		//r.HoursOfPrecipitation = s.HoursOfPrecipitation
		//r.HoursOfRain = s.HoursOfRain
		//r.HoursOfSnow = s.HoursOfSnow
		//r.HoursOfIce = s.HoursOfIce
		r.CloudCover = s.CloudCover
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		//r.TotalLiquid = s.TotalLiquid.Value
		r.Rain = s.Rain.Value
		r.Snow = s.Snow.Value
		r.Ice = s.Ice.Value
		return r, nil

	case "1.5e":
		fallthrough
	case "1.5":
		var r NullableDayNightDataFormatV1p3
		r.Icon = DisplayIconByAccuIcon(s.Icon)
		r.PrecipitationProbability = s.PrecipitationProbability
		r.ThunderstormProbability = s.ThunderstormProbability
		r.RainProbability = s.RainProbability
		r.SnowProbability = s.SnowProbability
		r.IceProbability = s.IceProbability
		//r.HoursOfPrecipitation = s.HoursOfPrecipitation
		//r.HoursOfRain = s.HoursOfRain
		//r.HoursOfSnow = s.HoursOfSnow
		//r.HoursOfIce = s.HoursOfIce
		r.CloudCover = s.CloudCover
		r.WindSpeed = s.Wind.Speed.Value
		r.WindHeading = s.Wind.Direction.Degrees
		r.WindGust = s.WindGust.Speed.Value
		//r.TotalLiquid = s.TotalLiquid.Value
		r.Rain = s.Rain.Value
		r.Snow = s.Snow.Value
		r.Ice = s.Ice.Value
		return r, nil

	default:
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableDayNightDataFormatV1p3) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableDayNightDataFormatV1p3")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableDayNightDataFormatV1p2) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableDayNightDataFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableAccuCurrentForecastResponse) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuCurrentForecastResponse")
	switch version {
	case "1.1e":
		fallthrough
	case "1.1":
		return s, nil

	case "1.2e":
		fallthrough
	case "1.2":
		var r NullableAccuCurrentForecastResponseFormatV1p2
		r.Epoch = s.EpochTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)
		r.IsDayTime = s.IsDayTime
		r.Temperature = s.Temperature.Metric.Value
		r.TornadoProbability = s.TornadoProbability
		r.HailProbability = s.HailProbability
		r.TornadoProbabilityDeprecated = s.TornadoProbability
		r.HailProbabilityDeprecated = s.HailProbability
		return r, nil

	case "1.3e":
		fallthrough
	case "1.3":
		var r NullableAccuCurrentForecastResponseFormatV1p2
		r.Epoch = s.EpochTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)
		r.IsDayTime = s.IsDayTime
		r.Temperature = s.Temperature.Metric.Value
		r.TornadoProbability = s.TornadoProbability
		r.HailProbability = s.HailProbability
		r.TornadoProbabilityDeprecated = s.TornadoProbability
		r.HailProbabilityDeprecated = s.HailProbability
		return r, nil

	case "1.4e":
		fallthrough
	case "1.4":
		var r NullableAccuCurrentForecastResponseFormatV1p3
		r.Epoch = s.EpochTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)
		r.IsDayTime = s.IsDayTime
		r.Temperature = s.Temperature.Metric.Value
		r.TornadoProbability = s.TornadoProbability
		r.HailProbability = s.HailProbability
		return r, nil

	case "1.5e":
		fallthrough
	case "1.5":
		var r NullableAccuCurrentForecastResponseFormatV1p3
		r.Epoch = s.EpochTime
		r.WeatherIcon = DisplayIconByAccuIcon(s.WeatherIcon)
		r.IsDayTime = s.IsDayTime
		r.Temperature = s.Temperature.Metric.Value
		r.TornadoProbability = s.TornadoProbability
		r.HailProbability = s.HailProbability
		return r, nil

	default:
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableAccuCurrentForecastResponseFormatV1p3) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuCurrentForecastResponseFormatV1p3")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableAccuCurrentForecastResponseFormatV1p2) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableAccuCurrentForecastResponseFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableUniversalForecast) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableUniversalForecast")
	switch version {
	case "1.1e":
		fallthrough
	case "1.1":
		return s, nil

	case "1.2e":
		fallthrough
	case "1.2":
		var r NullableUniversalForecastFormatV1p2
		r.Date = s.Date
		r.Time = s.Time
		r.GmtOffset = s.GmtOffset
		r.ForecastTime = s.ForecastTime
		r.Category = s.Category
		if s.Today != nil {
			r.Today, _ = s.Today.ResponseFormat(version)
		} else {
			r.Today = nil
		}

		if s.Current != nil {
			r.Current, _ = s.Current.ResponseFormat(version)
		} else {
			r.Current = nil
		}

		var daily []ApiResponseInterface
		if s.Daily != nil {
			for i := 0; i < len(*s.Daily); i++ {

				if !((*s.Daily)[i]).Moon.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i-1]).Moon.Rise
					} else {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i+1]).Moon.Rise
					}
				}
				if !((*s.Daily)[i]).Moon.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i-1]).Moon.Set
					} else {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i+1]).Moon.Set
					}
				}
				if !((*s.Daily)[i]).Sun.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i-1]).Sun.Rise
					} else {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i+1]).Sun.Rise
					}
				}
				if !((*s.Daily)[i]).Sun.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i-1]).Sun.Set
					} else {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i+1]).Sun.Set
					}
				}

				r, _ := ((*s.Daily)[i]).ResponseFormat(version)
				daily = append(daily, r)
			}
		}
		r.Daily = daily

		var hourly []ApiResponseInterface
		if s.Hourly != nil {
			for i := 0; i < len(*s.Hourly); i++ {

				// work around for null incoming temperature
				if i == 0 && !((*s.Hourly)[i]).Temperature.Value.Valid {
					((*s.Hourly)[i]).Temperature.Value = ((*s.Hourly)[i+1]).Temperature.Value
				}

				r, _ := ((*s.Hourly)[i]).ResponseFormat(version)
				hourly = append(hourly, r)
			}
		}
		r.Hourly = hourly

		r.NWSForecast = s.NWSForecast

		r.FlowControl = s.FlowControl
		return r, nil

	case "1.3e":
		fallthrough
	case "1.3":
		var r NullableUniversalForecastFormatV1p2
		r.Date = s.Date
		r.Time = s.Time
		r.GmtOffset = s.GmtOffset
		r.ForecastTime = s.ForecastTime
		r.Category = s.Category
		if s.Today != nil {
			r.Today, _ = s.Today.ResponseFormat(version)
		} else {
			r.Today = nil
		}

		if s.Current != nil {
			r.Current, _ = s.Current.ResponseFormat(version)
		} else {
			r.Current = nil
		}

		var daily []ApiResponseInterface
		if s.Daily != nil {
			for i := 0; i < len(*s.Daily); i++ {

				if !((*s.Daily)[i]).Moon.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i-1]).Moon.Rise
					} else {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i+1]).Moon.Rise
					}
				}
				if !((*s.Daily)[i]).Moon.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i-1]).Moon.Set
					} else {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i+1]).Moon.Set
					}
				}
				if !((*s.Daily)[i]).Sun.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i-1]).Sun.Rise
					} else {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i+1]).Sun.Rise
					}
				}
				if !((*s.Daily)[i]).Sun.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i-1]).Sun.Set
					} else {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i+1]).Sun.Set
					}
				}

				r, _ := ((*s.Daily)[i]).ResponseFormat(version)
				daily = append(daily, r)
			}
		}
		r.Daily = daily

		var hourly []ApiResponseInterface
		if s.Hourly != nil {
			for i := 0; i < len(*s.Hourly); i++ {

				// work around for null incoming temperature
				if i == 0 && !((*s.Hourly)[i]).Temperature.Value.Valid {
					((*s.Hourly)[i]).Temperature.Value = ((*s.Hourly)[i+1]).Temperature.Value
				}

				r, _ := ((*s.Hourly)[i]).ResponseFormat(version)
				hourly = append(hourly, r)
			}
		}
		r.Hourly = hourly

		r.NWSForecast = s.NWSForecast

		r.FlowControl = s.FlowControl
		return r, nil

	case "1.4e":
		fallthrough
	case "1.4":
		var r NullableUniversalForecastFormatV1p3
		r.Date = s.Date
		r.Time = s.Time
		r.GmtOffset = s.GmtOffset
		r.ForecastTime = s.ForecastTime
		r.Category = s.Category
		if s.Today != nil {
			r.Today, _ = s.Today.ResponseFormat(version)
		} else {
			r.Today = nil
		}

		if s.Current != nil {
			r.Current, _ = s.Current.ResponseFormat(version)
		} else {
			r.Current = nil
		}

		var daily []ApiResponseInterface
		if s.Daily != nil {
			for i := 0; i < len(*s.Daily); i++ {

				if !((*s.Daily)[i]).Moon.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i-1]).Moon.Rise
					} else {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i+1]).Moon.Rise
					}
				}
				if !((*s.Daily)[i]).Moon.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i-1]).Moon.Set
					} else {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i+1]).Moon.Set
					}
				}
				if !((*s.Daily)[i]).Sun.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i-1]).Sun.Rise
					} else {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i+1]).Sun.Rise
					}
				}
				if !((*s.Daily)[i]).Sun.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i-1]).Sun.Set
					} else {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i+1]).Sun.Set
					}
				}

				r, _ := ((*s.Daily)[i]).ResponseFormat(version)
				daily = append(daily, r)
			}
		}
		r.Daily = daily

		var hourly []ApiResponseInterface
		if s.Hourly != nil {
			for i := 0; i < len(*s.Hourly); i++ {

				if i == 0 && !((*s.Hourly)[i]).Temperature.Value.Valid {
					((*s.Hourly)[i]).Temperature.Value = ((*s.Hourly)[i+1]).Temperature.Value
				}

				r, _ := ((*s.Hourly)[i]).ResponseFormat(version)
				hourly = append(hourly, r)
			}
		}
		r.Hourly = hourly

		r.NWSForecast = s.NWSForecast

		r.FlowControl = s.FlowControl
		return r, nil

	case "1.5e":
		fallthrough
	case "1.5":
		var r NullableUniversalForecastFormatV1p3
		r.Date = s.Date
		r.Time = s.Time
		r.GmtOffset = s.GmtOffset
		r.ForecastTime = s.ForecastTime
		r.Category = s.Category
		if s.Today != nil {
			r.Today, _ = s.Today.ResponseFormat(version)
		} else {
			r.Today = nil
		}

		if s.Current != nil {
			r.Current, _ = s.Current.ResponseFormat(version)
		} else {
			r.Current = nil
		}

		var daily []ApiResponseInterface
		if s.Daily != nil {
			for i := 0; i < len(*s.Daily); i++ {

				if !((*s.Daily)[i]).Moon.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i-1]).Moon.Rise
					} else {
						((*s.Daily)[i]).Moon.Rise = ((*s.Daily)[i+1]).Moon.Rise
					}
				}
				if !((*s.Daily)[i]).Moon.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i-1]).Moon.Set
					} else {
						((*s.Daily)[i]).Moon.Set = ((*s.Daily)[i+1]).Moon.Set
					}
				}
				if !((*s.Daily)[i]).Sun.Rise.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i-1]).Sun.Rise
					} else {
						((*s.Daily)[i]).Sun.Rise = ((*s.Daily)[i+1]).Sun.Rise
					}
				}
				if !((*s.Daily)[i]).Sun.Set.Valid {
					if i > 0 {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i-1]).Sun.Set
					} else {
						((*s.Daily)[i]).Sun.Set = ((*s.Daily)[i+1]).Sun.Set
					}
				}

				r, _ := ((*s.Daily)[i]).ResponseFormat(version)
				daily = append(daily, r)
			}
		}
		r.Daily = daily

		var hourly []ApiResponseInterface
		if s.Hourly != nil {
			for i := 0; i < len(*s.Hourly); i++ {

				if i == 0 && !((*s.Hourly)[i]).Temperature.Value.Valid {
					((*s.Hourly)[i]).Temperature.Value = ((*s.Hourly)[i+1]).Temperature.Value
				}

				r, _ := ((*s.Hourly)[i]).ResponseFormat(version)
				hourly = append(hourly, r)
			}
		}
		r.Hourly = hourly

		r.NWSForecast = s.NWSForecast

		r.FlowControl = s.FlowControl
		return r, nil

	default:
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableUniversalForecastFormatV1p3) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableUniversalForecastFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (s NullableUniversalForecastFormatV1p2) ResponseFormat(version string) (ApiResponseInterface, error) {
	//log.Println("[ResponseFormat] NullableUniversalForecastFormatV1p2")
	if version == "1.2" || version == "1.3" || version == "1.4" || version == "1.5" || version == "1.2e" || version == "1.3e" || version == "1.4e" || version == "1.5e" {
		return s, nil
	} else {
		return s, errors.New("unsupported version")
	}
}

//==============================================
// Protocols - NullableGetWeatherForecastJson
//==============================================

//----------------------------------------------
// Get weather forecast for category one devices (v2)
//----------------------------------------------
/**
 * @brief
 */
func (accuLocation PostalCodeResponse) NullableGetWeatherForecastJsonExtended(category string, deviceID string, firmwareVersion string, callSubVersion string, includeToday bool, includeDaily bool, includeHourly bool, includeCurrent bool) ApiResponseInterface {
	//log.Printf("getWeatherForecast location key : %s, Timezone:%s, Device Category: %s, Device ID: %s", accuLocation.Key, accuLocation.TimeZone.Name, category, deviceID)

	var dailyClip = 7
	var hourlyClip = 12
	if callSubVersion == "3" || callSubVersion == "4" || callSubVersion == "5" {
		dailyClip = 8
		hourlyClip = 15
	}

	// Setup Forecast
	forecast := NullableUniversalForecast{}

	// Get Extended Info
	extendedInfo, err := device.GetExtendedDeviceInfo(deviceID)
	if err == nil {
		if extendedInfo.HasDateTimeBug {
			forecast.FlowControl = null.NewInt(ExceptionModeFlowCommand, true)
		} else {
			forecast.FlowControl = null.NewInt(DefaultModeFlowCommand, true)
		}
	} else {
		// default to exception mode to prevent accidental wipe.
		forecast.FlowControl = null.NewInt(DefaultModeFlowCommand, true)
	}
	forecast.ExtendedDeviceInfo = extendedInfo

	// Grab Weather Time
	weatherTime, err := GetLocalDateAndHourV2(accuLocation.TimeZone.Name, &extendedInfo)
	if err != nil {
		return ApiString("Could not get time from the timeZone")
	}

	// Time Range to Disable Flow (Temp Cut Off.)
	if weatherTime.DateTime.Month() < 9 || weatherTime.DateTime.Month() == 9 && weatherTime.DateTime.Day() < 20 {
		forecast.FlowControl = null.NewInt(DefaultModeFlowCommand, true)
	}

	daily, _ := JsonQueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
	forecast.Headline = &daily.Headline
	forecast.Today = daily.Today

	// Time formatting
	forecast.Time = null.NewString(weatherTime.LocalTime, true)
	forecast.GmtOffset = null.NewFloat(accuLocation.TimeZone.GmtOffset, true)

	//forecast.Date = null.NewString(weatherTime.LocalDate, true)
	forecast.Date = null.NewString(weatherTime.Iso8601, true)
	if extendedInfo.TimeZoneOverride.Enabled {
		forecast.GmtOffset = null.NewFloat(float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset)+(float64(extendedInfo.TimeZoneOverride.MinuteOffset)/60.0), true)
	}

	// Load Current & NWSForecast
	forecast.NWSForecast = getNWSInfoV2(accuLocation)
	current, _ := NullablequeryAccuCurrentForecastAPI(accuLocation.Key, weatherTime, forecast.NWSForecast)

	// @todo better error handling
	forecast.Current = &current

	// Set Category
	switch category {
	case device.CAT1:
		forecast.Category = null.NewInt(1, true)
	case device.CAT2:
		forecast.Category = null.NewInt(2, true)
	case device.CAT3:
		forecast.Category = null.NewInt(3, true)
	}

	//--------------------------------------------------------------------
	// Load: Seven Day Forecast, 12 Hour Forecast and ForecastTime
	//--------------------------------------------------------------------
	var sevenDayForecast []NullableAccuDailyForecast

	for i := 0; i < len(daily.DailyForecasts) && i < dailyClip; i++ {
		sevenDayForecast = append(sevenDayForecast, daily.DailyForecasts[i])
	}
	forecast.Daily = &sevenDayForecast

	// Hourly
	hourly := NullableQueryAccuHourForecastAPI(accuLocation.Key, "24hour", weatherTime)
	now := time.Unix(time.Now().Unix(), 0)
	var futureHourly []NullableAccuHourlyForecast
	var i int
	for i = 0; len(futureHourly) < hourlyClip && i < len(hourly); i++ {
		if !now.After(time.Unix(int64(hourly[i].EpochDateTime.Int64), 0)) {
			futureHourly = append(futureHourly, hourly[i])
		}
	}
	forecast.Hourly = &futureHourly

	// Note no validation of populated futureHourly is performed here.
	unixEpoch := time.Unix(int64(futureHourly[0].EpochDateTime.Int64), 0)
	forecast.ForecastTime = null.NewString(unixEpoch.Add(time.Minute*time.Duration(forecast.GmtOffset.Float64*60)).Format("2006-01-02T15:04:05-0700"), true)
	//--------------------------------------------------------------------
	// End Load: Seven Day Forecast, 12 Hour Forecast and ForecastTime
	//--------------------------------------------------------------------

	// Saving or updating keys which will be updated by the cache updater
	toUpdateKey := "activelocations:" + accuLocation.Key
	toUpdateVal := accuLocation.TimeZone.Name + ":" + category
	common.RedisInstance.SaveRedisData([]byte(toUpdateVal), toUpdateKey, 720*time.Hour)

	//------------------------
	// Crude override to drop some results in api response.
	//------------------------
	if !includeDaily {
		forecast.Daily = nil
	}
	if !includeHourly {
		forecast.Hourly = nil
	}
	if !includeCurrent {
		forecast.Current = nil
	}
	if !includeToday {
		forecast.Today = nil
	}

	return forecast
}

func (accuLocation PostalCodeResponse) NullableGetWeatherForecastJson(category string, deviceID string, firmwareVersion string, callSubVersion string) ApiResponseInterface {
	return accuLocation.NullableGetWeatherForecastJsonExtended(category, deviceID, firmwareVersion, callSubVersion, true, true, true, true)
}

//==============================================
// Protocols - GetWeatherForecastTest
//==============================================

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func (accuLocation PostalCodeResponse) GetWeatherForecastTest(category string, deviceID string) string {

	//log.Printf("GetWeatherForecastTest location key : %s, Timezone:%s, Device Category: %s, Device ID: %s", accuLocation.Key, accuLocation.TimeZone.Name, category, deviceID)
	weatherTime, err := GetLocalDateAndHour(accuLocation.TimeZone.Name)
	if err != nil {
		return "Could not get time from the timeZone"
	}

	if strings.TrimSpace(accuLocation.Key) == "" {
		return "Location Not found (L1)"
	}

	var apiResult interface{}
	var templateFile string

	switch category {
	case device.CAT1:
		ats := AccuTemplateCat1Struct{}

		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime
		ats.GmtOffset = accuLocation.TimeZone.GmtOffset

		apiResult = ats
		templateFile = "templateTestDeviceCat3"
	case device.CAT2:
		ats := AccuTemplateCat2Struct{}
		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime
		ats.GmtOffset = accuLocation.TimeZone.GmtOffset

		apiResult = ats
		templateFile = "templateTestDeviceCat2"
	case device.CAT3:
		ats := AccuTemplateCat3Struct{}

		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime

		apiResult = ats
		templateFile = "templateTestDevice"
	}

	// Loading the template file, which depends on the category
	/*
		body, err := ioutil.ReadFile(path.Join("/templates", templateFile))
		if err != nil {
			log.Printf(err.Error())
			return err.Error()
		}
	*/
	body, err := common.RedisInstance.GetCachedFile("/templates", templateFile, time.Hour)

	tmpl, err := template.New("body").Funcs(template.FuncMap{
		"getWifiIcon": func(accuIcon int) int {
			return IconMap[accuIcon]
		},
		"getFormatedDateFromEpoch": func(dateStr string) string {
			layout := "2006-01-02T15:04:05-07:00"
			t, err := time.Parse(layout, dateStr)

			if err != nil {
				log.Printf("GetWeatherForecastTest - Error parsing the time : %s %s", dateStr, err.Error())
			}

			return t.Format("15:04")
		},
		"getMoonPhrase": func(accuPhrase string) int {
			return MoonPhraseMap[accuPhrase]
		},
		"convertValueToEnum": func(value string) int {
			return ValueToEnumMap[value]
		},
	}).Parse(string(body))
	var forecast bytes.Buffer
	if tmpl == nil {
		log.Printf("Template is null -- Error")
		return "Problem loading template"
	}
	err = tmpl.Execute(&forecast, apiResult)
	if err != nil {
		log.Printf(err.Error())
		return err.Error()
	}

	// This hack is just because go has issues with interpreting "<" and ">" it's trying to render it as html block. TODO Fix it to properly set the template
	fc := strings.Replace(forecast.String(), "##", "<", -1)
	fc = strings.Replace(fc, "!!", ">", -1)

	return fc
}

//==============================================
// Protocols - GetWeatherForecastV2
//==============================================

//----------------------------------------------
// Get weather forecast for category one devices
//----------------------------------------------
/**
 * @brief
 */
func (accuLocation PostalCodeResponse) GetWeatherForecastV2(category string, deviceID string, firmwareVersion string) string {
	// Get Extended Info
	extendedInfo, err := device.GetExtendedDeviceInfo(deviceID)
	var flow int
	if err == nil {
		if extendedInfo.HasDateTimeBug {
			flow = ExceptionModeFlowCommand
		} else {
			flow = DefaultModeFlowCommand
		}
	} else {
		// default to exception mode to prevent accidental wipe.
		flow = DefaultModeFlowCommand
	}

	// @TODO - return anonymous if inside exception period and device is on previous API
	// @TODO - time loops and compression if specified in extendedInfo

	//log.Printf("getWeatherForecastV2 location key : %s, Timezone:%s, Device Category: %s, Device ID: %s", accuLocation.Key, accuLocation.TimeZone.Name, category, deviceID)
	weatherTime, err := GetLocalDateAndHourV2(accuLocation.TimeZone.Name, &extendedInfo)

	// Time Range to Disable Flow (Temp Cut Off.)
	if weatherTime.DateTime.Month() < 9 || weatherTime.DateTime.Month() == 9 && weatherTime.DateTime.Day() < 20 {
		flow = DefaultModeFlowCommand
	} else {
		if extendedInfo.HasDateTimeBug {
			// Force Anonymous response to pre patched firmway
			if firmwareVersion == "" {
				return "<anonymous:true>"
			}
		}
	}

	if err != nil {
		return "Could not get time from the timeZone"
	}

	if strings.TrimSpace(accuLocation.Key) == "" {
		return "Location Not found (L2)"
	}

	var apiResult interface{}
	var templateFile string

	switch category {
	case device.CAT1:
		ats := AccuTemplateCat1Struct{}
		ats.FlowControl = flow
		accu1dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
		ats.Headline = &accu1dForecast.Headline
		ats.DailyForecast = &accu1dForecast.DailyForecasts[0]
		ats.DailyForecast.Actual.Wind.Speed.ValueRound = Round(ats.DailyForecast.Actual.Wind.Speed.Value)
		ats.DailyForecast.Actual.WindGust.Speed.ValueRound = Round(ats.DailyForecast.Actual.WindGust.Speed.Value)
		ats.DailyForecast.Temperature.Maximum.ValueRound = Round(ats.DailyForecast.Temperature.Maximum.Value)
		ats.DailyForecast.Temperature.Minimum.ValueRound = Round(ats.DailyForecast.Temperature.Minimum.Value)
		ats.DailyForecast.Actual.Rain.ValueRound = Round(ats.DailyForecast.Actual.Rain.Value)
		ats.DailyForecast.Actual.Snow.ValueRound = Round(ats.DailyForecast.Actual.Snow.Value)
		ats.DailyForecast.NWSevereComponentMap = getNWSInfo(accuLocation.PrimaryPostalCode)

		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime
		ats.GmtOffset = accuLocation.TimeZone.GmtOffset
		if extendedInfo.TimeZoneOverride.Enabled {
			ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
		}

		apiResult = ats
		templateFile = "templateCat1V2"
	case device.CAT2:
		ats := AccuTemplateCat2Struct{}
		ats.FlowControl = flow
		accu1dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
		accuCurrentForecast := queryAccuCurrentForecastAPI(accuLocation.Key, weatherTime)

		ats.Headline = &accu1dForecast.Headline
		ats.DailyForecast = &accu1dForecast.DailyForecasts[0]
		ats.DailyForecast.NWSevereComponentMap = getNWSInfo(accuLocation.PrimaryPostalCode)
		ats.CurrentForecast = &accuCurrentForecast[0]

		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime
		ats.GmtOffset = accuLocation.TimeZone.GmtOffset
		if extendedInfo.TimeZoneOverride.Enabled {
			ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
		}
		apiResult = ats
		templateFile = "templateCat2V2"
	case device.CAT3:
		ats := AccuTemplateCat3Struct{}
		ats.FlowControl = flow
		accu10dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
		accu24hForecast := QueryAccuHourForecastAPI(accuLocation.Key, "24hour", weatherTime)

		var a7f []AccuDailyForecast
		for i := 0; i < 7; i++ {
			a7f = append(a7f, accu10dForecast.DailyForecasts[i])
		}

		n := time.Unix(time.Now().Unix(), 0)
		x := 0
		var a12f []AccuHourlyForecastResponse
		var i int
		// this is incorrect we should be returning 15 hourly records.
		for i = 0; i < 12+x && i < len(accu24hForecast); i++ {
			if n.After(time.Unix(int64(accu24hForecast[i].EpochDateTime), 0)) {
				// Should we be including the current hour?
				x++
			} else {
				a12f = append(a12f, accu24hForecast[i])
			}
		}

		ats.Accu7d = &a7f
		ats.DailyForecast = &accu10dForecast.DailyForecasts[0]
		ats.DailyForecast.NWSevereComponentMap = getNWSInfo(accuLocation.PrimaryPostalCode)
		ats.Accu24h = &a12f

		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime
		utime := time.Unix(int64(accu24hForecast[i-12].EpochDateTime), 0)
		ats.GmtOffset = accuLocation.TimeZone.GmtOffset
		if extendedInfo.TimeZoneOverride.Enabled {
			ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
		}
		//ats.ForecastTime = utime.Format("06:01:02 15:04")
		ats.ForecastTime = utime.Add(time.Minute * time.Duration(ats.GmtOffset*60)).Format("06:01:02 15:04")

		apiResult = ats
		templateFile = "templateCat3V2"
	}

	// Loading the template file, which depends on the category
	/*
		body, err := ioutil.ReadFile(path.Join("/templates", templateFile))
		if err != nil {
			log.Printf(err.Error())
			return err.Error()
		}
	*/
	body, err := common.RedisInstance.GetCachedFile("/templates", templateFile, time.Hour)

	tmpl, err := template.New("body").Funcs(template.FuncMap{
		"getWifiIcon": func(accuIcon int) int {
			return IconMap[accuIcon]
		},
		"getFormatedDateFromEpoch": func(dateStr string) string {
			layout := "2006-01-02T15:04:05-07:00"
			t, err := time.Parse(layout, dateStr)

			if err != nil {
				log.Printf("GetWeatherForecastV2 - Error parsing the time : %s %s", dateStr, err.Error())
			}

			return t.Format("15:04")
		},
		"getMoonPhrase": func(accuPhrase string) int {
			return MoonPhraseMap[accuPhrase]
		},
		"convertValueToEnum": func(value string) int {
			return ValueToEnumMap[value]
		},
	}).Parse(string(body))
	var forecast bytes.Buffer
	if tmpl == nil {
		log.Printf("Template is null -- Error")
		return "Problem loading template"
	}
	err = tmpl.Execute(&forecast, apiResult)
	if err != nil {
		log.Printf(err.Error())
		return err.Error()
	}

	// This hack is just because go has issues with interpreting "<" and ">" it's trying to render it as html block. TODO Fix it to properly set the template
	fc := strings.Replace(forecast.String(), "##", "<", -1)
	fc = strings.Replace(fc, "!!", ">", -1)

	// Saving or updating keys which will be updated by the cache updater
	toUpdateKey := "activelocations:" + accuLocation.Key
	toUpdateVal := accuLocation.TimeZone.Name + ":" + category

	common.RedisInstance.SaveRedisData([]byte(toUpdateVal), toUpdateKey, 720*time.Hour)
	return fc
}

//==============================================
// Protocols - GetWeatherForecast
//==============================================

//----------------------------------------------
// Get weather forecast for category one devices
//----------------------------------------------
/**
 * @brief
 */
func (accuLocation PostalCodeResponse) GetWeatherForecast(category string, deviceID string, forecastType string, firmwareVersion string) string {

	// Get Extended Info
	extendedInfo, err := device.GetExtendedDeviceInfo(deviceID)
	var flow int
	if err == nil {
		if extendedInfo.HasDateTimeBug {
			flow = ExceptionModeFlowCommand
		} else {
			flow = DefaultModeFlowCommand
		}
	} else {
		// default to exception mode to prevent accidental wipe.
		flow = DefaultModeFlowCommand
	}

	// @TODO - return anonymous if inside exception period and device is on previous API
	// @TODO - time loops and compression if specified in extendedInfo

	//log.Printf("getWeatherForecast location key : %s, Timezone:%s, Device Category: %s, Device ID: %s", accuLocation.Key, accuLocation.TimeZone.Name, category, deviceID)

	//log.Printf("getWeatherForecastV2 location key : %s, Timezone:%s, Device Category: %s, Device ID: %s", accuLocation.Key, accuLocation.TimeZone.Name, category, deviceID)
	weatherTime, err := GetLocalDateAndHourV2(accuLocation.TimeZone.Name, &extendedInfo)
	if err != nil {
		return "Could not get time from the timeZone"
	}

	// Time Range to Disable Flow
	if weatherTime.DateTime.Month() < 9 || weatherTime.DateTime.Month() == 9 && weatherTime.DateTime.Day() < 20 {
		flow = DefaultModeFlowCommand
	} else {
		if extendedInfo.HasDateTimeBug {
			// Force Anonymous response to pre patched firmway
			if firmwareVersion == "" {
				return "<anonymous:true>"
			}
		}
	}

	if strings.TrimSpace(accuLocation.Key) == "" {
		return "Location Not found: (L3)"
	}

	var apiResult interface{}
	var templateFile string

	if forecastType == ForecastTypeStreams {
		ats := AccuTemplateCat1Struct{}
		ats.FlowControl = flow
		accu1dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
		ats.Headline = &accu1dForecast.Headline
		ats.DailyForecast = &accu1dForecast.DailyForecasts[0]

		// Time formatting
		ats.DateStr = weatherTime.LocalDate
		ats.TimeStr = weatherTime.LocalTime
		ats.GmtOffset = accuLocation.TimeZone.GmtOffset
		if extendedInfo.TimeZoneOverride.Enabled {
			ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
		}
		apiResult = ats
		templateFile = "templateDatastreams"
	} else {

		switch category {
		case device.CAT1:
			ats := AccuTemplateCat1Struct{}
			ats.FlowControl = flow
			accu1dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
			ats.Headline = &accu1dForecast.Headline
			ats.DailyForecast = &accu1dForecast.DailyForecasts[0]
			ats.DailyForecast.Actual.Wind.Speed.ValueRound = Round(ats.DailyForecast.Actual.Wind.Speed.Value / 3.6)
			ats.DailyForecast.Actual.WindGust.Speed.ValueRound = Round(ats.DailyForecast.Actual.WindGust.Speed.Value / 3.6)
			ats.DailyForecast.Temperature.Maximum.ValueRound = Round(ats.DailyForecast.Temperature.Maximum.Value)
			ats.DailyForecast.Temperature.Minimum.ValueRound = Round(ats.DailyForecast.Temperature.Minimum.Value)
			ats.DailyForecast.Actual.Rain.ValueRound = Round(ats.DailyForecast.Actual.Rain.Value)
			ats.DailyForecast.Actual.Snow.ValueRound = Round(ats.DailyForecast.Actual.Snow.Value)

			// Time formatting
			ats.DateStr = weatherTime.LocalDate
			ats.TimeStr = weatherTime.LocalTime
			ats.GmtOffset = accuLocation.TimeZone.GmtOffset
			if extendedInfo.TimeZoneOverride.Enabled {
				ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
			}
			apiResult = ats
			templateFile = "templateCat1"
		case device.CAT2:
			ats := AccuTemplateCat2Struct{}
			ats.FlowControl = flow
			accu1dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
			accuCurrentForecast := queryAccuCurrentForecastAPI(accuLocation.Key, weatherTime)

			ats.Headline = &accu1dForecast.Headline
			ats.DailyForecast = &accu1dForecast.DailyForecasts[0]
			ats.CurrentForecast = &accuCurrentForecast[0]

			// Time formatting
			ats.DateStr = weatherTime.LocalDate
			ats.TimeStr = weatherTime.LocalTime
			ats.GmtOffset = accuLocation.TimeZone.GmtOffset
			if extendedInfo.TimeZoneOverride.Enabled {
				ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
			}
			apiResult = ats
			templateFile = "templateCat2"
		case device.CAT3:
			ats := AccuTemplateCat3Struct{}
			ats.FlowControl = flow
			accu10dForecast := QueryAccuDayForecastAPI(accuLocation.Key, accuLocation.TimeZone.Name, "10day", weatherTime)
			accu24hForecast := QueryAccuHourForecastAPI(accuLocation.Key, "24hour", weatherTime)

			var a7f []AccuDailyForecast
			for i := 0; i < 7; i++ {
				a7f = append(a7f, accu10dForecast.DailyForecasts[i])
			}

			n := time.Unix(time.Now().Unix(), 0)
			x := 0
			var a12f []AccuHourlyForecastResponse
			var i int
			for i = 0; i < 12+x && i < len(accu24hForecast); i++ {
				if n.After(time.Unix(int64(accu24hForecast[i].EpochDateTime), 0)) {
					// What about the current hour?, should this be n.after .. + EpochDateTime + 60 minutes)
					x++
				} else {
					a12f = append(a12f, accu24hForecast[i])
				}
			}

			ats.Accu7d = &a7f
			ats.DailyForecast = &accu10dForecast.DailyForecasts[0]
			ats.Accu24h = &a12f

			// Time formatting
			ats.DateStr = weatherTime.LocalDate
			ats.TimeStr = weatherTime.LocalTime
			utime := time.Unix(int64(accu24hForecast[i-12].EpochDateTime), 0)
			ats.GmtOffset = accuLocation.TimeZone.GmtOffset
			if extendedInfo.TimeZoneOverride.Enabled {
				ats.GmtOffset = float64(extendedInfo.TimeZoneOverride.Sign*extendedInfo.TimeZoneOverride.HourOffset) + (float64(extendedInfo.TimeZoneOverride.MinuteOffset) / 60.0)
			}
			//ats.ForecastTime = utime.Format("06:01:02 15:04")
			ats.ForecastTime = utime.Add(time.Minute * time.Duration(ats.GmtOffset*60)).Format("06:01:02 15:04")

			apiResult = ats
			templateFile = "templateCat3"
		}
	}

	// Loading the template file, which depends on the category
	/*
		body, err := ioutil.ReadFile(path.Join("/templates", templateFile))
		if err != nil {
			log.Printf(err.Error())
			return err.Error()
		}
	*/
	body, err := common.RedisInstance.GetCachedFile("/templates", templateFile, time.Hour)

	tmpl, err := template.New("body").Funcs(template.FuncMap{
		"getWifiIcon": func(accuIcon int) int {
			return IconMap[accuIcon]
		},
		"getFormatedDateFromEpoch": func(dateStr string) string {
			layout := "2006-01-02T15:04:05-07:00"
			t, err := time.Parse(layout, dateStr)

			if err != nil {
				log.Printf("GetWeatherForecast - Error parsing the time : %s %s", dateStr, err.Error())
			}

			return t.Format("15:04")
		},
		"getMoonPhrase": func(accuPhrase string) int {
			return MoonPhraseMap[accuPhrase]
		},
		"convertValueToEnum": func(value string) int {
			return ValueToEnumMap[value]
		},
	}).Parse(string(body))
	var forecast bytes.Buffer
	if tmpl == nil {
		log.Printf("Template is null -- Error")
		return "Problem loading template"
	}
	err = tmpl.Execute(&forecast, apiResult)
	if err != nil {
		log.Printf(err.Error())
		return err.Error()
	}

	// This hack is just because go has issues with interpreting "<" and ">" it's trying to render it as html block. TODO Fix it to properly set the template
	fc := strings.Replace(forecast.String(), "##", "<", -1)
	fc = strings.Replace(fc, "!!", ">", -1)

	// Saving or updating keys which will be updated by the cache updater
	toUpdateKey := "activelocations:" + accuLocation.Key
	toUpdateVal := accuLocation.TimeZone.Name + ":" + category

	common.RedisInstance.SaveRedisData([]byte(toUpdateVal), toUpdateKey, 720*time.Hour)
	return fc
}

//==============================================
// Functions
//==============================================

//----------------------------------------------
// @GetLocationFromPC
//----------------------------------------------
/**
 * @brief
 */
func GetLocationFromPC(accuPostalCode string) (PostalCodeResponse, error) {

	//log.Printf("GetLocationFromPC accuweather code: %s", accuPostalCode)
	pckey := "postalcode:" + accuPostalCode
	data, err := common.RedisInstance.GetCachedData(pckey)

	if err == nil {
		postalCodeResponse := PostalCodeResponse{}
		err := json.Unmarshal(data, &postalCodeResponse)
		if err != nil {
			log.Printf("Ummarshable - Purging Cache %s", pckey)
			//common.RedisInstance.RemoveKeyFromCache(pckey);
			return PostalCodeResponse{}, err
		}

		if strings.TrimSpace(postalCodeResponse.Key) == "" && postalCodeResponse.Code == "ServiceUnavailable" {
			// Temporary
			log.Printf("TEMPORARY - Purging Cache %s", pckey)
			common.RedisInstance.RemoveKeyFromCache(pckey)
		}

		if strings.TrimSpace(postalCodeResponse.Key) == "" && postalCodeResponse.Code == "" {
			log.Printf("TEMPORARY - Purging Empty Cache %s", pckey)
			common.RedisInstance.RemoveKeyFromCache(pckey)
		}

		return postalCodeResponse, nil

	} else {
		var Url *url.URL
		Url, _ = url.Parse("http://dataservice.accuweather.com")

		Url.Path += "/locations/v1/" + accuPostalCode
		parameters := url.Values{}
		parameters.Add("apikey", AccuApiKey)
		Url.RawQuery = parameters.Encode()

		resp, err := http.Get(Url.String())
		if err != nil {
			return PostalCodeResponse{}, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		var pc PostalCodeResponse
		err = json.Unmarshal(body, &pc)

		// No result returned
		if err != nil {
			err = errors.New("No location returned by accuweather for postal key : " + accuPostalCode)
			out := PostalCodeResponse{}
			out.Code = "NoEntries"
			dataBytes, _ := json.Marshal(out)
			common.RedisInstance.SaveRedisData(dataBytes, pckey, 60*time.Minute)
			return PostalCodeResponse{}, err
		}

		common.RedisInstance.SaveRedisData(body, pckey, 0)
		return pc, nil
	}
}

//----------------------------------------------
// @GetLocation
//----------------------------------------------
/**
 * @brief
 */
func GetLocation(postalCode string, countryCode string) (PostalCodeResponse, error) {
	var zipkey string

	// Fetch Redis
	if countryCode == "" {
		countryCode = DefaultCountryCode
		zipkey = "zip:" + postalCode + "_" + DefaultCountryCode
	} else {
		zipkey = "zip:" + postalCode + "_" + countryCode
	}

	data, err := common.RedisInstance.GetCachedData(zipkey)
	if err == nil {
		postalCodeResponse := PostalCodeResponse{}
		err := json.Unmarshal(data, &postalCodeResponse)
		if err != nil {
			log.Printf("Ummarshable - Purging Cache %s", zipkey)
			//common.RedisInstance.RemoveKeyFromCache(zipkey);
			return PostalCodeResponse{}, err
		}

		if strings.TrimSpace(postalCodeResponse.Key) == "" && postalCodeResponse.Code == "ServiceUnavailable" {
			// Temporary
			log.Printf("TEMPORARY - Purging Cache %s", zipkey)
			common.RedisInstance.RemoveKeyFromCache(zipkey)
		}

		if strings.TrimSpace(postalCodeResponse.Key) == "" && postalCodeResponse.Code == "" {
			log.Printf("TEMPORARY - Purging Empty Cache %s", zipkey)
			common.RedisInstance.RemoveKeyFromCache(zipkey)
		}

		if strings.TrimSpace(postalCodeResponse.Key) == "" {
			return PostalCodeResponse{}, err
		}

		return postalCodeResponse, nil

	} else {
		var Url *url.URL
		Url, _ = url.Parse("http://dataservice.accuweather.com")

		Url.Path += "/locations/v1/postalcodes/search"
		parameters := url.Values{}
		parameters.Add("apikey", AccuApiKey)
		parameters.Add("q", postalCode)
		Url.RawQuery = parameters.Encode()

		resp, err := http.Get(Url.String())
		if err != nil {
			return PostalCodeResponse{}, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		var pc []PostalCodeResponse
		json.Unmarshal(body, &pc)

		// No result returned
		if len(pc) == 0 {
			err = errors.New("No location returned by accuweather for postal code : " + postalCode + " and country code : " + countryCode)
			out := PostalCodeResponse{}
			out.Code = "NoEntries"
			dataBytes, _ := json.Marshal(out)
			common.RedisInstance.SaveRedisData(dataBytes, zipkey, 60*time.Minute)
			return PostalCodeResponse{}, err
		}

		// codeIndex variable is to save the index in the response array where we will get the correct location details
		codeIndex := -1
		for index, pcr := range pc {

			// Get the country code and write the data to cache
			c := strings.ToUpper(pc[index].Country.ID)
			countryzip := "zip:" + postalCode + "_" + c

			// Save data to redis and return
			dataBytes, _ := json.Marshal(pcr)
			common.RedisInstance.SaveRedisData(dataBytes, countryzip, 0)

			if (countryCode == "" && c == DefaultCountryCode) || (countryCode == c) {
				codeIndex = index
			}
		}

		if codeIndex == -1 {
			//err = errors.New("Location not found for postal code : " + postalCode + " and country code : " + countryCode)
			//return PostalCodeResponse{}, err
			// We couldn't get the location for the postal code defined for the device, since there is some values in the array, we will use the first item as location for that zip
			log.Printf("Location not found for postal code : " + postalCode + " and country code : " + countryCode + " . Using the default location returned")
			return pc[0], nil
		}

		return pc[codeIndex], nil
	}
}

//----------------------------------------------
// @SearchAllLocationsPerCountry
//----------------------------------------------
/**
 * @brief
 */
func SearchAllLocationsPerCountry(postalCode string, countryCode string) ([]byte, error) {

	locationKey := "location:lookup:" + countryCode + ":" + postalCode
	data, err := common.RedisInstance.GetCachedData(locationKey)
	if err == nil {
		return data, nil
	} else {
		var Url *url.URL
		Url, _ = url.Parse("http://dataservice.accuweather.com")

		Url.Path += "/locations/v1/postalcodes/" + countryCode + "/search"
		parameters := url.Values{}
		parameters.Add("apikey", AccuApiKey)
		parameters.Add("q", postalCode)
		Url.RawQuery = parameters.Encode()

		resp, err := http.Get(Url.String())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		var pc []PostalCodeResponse
		json.Unmarshal(body, &pc)

		// No result returned
		if len(pc) == 0 {
			err = errors.New("No location returned by accuweather for postal code : " + postalCode + " and country code : " + countryCode)
			locationDict := LocationLookupDict{}
			dataBytes, _ := json.Marshal(locationDict)
			// Save the empty data in redis for 10days it's not likely that the data will be available in accuweather anytime soon
			common.RedisInstance.SaveRedisData(dataBytes, locationKey, 240*time.Hour)
			return dataBytes, err
		}

		// Build postal code response
		locationList := make([]LocationLookupResponse, 0)
		for _, pcr := range pc {
			llr := LocationLookupResponse{}
			llr.CityName = pcr.LocalizedName
			llr.CountryCode = pcr.Country.ID
			llr.AcwKey = pcr.Key
			llr.AdministrativeArea = pcr.AdministrativeArea.ID
			llr.TzName = pcr.TimeZone.Name

			locationList = append(locationList, llr)
		}

		locationDict := LocationLookupDict{
			Locations: locationList,
		}

		dataBytes, _ := json.Marshal(locationDict)
		common.RedisInstance.SaveRedisData(dataBytes, locationKey, 0)

		return dataBytes, nil
	}
}

//----------------------------------------------
// @SearchAllLocationsPerPCorCity
//----------------------------------------------
/**
 * @brief
 */
func SearchAllLocationsPerPCorCity(pcOrCity string, countryCode string) ([]byte, error) {

	locationKey := "location:lookup:pcorcity:" + countryCode + ":" + pcOrCity
	data, err := common.RedisInstance.GetCachedData(locationKey)
	if err == nil {
		return data, nil
	} else {
		path := "/locations/v1/postalcodes/" + countryCode + "/search"
		pc := httpAccuGetLocation(path, pcOrCity)

		path = "/locations/v1/cities/" + countryCode + "/search"
		locationFromCities := httpAccuGetLocation(path, pcOrCity)

		for _, acwLocation := range locationFromCities {
			pc = append(pc, acwLocation)
		}

		// No result returned
		if len(pc) == 0 {
			err = errors.New("No location returned by accuweather for postal code or city : " + pcOrCity + " and country code : " + countryCode)
			locationDict := LocationLookupDict{}
			dataBytes, _ := json.Marshal(locationDict)
			// Save the empty data in redis for 10days it's not likely that the data will be available in accuweather anytime soon
			common.RedisInstance.SaveRedisData(dataBytes, locationKey, 240*time.Hour)
			return dataBytes, err
		}

		// Build postal code response
		locationList := make([]LocationLookupResponse, 0)
		for _, pcr := range pc {
			if findACKeyFromList(locationList, pcr.LocalizedName, pcr.Country.ID, pcr.AdministrativeArea.ID, pcr.TimeZone.Name) == true {
				continue
			}

			llr := LocationLookupResponse{}
			llr.CityName = pcr.LocalizedName
			llr.CountryCode = pcr.Country.ID
			llr.AcwKey = pcr.Key
			llr.AdministrativeArea = pcr.AdministrativeArea.ID
			llr.TzName = pcr.TimeZone.Name

			locationList = append(locationList, llr)
		}

		locationDict := LocationLookupDict{
			Locations: locationList,
		}

		dataBytes, _ := json.Marshal(locationDict)
		common.RedisInstance.SaveRedisData(dataBytes, locationKey, 0)

		return dataBytes, nil
	}
}

//----------------------------------------------
//
//----------------------------------------------
// Hour api forecast query to Accuweather
/**
 * @brief
 */
func QueryAccuHourForecastAPI(locationKey string, period string, weatherTime WeatherTime) []AccuHourlyForecastResponse {

	key := "forecast:" + locationKey + ":24hour:" + weatherTime.HourRange + "_" + weatherTime.LocalDate
	//redisInstance := &cache.RedisInstance{RedisSession: main.RedisClient}
	data, err := common.RedisInstance.GetCachedData(key)
	path := "/forecasts/v1/hourly/24hour/" + locationKey

	if err != nil {
		log.Printf("KeyNotFound : %s", key)
		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)

		// Persist Last cache update
		nowStr := time.Now().Format("02:01:2006 15:04:05")
		common.RedisInstance.SaveRedisData([]byte(nowStr), "forecastupdate:"+period+":"+locationKey, 0)
	}

	var accuForecast []AccuHourlyForecastResponse
	retryCount := 0
	json.Unmarshal(data, &accuForecast)
	for len(accuForecast) == 0 && retryCount < MaxRetries {
		retryCount = retryCount + 1
		log.Printf("Could not correctly fetch the weather forecast, retry number : %d path: %s", retryCount, path)
		time.Sleep(time.Duration(retryCount) * time.Second)

		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)
		json.Unmarshal(data, &accuForecast)
	}

	if retryCount == MaxRetries {

		return accuForecast
	}

	return accuForecast
}

//----------------------------------------------
// Hour api forecast query to Accuweather
//----------------------------------------------
/**
 * @brief
 */
func NullableQueryAccuHourForecastAPI(locationKey string, period string, weatherTime WeatherTime) []NullableAccuHourlyForecast {

	key := "forecast:" + locationKey + ":24hour:" + weatherTime.HourRange + "_" + weatherTime.LocalDate
	//redisInstance := &cache.RedisInstance{RedisSession: main.RedisClient}
	data, err := common.RedisInstance.GetCachedData(key)
	path := "/forecasts/v1/hourly/24hour/" + locationKey

	if err != nil {
		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)

		// Persist Last cache update
		nowStr := time.Now().Format("02:01:2006 15:04:05")
		common.RedisInstance.SaveRedisData([]byte(nowStr), "forecastupdate:"+period+":"+locationKey, 0)
	}

	var accuForecast []NullableAccuHourlyForecast
	retryCount := 0
	json.Unmarshal(data, &accuForecast)
	for len(accuForecast) == 0 && retryCount < MaxRetries {
		retryCount = retryCount + 1
		log.Printf("Could not correctly fetch the weather forecast, retry number : %s", retryCount)
		time.Sleep(time.Duration(retryCount) * time.Second)

		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)
		json.Unmarshal(data, &accuForecast)
	}

	if retryCount == MaxRetries {

		return accuForecast
	}

	return accuForecast
}

//----------------------------------------------
// Day api forecast query to Accuweather
//----------------------------------------------
/**
 * @brief
 */
func QueryAccuDayForecastAPI(locationKey string, timeZone string, period string, weatherTime WeatherTime) DailyForecast {

	key := "forecast:" + locationKey + ":" + period + ":" + weatherTime.HourRange + "_" + weatherTime.LocalDate

	//redisInstance := &cache.RedisInstance{RedisSession: main.RedisClient}
	data, err := common.RedisInstance.GetCachedData(key)
	path := "/forecasts/v1/daily/" + period + "/" + locationKey

	if err != nil {
		log.Printf("KeyNotFound : %s", key)
		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)
		// Persist Last cache update
		nowStr := time.Now().Format("02:01:2006 15:04:05")
		common.RedisInstance.SaveRedisData([]byte(nowStr), "forecastupdate:"+period+":"+locationKey, 0)
	}

	var accuForecast DailyForecast
	retryCount := 0
	json.Unmarshal(data, &accuForecast)
	for len(accuForecast.DailyForecasts) == 0 && retryCount < MaxRetries {
		retryCount = retryCount + 1
		log.Printf("Could not correctly fetch the weather forecast, retry number : %d path: %s", retryCount, path)
		time.Sleep(time.Duration(retryCount) * time.Second)

		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)
		json.Unmarshal(data, &accuForecast)
	}

	if retryCount == MaxRetries {

		return accuForecast
	}

	if weatherTime.DayInfo == DayInfoDay {
		accuForecast.DailyForecasts[0].Actual = &accuForecast.DailyForecasts[0].Day
	} else {
		accuForecast.DailyForecasts[0].Actual = &accuForecast.DailyForecasts[0].Night
	}

	// Populate Air and Pollen Map
	for index, forecast := range accuForecast.DailyForecasts {
		accuForecast.DailyForecasts[index].AirAndPollenMap = make(map[string]int)
		accuForecast.DailyForecasts[index].AirAndPollenCategoryMap = make(map[string]string)
		for _, airAndPollen := range forecast.AirAndPollen {
			accuForecast.DailyForecasts[index].AirAndPollenMap[airAndPollen.Name] = airAndPollen.Value
			accuForecast.DailyForecasts[index].AirAndPollenCategoryMap[airAndPollen.Name] = airAndPollen.Category
		}
	}

	// Add formatted Set and Rise

	return accuForecast
}

//----------------------------------------------
// Day api forecast query to Accuweather
//----------------------------------------------
/**
 * @brief
 */
func JsonQueryAccuDayForecastAPI(locationKey string, timeZone string, period string, weatherTime WeatherTime) (NullableDailyForecast, error) {
	return getNullableDailyForecast(locationKey, timeZone, period, weatherTime)
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func GetLocalDateAndHourV2(timeZone string, extendedInfo *device.ExtendedDeviceInfo) (WeatherTime, error) {
	weatherTime := WeatherTime{}

	// Load the specific location
	var loc *time.Location
	var err error
	LocationMapMutex.Lock()
	loc, keyFound := LocationMap[timeZone]
	LocationMapMutex.Unlock()

	if !keyFound {
		loc, err = time.LoadLocation(timeZone)
		if err != nil {
			log.Printf("Error Loading the location for timezone %s : %s", timeZone, err.Error())
			return weatherTime, err
		}
		LocationMap[timeZone] = loc
	}

	//set Location
	baseTime := time.Now()

	// 1. Apply Time Acceleration Option
	if extendedInfo.TimeCompression.Enabled {
		t := baseTime.Unix()
		elapsed := t - extendedInfo.TimeCompression.StartTime
		adjusted := float64(elapsed) * extendedInfo.TimeCompression.AccelerationRate
		final := extendedInfo.TimeCompression.StartTime + int64(adjusted) + extendedInfo.TimeCompression.TimeOffset
		baseTime = time.Unix(final, 0)
	}

	nowLocal := baseTime.In(loc)

	// 2. Apply Looping Constraint
	if extendedInfo.TimeLoop.Enabled {
		// Apply Offset if any
		nowLocal = nowLocal.Add(time.Duration(int64(time.Second) * extendedInfo.TimeLoop.LoopOffset))

		switch extendedInfo.TimeLoop.Mode {

		case device.TimeLoopSeptemberOctober:
			rangeStart := time.Date(nowLocal.Year(), 9, 15, 23, 45, 0, 0, loc)
			rangeEnd := time.Date(nowLocal.Year()+1, 3, 15, 0, 15, 0, 0, loc)
			elapsed := nowLocal.Sub(rangeStart)
			interval := rangeEnd.Sub(rangeStart)
			delta := elapsed % interval
			nowLocal = rangeStart.Add(delta)

		case device.TimeLoopDecemberJanuary:
			rangeStart := time.Date(nowLocal.Year(), 12, 31, 23, 45, 0, 0, loc)
			rangeEnd := time.Date(nowLocal.Year()+1, 1, 1, 0, 15, 0, 0, loc)
			elapsed := nowLocal.Sub(rangeStart)
			interval := rangeEnd.Sub(rangeStart)
			delta := elapsed % interval
			nowLocal = rangeStart.Add(delta)

		case device.TimeLoopDstIn:
			// NYI
		case device.TimeLoopDstOut:
			// NYI

		case device.TimeLoopCustom:
			rangeStart := time.Unix(extendedInfo.TimeLoop.LoopStart, 0)
			rangeEnd := time.Unix(extendedInfo.TimeLoop.LoopEnd, 0)
			elapsed := nowLocal.Sub(rangeStart)
			interval := rangeEnd.Sub(rangeStart)
			delta := elapsed % interval
			nowLocal = rangeStart.Add(delta)

		case device.TimeLoopBackendDriver:
			// @todo make call to elixir backend to determine reported time.
		}
	}

	weatherTime.DateTime = nowLocal

	// Getting the hour in that timezone
	localHour, _ := strconv.Atoi(nowLocal.Format("15"))
	if localHour >= 7 && localHour < 19 {
		weatherTime.DayInfo = DayInfoDay
	} else {
		weatherTime.DayInfo = DayInfoNight
	}

	if localHour >= 0 && localHour < 6 {
		weatherTime.HourRange = "00"
	} else if localHour >= 6 && localHour < 12 {
		weatherTime.HourRange = "01"
	} else if localHour >= 12 && localHour < 18 {
		weatherTime.HourRange = "02"
	} else {
		weatherTime.HourRange = "03"
	}

	weatherTime.LocalDate = nowLocal.Format("02:01:2006")
	weatherTime.LocalTime = nowLocal.Format("15:04:05")
	weatherTime.Iso8601 = nowLocal.Format("2006-01-02T15:04:05-0700")
	return weatherTime, nil
}

//----------------------------------------------
// Round function, will use this as we are using go 1.7 and round is included in go 1.10
// After we upgrade to newer version of Go this needs to be removed
//----------------------------------------------
/**
 * @brief
 */
func Round(x float64) int {
	t := math.Trunc(x)
	if math.Abs(x-t) >= 0.5 {
		return int(t + math.Copysign(1, x))
	}

	return int(t)
}

//----------------------------------------------
// Local Funcs
//----------------------------------------------

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func findACKeyFromList(ll []LocationLookupResponse, city string, country string, adminArea string, tz string) bool {

	found := false
	for _, loc := range ll {
		if loc.CityName == city && loc.CountryCode == country && loc.AdministrativeArea == adminArea && loc.TzName == tz {
			found = true
			break
		}
	}

	return found
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func httpAccuGetLocation(path string, queryString string) []PostalCodeResponse {

	var pc []PostalCodeResponse
	var Url *url.URL
	Url, _ = url.Parse(AccuBaseUrl)

	Url.Path += path
	parameters := url.Values{}
	parameters.Add("apikey", AccuApiKey)
	parameters.Add("q", queryString)
	Url.RawQuery = parameters.Encode()

	resp, err := http.Get(Url.String())
	if err != nil {
		log.Printf(err.Error())
		return pc
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &pc)

	return pc

}

//----------------------------------------------
// Day api forecast query to Accuweather
//----------------------------------------------
/**
 * @brief
 */
func getNullableDailyForecast(locationKey string, timeZone string, period string, weatherTime WeatherTime) (NullableDailyForecast, error) {
	key := "forecast:" + locationKey + ":" + period + ":" + weatherTime.HourRange + "_" + weatherTime.LocalDate
	path := "/forecasts/v1/daily/" + period + "/" + locationKey

	var response NullableDailyForecast
	data, err := common.RedisInstance.GetCachedData(key)
	if err != nil {
		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)
		nowStr := time.Now().Format("02:01:2006 15:04:05")
		common.RedisInstance.SaveRedisData([]byte(nowStr), "forecastupdate:"+period+":"+locationKey, 0)
	}

	retryCount := 0
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Printf("Json Error raised %v", err)
	}
	for len(response.DailyForecasts) == 0 && retryCount < MaxRetries {
		// @TODO Note this logic may result in many processes slamming Accuweather at once.
		retryCount = retryCount + 1
		log.Printf("Could not correctly fetch the weather forecast, retry number : %s", retryCount)
		time.Sleep(time.Duration(retryCount) * time.Second)
		data, _ = httpAccuGetAndCache(path, key, ForecastExpireHours*time.Hour)
		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Printf("Json Error raised %v", err)
		}
	}

	if len(response.DailyForecasts) == 0 {
		if err == nil {
			err = errors.New("incomplete data")
		}
		return response, err
	}

	if weatherTime.DayInfo == DayInfoDay {
		response.Today = &response.DailyForecasts[0].Day
	} else {
		response.Today = &response.DailyForecasts[0].Night
	}

	// Populate Air and Pollen Map
	for index, forecast := range response.DailyForecasts {
		response.DailyForecasts[index].AirAndPollenMap = make(map[string]int)
		response.DailyForecasts[index].AirAndPollenCategoryMap = make(map[string]string)
		for _, airAndPollen := range forecast.AirAndPollen {
			if airAndPollen.Value.Valid {
				response.DailyForecasts[index].AirAndPollenMap[airAndPollen.Name.String] = int(airAndPollen.Value.Int64)
			}
			if airAndPollen.Category.Valid {
				response.DailyForecasts[index].AirAndPollenCategoryMap[airAndPollen.Name.String] = airAndPollen.Category.String
			}
		}
	}

	return response, err
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func getNWSInfo(zip string) map[string]string {
	var severeComponentMap map[string]string
	if strings.TrimSpace(zip) != "" {
		key := "nwsforecast:" + zip
		data, err := common.RedisInstance.GetCachedData(key)
		if err != nil {
			severeComponentMap = nws.GetSevereComponentMap(zip)
			dataBytes, _ := json.Marshal(severeComponentMap)
			common.RedisInstance.SaveRedisData(dataBytes, key, 24*time.Hour)
		} else {
			err = json.Unmarshal(data, &severeComponentMap)
			if err != nil {
				log.Printf("NWS Unmarshall Error: %v", err)
			}
		}
	}
	return severeComponentMap
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func getNWSInfoV2(accuLocation PostalCodeResponse) map[string]string {
	var severeComponentMap map[string]string
	if strings.TrimSpace(accuLocation.PrimaryPostalCode) != "" && accuLocation.Country.ID == "US" {
		key := "nwsforecast:" + accuLocation.PrimaryPostalCode
		data, err := common.RedisInstance.GetCachedData(key)
		if err != nil {
			severeComponentMap = nws.GetSevereComponentMap(accuLocation.PrimaryPostalCode)
			dataBytes, _ := json.Marshal(severeComponentMap)
			common.RedisInstance.SaveRedisData(dataBytes, key, 24*time.Hour)
		} else {
			err = json.Unmarshal(data, &severeComponentMap)
			if err != nil {
				log.Printf("NWS Unmarshall Error: %v", err)
			}
		}
	}
	return severeComponentMap
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func queryAccuCurrentForecastAPI(locationKey string, weatherTime WeatherTime) []AccuCurrentForecastResponse {

	key := locationKey + "_AccuCurrentForecast_" + weatherTime.LocalDate

	path := "/currentconditions/v1/" + locationKey
	data, err := common.RedisInstance.GetCachedData(key)

	if err != nil {
		log.Printf("KeyNotFound : %s", key)
		data, _ = httpAccuGetAndCache(path, key, 3*time.Hour)

		// Persist Last cache update
		nowStr := time.Now().Format("02:01:2006 15:04:05")
		common.RedisInstance.SaveRedisData([]byte(nowStr), "forecastupdate:current:"+locationKey, 0)
	}

	var accuCurrentForecastResponse []AccuCurrentForecastResponse
	retryCount := 0
	json.Unmarshal(data, &accuCurrentForecastResponse)

	for len(accuCurrentForecastResponse) == 0 && retryCount < MaxRetries {
		retryCount = retryCount + 1
		log.Printf("Could not correctly fetch the weather forecast, retry number : %d path: %s", retryCount, path)
		time.Sleep(time.Duration(retryCount) * time.Second)

		data, _ = httpAccuGetAndCache(path, key, 3*time.Hour)
		json.Unmarshal(data, &accuCurrentForecastResponse)
	}

	if retryCount == MaxRetries {
		return accuCurrentForecastResponse
	}

	return accuCurrentForecastResponse
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func NullablequeryAccuCurrentForecastAPI(locationKey string, weatherTime WeatherTime, nws map[string]string) (NullableAccuCurrentForecastResponse, error) {
	// @note LocalDate part of key is DD:MM:YYYY string.
	// @note cache invalidation logic could be better here.
	key := locationKey + "_AccuCurrentForecast_" + weatherTime.LocalDate
	path := "/currentconditions/v1/" + locationKey
	data, cacheErr := common.RedisInstance.GetCachedData(key)
	if cacheErr != nil {
		data, _ = httpAccuGetAndCache(path, key, 3*time.Hour)
		// Persist Last cache update
		nowStr := time.Now().Format("02:01:2006 15:04:05")
		common.RedisInstance.SaveRedisData([]byte(nowStr), "forecastupdate:current:"+locationKey, 0)
	}

	var accuCurrentForecastResponse []NullableAccuCurrentForecastResponse
	retryCount := 0

	err := json.Unmarshal(data, &accuCurrentForecastResponse)
	for len(accuCurrentForecastResponse) == 0 && retryCount < MaxRetries {
		retryCount = retryCount + 1
		log.Printf("Could not correctly fetch the weather forecast, retry number : %s", retryCount)
		time.Sleep(time.Duration(retryCount) * time.Second)

		data, _ = httpAccuGetAndCache(path, key, 3*time.Hour)
		err = json.Unmarshal(data, &accuCurrentForecastResponse)
	}

	if len(accuCurrentForecastResponse) > 0 {
		// Adapter: Hail & Tornado Probability
		if accuCurrentForecastResponse != nil {
			// @todo deal with parse failure.
			// Hail Probability
			if val, ok := nws["hail"]; ok {
				if ok {
					p, e := strconv.Atoi(val)
					if e == nil {
						accuCurrentForecastResponse[0].HailProbability = null.NewInt(int64(p), true)
					}
				}
			}
			// Tornado Probability
			if val, ok := nws["tornadoes"]; ok {
				if ok {
					p, e := strconv.Atoi(val)
					if e == nil {
						accuCurrentForecastResponse[0].TornadoProbability = null.NewInt(int64(p), true)
					}
				}
			}
		}
		// End Adapter: Hail & Tornado Probability
		return accuCurrentForecastResponse[0], nil
	} else {
		nullResponse := NullableAccuCurrentForecastResponse{}
		return nullResponse, err
	}
}

//----------------------------------------------
//
//----------------------------------------------
/**
 * @brief
 */
func GetLocalDateAndHour(timeZone string) (WeatherTime, error) {

	weatherTime := WeatherTime{}

	// Load the specific location
	var loc *time.Location
	var err error
	LocationMapMutex.Lock()
	loc, keyFound := LocationMap[timeZone]
	LocationMapMutex.Unlock()

	if !keyFound {
		loc, err = time.LoadLocation(timeZone)
		if err != nil {
			log.Printf("Error Loading the location for timezone %s : %s", timeZone, err.Error())
			return weatherTime, err
		}
		LocationMap[timeZone] = loc
	}

	//set Location
	nowLocal := time.Now().In(loc)

	// Getting the hour in that timezone
	localHour, _ := strconv.Atoi(nowLocal.Format("15"))
	if localHour >= 7 && localHour < 19 {
		weatherTime.DayInfo = DayInfoDay
	} else {
		weatherTime.DayInfo = DayInfoNight
	}

	if localHour >= 0 && localHour < 6 {
		weatherTime.HourRange = "00"
	} else if localHour >= 6 && localHour < 12 {
		weatherTime.HourRange = "01"
	} else if localHour >= 12 && localHour < 18 {
		weatherTime.HourRange = "02"
	} else {
		weatherTime.HourRange = "03"
	}

	weatherTime.LocalDate = nowLocal.Format("02:01:2006")
	weatherTime.LocalTime = nowLocal.Format("15:04:05")
	weatherTime.Iso8601 = nowLocal.Format("2006-01-02T15:04:05-0700")
	return weatherTime, nil
}

//----------------------------------------------
// Generic Code to call Accuweather API service. All api calls to Accu should be handled by this method
// In cache of no error, the raw data is immediately cached
//----------------------------------------------
/**
 * @brief
 */
func httpAccuGetAndCache(path string, cacheKey string, cacheDuration time.Duration) ([]byte, error) {
	var Url *url.URL
	Url, _ = url.Parse(AccuBaseUrl)

	Url.Path += path
	parameters := url.Values{}
	parameters.Add("apikey", AccuApiKey)
	parameters.Add("details", "true")
	parameters.Add("metric", "true")
	Url.RawQuery = parameters.Encode()

	log.Println("[Get Forecast] " + Url.String())

	resp, err := http.Get(Url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//redisInstance := &cache.RedisInstance{RedisSession: main.RedisClient}
	common.RedisInstance.SaveRedisData(body, cacheKey, cacheDuration)
	return body, nil
}
