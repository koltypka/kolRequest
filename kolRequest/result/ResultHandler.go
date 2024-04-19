package result

import (
	"encoding/json"
	"encoding/xml"
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

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
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

	json.Unmarshal(r.Body, &jsonResult)

	result = jsonResult.(map[string]interface{})

	return result, r.err
}

func (r *ResultHandler) ToStructRss() (result Channel, err error) {
	defer func() { r.err = myErr.Handle("Request error", r.err) }()

	result = Channel{}

	if !r.isSuccess {
		return result, r.err
	}

	xml.Unmarshal(r.Body, &result)

	return result, r.err
}
