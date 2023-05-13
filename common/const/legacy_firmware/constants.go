package legacy_firmware

//================================================================
//================================================================
// Globals
//================================================================
//================================================================

//==============================================
// Globals - Constants
//==============================================

//----------------------------------------------
// Globals - Constants - Moonphase
//----------------------------------------------

//----------------------------------------------
// Globals - Constants - Moonphase
//----------------------------------------------
const (
	MoonPhaseNew            = 1
	MoonPhaseWaxingCrescent = 2
	MoonPhaseFirst          = 3
	MoonPhaseWaxingGibbous  = 4
	MoonPhaseFull           = 5
	MoonPhaseWaningGibbous  = 6
	MoonPhaseLast           = 7
	MoonPhaseWaningCrescent = 8
)

//----------------------------------------------
// Globals - Constants - Weather Category
//----------------------------------------------
const (
	WeatherCategoryLow           = 1
	WeatherCategoryHigh          = 2
	WeatherCategoryGood          = 3
	WeatherCategoryModerate      = 4
	WeatherCategoryUnhealthy     = 5
	WeatherCategoryHazardous     = 6
	WeatherCategoryPoor          = WeatherCategoryModerate
	WeatherCategoryFair          = WeatherCategoryGood
	WeatherCategoryVeryGood      = WeatherCategoryGood
	WeatherCategoryExcellent     = WeatherCategoryGood
	WeatherCategoryVeryHigh      = WeatherCategoryGood
	WeatherCategoryExtreme       = WeatherCategoryHazardous
	WeatherCategoryUnlikely      = WeatherCategoryLow
	WeatherCategoryWatch         = WeatherCategoryHigh
	WeatherCategoryAdvisory      = WeatherCategoryHigh
	WeatherCategoryWarning       = WeatherCategoryHazardous
	WeatherCategoryEmergency     = WeatherCategoryHazardous
	WeatherCategoryBeneficial    = WeatherCategoryGood
	WeatherCategoryNeutral       = WeatherCategoryGood
	WeatherCategoryAtRisk        = WeatherCategoryHazardous
	WeatherCategoryAtHighRisk    = WeatherCategoryHazardous
	WeatherCategoryAtExtremeRisk = WeatherCategoryHazardous
	WeatherCategoryVeryLikely    = WeatherCategoryHigh
	WeatherCategoryLikely        = WeatherCategoryModerate
	WeatherCategoryPossibly      = WeatherCategoryLow
	WeatherCategoryVeryUnlikely  = WeatherCategoryLow
)
