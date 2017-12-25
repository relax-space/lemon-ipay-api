package ipay

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/relax-space/go-kit/httpreq"

	"github.com/relax-space/go-kit/model"
)

type DefaultBmappingService struct {
}

func (a *DefaultBmappingService) Get(bmappingUrl, serviceType string, serviceContent interface{}) (eId int64, httpStatus int, err error) {

	baseServiceDto, ok := serviceContent.(BaseServiceDto)
	if !ok {
		httpStatus = http.StatusBadRequest
		err = errors.New("pay_type and e_id is required when servie_type is default")
		return
	}
	httpStatus = http.StatusOK
	eId = baseServiceDto.EId
	return
}

type BrandABmappingService struct {
}

func (a *BrandABmappingService) Get(bmappingUrl, serviceType string, serviceContent interface{}) (eId int64, httpStatus int, err error) {
	type BrandABmappingContent struct {
		GroupCode string `json:"group_code"`
		Code      string `json:"code"`
		CountryId int64  `json:"country_id"`
		EType     int64  `json:"e_type"`
	}

	dto, ok := serviceContent.(BrandABmappingContent)
	if !ok {
		httpStatus = http.StatusBadRequest
		err = errors.New("group_code、code、country_id、e_type is required when servie_type is branda")
		return
	}

	type EIdDto struct {
		EId int64 `json:"e_id"`
	}
	var result struct {
		Success bool        `json:"success"`
		EIdDto  EIdDto      `json:"result"`
		Error   model.Error `json:"error"`
	}
	url := fmt.Sprintf("%v/storeEPay/%v/%v/%v/%v", bmappingUrl, dto.GroupCode, dto.Code, dto.CountryId, dto.EType)

	_, err = httpreq.GET("", url, &result)
	if err != nil {
		httpStatus = http.StatusInternalServerError
		return
	}
	eId = result.EIdDto.EId
	return

}
