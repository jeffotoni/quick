package quick

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
)

type (
	// QuickMockCtx defines the interface for mocking HTTP methods in Quick.
	QuickMockCtx interface {
		Get(URI string) error
		Post(URI string, body []byte) error
		Put(URI string, body []byte) error
		Delete(URI string) error
	}

	// quickMockCtxJSON implements QuickMockCtx for JSON-based requests.
	quickMockCtxJSON struct {
		Ctx    *Ctx
		Params map[string]string
	}

	// quickMockCtxXML implements QuickMockCtx for XML-based requests with optional content type.
	quickMockCtxXML struct {
		Ctx         *Ctx
		Params      map[string]string
		ContentType string
	}
)

// QuickMockCtxJSON creates a new mock context for JSON content type.
func QuickMockCtxJSON(ctx *Ctx, params map[string]string) QuickMockCtx {
	return &quickMockCtxJSON{
		Ctx:    ctx,
		Params: params,
	}
}

// QuickMockCtxXML creates a new mock context for XML content type.
func QuickMockCtxXML(ctx *Ctx, params map[string]string, contentType string) QuickMockCtx {
	return &quickMockCtxXML{
		Ctx:         ctx,
		Params:      params,
		ContentType: contentType,
	}
}

// Get simulates a GET request for JSON.
func (m quickMockCtxJSON) Get(URI string) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}
	queryMap := make(map[string]string)

	req := httptest.NewRequest("GET", URI, nil)
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", ContentTypeAppJSON)
	m.Ctx.Params = m.Params

	// Parse query parameters into map
	query := req.URL.Query()
	for _, pair := range strings.Split(query.Encode(), "&") {
		spltVal := strings.Split(pair, "=")
		if len(spltVal) > 1 {
			queryMap[spltVal[0]] = spltVal[1]
		}
	}
	m.Ctx.Query = queryMap
	return nil
}

// Post simulates a POST request for JSON.
func (m quickMockCtxJSON) Post(URI string, body []byte) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	req := httptest.NewRequest("POST", URI, io.NopCloser(bytes.NewBuffer(body)))
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", ContentTypeAppJSON)
	m.Ctx.Params = m.Params
	m.Ctx.bodyByte = body
	return nil
}

// Put simulates a PUT request for JSON.
func (m quickMockCtxJSON) Put(URI string, body []byte) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	req := httptest.NewRequest("PUT", URI, io.NopCloser(bytes.NewBuffer(body)))
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", ContentTypeAppJSON)
	m.Ctx.Params = m.Params
	m.Ctx.bodyByte = body
	return nil
}

// Delete simulates a DELETE request for JSON.
func (m quickMockCtxJSON) Delete(URI string) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	req := httptest.NewRequest("DELETE", URI, nil)
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", ContentTypeAppJSON)
	m.Ctx.Params = m.Params
	return nil
}

// Get simulates a GET request for XML with optional content type.
func (m quickMockCtxXML) Get(URI string) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}
	queryMap := make(map[string]string)

	contentT := ContentTypeTextXML
	if len(m.ContentType) != 0 {
		contentT = m.ContentType
	}

	req := httptest.NewRequest("GET", URI, nil)
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", contentT)
	m.Ctx.Params = m.Params

	// Parse query parameters into map
	query := req.URL.Query()
	for _, pair := range strings.Split(query.Encode(), "&") {
		spltVal := strings.Split(pair, "=")
		if len(spltVal) > 1 {
			queryMap[spltVal[0]] = spltVal[1]
		}
	}
	m.Ctx.Query = queryMap
	return nil
}

// Post simulates a POST request for XML with optional content type.
func (m quickMockCtxXML) Post(URI string, body []byte) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	contentT := ContentTypeTextXML
	if len(m.ContentType) != 0 {
		contentT = m.ContentType
	}

	req := httptest.NewRequest("POST", URI, io.NopCloser(bytes.NewBuffer(body)))
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", contentT)
	m.Ctx.Params = m.Params
	m.Ctx.bodyByte = body
	return nil
}

// Put simulates a PUT request for XML with optional content type.
func (m quickMockCtxXML) Put(URI string, body []byte) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	contentT := ContentTypeTextXML
	if len(m.ContentType) != 0 {
		contentT = m.ContentType
	}

	req := httptest.NewRequest("PUT", URI, io.NopCloser(bytes.NewBuffer(body)))
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", contentT)
	m.Ctx.Params = m.Params
	m.Ctx.bodyByte = body
	return nil
}

// Delete simulates a DELETE request for XML with optional content type.
func (m quickMockCtxXML) Delete(URI string) error {
	if m.Ctx == nil {
		return errors.New("ctx is null")
	}

	contentT := ContentTypeTextXML
	if len(m.ContentType) != 0 {
		contentT = m.ContentType
	}

	req := httptest.NewRequest("DELETE", URI, nil)
	m.Ctx.Request = req
	m.Ctx.Request.Header.Set("Content-Type", contentT)
	m.Ctx.Params = m.Params
	return nil
}
