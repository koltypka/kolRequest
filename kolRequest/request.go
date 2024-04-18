package request

import (
	"net/http"
	"net/url"

	myResult "github.com/koltypka/kolRequest/kolRequest/result"
	"github.com/koltypka/kolRequest/myErr"
)

type Request struct {
	rawUrl     string
	client     http.Client
	header     map[string]string
	parameters map[string]string //TODO написать функцию, очищающую параметры
}

func New(rawUrl string) Request {
	return Request{
		rawUrl:     rawUrl,
		client:     http.Client{},
		header:     make(map[string]string),
		parameters: make(map[string]string)}
}

func (r *Request) AddParam(key, value string) {
	r.parameters[key] = value
}

func (r *Request) AddHeader(key, value string) {
	r.header[key] = value
}

func (r *Request) Get(method string) (data *myResult.ResultHandler, err error) {
	return r.run(method, http.MethodGet)
}

func (r *Request) Head(method string) (data *myResult.ResultHandler, err error) {
	return r.run(method, http.MethodHead)
}

func (r *Request) Put(method string) (data *myResult.ResultHandler, err error) {
	return r.run(method, http.MethodPut)
}

func (r *Request) Post(method string) (data *myResult.ResultHandler, err error) {
	return r.run(method, http.MethodPost)
}

func (r *Request) Delete(method string) (data *myResult.ResultHandler, err error) {
	return r.run(method, http.MethodDelete)
}

func (r *Request) run(method, httpMethod string) (data *myResult.ResultHandler, err error) {
	defer func() { err = myErr.Handle("Request error", err) }()

	parsedURL, err := url.Parse(r.rawUrl)

	if len(parsedURL.Path) > 0 {
		method = parsedURL.Path + method
	}

	myUrl := url.URL{Scheme: parsedURL.Scheme, Host: parsedURL.Host, Path: method}

	req, err := http.NewRequest(httpMethod, myUrl.String(), nil)

	r.prepareHeaders(req)
	req.URL.RawQuery = r.prepareParams().Encode() //подставляем параметры

	result, err := r.client.Do(req)

	defer func() {
		if result != nil {
			_ = result.Body.Close()
		}
	}()

	body := myResult.New(result)

	return body, nil
}

func (r *Request) prepareHeaders(req *http.Request) {
	for key, value := range r.header {
		req.Header.Add(key, value)
	}
}

func (r *Request) prepareParams() url.Values {
	result := url.Values{}

	for key, value := range r.parameters {
		result.Set(key, value)
	}

	return result
}
