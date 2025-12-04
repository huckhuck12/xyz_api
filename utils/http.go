package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 15,
}

func Request(url, method string, body map[string]any, headers map[string]string) (*http.Response, int, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request body: %v", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %v", err)
	}

	if headers != nil {
		// 请求头
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		fmt.Println("=========[REQUEST INFO]=========")
		fmt.Println("Request Url:", url)
		fmt.Println("Request Method:", method)

		if req.Body != nil {
			requestBody := req.Body
			var requestBodyBytes []byte
			if requestBody != nil {
				requestBodyBytes, _ = io.ReadAll(requestBody)
			}

			fmt.Println("Request Body:", string(requestBodyBytes))

			req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}

		fmt.Println("Request Headers:")
		for key, values := range req.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}

		fmt.Println("========================")
	}

	var resp *http.Response
	var retryErr error
	//for i := 0; i < 3; i++ {
	resp, retryErr = client.Do(req)
	//if retryErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
	//	break
	//}

	//time.Sleep(time.Second * 2)
	//}

	//defer resp.Body.Close()

	if retryErr != nil {
		return nil, 0, fmt.Errorf("failed to send request: %v", retryErr)
	}

	// 读取响应内容，便于在错误时输出真实错误信息
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 尽量读取 body，但不要因为读取失败而丢失状态码信息
		var bodyBytes []byte
		if resp.Body != nil {
			bodyBytes, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, resp.StatusCode, fmt.Errorf("received 401 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
		}

		return nil, resp.StatusCode, fmt.Errorf("received non-200 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, resp.StatusCode, nil
}
