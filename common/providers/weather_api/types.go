package weather_api

//==============================================
// CopyRight 2020 La Crosse Technology, LTD.
//==============================================

//==============================================
// Imports
//==============================================

type (

	/**
	 * @brief Api Response Interface
	 */
	ApiResponseInterface interface {
		JsonResponse(string) (string, error)
		ResponseFormat(string) (ApiResponseInterface, error)
	}

	/**
	 * @brief Weather Icon Record
	 */
	AccuIcon struct {
		IconNumber  int
		Icon        string
		Day         bool
		Night       bool
		Text        string
		DisplayIcon int
	}


)