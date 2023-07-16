package request

import (
	"io"
	"net/http"
	"net/url"

	"github.com/koltypka/kolRequest/myErr"
)

type Request struct {
	rawUrl     string
	client     http.Client
	header     map[string]string
	parameters map[string]string
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

func (r *Request) Get(method string) (data []byte, err error) {
	return r.run(method, http.MethodGet)
}

func (r *Request) Post(method string) (data []byte, err error) {
	return r.run(method, http.MethodPost)
}

func (r *Request) run(method, httpMethod string) (data []byte, err error) {
	defer func() { err = myErr.Handle("Request error", err) }()

	parsedURL, err := url.Parse(r.rawUrl)

	if err != nil {
		return nil, err
	}

	myUrl := url.URL{Scheme: parsedURL.Scheme, Host: parsedURL.Host, Path: method}

	req, err := http.NewRequest(httpMethod, myUrl.String(), nil)

	if err != nil {
		return nil, err
	}

	r.prepareHeaders(req)
	req.URL.RawQuery = r.prepareParams().Encode() //подставляем параметры

	result, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { _ = result.Body.Close() }()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

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
