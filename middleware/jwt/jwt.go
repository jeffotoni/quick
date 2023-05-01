package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// var mapEncoding = map[string]string{
// 	"HS256": "sha256",
// 	"HS512": "sha512",
// }

type (
	Config struct {
		Header Header `json:"header"`

		Payload Payload `json:"payload"`

		SecretKey string `json:"secretKey"`

		ExpiresIn string
	}

	Header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	Payload struct {
		Login    string `json:"login"`
		Password string `json:"payload"`
	}
)

var defaultJwtConfig = Config{
	Header: Header{
		Alg: "HS256",
		Typ: "JWT",
	},
	Payload: Payload{
		Login:    "quick",
		Password: "quick.com.br",
	},
	SecretKey: "quick-is-awesome!",
	ExpiresIn: "500s",
}

func New(config ...Config) func(http.Handler) http.Handler {
	c := defaultJwtConfig
	if len(config) > 0 {
		c = config[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signature := genJwtSignature(c)
			log.Printf("data -> %v", signature)
			w.Header().Set("jwt-token", signature)
		})
	}
}

func genJwtSignature(config Config) string {
	hm, errM := json.Marshal(config.Header)

	if errM != nil {
		log.Printf("error: %v", errM)
		return ""
	}

	pm, errM := json.Marshal(config.Payload)

	if errM != nil {
		log.Printf("error: %v", errM)
		return ""
	}

	var (
		b64Header  = strings.TrimRight(base64.URLEncoding.EncodeToString(hm), "=")
		b64Payload = strings.TrimRight(base64.URLEncoding.EncodeToString(pm), "=")
	)

	data := b64Header + "." + b64Payload

	switch config.Header.Alg {
	case "HS256":
		mac := hmac.New(sha256.New, []byte(config.SecretKey))
		_, err := mac.Write([]byte(data))
		if err != nil {
			log.Printf("error: %v", err)
			return ""
		}

		s := strings.TrimRight(base64.URLEncoding.EncodeToString(mac.Sum(nil)), "=")

		cleanResult(&s)

		return data + "." + s
	case "HS512":
		mac := hmac.New(sha512.New, []byte(config.SecretKey))
		_, err := mac.Write([]byte(data))
		if err != nil {
			log.Printf("error: %v", err)
			return ""
		}

		return string(mac.Sum(nil))
	}

	return ""
}

func cleanResult(hash *string) {
	*hash = strings.ReplaceAll(*hash, "+", "-")
	*hash = strings.ReplaceAll(*hash, "/+", "-")
	*hash = strings.ReplaceAll(*hash, "/", "_")
	*hash = strings.ReplaceAll(*hash, "\\", "_")
	*hash = strings.ReplaceAll(*hash, "=", "")
}
