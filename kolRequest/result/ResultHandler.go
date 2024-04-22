package result

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/koltypka/kolRequest/kolRequest/result/structures"
	"github.com/koltypka/kolRequest/myErr"
)

type ResultHandler struct {
	Body      []byte
	Header    http.Header
	isSuccess bool
	err       error
}

func New(Response *http.Response) *ResultHandler {
	if Response == nil {
		return &ResultHandler{
			Body:      make([]byte, 0),
			Header:    nil,
			isSuccess: false,
		}
	}
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

	err = json.Unmarshal(r.Body, &jsonResult)

	if err != nil {
		return nil, err
	}

	result = jsonResult.(map[string]interface{})

	return result, r.err
}

func (r *ResultHandler) ToStruct(StructType string) (result map[string]interface{}, err error) {
	defer func() { r.err = myErr.Handle("Request error", r.err) }()

	if !r.isSuccess {
		return nil, r.err
	}

	tmpResult := structures.Handler(StructType, r.Body)

	if tmpResult == nil {
		r.err = errors.New("Empty Result")
		return nil, r.err
	}

	result = tmpResult.(map[string]interface{})

	return result, r.err
}
