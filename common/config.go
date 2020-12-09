// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server example. This is by no means a complete server.
//
// If you'd like to test diameter over SSL, generate SSL certificates:
//   go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
//
// And start the server with `-cert_file cert.pem -key_file key.pem`.
//
// By default this server runs in a single OS thread. If you want to
// make it run on more, set the GOMAXPROCS=n environment variable.
// See Go's FAQ for details: http://golang.org/doc/faq#Why_no_multi_CPU

package common

import (
	"os"

	"github.com/kelseyhightower/envconfig"

	log "github.com/sirupsen/logrus"

	be "gitlab.smartfren.com/paggr/libraries"
)

//Config stores global configuration loaded from json file
type Configuration struct {
	ListenPort      string            `required:"true" split_words:"true"`
	RootURL         string            `required:"true" split_words:"true"`
	OriginHost      string            `required:"true" split_words:"true"`
	CertificateFile string            `split_words:"true"`
	KeyFile         string            `split_words:"true"`
	PrivateKey      string            `required:"true" split_words:"true"`
	BackendKey      string            `split_words:"true"`
	BackendURL      string            `required:"true" split_words:"true"`
	MerchantID      map[string]string `required:"true" default:"recharge:990001,billpay:990007,package:990006,starterpack:990009,esim:990008" split_words:"true"`
	ClientKey       string            ` split_words:"true"`
	ServerKey       string            `split_words:"true"`
	ProxyURL        string            `split_words:"true"`
	SnapURL         string            `required:"true" split_words:"true"`
	TokenURL        string            `required:"true" split_words:"true"`
	PushInvoiceURL  string            `required:"true" split_words:"true"`
	PushStatusURL   string            `required:"true" split_words:"true"`
	PullStatusURL   string            `required:"true" split_words:"true"`
	FinishURL       string            `split_words:"true"`
	Timeout         int64             `split_words:"true"`
	MerchantName    string            `split_words:"true"`
	BillerCode      string            `required:"true" split_words:"true"`
	IbmClientId     string            `split_words:"true"`
	TxType          string            `split_words:"true"`
	GrantType       string            `split_words:"true""`
	ClientId        string            `split_words:"true""`
	ClientSecret    string            `split_words:"true""`
	Scope           string            `split_words:"true""`
	RefundURL   	string            `required:"true" split_words:"true"`
}

var Config Configuration
var Signer *be.Signer
var BackendSigner *be.Unsigner
var logger *log.Entry

func LoadConfig() {

	if err := envconfig.Process("SIMOBI", &Config); err != nil {
		log.Infof("Loaded configs: %+v", Config)
		log.Error(err)
		os.Exit(1)
	}
	log.Infof("Loaded configs: %+v", Config)

	var err error
	if Signer, err = be.NewSignerFromFile(Config.PrivateKey); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if BackendSigner, err = be.NewUnsignerFromFile(Config.BackendKey); err != nil {
		log.Error(err)
		os.Exit(1)
	}

}

