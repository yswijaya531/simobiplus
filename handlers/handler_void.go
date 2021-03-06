package handlers

import (
	cm "github.com/yswijaya531/simobiplus/common"

	log "github.com/sirupsen/logrus"
	ex "github.com/wolvex/go/error"
	be "gitlab.smartfren.com/paggr/libraries"
)


func VoidHandler(req be.Message) (res be.Message, errs error) {
	var err *ex.AppError

	defer panicRecovery()
	
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught")
		} else {
			err = ex.Errorc(be.ERR_SUCCESS).Rem("success")
		}
		res.Response.Result = buildResponse(err)
	}()
	
	if req.Request == nil || req.Request.Order == nil || req.Request.Void == nil ||
		req.Request.Void.Partner == nil || req.Request.Void.Resource == nil || req.Request.Order.Goods == nil {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing mandatory parameter")
		return
	}

	var reqs be.Message
	//initialize response, echoed from request
	res = initResponseFromRequest(reqs)
	res.Response.Order = req.Request.Order
	res.Response.Void = req.Request.Void
	log.WithField("info", req.Request.Order.InvoiceNo).Info("Request.Order.InvoiceNo")

	invoiceNo := req.Request.Order.InvoiceNo
	
	var snap *HttpClient

	if snap, err = NewSnapClient(); err != nil {
		return
	}

	var response *cm.SimobiCallBack
	if response, err = snap.PushRefund(invoiceNo, req.Request.Void.ApprovalCode, req.Request.Order.Goods[0].Category, int64(req.Request.Void.Resource[0].Value)); err != nil {
		return
	}

	log.WithField("info", response).Info("response")

	return res, errs
}