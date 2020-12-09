package middleware

import (
	cm "simobiplus/common"

	ex "github.com/wolvex/go/error"
	be "github.com/wolvex/paymentaggregator"
) 

func Authorize(msg be.Message) *ex.AppError {
	//validating signature sent by backend
	if e := cm.BackendSigner.CheckRequest(msg); e != nil {
		return ex.Error(e, be.ERR_INVALID_SIGNATURE).Rem("Signature invalid")
	}
	return nil
}
