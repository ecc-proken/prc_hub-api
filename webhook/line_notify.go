package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type LineNotifyApiResponse struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

func LineNotify(token string, msg string) (err error) {
	form := url.Values{}
	form.Add("message", msg)
	body := strings.NewReader(form.Encode())

	// LINE notify API用リクエストの作成
	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", body)
	if err != nil {
		return
	}
	// ヘッダを追加
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// LINE notify APIにリクエスト送信
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// レスポンスボディを読込
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	// レスポンスJSONを解析
	var resBody LineNotifyApiResponse
	err = json.Unmarshal(bodyBytes, &resBody)
	if err != nil {
		return
	}
	// エラーメッセージ
	if resBody.Status != http.StatusOK {
		err = errors.New(resBody.Message)
	}

	return
}
