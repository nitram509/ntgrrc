package ntgrrc

import (
	"fmt"
)

// PrintDebugReport connects to multiple endpoints of the switch and prints their response to the console (STDOUT)
// this is only useful for debugging use cases
func PrintDebugReport(address string) error {
	args := NewSession()
	args.PrintVerbose = true
	// FIXME: sessions should not be loaded from disk
	model, _, err := readTokenAndModel2GlobalOptions(args, address)
	if err != nil {
		println("Warning, prior error: " + err.Error())
		printDebugNotLoggedIn(args, address, err)
	}
	printDebugLoggedIn(args, model, address)
	return nil
}

func printDebugNotLoggedIn(args *NtgrrcSession, host string, err error) {
	println("---[DEBUG: not logged in]---")
	println(fmt.Sprintf("Not logged in error: %s", err))
	println("Please try to login and run `debug-report` command again, in order to detect the model and get even more debug information")
	reqUrls := []string{
		fmt.Sprintf("http://%s/", host),
		fmt.Sprintf("http://%s/login.cgi", host),
		fmt.Sprintf("http://%s/wmi/login", host),
		fmt.Sprintf("http://%s/redirect.html", host),
	}
	for _, reqUrl := range reqUrls {
		body, err := doUnauthenticatedHttpRequestAndReadResponse(args, "GET", reqUrl, "")
		println(fmt.Sprintf("---[RESPONSE: %s]---", reqUrl))
		if err != nil {
			println("ERROR: " + err.Error())
		} else {
			println(body)
		}
		println("---[/RESPONSE]---")
	}
	println("---[/DEBUG]---")
}

func printDebugLoggedIn(args *NtgrrcSession, model NetgearModel, host string) {
	var reqUrls []string
	if !isModel30x(model) {
		reqUrls = append(reqUrls,
			// GS316xx
			fmt.Sprintf("http://%s/iss/specific/poe.html", host),
			fmt.Sprintf("http://%s/iss/specific/poePortConf.html", host),
			fmt.Sprintf("http://%s/iss/specific/poePortStatus.html", host),
			fmt.Sprintf("http://%s/iss/specific/poePortStatus.html?GetData=TRUE", host),
			fmt.Sprintf("http://%s/iss/specific/getPortRate.html", host),
			fmt.Sprintf("http://%s/iss/specific/dashboard.html", host),
			fmt.Sprintf("http://%s/iss/specific/homepage.html", host),
		)
	}
	if !isModel316(model) {
		reqUrls = append(reqUrls,
			// GS30xxx
			fmt.Sprintf("http://%s/getPoePortStatus.cgi", host),
			fmt.Sprintf("http://%s/PoEPortConfig.cgi", host),
			fmt.Sprintf("http://%s/port_status.cgi", host),
			fmt.Sprintf("http://%s/dashboard.cgi", host),
		)
	}
	if len(reqUrls) > 0 {
		println(fmt.Sprintf("---[DEBUG: model '%s']---", model))
		for _, reqUrl := range reqUrls {
			body, err := doHttpRequestAndReadResponse(args, "GET", host, reqUrl, "")
			println(fmt.Sprintf("---[RESPONSE: %s]---", reqUrl))
			if err != nil {
				println("ERROR: " + err.Error())
			} else if checkIsLoginRequired(body) {
				println("WARN: it seems the session token expired, please re-login")
			} else {
				println(body)
			}
			println("---[/RESPONSE]---")
		}
		println("---[/DEBUG]---")
	}
}
