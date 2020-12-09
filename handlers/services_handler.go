package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	ex "github.com/wolvex/go/error"
	be "github.com/wolvex/paymentaggregator"
	cm "github.com/yswijaya531/simobiplus/common"
)

type HttpClient struct {
	ServerKey  string
	Session    *http.Client
	ExpiryTime int64
}

func isTimeout(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	} else {
		return false
	}
}

func getSimobiToken(c *HttpClient) (err *ex.AppError) {

	form := url.Values{}
	form.Add("grant_type", cm.Config.GrantType)
	form.Add("client_id", cm.Config.ClientId)
	form.Add("client_secret", cm.Config.ClientSecret)
	form.Add("scope", cm.Config.Scope)
	form.Add("x-ibm-client-id", cm.Config.IbmClientId)

	var (
		e         error
		stringURL = cm.Config.SnapURL + cm.Config.TokenURL
	)
	resp, e := c.Session.PostForm(stringURL, form)

	if e != nil {
		err = ex.Error(e, be.ERR_SYSTEM_ERROR).Rem("Error, cant reach url destination %s", stringURL)
		return
	}

	//dump response for logging
	if dump, e := httputil.DumpResponse(resp, true); e != nil {
		log.Info(string(dump))
	} else {
		statusCode := resp.StatusCode
		log.WithFields(log.Fields{
			"request" : string(dump),
			"statusCode" :     statusCode,
		  }).Info("Sending HTTP request")
	}

	if resp.StatusCode != 200 {
		err = ex.Error(e, be.ERR_UNAUTHORIZED).Rem("Unable to get token")
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	mdl := cm.TokenResponse{}
	xerr := json.Unmarshal(body, &mdl)
	if xerr != nil {
		err = ex.Error(xerr, be.ERR_OTHERS).Rem("Error UnMarshall response from Auth TokenSimobi")
		return
	}
	c.ServerKey = mdl.AccessToken

	return
}

//PullInvoice is ...
func (c *HttpClient) PullInvoice(invoiceNo string, reqDate string) (response *cm.SimobiCallBack, err *ex.AppError) {

	msg := &cm.SimobiPull{
		BillerCode:  cm.Config.BillerCode,
		TxID:        invoiceNo,
		RequestDate: reqDate,
	}

	if response, err = pullSimobiAPI(c, cm.Config.PushStatusURL, msg); err != nil {
		return
	}

	remark := response.ResponseCode + ":" + response.ResponseMessage

	switch response.ResponseCode {
	case "00":
		return
	case "09":
		err = ex.Errorc(be.ERR_TRX_UNAUTHORIZED).Rem(remark)
	default:
		err = ex.Errorc(be.ERR_OTHERS).Rem(remark)
	}
	return
}


func pullSimobiAPI(c *HttpClient, url string, msg *cm.SimobiPull) (response *cm.SimobiCallBack, err *ex.AppError) {
	var (
		body       []byte
		e          error
		statusCode int		
	)

	if strings.HasPrefix(url, "/") {
		url = fmt.Sprintf("%s%s", cm.Config.SnapURL, url)
	} else {
		url = fmt.Sprintf("%s/%s", cm.Config.SnapURL, url)
	}

	if body, e = json.Marshal(msg); e != nil {
		err = ex.Error(e, be.ERR_INVALID_FORMAT).Rem("Unable to marshal request to json format")
		return
	}

	//initiliaze request
	if post, e := http.NewRequest("POST", strings.TrimSpace(url), bytes.NewBuffer(body)); e != nil {
		err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to create new http request")
		return

	} else {
		//assign headers
		if c.ServerKey == "" {
			if err = getSimobiToken(c); err != nil {
				return
			}
		}

		post.Header.Add("Accept", "application/json; charset=utf-8")
		post.Header.Add("Content-Type", "application/json")
		post.Header.Add("x-ibm-client-id", cm.Config.IbmClientId)
		post.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ServerKey))

		//dump message for logging
		if dump, e := httputil.DumpRequestOut(post, true); e != nil {
			err = ex.Error(e, be.ERR_OTHERS).Rem("Error in dump request")
			return
		} else {			
			log.WithFields(log.Fields{
				"request": string(dump),
				"url" :     url,
			  }).Info("Sending HTTP request")
		}

		if res, e := c.Session.Do(post); e != nil {
			if isTimeout(e) {
				err = ex.Errorc(be.ERR_TIMEOUT).Rem("Timeout detected")
			} else {
				err = ex.Errorc(be.ERR_OTHERS).Rem("Unable to send POST to %s", url)
			}
			return

		} else {
			defer res.Body.Close()

			//dump response for logging
			if dump, e := httputil.DumpResponse(res, true); e != nil {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Error in dump response")
				return
			} else {
				statusCode = res.StatusCode				
				log.WithFields(log.Fields{
					"request": string(dump),
					"url" :     url,
				  }).Info("Sending HTTP request")
			}

			if res.StatusCode == http.StatusNotFound {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Service not available")
				return
			}

			if body, e = ioutil.ReadAll(res.Body); err != nil {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to fetch response body")
				return
			}

			if e = json.Unmarshal(body, &response); e != nil {
				err = ex.Error(e, be.ERR_INVALID_FORMAT).Rem("Unable to decode response")
				return
			}

			response.HTTPStatus = statusCode
		}
	}

	return
}

func (c *HttpClient) PushInvoice(invoiceNo string, mobileNo string, itemName string, qty, amount int64, txtype string, category string) (response *cm.SimobiResponse, err *ex.AppError) {
	curTime := time.Now()
	ReqDate := curTime.Format("02-01-2006 15:04:05")

	msg := &cm.SimobiRequest{
		BillerCode:   cm.Config.MerchantID[category],
		MobileNumber: mobileNo,
		TxID:         invoiceNo,
		ItemName:     fmt.Sprintf("%s", itemName),
		TxType:       txtype,
		Qty:          fmt.Sprintf("%d", qty),
		Amount:       fmt.Sprintf("%d", amount),
		TxDate:       ReqDate,
	}

	if response, err = submitSimobiAPI(c, cm.Config.PushInvoiceURL, msg); err != nil {
		return
	}

	remark := response.ResponseCode + ":" + response.ResponseMessage

	switch response.ResponseCode {
	case "00":
		return
	case "09":
		err = ex.Errorc(be.ERR_TRX_UNAUTHORIZED).Rem(remark)
	case "06":
		err = ex.Errorc(be.ERR_ACCOUNT_NOT_FOUND).Rem(remark)
	default:
		err = ex.Errorc(be.ERR_OTHERS).Rem(remark)
	}
	return
}

func submitSimobiAPI(c *HttpClient, url string, msg *cm.SimobiRequest) (response *cm.SimobiResponse, err *ex.AppError) {
	var (
		body       []byte
		e          error
		statusCode int
	)

	if strings.HasPrefix(url, "/") {
		url = fmt.Sprintf("%s%s", cm.Config.SnapURL, url)
	} else {
		url = fmt.Sprintf("%s/%s", cm.Config.SnapURL, url)
	}

	if body, e = json.Marshal(msg); e != nil {
		err = ex.Error(e, be.ERR_INVALID_FORMAT).Rem("Unable to marshal request to json format")
		return
	}

	//initiliaze request
	if post, e := http.NewRequest("POST", strings.TrimSpace(url), bytes.NewBuffer(body)); e != nil {
		err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to create new http request")
		return

	} else {
		//assign headers
		if c.ServerKey == "" {
			if err = getSimobiToken(c); err != nil {
				return
			}
		}

		post.Header.Add("Accept", "application/json; charset=utf-8")
		post.Header.Add("Content-Type", "application/json")
		post.Header.Add("x-ibm-client-id", cm.Config.IbmClientId)
		post.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ServerKey))

		//dump message for logging
		if dump, e := httputil.DumpRequestOut(post, true); e != nil {
			err = ex.Error(e, be.ERR_OTHERS).Rem("Error in dump request")
			return
		} else {
			log.WithFields(log.Fields{
				"request": string(dump),
				"url":     url,
			}).Info("Sending HTTP request")
		}

		if res, e := c.Session.Do(post); e != nil {
			if isTimeout(e) {
				err = ex.Errorc(be.ERR_TIMEOUT).Rem("Timeout detected")
			} else {
				err = ex.Errorc(be.ERR_OTHERS).Rem("Unable to send POST to %s", url)
			}
			return

		} else {
			defer res.Body.Close()

			//dump response for logging
			if dump, e := httputil.DumpResponse(res, true); e != nil {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Error in dump response")
				return
			} else {
				statusCode = res.StatusCode
				log.WithFields(log.Fields{
					"status":   statusCode,
					"response": string(dump),
				}).Info("Receiving HTTP response")
			}

			if res.StatusCode == http.StatusNotFound {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Service not available")
				return
			}

			if body, e = ioutil.ReadAll(res.Body); err != nil {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to fetch response body")
				return
			}

			if e = json.Unmarshal(body, &response); e != nil {
				err = ex.Error(e, be.ERR_INVALID_FORMAT).Rem("Unable to decode response")
				return
			}

			response.HTTPStatus = statusCode
		}
	}

	return
}

func (c *HttpClient) PushRefund( invoiceNo string, authCode string, category string, amount int64) (response *cm.SimobiCallBack, err *ex.AppError) {
	curTime := time.Now()
	ReqDate := curTime.Format("02-01-2006 15:04:05")

	msg := &cm.SimobiRequest{
		TxID:         invoiceNo,
		TxDate:       ReqDate,
		BillerCode:   cm.Config.MerchantID[category],
		AuthCode:     authCode,
		Amount: 	  fmt.Sprintf("%d", amount),
	}

	if response, err = refundSimobiAPI(c, cm.Config.RefundURL, msg); err != nil {
		return
	}

	remark := response.ResponseCode + ":" + response.ResponseMessage

	switch response.ResponseCode {
	case "00":
		return
	case "09":
		err = ex.Errorc(be.ERR_TRX_UNAUTHORIZED).Rem(remark)
	case "06":
		err = ex.Errorc(be.ERR_ACCOUNT_NOT_FOUND).Rem(remark)
	default:
		err = ex.Errorc(be.ERR_OTHERS).Rem(remark)
	}
	return
}

func refundSimobiAPI(c *HttpClient, url string, msg *cm.SimobiRequest) (response *cm.SimobiCallBack, err *ex.AppError) {
	var (
		body       []byte
		e          error
		statusCode int
	)

	if strings.HasPrefix(url, "/") {
		url = fmt.Sprintf("%s%s", cm.Config.SnapURL, url)
	} else {
		url = fmt.Sprintf("%s/%s", cm.Config.SnapURL, url)
	}

	if body, e = json.Marshal(msg); e != nil {
		err = ex.Error(e, be.ERR_INVALID_FORMAT).Rem("Unable to marshal request to json format")
		return
	}

	//initiliaze request
	
	if post, e := http.NewRequest("POST", strings.TrimSpace(url), bytes.NewBuffer(body)); e != nil {
		err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to create new http request")
		return

	} else {
		//assign headers
		if c.ServerKey == "" {
			if err = getSimobiToken(c); err != nil {
				return
			}
		}

		post.Header.Add("Accept", "application/json; charset=utf-8")
		post.Header.Add("Content-Type", "application/json")
		post.Header.Add("x-ibm-client-id", cm.Config.IbmClientId)
		post.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ServerKey))

		//dump message for logging
		if dump, e := httputil.DumpRequestOut(post, true); e != nil {
			err = ex.Error(e, be.ERR_OTHERS).Rem("Error in dump request")
			return
		} else {
			log.WithFields(log.Fields{
				"request": string(dump),
				"url":     url,
			}).Info("Sending HTTP request")
		}

		if res, e := c.Session.Do(post); e != nil {
			if isTimeout(e) {
				err = ex.Errorc(be.ERR_TIMEOUT).Rem("Timeout detected")
			} else {
				err = ex.Errorc(be.ERR_OTHERS).Rem("Unable to send POST to %s", url)
			}
			return

		} else {
			defer res.Body.Close()

			//dump response for logging
			if dump, e := httputil.DumpResponse(res, true); e != nil {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Error in dump response")
				return
			} else {
				statusCode = res.StatusCode
				log.WithFields(log.Fields{
					"status":   statusCode,
					"response": string(dump),
				}).Info("Receiving HTTP response")
			}

			if res.StatusCode == http.StatusNotFound {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Service not available")
				return
			}

			if body, e = ioutil.ReadAll(res.Body); err != nil {
				err = ex.Error(e, be.ERR_OTHERS).Rem("Unable to fetch response body")
				return
			}

			if e = json.Unmarshal(body, &response); e != nil {
				err = ex.Error(e, be.ERR_INVALID_FORMAT).Rem("Unable to decode response ssss")
				return
			}

			
		}
	}

	return
}