package quick

import (
	"crypto/tls"
	"io"
	"net/http"

	p "github.com/jeffotoni/quick/internal/print"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	ClientGetInsecure HTTPClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        100,
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	ClientGetSecure HTTPClient = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        100,
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 100,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func (q *Quick) GetRequest(path string, handlerFunc HandleFunc) {

	path, params, partternExist := extractParamsPattern(path)

	route := Route{
		Pattern: partternExist,
		Path:    path,
		Params:  params,
		handler: extractCustomParamsGet(q, path, params, handlerFunc),
		Method:  MethodGet,
	}

	q.appendRoute(&route)
	q.mux.HandleFunc(path, route.handler)
}

func extractCustomParamsGet(q *Quick, pathTmp, paramsPath string, handlerFunc HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		v := req.Context().Value(myContextKey)
		if v == nil {
			http.NotFound(w, req)
			return
		}

		cval := v.(ctxServeHttp)
		querys := make(map[string]string)
		queryParams := req.URL.Query()
		for key, values := range queryParams {
			querys[key] = values[0]
		}
		headersMap := extractHeaders(*req)

		resp, err := ClientGetInsecure.Do(req)

		if err != nil {
			p.Stdout("\033[0;33merror:", err.Error(), "\033[0m\n")
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			p.Stdout("\033[0;33merror:", err.Error(), "\033[0m\n")
		}

		w.Write(b)

		c := &Ctx{
			Response: w,
			Request:  req,
			Params:   cval.ParamsMap,
			Query:    querys,
			//bodyByte: extractBodyBytes(req.Body),
			//bodyByte: extractBodyBytes(req.Body),
			Headers:      headersMap,
			MoreRequests: q.config.MoreRequests,
		}
		execHandleFunc(c, handlerFunc)
	}
}
