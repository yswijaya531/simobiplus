package handlers

import (
	"fmt"
	"strconv"

	cm "github.com/yswijaya531/simobiplus/common"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	ex "github.com/wolvex/go/error"
	be "github.com/wolvex/paymentaggregator"
	"github.com/xid"
	"gopkg.in/go-playground/validator.v9"
)


func CallBackHandler(e echo.Context) (res cm.SimobiCallBack, errs error) {
	var err *ex.AppError

	defer panicRecovery()
	
	defer func() {

		if err != nil {
			log.WithField("error", err.Dump()).Error("Exception caught")
			res.ResponseCode = strconv.Itoa(err.ErrCode)
			res.ResponseMessage = err.Remark

		}
	}()

	//initialize callback struct
	req := new(cm.SimobiCallBack)
		
	if errs = e.Bind(req); err != nil {
		return res, errs
	}

	validate := validator.New()
	if e := validate.Struct(req); e != nil {
		err = ex.Errorc(be.ERR_INVALID_FORMAT).Rem("One or more parameter has invalid format")
		return
	}

	result := &be.Result{}
	
	var request *cm.SimobiCallBack
	
	if err = InspectResponseCode(request); err == nil {
		result.Code = be.ERR_SUCCESS
		result.Remark = "Success"
	} else {
		if err.ErrCode == be.ERR_IN_PROGRESS {
			return
		}
		result.Code = err.ErrCode
		result.Remark = err.Remark
	}
	
	if req.DataStatus.Amount == "" {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing paramter Amount")
		return
	}

	if req.DataStatus.TxID == "" {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing paramter TxID")
		return
	}

	if req.DataStatus.TransRefNum == "" {
		err = ex.Errorc(be.ERR_PARAM_MISSING).Rem("Missing paramter TransRefNum")
		return
	}

	var amount float64
	if amount, errs = strconv.ParseFloat(req.DataStatus.Amount, 64); errs != nil {
		err = ex.Errorc(be.ERR_INVALID_FORMAT).Rem("Amount parameter has invalid format")
		return
	}

	errx, _ := strconv.Atoi(req.ResponseCode)
	err = checkErrorCodes(errx, req.ResponseMessage)

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

	callbackMsg := &be.Message{
		Request: &be.RequestMessage{
			Order: &be.Order{
				InvoiceNo: req.DataStatus.TxID,
				TotalPrice: &be.Amount{
					Currency: "IDR",
					Value:    amount,
				},
				Goods: goods,
			},
			Payment: &be.Payment{
				ApprovalCode: req.DataStatus.AuthCode,
				Method:   "SIMOBI",
				Resource: resources,
			},
			Result: result,
		},
		MsgID: xid.New().String(),
	}

	if req.DataStatus.TransRefNum != "" {
		callbackMsg.Request.Payment.ApprovalCode = req.DataStatus.AuthCode
	}

	backend := be.NewClient(fmt.Sprintf("%s/callback", cm.Config.BackendURL), cm.Config.OriginHost, cm.Signer, nil, cm.Config.Timeout)

	var msg *be.Message
	msg, err = backend.Post(callbackMsg)
	res.ResponseCode = strconv.Itoa(msg.Response.Result.Code)
	res.ResponseMessage = "Success"

	if err == nil {
		if msg.Response.Result.Code != be.ERR_SUCCESS {
			err = ex.Errorc(msg.Response.Result.Code).Rem(msg.Response.Result.Remark)
		}
	}
	
	return res, errs
}