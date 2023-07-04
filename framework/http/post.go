package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Post(url string, data interface{}, result interface{}, header map[string]string) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}
	var client = http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	errs := json.Unmarshal(body, &result)
	if err != nil {
		return errs
	}
	fmt.Println(result)
	return nil
}

func Post2(url string, data string, result interface{}, header map[string]string) error {
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	if err != nil {
		return err
	}
	//request.Header.Add("Content-Type", "application/json")
	for k, v := range header {
		request.Header.Add(k, v)
	}
	var client = http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	errs := json.Unmarshal(body, &result)
	if err != nil {
		return errs
	}
	fmt.Println(result)
	return nil
}
