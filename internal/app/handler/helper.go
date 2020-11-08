package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vincentwijaya/go-pkg/log"
	"github.com/vincentwijaya/go-template/constant/errs"
	"github.com/vincentwijaya/go-template/internal/entity"
)

func writeResponse(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := &entity.Response{
		Success: false,
	}

	if err != nil {
		response.Message, response.Code = classifyError(err)
	} else {
		response.Success = true
		response.Code = errs.Success
		response.Message = "OK"
	}

	if data != nil && response.Success {
		response.Data = data
	}

	log.Infof("Response: %+v", response)
	body, _ := json.Marshal(response)
	_, _ = w.Write(body)
}

//maps
type errResponse struct {
	message string
	code    string
	apiFail bool
}

var (
	errToResponse = map[error]errResponse{
		errs.BadRequest: {
			message: errs.GeneralErrorMessage,
			code:    errs.BadRequestCode,
			apiFail: true,
		},
		errs.NoData: {
			message: errs.NoDataMessage,
			code:    errs.NoDataCode,
			apiFail: false,
		},
		errs.Unauthorized: {
			message: errs.UnauthorizedMessage,
			code:    errs.UnauthorizedErrorCode,
			apiFail: false,
		},
	}
)

func classifyError(err error) (string, string) {
	val, exist := errToResponse[err]
	if !exist {
		log.Infof("Unmapped error:%v", err.Error())
		return errs.GeneralErrorMessage, errs.UndefinedErrorCode
	}
	if val.apiFail {
		// on api fail, return general error message
		// log the error on log
		log.Errorf("Error on API code [%s]:%s", val.code, err.Error())
	}
	return val.message, val.code
}
