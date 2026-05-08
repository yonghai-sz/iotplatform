package httpclient

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Interval           int
	ProxyUrl           string
	InsecureSkipVerify bool
}

func createHttpClient(config *Config) *http.Client {

	dialContext := (&net.Dialer{
		Timeout:   time.Second * 30,
		KeepAlive: time.Second * 30,
	}).DialContext

	tr := &http.Transport{
		DialContext:         dialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     time.Second * 60,
	}
	if config.ProxyUrl != "" {
		proxy := func(_ *http.Request) (*url.URL, error) { return url.Parse(config.ProxyUrl) }
		tr.Proxy = proxy
	}
	if config.InsecureSkipVerify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Transport: tr,
	}
	if config.Interval > 0 {
		timeout := time.Second * time.Duration(config.Interval)
		client.Timeout = timeout
	}
	return client
}

/*
 *
 *
 *
 */

type myHttpClient struct {
	client *http.Client
}

func NewHttpClient() *myHttpClient {
	return &myHttpClient{
		client: createHttpClient(&Config{Interval: 15}),
	}
}

func (m *myHttpClient) Req(url string, method string, data []byte,
	setters ...func(*http.Request)) (body []byte, header http.Header, statusCode int, err error) {

	var request *http.Request
	request, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	for _, setter := range setters {
		setter(request)
	}

	var response *http.Response
	response, err = m.client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	header = response.Header
	statusCode = response.StatusCode
	return
}
