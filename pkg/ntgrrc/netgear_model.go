package ntgrrc

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NetgearModel string

const (
	GS30xEPx NetgearModel = "GS30xEPx"
	GS305EP  NetgearModel = "GS305EP"
	GS305EPP NetgearModel = "GS305EPP"
	GS308EP  NetgearModel = "GS308EP"
	GS308EPP NetgearModel = "GS308EPP"
	GS316EP  NetgearModel = "GS316EP"
	GS316EPP NetgearModel = "GS316EPP"
)

func (session *NtgrrcSession) DetectNetgearModel(host string) (NetgearModel, error) {
	url := fmt.Sprintf("http://%s/", host)
	if session.PrintVerbose {
		println("detecting Netgear switch model: " + url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if session.PrintVerbose {
		println(fmt.Sprintf("HTTP response code %d", resp.StatusCode))
	}
	if resp.StatusCode != 200 {
		println(fmt.Sprintf("Warning: response code was not 200; unusual, but will attempt detection anyway"))
	}
	responseBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return "", err
	}
	model := detectNetgearModelFromResponse(string(responseBody))
	if model == "" {
		return "", errors.New("Can't auto-detect Netgear model from response. You may try using --model parameter ")
	}
	if session.PrintVerbose {
		println(fmt.Sprintf("Detected model %s", model))
	}
	return model, nil
}

func detectNetgearModelFromResponse(body string) NetgearModel {
	if strings.Contains(strings.ToLower(body), "<title>") && strings.Contains(body, "GS316EPP") {
		return GS316EPP
	}
	if strings.Contains(strings.ToLower(body), "<title>") && strings.Contains(body, "GS316EP") {
		return GS316EP
	}
	if strings.Contains(strings.ToLower(body), "<title>") && strings.Contains(body, "Redirect to Login") {
		return GS30xEPx
	}
	return ""
}

func isModel30x(nm NetgearModel) bool {
	return nm == GS305EP || nm == GS305EPP || nm == GS308EP || nm == GS308EPP || nm == GS30xEPx
}

func isModel316(nm NetgearModel) bool {
	return nm == GS316EP || nm == GS316EPP
}

func isSupportedModel(modelName string) bool {
	return isModel30x(NetgearModel(modelName)) || isModel316(NetgearModel(modelName))
}
