package middleware

import (
	cm "github.com/yswijaya531/simobiplus/common"

	ex "github.com/wolvex/go/error"
	be "gitlab.smartfren.com/paggr/libraries"
) 

func Authorize(msg be.Message) *ex.AppError {
	//validating signature sent by backend
	if e := cm.BackendSigner.CheckRequest(msg); e != nil {
		return ex.Error(e, be.ERR_INVALID_SIGNATURE).Rem("Signature invalid")
	}
	return nil
}

func BuildResponse(req be.Message, err int, remark string) be.Message {
	return be.Message{
		OriginHost: cm.Config.OriginHost,
		Version:    "1.0",
		MsgID:      req.MsgID,
		Response: &be.ResponseMessage{
			Result: &be.Result{
				Code:   err,
				Remark: remark,
			},
		},
	}
}