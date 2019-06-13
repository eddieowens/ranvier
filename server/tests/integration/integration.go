package integration

import (
	"bytes"
	"fmt"
	"github.com/eddieowens/ranvier/server/tests"
	json "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
)

type Integration struct {
	tests.TestSuite
}

type Request struct {
	url         string
	Body        interface{}
	PathParams  map[string]string
	QueryParams map[string]string
}

func (i *Integration) Request(request Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	reader := bytes.NewReader(nil)
	if request.Body != nil {
		d, err := json.Marshal(request.Body)
		if err != nil {
			panic(err)
		}
		reader = bytes.NewReader(d)
	}

	names := make([]string, len(request.PathParams))
	values := make([]string, len(request.PathParams))
	j := 0
	for k, v := range request.PathParams {
		names[j] = k
		values[j] = v
		j++
	}

	init := false
	request.url = "/"
	for k, v := range request.QueryParams {
		if !init {
			request.url += fmt.Sprintf("?%s=%s", k, v)
			init = true
		} else {
			request.url += fmt.Sprintf("&%s=%s", k, v)
		}
	}

	req := httptest.NewRequest(http.MethodGet, request.url, reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	for k, v := range request.QueryParams {
		ctx.Set(k, v)
	}
	ctx.SetParamNames(names...)
	ctx.SetParamValues(values...)

	return ctx, rec
}
