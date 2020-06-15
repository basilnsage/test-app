package shared

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func parseRespBody(resp *http.Response) (*[]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()
	return &body, err
}

func RespErrorCheck(resp *http.Response, e error) (int, error) {
	if e != nil {
		return http.StatusInternalServerError, e
	}
	if resp.StatusCode - 400 >= 0 {
		body, err := ioutil.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()
		if err != nil {
			return resp.StatusCode, err
		} else {
			return resp.StatusCode, errors.New(string(body))
		}
	}
	return resp.StatusCode, nil
}
