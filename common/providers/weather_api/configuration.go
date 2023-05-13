package weather_api

//==============================================
// CopyRight 2020 La Crosse Technology, LTD.
//==============================================

//==============================================
// Imports
//==============================================
import (

)

//================================================================
//================================================================
// Globals
//================================================================
//================================================================

//==============================================
// Globals - Constants
//==============================================

/**
 * @brief
 */
const (
	ForecastTypeStreams       = "DATASTREAMS"
	DayInfoNight              = "Night"
	DayInfoDay                = "Day"
	ForecastExpireHours       = 12
	AccuBaseUrl               = "http://dataservice.accuweather.com"
	DefaultCountryCode        = "US"
	MaxRetries                = 5
	LocationTypeCity          = "City"
	ExceptionModeFlowCommand  = 0
	DefaultModeFlowCommand    = 2
	GEO_REFRESH_API_HIT_AMOUNT = 3 // Amount of time Keith's refresh api is called
)
