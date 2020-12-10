SET SIMOBI_LISTEN_PORT=:9662
SET SIMOBI_CERTIFICATE_FILE=simobi-ssl.crt
SET SIMOBI_KEY_FILE=simobi-ssl.key
SET SIMOBI_ROOT_URL=/paggr-simobi
SET SIMOBI_ORIGIN_HOST=pay-aggr-simobi
SET SIMOBI_PRIVATE_KEY=simobi.key
SET SIMOBI_BACKEND_KEY=backend.key.pub
SET SIMOBI_MERCHANT_ID=recharge:990001,billpay:990007,package:990006,starterpack:990009,esim:990008
SET SIMOBI_BILLER_CODE=990001
SET SIMOBI_CLIENT_KEY=SB-Mid-client-_yCZ7Mjfd-iBtX6Q
SET SIMOBI_SERVER_KEY=SB-Mid-server-IMrP4nx5yR6tq69LZDl-lqCr
SET SIMOBI_MERCHANT_NAME=Smartfren
SET SIMOBI_BRAND_NAME=Smartfren
SET SIMOBI_BACKEND_URL=https://localhost:9610/paggr-be
SET SIMOBI_SNAP_URL=http://localhost
SET SIMOBI_TOKEN_URL=/labs/sb/oauthdev/oauth2/token/index.php 
SET SIMOBI_PUSH_INVOICE_URL=/simobi/push-invoice/index.php
SET SIMOBI_PUSH_STATUS_URL=/simobi/push/index.php
SET SIMOBI_PULL_STATUS_URL=/simobi/pull/index.php
SET SIMOBI_REFUND_URL=/simobi/refund-invoice/index.php
SET SIMOBI_TXTTYPE=0
SET SIMOBI_IBM_CLIENT_ID=f54329c4-2422-46fb-9bd0-5132364d3ae9 
SET SIMOBI_GRANT_TYPE=client_credentials
SET SIMOBI_CLIENT_ID=f54329c4-2422-46fb-9bd0-5132364d3ae9
SET SIMOBI_CLIENT_SECRET=pW8yQ2rF8gD3gT2dR6wN5aF2gD3bF1eO4tF4eI4eB6mH0eY2vP
SET SIMOBI_SCOPE=BSIM-Dev 
go run main.go


