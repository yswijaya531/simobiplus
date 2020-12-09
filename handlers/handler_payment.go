package handlers

import (
	cm "github.com/yswijaya531/simobiplus/common"

	log "github.com/Sirupsen/logrus"
	ex "github.com/wolvex/go/error"
	be "github.com/wolvex/paymentaggregator"
)


func PaymentHandler(req be.Message) (res be.Message, errs error) {
	var err *ex.AppError

	defer panicRecovery()
	
	defer func() {
		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught")
		} else {
			err = ex.Errorc(be.ERR_IN_PROGRESS).Rem("Payment in progress")
		}
		res.Response.Result = buildResponse(err)
	}()
	
	if req.Request == nil || req.Request.Order == nil || req.Request.Payment == nil ||
		req.Request.Payment.Partner == nil || req.Request.Payment.Resource == nil {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing mandatory parameter")
		return
	}

	var reqs be.Message

	//initialize response, echoed from request
	res = initResponseFromRequest(reqs)
	res.Response.Order = req.Request.Order
	res.Response.Payment = req.Request.Payment
	log.WithField("info", req.Request.Order.InvoiceNo).Info("Request.Order.InvoiceNo")

	invoiceNo := req.Request.Order.InvoiceNo
	if len(req.Request.Order.InvoiceNo) > 12 {
		invoiceNo = invoiceNo[0:12]
	}

	var snap *HttpClient

	if snap, err = NewSnapClient(); err != nil {
		return
	}

	var response *cm.SimobiResponse
	if response, err = snap.PushInvoice(invoiceNo, NormalizeMDN(req.Request.Payment.Account.ID),
		req.Request.Order.Title, 1, int64(req.Request.Payment.Resource[0].Value), req.Request.Order.Goods[0].Code, req.Request.Order.Goods[0].Category); err != nil {
		return
	}

	res.Response.Payment.ApprovalCode = response.DataStatus.TransRefNum

	return res, errs
}
