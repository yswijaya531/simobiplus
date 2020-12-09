package handlers

import (
	"time"

	log "github.com/sirupsen/logrus"
	cm "github.com/yswijaya531/simobiplus/common"

	ex "github.com/wolvex/go/error"
	be "gitlab.smartfren.com/paggr/libraries"
)

func AdviseHandler(req be.Message) (res be.Message, errs error) {

	defer panicRecovery()
	
	var err *ex.AppError

	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught")
			res.Response.Result = buildResponse(err)
		} else {
			err = ex.Errorc(be.ERR_SUCCESS).Rem("Success")
		}
	}()

	//init Response from Request
	res = initResponseFromRequest(req)
	res.Response.Order = req.Request.Order
	res.Response.Order = req.Request.Order

	invoiceNo := req.Request.Order.InvoiceNo
	// if len(invoiceNo) > 12 {
	// 	invoiceNo = invoiceNo[0:12]
	// }

	//Calling Simobi
	var snap *HttpClient
	if snap, err = NewSnapClient(); err != nil {
		return
	}

	curTime := time.Now()
	ReqDate := curTime.Format("02-01-2006 15:04:05")

	var response *cm.SimobiCallBack
	if response, err = snap.PullInvoice(invoiceNo, ReqDate); err != nil {
		return
	}

	result := &be.Result{}
	
	if err = InspectResponseCode(response); err == nil {
		result.Code = be.ERR_SUCCESS
		result.Remark = "Success"
	} else {
		if err.ErrCode == be.ERR_IN_PROGRESS {
			return
		}
		result.Code = err.ErrCode
		result.Remark = err.Remark
	}
	
	return
}