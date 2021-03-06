package handlers

import (
	"fmt"
	"strconv"
	"time"

	cm "github.com/yswijaya531/simobiplus/common"

	log "github.com/sirupsen/logrus"
	ex "github.com/wolvex/go/error"
	be "gitlab.smartfren.com/paggr/libraries"
)

var err *ex.AppError

func NotifyHandler(req be.Message) (res be.Message, errs error) {
	
	defer panicRecovery()

	var err *ex.AppError
	
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught")
		}
	}()
		
	
	res = initResponseFromRequest(req)
	res.Response.Order = req.Request.Order	
	res.Response.Result = &be.Result{
		Code:   be.ERR_SUCCESS,
		Remark: "Success",
	}
	
	go processNotify(req)

	return res, errs
	
}

func processNotify(req be.Message) {
	
	var e error	

	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught")
		}
	}()

	curTime := time.Now()
	reqDate := curTime.Format("02-01-2006 15:04:05")
	invoiceNo := req.Request.Order.InvoiceNo
	if len(req.Request.Order.InvoiceNo) > 12 {
		invoiceNo = invoiceNo[0:12]
	}

	//Calling Simobi
	var snap *HttpClient
	if snap, err = NewSnapClient(); err != nil {
		return
	}

	var response *cm.SimobiCallBack
	if response, err = snap.PullInvoice(invoiceNo, reqDate); err != nil {
		return
	}

	if response.DataStatus.Amount == "" {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing parameter Amount")
		return
	}

	if response.DataStatus.TxID == "" {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing paramter TxID")
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

	var amount float64
	if amount, e = strconv.ParseFloat(response.DataStatus.Amount, 64); e != nil {
		err = ex.Errorc(be.ERR_INVALID_FORMAT).Rem("Amount parameter has invalid format")
		return
	}

	goods := []*be.GoodsItem{
		&be.GoodsItem{
			Price: &be.Amount{
				Currency: "IDR",
				Value:    amount,
			},
		},
	}

	resources := []*be.Amount{
		&be.Amount{
			Currency: "IDR",
			Value:    amount,
		},
	}

	callback := &be.Message{
		Request: &be.RequestMessage{
			Order: &be.Order{
				InvoiceNo: response.DataStatus.TxID,
				TotalPrice: &be.Amount{
					Currency: "IDR",
					Value:    amount,
				},
				Goods: goods,
			},
			Payment: &be.Payment{
				Method:       "SIMOBI",
				Resource:     resources,
				ApprovalCode: response.DataStatus.AuthCode,
			},
			Result: result,
		},
		MsgID: req.MsgID,
	}

	backend := be.NewClient(fmt.Sprintf("%s/callback", cm.Config.BackendURL), cm.Config.OriginHost, cm.Signer, nil, cm.Config.Timeout)
	_, err = backend.Post(callback)
	
	return 
}	

