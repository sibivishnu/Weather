package weather_api

//==============================================
// CopyRight 2020 La Crosse Technology, LTD.
//==============================================


//==============================================
// Imports
//==============================================
import (
	"../../const/accuweather"
	"../../const/firmware"
	"../../const/legacy_firmware"
	"../../const/mcu"
	"sync"
	"time"
)


var (

	/**
	 * @brief Maps Accuweather Icons into Display Icons.
	 */
	AccuweatherIcons = map[int]AccuIcon{
	const_accuweather.IconError:                              {IconNumber: const_accuweather.IconError, DisplayIcon: 0, Icon: "", Day: false, Night: false, Text: "Error"},
	const_accuweather.IconDaySunny:                           {IconNumber: const_accuweather.IconDaySunny, DisplayIcon: const_mcu.IconSunny, Icon: "https://developer.accuweather.com/sites/default/files/01-s.png", Day: true, Night: false, Text: "Sunny"},
	const_accuweather.IconDayMostlySunny:                     {IconNumber: const_accuweather.IconDayMostlySunny, DisplayIcon: const_mcu.IconSunny, Icon: "https://developer.accuweather.com/sites/default/files/02-s.png", Day: true, Night: false, Text: "Mostly Sunny"},
	const_accuweather.IconDayPartlySunny:                     {IconNumber: const_accuweather.IconDayPartlySunny, DisplayIcon: const_mcu.IconPartlyCloudy, Icon: "https://developer.accuweather.com/sites/default/files/03-s.png", Day: true, Night: false, Text: "Partly Sunny"},
	const_accuweather.IconDayIntermittentClouds:              {IconNumber: const_accuweather.IconDayIntermittentClouds, DisplayIcon: const_mcu.IconPartlyCloudy, Icon: "https://developer.accuweather.com/sites/default/files/04-s.png", Day: true, Night: false, Text: "Intermittent Clouds"},
	const_accuweather.IconDayHazySunshine:                    {IconNumber: const_accuweather.IconDayHazySunshine, DisplayIcon: const_mcu.IconPartlyCloudy, Icon: "https://developer.accuweather.com/sites/default/files/05-s.png", Day: true, Night: false, Text: "Hazy Sunshine"},
	const_accuweather.IconDayMostlyCloudy:                    {IconNumber: const_accuweather.IconDayMostlyCloudy, DisplayIcon: const_mcu.IconCloudy, Icon: "https://developer.accuweather.com/sites/default/files/06-s.png", Day: true, Night: false, Text: "Mostly Cloudy"},
	const_accuweather.IconCloudy:                             {IconNumber: const_accuweather.IconCloudy, DisplayIcon: const_mcu.IconCloudy, Icon: "https://developer.accuweather.com/sites/default/files/07-s.png", Day: true, Night: true, Text: "Cloudy"},
	const_accuweather.IconDreary:                             {IconNumber: const_accuweather.IconDreary, DisplayIcon: const_mcu.IconCloudy, Icon: "https://developer.accuweather.com/sites/default/files/08-s.png", Day: true, Night: true, Text: "Dreary (Overcast)"},
	const_accuweather.IconFog:                                {IconNumber: const_accuweather.IconFog, DisplayIcon: const_mcu.IconFog, Icon: "https://developer.accuweather.com/sites/default/files/11-s.png", Day: true, Night: true, Text: "Fog"},
	const_accuweather.IconShowers:                            {IconNumber: const_accuweather.IconShowers, DisplayIcon: const_mcu.IconRain, Icon: "https://developer.accuweather.com/sites/default/files/12-s.png", Day: true, Night: true, Text: "Showers"},
	const_accuweather.IconDayMostlyCloudyWithShowers:         {IconNumber: const_accuweather.IconDayMostlyCloudyWithShowers, DisplayIcon: const_mcu.IconRain, Icon: "https://developer.accuweather.com/sites/default/files/13-s.png", Day: true, Night: false, Text: "Mostly Cloudy w/ Showers"},
	const_accuweather.IconDayPartlySunnyWithShowers:          {IconNumber: const_accuweather.IconDayPartlySunnyWithShowers, DisplayIcon: const_mcu.IconPartlyCloudyWithRain, Icon: "https://developer.accuweather.com/sites/default/files/14-s.png", Day: true, Night: false, Text: "Partly Sunny w/ Showers"},
	const_accuweather.IconThunderstorms:                      {IconNumber: const_accuweather.IconThunderstorms, DisplayIcon: const_mcu.IconStorm, Icon: "https://developer.accuweather.com/sites/default/files/15-s.png", Day: true, Night: true, Text: "T-Storms"},
	const_accuweather.IconDayMostlyCloudyWithThunderStorms:   {IconNumber: const_accuweather.IconDayMostlyCloudyWithThunderStorms, DisplayIcon: const_mcu.IconStorm, Icon: "https://developer.accuweather.com/sites/default/files/16-s.png", Day: true, Night: false, Text: "Mostly Cloudy w/ T-Storms"},
	const_accuweather.IconDayPartlySunnyWithThunderstorms:    {IconNumber: const_accuweather.IconDayPartlySunnyWithThunderstorms, DisplayIcon: const_mcu.IconStorm, Icon: "https://developer.accuweather.com/sites/default/files/17-s.png", Day: true, Night: false, Text: "Partly Sunny w/ T-Storms"},
	const_accuweather.IconRain:                               {IconNumber: const_accuweather.IconRain, DisplayIcon: const_mcu.IconRain, Icon: "https://developer.accuweather.com/sites/default/files/18-s.png", Day: true, Night: true, Text: "Rain"},
	const_accuweather.IconFlurries:                           {IconNumber: const_accuweather.IconFlurries, DisplayIcon: const_mcu.IconFlurries, Icon: "https://developer.accuweather.com/sites/default/files/19-s.png", Day: true, Night: true, Text: "Flurries"},
	const_accuweather.IconDayMostlyCloudyWithFlurries:        {IconNumber: const_accuweather.IconDayMostlyCloudyWithFlurries, DisplayIcon: const_mcu.IconFlurries, Icon: "https://developer.accuweather.com/sites/default/files/20-s.png", Day: true, Night: false, Text: "Mostly Cloudy w/ Flurries"},
	const_accuweather.IconDayPartlySunnyWithFlurries:         {IconNumber: const_accuweather.IconDayPartlySunnyWithFlurries, DisplayIcon: const_mcu.IconFlurries, Icon: "https://developer.accuweather.com/sites/default/files/21-s.png", Day: true, Night: false, Text: "Partly Sunny w/ Flurries"},
	const_accuweather.IconSnow:                               {IconNumber: const_accuweather.IconSnow, DisplayIcon: const_mcu.IconSnow, Icon: "https://developer.accuweather.com/sites/default/files/22-s.png", Day: true, Night: true, Text: "Snow"},
	const_accuweather.IconDayMostlyCloudyWithSnow:            {IconNumber: const_accuweather.IconDayMostlyCloudyWithSnow, DisplayIcon: const_mcu.IconSnow, Icon: "https://developer.accuweather.com/sites/default/files/23-s.png", Day: true, Night: false, Text: "Mostly Cloudy w/ Snow"},
	const_accuweather.IconIce:                                {IconNumber: const_accuweather.IconIce, DisplayIcon: const_mcu.IconIce, Icon: "https://developer.accuweather.com/sites/default/files/24-s.png", Day: true, Night: true, Text: "Ice"},
	const_accuweather.IconSleet:                              {IconNumber: const_accuweather.IconSleet, DisplayIcon: const_mcu.IconFreezingRain, Icon: "https://developer.accuweather.com/sites/default/files/25-s.png", Day: true, Night: true, Text: "Sleet"},
	const_accuweather.IconFreezingRain:                       {IconNumber: const_accuweather.IconFreezingRain, DisplayIcon: const_mcu.IconFreezingRain, Icon: "https://developer.accuweather.com/sites/default/files/26-s.png", Day: true, Night: true, Text: "Freezing Rain"},
	const_accuweather.IconRainAndSnow:                        {IconNumber: const_accuweather.IconRainAndSnow, DisplayIcon: const_mcu.IconFreezingRain, Icon: "https://developer.accuweather.com/sites/default/files/29-s.png", Day: true, Night: true, Text: "Rain and Snow"},
	const_accuweather.IconHot:                                {IconNumber: const_accuweather.IconHot, DisplayIcon: const_mcu.IconSunny, Icon: "https://developer.accuweather.com/sites/default/files/30-s.png", Day: true, Night: true, Text: "Hot"},
	const_accuweather.IconCold:                               {IconNumber: const_accuweather.IconCold, DisplayIcon: const_mcu.IconSunny, Icon: "https://developer.accuweather.com/sites/default/files/31-s.png", Day: true, Night: true, Text: "Cold"},
	const_accuweather.IconWindy:                              {IconNumber: const_accuweather.IconWindy, DisplayIcon: const_mcu.IconWind, Icon: "https://developer.accuweather.com/sites/default/files/32-s.png", Day: true, Night: true, Text: "Windy"},
	const_accuweather.IconNightClear:                         {IconNumber: const_accuweather.IconNightClear, DisplayIcon: const_mcu.IconSunny, Icon: "https://developer.accuweather.com/sites/default/files/33-s.png", Day: false, Night: true, Text: "Clear"},
	const_accuweather.IconNightMostlyClear:                   {IconNumber: const_accuweather.IconNightMostlyClear, DisplayIcon: const_mcu.IconSunny, Icon: "https://developer.accuweather.com/sites/default/files/34-s.png", Day: false, Night: true, Text: "Mostly Clear"},
	const_accuweather.IconNightPartlyCloudy:                  {IconNumber: const_accuweather.IconNightPartlyCloudy, DisplayIcon: const_mcu.IconPartlyCloudy, Icon: "https://developer.accuweather.com/sites/default/files/35-s.png", Day: false, Night: true, Text: "Partly Cloudy"},
	const_accuweather.IconNightIntermittentClouds:            {IconNumber: const_accuweather.IconNightIntermittentClouds, DisplayIcon: const_mcu.IconPartlyCloudy, Icon: "https://developer.accuweather.com/sites/default/files/36-s.png", Day: false, Night: true, Text: "Intermittent Clouds"},
	const_accuweather.IconNightHazyMoonlight:                 {IconNumber: const_accuweather.IconNightHazyMoonlight, DisplayIcon: const_mcu.IconPartlyCloudy, Icon: "https://developer.accuweather.com/sites/default/files/37-s.png", Day: false, Night: true, Text: "Hazy Moonlight"},
	const_accuweather.IconNightMostlyCloudy:                  {IconNumber: const_accuweather.IconNightMostlyCloudy, DisplayIcon: const_mcu.IconCloudy, Icon: "https://developer.accuweather.com/sites/default/files/38-s.png", Day: false, Night: true, Text: "Mostly Cloudy"},
	const_accuweather.IconNightPartlyCloudyWithShowers:       {IconNumber: const_accuweather.IconNightPartlyCloudyWithShowers, DisplayIcon: const_mcu.IconPartlyCloudyWithRain, Icon: "https://developer.accuweather.com/sites/default/files/39-s.png", Day: false, Night: true, Text: "Partly Cloudy w/ Showers"},
	const_accuweather.IconNightMostlyCloudyWithShowers:       {IconNumber: const_accuweather.IconNightMostlyCloudyWithShowers, DisplayIcon: const_mcu.IconRain, Icon: "https://developer.accuweather.com/sites/default/files/40-s.png", Day: false, Night: true, Text: "Mostly Cloudy w/ Showers"},
	const_accuweather.IconNightPartlyCloudyWithThunderstorms: {IconNumber: const_accuweather.IconNightPartlyCloudyWithThunderstorms, DisplayIcon: const_mcu.IconStorm, Icon: "https://developer.accuweather.com/sites/default/files/41-s.png", Day: false, Night: true, Text: "Partly Cloudy w/ T-Storms"},
	const_accuweather.IconNightMostlyCloudyWithThunderStorms: {IconNumber: const_accuweather.IconNightMostlyCloudyWithThunderStorms, DisplayIcon: const_mcu.IconStorm, Icon: "https://developer.accuweather.com/sites/default/files/42-s.png", Day: false, Night: true, Text: "Mostly Cloudy w/ T-Storms"},
	const_accuweather.IconNightMostlyCloudyWithFlurries:      {IconNumber: const_accuweather.IconNightMostlyCloudyWithFlurries, DisplayIcon: const_mcu.IconFlurries, Icon: "https://developer.accuweather.com/sites/default/files/43-s.png", Day: false, Night: true, Text: "Mostly Cloudy w/ Flurries"},
	const_accuweather.IconNightMostlyCloudyWithSnow:          {IconNumber: const_accuweather.IconNightMostlyCloudyWithSnow, DisplayIcon: const_mcu.IconSnow, Icon: "https://developer.accuweather.com/sites/default/files/44-s.png", Day: false, Night: true, Text: "Mostly Cloudy w/ Snow "},
}
)

//----------------------------------------------
// Globals - Tables - Weather Icon Details
//----------------------------------------------
var (

	//----------------------------------------------
	// Globals - Tables - Accuweather to Lax Icon Codes
	//----------------------------------------------

	/**
	 * @brief Maps accuweather icons into display icons with fewer details than AccuweatherIcons table.
	 */
	IconMap = map[int]int{
		const_accuweather.IconDaySunny:                           AccuweatherIcons[const_accuweather.IconDaySunny].DisplayIcon,
		const_accuweather.IconDayMostlySunny:                     AccuweatherIcons[const_accuweather.IconDayMostlySunny].DisplayIcon,
		const_accuweather.IconDayPartlySunny:                     AccuweatherIcons[const_accuweather.IconDayPartlySunny].DisplayIcon,
		const_accuweather.IconDayIntermittentClouds:              AccuweatherIcons[const_accuweather.IconDayIntermittentClouds].DisplayIcon,
		const_accuweather.IconDayHazySunshine:                    AccuweatherIcons[const_accuweather.IconDayHazySunshine].DisplayIcon,
		const_accuweather.IconDayMostlyCloudy:                    AccuweatherIcons[const_accuweather.IconDayMostlyCloudy].DisplayIcon,
		const_accuweather.IconCloudy:                             AccuweatherIcons[const_accuweather.IconCloudy].DisplayIcon,
		const_accuweather.IconDreary:                             AccuweatherIcons[const_accuweather.IconDreary].DisplayIcon,
		const_accuweather.IconFog:                                AccuweatherIcons[const_accuweather.IconFog].DisplayIcon,
		const_accuweather.IconShowers:                            AccuweatherIcons[const_accuweather.IconShowers].DisplayIcon,
		const_accuweather.IconDayMostlyCloudyWithShowers:         AccuweatherIcons[const_accuweather.IconDayMostlyCloudyWithShowers].DisplayIcon,
		const_accuweather.IconDayPartlySunnyWithShowers:          AccuweatherIcons[const_accuweather.IconDayPartlySunnyWithShowers].DisplayIcon,
		const_accuweather.IconThunderstorms:                      AccuweatherIcons[const_accuweather.IconThunderstorms].DisplayIcon,
		const_accuweather.IconDayMostlyCloudyWithThunderStorms:   AccuweatherIcons[const_accuweather.IconDayMostlyCloudyWithThunderStorms].DisplayIcon,
		const_accuweather.IconDayPartlySunnyWithThunderstorms:    AccuweatherIcons[const_accuweather.IconDayPartlySunnyWithThunderstorms].DisplayIcon,
		const_accuweather.IconRain:                               AccuweatherIcons[const_accuweather.IconRain].DisplayIcon,
		const_accuweather.IconFlurries:                           AccuweatherIcons[const_accuweather.IconFlurries].DisplayIcon,
		const_accuweather.IconDayMostlyCloudyWithFlurries:        AccuweatherIcons[const_accuweather.IconDayMostlyCloudyWithFlurries].DisplayIcon,
		const_accuweather.IconDayPartlySunnyWithFlurries:         AccuweatherIcons[const_accuweather.IconDayPartlySunnyWithFlurries].DisplayIcon,
		const_accuweather.IconSnow:                               AccuweatherIcons[const_accuweather.IconSnow].DisplayIcon,
		const_accuweather.IconDayMostlyCloudyWithSnow:            AccuweatherIcons[const_accuweather.IconDayMostlyCloudyWithSnow].DisplayIcon,
		const_accuweather.IconIce:                                AccuweatherIcons[const_accuweather.IconIce].DisplayIcon,
		const_accuweather.IconSleet:                              AccuweatherIcons[const_accuweather.IconSleet].DisplayIcon,
		const_accuweather.IconFreezingRain:                       AccuweatherIcons[const_accuweather.IconFreezingRain].DisplayIcon,
		const_accuweather.IconRainAndSnow:                        AccuweatherIcons[const_accuweather.IconRainAndSnow].DisplayIcon,
		const_accuweather.IconHot:                                AccuweatherIcons[const_accuweather.IconHot].DisplayIcon,
		const_accuweather.IconCold:                               AccuweatherIcons[const_accuweather.IconCold].DisplayIcon,
		const_accuweather.IconWindy:                              AccuweatherIcons[const_accuweather.IconWindy].DisplayIcon,
		const_accuweather.IconNightClear:                         AccuweatherIcons[const_accuweather.IconNightClear].DisplayIcon,
		const_accuweather.IconNightMostlyClear:                   AccuweatherIcons[const_accuweather.IconNightMostlyClear].DisplayIcon,
		const_accuweather.IconNightPartlyCloudy:                  AccuweatherIcons[const_accuweather.IconNightPartlyCloudy].DisplayIcon,
		const_accuweather.IconNightIntermittentClouds:            AccuweatherIcons[const_accuweather.IconNightIntermittentClouds].DisplayIcon,
		const_accuweather.IconNightHazyMoonlight:                 AccuweatherIcons[const_accuweather.IconNightHazyMoonlight].DisplayIcon,
		const_accuweather.IconNightMostlyCloudy:                  AccuweatherIcons[const_accuweather.IconNightMostlyCloudy].DisplayIcon,
		const_accuweather.IconNightPartlyCloudyWithShowers:       AccuweatherIcons[const_accuweather.IconNightPartlyCloudyWithShowers].DisplayIcon,
		const_accuweather.IconNightMostlyCloudyWithShowers:       AccuweatherIcons[const_accuweather.IconNightMostlyCloudyWithShowers].DisplayIcon,
		const_accuweather.IconNightPartlyCloudyWithThunderstorms: AccuweatherIcons[const_accuweather.IconNightPartlyCloudyWithThunderstorms].DisplayIcon,
		const_accuweather.IconNightMostlyCloudyWithThunderStorms: AccuweatherIcons[const_accuweather.IconNightMostlyCloudyWithThunderStorms].DisplayIcon,
		const_accuweather.IconNightMostlyCloudyWithFlurries:      AccuweatherIcons[const_accuweather.IconNightMostlyCloudyWithFlurries].DisplayIcon,
		const_accuweather.IconNightMostlyCloudyWithSnow:          AccuweatherIcons[const_accuweather.IconNightMostlyCloudyWithSnow].DisplayIcon,
	}

	//----------------------------------------------
	// Globals - Tables - Moonphase Map
	//----------------------------------------------
	/**
	 * @brief MoonPhase Mappings between accuweather and MCU values.
	 */
	MoonPhraseMap = map[string]int{
		const_accuweather.MoonPhaseNew:            legacy_const_firmware.MoonPhaseNew,
		const_accuweather.MoonPhaseWaxingCrescent: legacy_const_firmware.MoonPhaseWaxingCrescent,
		const_accuweather.MoonPhaseFirst:          legacy_const_firmware.MoonPhaseFirst,
		const_accuweather.MoonPhaseWaxingGibbous:  legacy_const_firmware.MoonPhaseWaxingGibbous,
		const_accuweather.MoonPhaseFull:           legacy_const_firmware.MoonPhaseFull,
		const_accuweather.MoonPhaseWaningGibbous:  legacy_const_firmware.MoonPhaseWaningGibbous,
		const_accuweather.MoonPhaseLast:           legacy_const_firmware.MoonPhaseLast,
		const_accuweather.MoonPhaseWaningCrescent: legacy_const_firmware.MoonPhaseWaningCrescent,
	}

	//----------------------------------------------
	// Globals - Tables - Weather Category Map
	//----------------------------------------------
	/**
	 * @brief WeatherCategory Mappings
	 */
	ValueToEnumMap = map[string]int{
		const_accuweather.WeatherCategoryLow:           legacy_const_firmware.WeatherCategoryLow,
		const_accuweather.WeatherCategoryHigh:          legacy_const_firmware.WeatherCategoryHigh,
		const_accuweather.WeatherCategoryGood:          legacy_const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryModerate:      legacy_const_firmware.WeatherCategoryModerate,
		const_accuweather.WeatherCategoryUnhealthy:     legacy_const_firmware.WeatherCategoryUnhealthy,
		const_accuweather.WeatherCategoryHazardous:     legacy_const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryPoor:          legacy_const_firmware.WeatherCategoryPoor,
		const_accuweather.WeatherCategoryFair:          legacy_const_firmware.WeatherCategoryFair,
		const_accuweather.WeatherCategoryVeryGood:      legacy_const_firmware.WeatherCategoryVeryGood,
		const_accuweather.WeatherCategoryExcellent:     legacy_const_firmware.WeatherCategoryExcellent,
		const_accuweather.WeatherCategoryVeryHigh:      legacy_const_firmware.WeatherCategoryVeryHigh,
		const_accuweather.WeatherCategoryExtreme:       legacy_const_firmware.WeatherCategoryExtreme,
		const_accuweather.WeatherCategoryUnlikely:      legacy_const_firmware.WeatherCategoryUnlikely,
		const_accuweather.WeatherCategoryWatch:         legacy_const_firmware.WeatherCategoryWatch,
		const_accuweather.WeatherCategoryAdvisory:      legacy_const_firmware.WeatherCategoryAdvisory,
		const_accuweather.WeatherCategoryWarning:       legacy_const_firmware.WeatherCategoryWarning,
		const_accuweather.WeatherCategoryEmergency:     legacy_const_firmware.WeatherCategoryEmergency,
		const_accuweather.WeatherCategoryBeneficial:    legacy_const_firmware.WeatherCategoryBeneficial,
		const_accuweather.WeatherCategoryNeutral:       legacy_const_firmware.WeatherCategoryNeutral,
		const_accuweather.WeatherCategoryAtRisk:        legacy_const_firmware.WeatherCategoryAtRisk,
		const_accuweather.WeatherCategoryAtHighRisk:    legacy_const_firmware.WeatherCategoryAtHighRisk,
		const_accuweather.WeatherCategoryAtExtremeRisk: legacy_const_firmware.WeatherCategoryAtExtremeRisk,
		const_accuweather.WeatherCategoryVeryLikely:    legacy_const_firmware.WeatherCategoryVeryLikely,
		const_accuweather.WeatherCategoryLikely:        legacy_const_firmware.WeatherCategoryLikely,
		const_accuweather.WeatherCategoryPossibly:      legacy_const_firmware.WeatherCategoryPossibly,
		const_accuweather.WeatherCategoryVeryUnlikely:  legacy_const_firmware.WeatherCategoryVeryUnlikely,
	}

	//----------------------------------------------
	// Globals - Tables - Moonphase Map
	//----------------------------------------------


	// Moon Phases and Weather Categories must match Enum Values from Firmware trie.h
	/**
	 * @brief
	 */
	MoonPhaseToFirmwareMap = map[string]int{
		const_accuweather.MoonPhaseNew:            const_firmware.MoonPhaseNew,
		const_accuweather.MoonPhaseWaxingCrescent: const_firmware.MoonPhaseWaxingCrescent,
		const_accuweather.MoonPhaseFirst:          const_firmware.MoonPhaseFirst,
		const_accuweather.MoonPhaseWaxingGibbous:  const_firmware.MoonPhaseWaxingGibbous,
		const_accuweather.MoonPhaseFull:           const_firmware.MoonPhaseFull,
		const_accuweather.MoonPhaseWaningGibbous:  const_firmware.MoonPhaseWaningGibbous,
		const_accuweather.MoonPhaseLast:           const_firmware.MoonPhaseLast,
		const_accuweather.MoonPhaseWaningCrescent: const_firmware.MoonPhaseWaningCrescent,
	}

	//----------------------------------------------
	// Globals - Tables - Weather Category to Enum
	//----------------------------------------------
	/**
	 * @brief
	 */
	WeatherCategoryToFirmwareMap = map[string]int{
		const_accuweather.WeatherCategoryLow:       const_firmware.WeatherCategoryLow,
		const_accuweather.WeatherCategoryHigh:      const_firmware.WeatherCategoryHigh,
		const_accuweather.WeatherCategoryGood:      const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryModerate:  const_firmware.WeatherCategoryModerate,
		const_accuweather.WeatherCategoryUnhealthy: const_firmware.WeatherCategoryUnhealthy,
		const_accuweather.WeatherCategoryHazardous: const_firmware.WeatherCategoryHazardous,

		// Backwards compatibility mapping
		const_accuweather.WeatherCategoryPoor:          const_firmware.WeatherCategoryModerate,
		const_accuweather.WeatherCategoryFair:          const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryVeryGood:      const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryExcellent:     const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryVeryHigh:      const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryExtreme:       const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryUnlikely:      const_firmware.WeatherCategoryLow,
		const_accuweather.WeatherCategoryWatch:         const_firmware.WeatherCategoryHigh,
		const_accuweather.WeatherCategoryAdvisory:      const_firmware.WeatherCategoryHigh,
		const_accuweather.WeatherCategoryWarning:       const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryEmergency:     const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryBeneficial:    const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryNeutral:       const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryAtRisk:        const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryAtHighRisk:    const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryAtExtremeRisk: const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryVeryLikely:    const_firmware.WeatherCategoryHigh,
		const_accuweather.WeatherCategoryLikely:        const_firmware.WeatherCategoryModerate,
		const_accuweather.WeatherCategoryPossibly:      const_firmware.WeatherCategoryLow,
		const_accuweather.WeatherCategoryVeryUnlikely:  const_firmware.WeatherCategoryLow,
	}

	//----------------------------------------------
	// Globals - Tables - Weather Category to Enum (extended set)
	//----------------------------------------------
	/**
	 * @brief
	 */
	WeatherCategoryToFirmwareMapExtended = map[string]int{
		const_accuweather.WeatherCategoryLow:           const_firmware.WeatherCategoryLow,
		const_accuweather.WeatherCategoryHigh:          const_firmware.WeatherCategoryHigh,
		const_accuweather.WeatherCategoryGood:          const_firmware.WeatherCategoryGood,
		const_accuweather.WeatherCategoryModerate:      const_firmware.WeatherCategoryModerate,
		const_accuweather.WeatherCategoryUnhealthy:     const_firmware.WeatherCategoryUnhealthy,
		const_accuweather.WeatherCategoryHazardous:     const_firmware.WeatherCategoryHazardous,
		const_accuweather.WeatherCategoryPoor:          const_firmware.WeatherCategoryPoor,
		const_accuweather.WeatherCategoryFair:          const_firmware.WeatherCategoryFair,
		const_accuweather.WeatherCategoryVeryGood:      const_firmware.WeatherCategoryVeryGood,
		const_accuweather.WeatherCategoryExcellent:     const_firmware.WeatherCategoryExcellent,
		const_accuweather.WeatherCategoryVeryHigh:      const_firmware.WeatherCategoryVeryHigh,
		const_accuweather.WeatherCategoryExtreme:       const_firmware.WeatherCategoryExtreme,
		const_accuweather.WeatherCategoryUnlikely:      const_firmware.WeatherCategoryUnlikely,
		const_accuweather.WeatherCategoryWatch:         const_firmware.WeatherCategoryWatch,
		const_accuweather.WeatherCategoryAdvisory:      const_firmware.WeatherCategoryAdvisory,
		const_accuweather.WeatherCategoryWarning:       const_firmware.WeatherCategoryWarning,
		const_accuweather.WeatherCategoryEmergency:     const_firmware.WeatherCategoryEmergency,
		const_accuweather.WeatherCategoryBeneficial:    const_firmware.WeatherCategoryBeneficial,
		const_accuweather.WeatherCategoryNeutral:       const_firmware.WeatherCategoryNeutral,
		const_accuweather.WeatherCategoryAtRisk:        const_firmware.WeatherCategoryAtRisk,
		const_accuweather.WeatherCategoryAtHighRisk:    const_firmware.WeatherCategoryAtHighRisk,
		const_accuweather.WeatherCategoryAtExtremeRisk: const_firmware.WeatherCategoryAtExtremeRisk,
		const_accuweather.WeatherCategoryVeryLikely:    const_firmware.WeatherCategoryVeryLikely,
		const_accuweather.WeatherCategoryLikely:        const_firmware.WeatherCategoryLikely,
		const_accuweather.WeatherCategoryPossibly:      const_firmware.WeatherCategoryPossibly,
		const_accuweather.WeatherCategoryVeryUnlikely:  const_firmware.WeatherCategoryVeryUnlikely,
	}

	//----------------------------------------------
	// Globals - Tables -
	//----------------------------------------------
	/**
	 * @brief
	 */
	LocationMap map[string]*time.Location

	//----------------------------------------------
	// Globals -
	//----------------------------------------------
	/**
	 * @brief
	 */
	AccuApiKey string

	//----------------------------------------------
	// Globals -
	//----------------------------------------------

	// Mutex to control access to the map
	/**
	 * @brief
	 */
	LocationMapMutex = sync.RWMutex{}
)


