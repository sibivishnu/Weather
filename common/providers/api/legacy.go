package api

import "time"

/*!
 * @brief Geographic location details.
 */
type PostalCodeResponse struct {
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

/*!
 * @brief Geographic Country
 */
type AccuCountry struct {
	ID            string
	LocalizedName string
	EnglishName   string
}

/*!
 * @brief Geographic Time Zone
 */
type AccuTimeZone struct {
	Code             string
	Name             string
	GmtOffset        float64
	IsDaylightSaving bool
	NextOffsetChange string
}

/*!
 * @brief Forecast Time Details
 */
type WeatherTime struct {
	LocalTime string
	LocalDate string
	DayInfo   string
	HourRange string
	Iso8601   string
	DateTime  time.Time
}

/*!
 * @brief Hourly Forecast Record
 */
type AccuHourlyForecastResponse struct {
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

/*!
 * @brief Current Forecast Record
 */
type AccuCurrentForecastResponse struct {
	LocalObservationDateTime string
	EpochTime                int
	WeatherText              string
	WeatherIcon              int
	IsDayTime                bool
	Temperature              CurrentTemp
}

/*!
 * @brief Current Temperature Record
 */
type CurrentTemp struct {
	Metric   Temperature
	Imperial Temperature
}

/*!
 * @brief Wind speed and direction record.
 */
type Wind struct {
	Speed     Temperature
	Direction Direction
}

/*!
 * @brief Angular Direction.
 */
type Direction struct {
	Degrees   int
	Localized string
	English   string
}

/*!
 * @brief Forecast Headline String Record
 */
type AccuHeadline struct {
	EffectiveDate      string
	EffectiveEpochDate int
	Severity           int
	Text               string
	Category           string
	EndDate            string
	EndEpochDate       int
}

/*!
 * @brief Sun/Moon rise/phase/set details.
 */
type Sun struct {
	Rise      string
	EpochRise int64
	Set       string
	EpochSet  int64
	Phase     string
	Age       int
}

/*!
 * @brief Air and Pollen Category
 */
type AirAndPollen struct {
	Name          string
	Value         int
	Category      string
	CategoryValue int
	Type          string
}

/*!
 * @brief Daily Forecast Struct
 */
type DailyForecast struct {
	Headline       AccuHeadline
	DailyForecasts []AccuDailyForecast
}

/*!
 * @brief Legacy Category 1 custom encoding template.
 */
type AccuTemplateCat1Struct struct {
	DateStr       string
	TimeStr       string
	GmtOffset     float64
	Accu1d        *DailyForecast
	Headline      *AccuHeadline
	DailyForecast *AccuDailyForecast
	FlowControl   int
}

/*!
 * @brief Legacy Category 2 custom encoding template.
 */
type AccuTemplateCat2Struct struct {
	DateStr         string
	TimeStr         string
	GmtOffset       float64
	Accu1d          *DailyForecast
	Headline        *AccuHeadline
	DailyForecast   *AccuDailyForecast
	CurrentForecast *AccuCurrentForecastResponse
	FlowControl     int
}

/*!
 * @brief Legacy Category 3 custom encoding template.
 */
type AccuTemplateCat3Struct struct {
	DateStr       string
	TimeStr       string
	ForecastTime  string
	GmtOffset     float64
	DailyForecast *AccuDailyForecast
	Accu7d        *[]AccuDailyForecast
	Accu24h       *[]AccuHourlyForecastResponse
	FlowControl   int
}

/*!
 * @brief High Low temperature record.
 */
type MinMaxTemperature struct {
	Minimum Temperature
	Maximum Temperature
}

/*!
 * @brief temperature summary
 * More details needed.
 */
type SummaryTemperature struct {
	Heating Temperature
	Cooling Temperature
}

/*!
 * @brief day and night forecast records of daily forecast entry.
 */
type DayNightData struct {
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

/*!
 * @brief Temperature record (value and units)
 */
type Temperature struct {
	Value      float64
	ValueRound int
	Unit       string
	UnitType   int
}

/*!
 * @brief Daily Forecast Record
 */
type AccuDailyForecast struct {
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
