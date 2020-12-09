package common

//Struct API
type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ConsentedOn int    `json:"consented_on"`
	Scope       string `json:"scope"`
}

type SimobiRequest struct {
	BillerCode   string `json:"billerCode,omitempty"`
	MobileNumber string `json:"mobileNumber,omitempty"`
	TxID         string `json:"txId,omitempty"`
	TxDate       string `json:"requestDate,omitempty"`
	TxType       string `json:"txTpe,omitempty"`
	ItemName     string `json:"itemName,omitempty"`
	Qty          string `json:"qty,omitempty"`
	Amount       string `json:"txAmt,omitempty"`
	AuthCode 	 string `json:"authCode,omitempty"`
}

type SimobiResponse struct {
	DataStatus struct {
		Amount       float64 `json:"amount"`
		MerchantName string  `json:"merchantName"`
		MobileNumber string  `json:"mobileNumber"`
		RequestDate  string  `json:"requestDate"`
		Status       string  `json:"status"`
		TransRefNum  string  `json:"transRefNum"`
		TxID         string  `json:"txId"`
		AuthCode	 string  `json:"authCode,omitempty"`
	} `json:"dataStatus"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	HTTPStatus      int    `json:"-,omitempty"`
}

type SimobiStatus struct {
	TxID            string `json:"txId,omitempty"`
	TransRefNum     string `json:"transRefNum,omitempty"`
	TxDate          string `json:"transactionDate,omitempty"`
	MerchantName    string `json:"merchantName,omitempty"`
	Amount          int    `json:"amount,omitempty"`
	MobileNo        string `json:"mobileNo,omitempty"`
	MobileNumber    string `json:"mobileNumber,omitempty"`
	ResponseCode    string `json:"responseCode,omitempty"`
	ResponseMessage string `json:"responseMessage,omitempty"`
	Status          string `json:"status,omitempty"`
}

type SimobiPull struct {
	BillerCode  string `json:"billerCode"`
	RequestDate string `json:"requestDate"`
	TxID        string `json:"txId"`
}

type SimobiPullResponse struct {
	DataStatus struct {
		DataStatus struct {
			Amount          string `json:"amount"`
			MerchantName    string `json:"merchantName"`
			MobileNumber    string `json:"mobileNumber"`
			Status          string `json:"status"`
			TransRefNum     string `json:"transRefNum"`
			TransactionDate string `json:"transactionDate"`
			TxID            string `json:"txId"`
			AuthCode 		string `json:"authCode"`
		} `json:"dataStatus"`
	} `json:"dataStatus"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	HTTPStatus      int    `json:"-"`
}

type SimobiCallBack struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	DataStatus      struct {
		TxID         string `json:"txId"`
		TransRefNum  string `json:"transRefNum"`
		MobileNumber string `json:"mobileNumber"`
		RequestDate  string `json:"requestDate"`
		MerchantName string `json:"merchantName"`
		Amount       string `json:"amount"`
		AuthCode     string `json:"authCode"`
		Status       string `json:"status"`
	} `json:"dataStatus"`
	HTTPStatus      int    `json:"-"`
}

//End Struct API
