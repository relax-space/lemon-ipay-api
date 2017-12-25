package ipay

type ReqDto struct {
	Sign           string      `json:"sign"`
	PayType        string      `json:"pay_type"`
	ServiceType    string      `json:"servcie_type"`
	ServiceContent interface{} `json:"service_content"`
	RawContent     interface{} `json:"raw_content"`
}

type PayReqDto struct {
	ReqDto
	AuthCode string `json:"auth_code"`
}
type NoPayReqDto struct {
	ReqDto
	OutTradeNo string `json:"out_trade_no"`
}

type RespDto struct {
	Sign          string      `json:"sign"`
	UnifiedResult interface{} `json:"unified_result"`
	RawResult     interface{} `json:"raw_result"`
}

type BaseServiceDto struct {
	PayType string `json:"pay_type"`
	EId     int64  `json:"e_id"`
}
