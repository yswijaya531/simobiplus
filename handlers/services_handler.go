package handlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	cm "github.com/yswijaya531/simobiplus/common"

	ex "github.com/wolvex/go/error"
	be "gitlab.smartfren.com/paggr/libraries"
)

func panicRecovery() {
	if r := recover(); r != nil {
		fmt.Printf("Recovering from panic: %v \n", r)
	}
}

func buildResponse(err *ex.AppError) *be.Result {
	if err != nil {
		return &be.Result{
			Code:   err.ErrCode,
			Remark: err.Remark,
		}
	}

	return &be.Result{
		Code:   be.ERR_SUCCESS,
		Remark: "Success",
	}
}

func initResponseFromRequest(req be.Message) be.Message {
	return be.Message{
		OriginHost: cm.Config.OriginHost,
		Version:    "1.0",
		MsgID:      req.MsgID,
		Response:   &be.ResponseMessage{},
	}
}

func NormalizeMDN(mdn string) string {
	if _, err := strconv.ParseFloat(mdn, 64); err != nil {
		return ""
	} else if len(mdn) > 20 {
		return ""
	}

	if strings.HasPrefix(mdn, "62") {
		return mdn
	} else if strings.HasPrefix(mdn, "+62") {
		return strings.Replace(mdn, "+62", "62", 1)
	} else if strings.HasPrefix(mdn, "0") {
		return strings.Replace(mdn, "0", "62", 1)
	}

	return mdn
}

func NewSnapClient() (client *HttpClient, err *ex.AppError) {
	transport := &http.Transport{}
	if strings.HasPrefix(cm.Config.SnapURL, "https") {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if cm.Config.ProxyURL != "" {
		if pUrl, e := url.Parse(cm.Config.ProxyURL); err != nil {
			err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to parse proxy url")
			return
		} else {
			transport.Proxy = http.ProxyURL(pUrl)
		}
	}

	client = &HttpClient{
		Session: &http.Client{
			Timeout:   time.Duration(cm.Config.Timeout) * time.Millisecond,
			Transport: transport,
		},
	}

	return
}

func checkErrorCodes(responseCode int, responseMessage string) *ex.AppError {

	strMsg := fmt.Sprintf("%d - %s", responseCode, responseMessage)

	switch int(responseCode) {
	case 0:
		return nil
	case 2:
		return ex.Errorc(be.ERR_TRX_UNAUTHORIZED).Rem(strMsg)

	case 7, 9:
		return ex.Errorc(be.ERR_INVALID_FORMAT).Rem(strMsg)

	case 6, 14:
		return ex.Errorc(be.ERR_ACCOUNT_NOT_FOUND).Rem(strMsg)

	case 8, 51, 61:
		return ex.Errorc(be.ERR_TRX_INVALID).Rem(strMsg)

	default:
		return ex.Errorc(be.ERR_OTHERS).Rem(strMsg)

	}
}

func InspectResponseCode(msg *cm.SimobiCallBack) *ex.AppError {
	
	switch  {
		case msg.DataStatus.Status == "submitted":			
			return ex.Errorc(be.ERR_IN_PROGRESS).Rem("in progress")
		case msg.DataStatus.Status == "paid":			
			return ex.Errorc(be.ERR_SUCCESS).Rem("Success")
		case msg.ResponseCode == "01" || msg.ResponseCode == "97" || msg.ResponseCode == "99":	
			return ex.Errorc(be.ERR_TRX_INVALID).Rem("Failed - Catch Error")
		case msg.ResponseCode == "02" || msg.ResponseCode == "98":	
			return ex.Errorc(be.ERR_PAYMENT_IN_PROGRESS).Rem("pending")
		case msg.ResponseCode == "06":	
			return ex.Errorc(be.ERR_ACCOUNT_NOT_FOUND).Rem("User not Found")
		case msg.ResponseCode == "08":	
			return ex.Errorc(be.ERR_TRX_DUPLICATE).Rem("txId is already Exist")
		case msg.ResponseCode == "51":	
			return ex.Errorc(be.ERR_PAYMENT_DECLINED).Rem("Insufficient balance")
		case msg.ResponseCode == "61":	
			return ex.Errorc(be.ERR_TRX_UNAUTHORIZED).Rem("Amount limit exceeded")								
		default:			
			return ex.Errorc(be.ERR_OTHERS).Rem(msg.ResponseCode, msg.ResponseMessage)
	}

}