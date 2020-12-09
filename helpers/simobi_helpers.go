package helpers

import (
	cm "simobiplus/common"

	ex "github.com/wolvex/go/error"
	be "github.com/wolvex/paymentaggregator"
)

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