package quick

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
)

type (
	QuickMockCtx interface {
		Get(URI string) error
		Post(URI string, body []byte) error
		Put(URI string, body []byte) error
		Delete(URI string) error
	}

	quickMockCtxJSON struct {
		Ctx    *Ctx
		Params map[string]string
	}

	quickMockCtxXML struct {
		Ctx    *Ctx
		Params map[string]string
	}
)

func QuickMockCtxJSON(ctx *Ctx, params map[string]string) QuickMockCtx {
	return &quickMockCtxJSON{
		Ctx:    ctx,
		Params: params,
	}
}

func (m quickMockCtxJSON) Get(URI string) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}
	queryMap := make(map[string]string)

	req := httptest.NewRequest("GET", URI, nil)
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", "application/json")
	m.Ctx.Params = m.Params
	query := req.URL.Query()
	spltQuery := strings.Split(query.Encode(), "&")

	for i := 0; i < len(spltQuery); i++ {
		spltVal := strings.Split(spltQuery[i], "=")
		if len(spltVal) > 1 {
			queryMap[spltVal[0]] = spltVal[1]
		}
	}

	m.Ctx.Query = queryMap
	return nil
}

func (m quickMockCtxJSON) Post(URI string, body []byte) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	req := httptest.NewRequest("POST", URI, io.NopCloser(bytes.NewBuffer(body)))
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", "application/json")
	m.Ctx.Params = m.Params
	m.Ctx.bodyByte = body
	return nil
}

func (m quickMockCtxJSON) Put(URI string, body []byte) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	req := httptest.NewRequest("PUT", URI, io.NopCloser(bytes.NewBuffer(body)))
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", "application/json")
	m.Ctx.Params = m.Params
	m.Ctx.bodyByte = body
	return nil
}

func (m quickMockCtxJSON) Delete(URI string) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	req := httptest.NewRequest("DELETE", URI, nil)
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", "application/json")
	m.Ctx.Params = m.Params
	return nil
}
