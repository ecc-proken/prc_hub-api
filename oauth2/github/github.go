package github

import "errors"

type Client struct {
	ClientId     string
	ClientSecret string
}

var client *Client = nil

// クライアントId, クライアントシークレット設定
func SetClient(clientId string, clientSecret string) error {
	// TODO: clientId, clientSecretの検証
	client = &Client{clientId, clientSecret}
	return nil
}

// クライアント情報取得
func GetClient() (*Client, error) {
	if client == nil {
		return nil, errors.New("clientId and clientSecret do not set")
	}
	return client, nil
}
