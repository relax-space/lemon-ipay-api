package ipay

const (
	BMAPPING_TYPE_DEFAULT = "default"
	BMAPPING_TYPE_BRANDA  = "branda"
)

func CreateBmappingService(serviceType string) (service IBmappingService) {
	switch serviceType {
	case "default":
		service =new(DefaultBmappingService)
	case "branda":
		service =new(BrandABmappingService)
	}
	return
}
