package frizzante

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var client http.Client

// HttpGet sends an http request using the GET verb.
func HttpGet(path string, header map[string]string) (res string, err error) {
	if nil == header {
		header = map[string]string{}
	}

	request, requestError := http.NewRequest("GET", path, nil)
	if requestError != nil {
		return "", requestError
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	response, doError := client.Do(request)
	if doError != nil {
		return "", doError
	}
	defer func(Body io.ReadCloser) { err = Body.Close() }(response.Body)

	bodyBytes, readError := io.ReadAll(response.Body)
	if readError != nil {
		return "", readError
	}
	result := string(bodyBytes)

	if response.StatusCode >= 300 {
		return result, fmt.Errorf("server responded with status code '%d'", response.StatusCode)
	}

	return result, nil
}

// HttpDelete sends an http request using the DELETE verb.
func HttpDelete(path string, header map[string]string) (err error) {
	if nil == header {
		header = map[string]string{}
	}

	request, requestError := http.NewRequest("DELETE", path, nil)
	if requestError != nil {
		return requestError
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	response, doError := client.Do(request)
	if doError != nil {
		return doError
	}
	defer func(Body io.ReadCloser) { err = Body.Close() }(response.Body)

	if 200 != response.StatusCode {
		return fmt.Errorf("server responded with status code '%d'", response.StatusCode)
	}

	return nil
}

// HttpPost sends an http request using the POST verb.
func HttpPost(path string, contents string, header map[string]string) (res string, err error) {
	if nil == header {
		header = map[string]string{}
	}

	request, requestError := http.NewRequest("POST", path, bytes.NewBuffer([]byte(contents)))
	if requestError != nil {
		return "", requestError
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	response, doError := client.Do(request)
	if doError != nil {
		return "", doError
	}
	defer func(Body io.ReadCloser) { err = Body.Close() }(response.Body)

	if 200 != response.StatusCode {
		return "", fmt.Errorf("server responded with status code '%d'", response.StatusCode)
	}

	bodyBytes, readError := io.ReadAll(response.Body)
	if readError != nil {
		return "", readError
	}

	result := string(bodyBytes)

	return result, nil
}

// HttpPut sends an http request using the PUT verb.
func HttpPut(path string, header map[string]string, contents string) (res string, err error) {
	if nil == header {
		header = map[string]string{}
	}

	request, requestError := http.NewRequest("PUT", path, bytes.NewBuffer([]byte(contents)))
	if requestError != nil {
		return "", requestError
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	response, doError := client.Do(request)
	if doError != nil {
		return "", doError
	}
	defer func(Body io.ReadCloser) { err = Body.Close() }(response.Body)

	if 200 != response.StatusCode {
		return "", fmt.Errorf("server responded with status code '%d'", response.StatusCode)
	}

	bodyBytes, readError := io.ReadAll(response.Body)
	if readError != nil {
		return "", readError
	}

	result := string(bodyBytes)

	return result, nil
}
