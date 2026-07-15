package waroom

import (
	"fmt"
	"net/http"
	"strings"
)

// APIError は Waroom API サーバーが返したエラーレスポンスを表す。
// トランスポート層のエラー(接続失敗・デコード失敗)は含まない。
type APIError struct {
	StatusCode int      // 400, 401, 402, 404, 422
	Messages   []string // レスポンスボディの error / errors を正規化したもの
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("API error %d", e.StatusCode)
	if text := http.StatusText(e.StatusCode); text != "" {
		msg += " " + text
	}
	if len(e.Messages) > 0 {
		msg += ": " + strings.Join(e.Messages, "; ")
	}
	return msg
}

// optMessages は OptString のエラーメッセージを []string に正規化する。未設定なら nil。
func optMessages(s OptString) []string {
	if v, ok := s.Get(); ok {
		return []string{v}
	}
	return nil
}
