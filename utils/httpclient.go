package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	cli        http.Client
	MaxRetries int
}

func (fc *HttpClient) ForwardTo(req *http.Request, jsonResp interface{}) error {
	resp, err := fc.do(req)
	if err != nil || resp == nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("response has status:%s and body:%q", resp.Status, rb)
	}

	if jsonResp != nil {
		return json.NewDecoder(resp.Body).Decode(jsonResp)
	}
	return nil
}

func (fc *HttpClient) do(req *http.Request) (resp *http.Response, err error) {
	if resp, err = fc.cli.Do(req); err == nil {
		return
	}

	maxRetries := fc.MaxRetries
	backoff := 10 * time.Millisecond

	for retries := 1; retries < maxRetries; retries++ {
		time.Sleep(backoff)
		backoff *= 2

		if resp, err = fc.cli.Do(req); err == nil {
			break
		}
	}
	return
}

func JsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(t); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
