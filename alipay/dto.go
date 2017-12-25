package alipay

import alpay "github.com/relax-space/lemon-alipay-sdk"

type ReqPayDto struct {
	*alpay.ReqPayDto
	EId int64 `json:"e_id"`
}
type ReqQueryDto struct {
	*alpay.ReqQueryDto
	EId int64 `json:"e_id"`
}
type ReqRefundDto struct {
	*alpay.ReqRefundDto
	EId int64 `json:"e_id"`
}
type ReqReverseDto struct {
	*alpay.ReqReverseDto
	EId int64 `json:"e_id"`
}

// type ReqRefundQueryDto struct {
// 	*alpay.ReqRefundQueryDto
// 	EId int64 `json:"e_id"`
// }
type ReqPrepayDto struct {
	*alpay.ReqPrepayDto
	EId int64 `json:"e_id"`
}
