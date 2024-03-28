package main

import (
	"io"
	"net/http"
	"strings"
)

func requestPage(args *GlobalOptions, host string, url string) (string, error) {
	return doHttpRequestAndReadResponse(args, http.MethodGet, host, url, "")
}

func postPage(args *GlobalOptions, host string, url string, requestBody string) (string, error) {
	return doHttpRequestAndReadResponse(args, http.MethodPost, host, url, requestBody)
}

func doHttpRequestAndReadResponse(args *GlobalOptions, httpMethod string, host string, requestUrl string, requestBody string) (string, error) {
	model, token, err := readTokenAndModel2GlobalOptions(args, host)
	if err != nil {
		return "", err
	}

	if args.Verbose {
		println("Fetching data from: " + requestUrl)
	}

	if isModel316(model) {
		requestUrl = requestUrl + "?Gambit=" + token
	}

	req, err := http.NewRequest(httpMethod, requestUrl, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	if isModel30x(model) {
		req.Header.Set("Cookie", "SID="+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if args.Verbose {
		println(resp.Status)
	}
	bytes, err := io.ReadAll(resp.Body)
	return string(bytes), err
}
