package result

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/koltypka/kolRequest/myErr"
)

type ResultHandler struct {
	Body      []byte
	Header    http.Header
	isSuccess bool
	err       error
}

func New(Response *http.Response) *ResultHandler {
	res, err := io.ReadAll(Response.Body)

	flag := true
	if err != nil {
		flag = false
	}

	return &ResultHandler{
		Body:      res,
		Header:    Response.Header,
		isSuccess: flag,
		err:       err,
	}
}

func (r *ResultHandler) ToString() (string, error) {
	defer func() { r.err = myErr.Handle("Request error", r.err) }()
	if !r.isSuccess {
		return "", r.err
	}

	result := string(r.Body[:])

	return result, r.err
}

func (r *ResultHandler) ToJson() (result map[string]interface{}, err error) {
	defer func() { r.err = myErr.Handle("Request error", r.err) }()

	if !r.isSuccess {
		return nil, r.err
	}

	var jsonResult interface{}

	json.Unmarshal(r.Body, &jsonResult)

	result = jsonResult.(map[string]interface{})

	return result, r.err
}
