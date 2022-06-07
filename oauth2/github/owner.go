package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Owner struct {
	Name      string `json:"login"`
	Id        uint64 `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

// GitHubの登録情報を取得
func (g *Client) GetOwner(token string) (o Owner, err error) {
	// githubAPI用リクエストの作成
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return
	}
	// Authorizationヘッダを追加
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	// githubAPIにリクエスト送信
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

	// httpステータスコードを確認
	if res.StatusCode != http.StatusOK {
		// エラー内容を受け渡し
		err = errors.New(string(bodyBytes))
		return
	}

	// レスポンスJSONを解析
	err = json.Unmarshal(bodyBytes, &o)
	if err != nil {
		return
	}

	return o, nil
}

type OwnerEmail struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visivility bool   `json:"visibility"`
}

// GitHubに登録されているemail一覧を取得
func (g *Client) getOwnerEmails(token string) (emails []OwnerEmail, err error) {
	// githubAPI用リクエストの作成
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	// Authorizationヘッダを追加
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	// githubAPIにリクエスト送信
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// レスポンスボディを読込
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// httpステータスコードを確認
	if res.StatusCode != http.StatusOK {
		// エラー内容を受け渡し
		return nil, errors.New(string(bodyBytes))
	}

	// レスポンスJSONを解析
	err = json.Unmarshal(bodyBytes, &emails)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

// GitHubに登録されているemailからprimaryに設定されているものを取得
func (g *Client) GetOwnerPrimaryEmail(token string) (e OwnerEmail, err error) {
	// GitHubに登録されているemail一覧を取得
	emails, err := g.getOwnerEmails(token)
	if err != nil {
		return
	}

	// primaryに設定されているemailを取得
	found := false
	for _, email := range emails {
		if email.Primary {
			e = email
			found = true
			break
		}
	}
	if !found {
		// primary emailが無い
		err = errors.New("primary email not found in body of response from \"api.github.com\"")
		return
	}

	return e, nil
}
