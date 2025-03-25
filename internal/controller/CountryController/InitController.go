package countrycontroller

import "ten_module/internal/service/countryservice"

func InitCountryControll() {
	countryservice.InitCountryService()
	InitCountryController()
}
