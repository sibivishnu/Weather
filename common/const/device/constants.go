package device
//----------------------------------------------
// CopyRight 2019 La Crosse Technology, LTD.
//----------------------------------------------

//----------------------------------------------
// Imports
//----------------------------------------------
import (
)

//----------------------------------------------
// Constants
//----------------------------------------------
const (
	CAT1 = "1"
	CAT2 = "2"
	CAT3 = "3"

	// Repeating Time Period
	// TimeLoopDisabled         = 0
	TimeLoopSeptemberOctober = 1
	TimeLoopDecemberJanuary  = 2
	TimeLoopDstIn            = 3
	TimeLoopDstOut           = 4
	TimeLoopCustom           = 5
	TimeLoopBackendDriver    = 6

	// Time Acceleration
	// TimeCompressionDisabled = 0

	// Control data sent to display
	ForecastScriptDisabled       = 0
	ForecastScriptStaticA        = 1
	// ForecastScriptStaticB        = 2
	// ForecastScriptStaticC        = 3
	// ForecastScriptExtremeWeather = 4
	ForecastScriptBackendDriver  = 5
)
