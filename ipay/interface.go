package ipay

type IPayService interface {
	Pay(reqDto *ReqDto) (result map[string]interface{}, httpStatus int, err error)
	Query(reqDto *ReqDto) (result map[string]interface{}, httpStatus int, err error)
	Refund(reqDto *ReqDto) (result map[string]interface{}, httpStatus int, err error)
	Reverse(reqDto *ReqDto) (result map[string]interface{}, httpStatus int, err error)
}

type IBmappingService interface {
	Get(bmappingUrl, serviceType string, serviceContent interface{}) (eId int64, httpStatus int, err error)
}
