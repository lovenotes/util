package netutil

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Request(method, url string, headers map[string]string, body io.Reader, timeout int) (*http.Response, error) {
	transport := &http.Transport{
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Millisecond * time.Duration(timeout),
		Transport: transport,
	}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	retry := 3
	var resp *http.Response

	for i := 0; i < retry; i++ {
		resp, err = client.Do(req)

		if err != nil {
			time.Sleep(100 * time.Millisecond)

			continue
		} else {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func HttpGet(url string, refer string, headers map[string]string, timeout int) (int, string, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	if refer != "" {
		headers["Referer"] = refer
	}

	resp, err := Request("GET", url, headers, nil, timeout)

	if err != nil {
		return -1, "", err
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	defer resp.Body.Close()
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(body), nil
}

func HttpGetLen(url string, refer string, headers map[string]string, timeout int) (int64, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	if refer != "" {
		headers["Referer"] = refer
	}

	resp, err := Request("GET", url, headers, nil, timeout)

	if err != nil {
		return 0, nil
	}

	contentLen := resp.Header.Get("Content-Length")

	length, err := strconv.ParseInt(contentLen, 10, 64)

	if err != nil {
		return 0, err
	}

	return length, nil
}

//Post
func HttpPost(url string, refer string, headers map[string]string, params string, timeout int) (int, string, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	if refer != "" {
		headers["Referer"] = refer
	}

	resp, err := Request("POST", url, headers, bytes.NewBuffer([]byte(params)), timeout)

	if err != nil {
		return -1, "", err
	}

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	defer resp.Body.Close()
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(body), nil
}

//PostFrom
func HttpPostFrom(url string, refer string, headers map[string]string, params url.Values, timeout int) (int, string, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	if refer != "" {
		headers["Referer"] = refer
	}

	resp, err := Request("POST", url, headers, bytes.NewBuffer([]byte(params.Encode())), timeout)

	if err != nil {
		return -1, "", err
	}

	var reader io.ReadCloser

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	defer resp.Body.Close()
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, string(body), nil
}
