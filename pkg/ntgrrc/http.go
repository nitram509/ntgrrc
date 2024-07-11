package ntgrrc

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func requestPage(args *NtgrrcSession, host string, url string) (string, error) {
	return doHttpRequestAndReadResponse(args, http.MethodGet, host, url, "")
}

func postPage(args *NtgrrcSession, host string, url string, requestBody string) (string, error) {
	return doHttpRequestAndReadResponse(args, http.MethodPost, host, url, requestBody)
}

func doHttpRequestAndReadResponse(args *NtgrrcSession, httpMethod string, host string, requestUrl string, requestBody string) (string, error) {
	model, token, err := readTokenAndModel2GlobalOptions(args, host)
	if err != nil {
		return "", err
	}

	if args.PrintVerbose {
		println("Fetching data from: " + requestUrl)
	}

	if isModel316(model) {
		if strings.Contains(requestUrl, "?") {
			splits := strings.Split(requestUrl, "?")
			requestUrl = splits[0] + "?Gambit=" + token + "&" + splits[1]
		} else {
			requestUrl = requestUrl + "?Gambit=" + token
		}
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
	if args.PrintVerbose {
		println(resp.Status)
	}
	bytes, err := io.ReadAll(resp.Body)
	return string(bytes), err
}

func doUnauthenticatedHttpRequestAndReadResponse(args *NtgrrcSession, httpMethod string, requestUrl string, requestBody string) (string, error) {
	if args.PrintVerbose {
		println("Fetching data from: " + requestUrl)
	}

	req, err := http.NewRequest(httpMethod, requestUrl, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if args.PrintVerbose {
		println(resp.Status)
		for name, values := range resp.Header {
			for _, value := range values {
				println(fmt.Sprintf("Response header: '%s' -- '%s'", name, value))
			}
		}
	}
	bytes, err := io.ReadAll(resp.Body)
	return string(bytes), err
}
