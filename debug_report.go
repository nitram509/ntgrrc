package main

import (
	"fmt"
)

type DebugReportCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (drc *DebugReportCommand) Run(args *GlobalOptions) error {
	args.Verbose = true
	model, _, err := readTokenAndModel2GlobalOptions(args, drc.Address)
	if err != nil {
		printDebugNotLoggedIn(args, drc.Address, err)
	}
	if isModel30x(model) {
		printDebugReportGs30xEPx(drc.Address)
	}
	if isModel316(model) {
		printDebugReportGs316EPx(drc.Address)
	}
	return nil
}

func printDebugNotLoggedIn(args *GlobalOptions, host string, err error) {
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

func printDebugLoggedIn(args *GlobalOptions, model NetgearModel, host string) {
	println("---[DEBUG: model " + model + "]---")

	var reqUrls []string
	if isModel30x(model) {
		reqUrls = []string{
			fmt.Sprintf("http://%s/getPoePortStatus.cgi", host),
			fmt.Sprintf("http://%s/PoEPortConfig.cgi", host),
			fmt.Sprintf("http://%s/port_status.cgi", host),
			fmt.Sprintf("http://%s/dashboard.cgi", host),
		}
	}
	if isModel316(model) {
		reqUrls = []string{
			fmt.Sprintf("http://%s/iss/specific/poe.html", host),
			fmt.Sprintf("http://%s/iss/specific/poePortConf.html", host),
			fmt.Sprintf("http://%s/iss/specific/poePortStatus.html", host),
			fmt.Sprintf("http://%s/iss/specific/dashboard.html", host),
			fmt.Sprintf("http://%s/iss/specific/homepage.html", host),
		}
	}

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

func printDebugReportGs30xEPx(string) {
	//var reqUrl string
	//reqUrl = fmt.Sprintf("http://%s/getPoePortStatus.cgi", drc.Address)
	//requestPage(args, drc.Address, reqUrl)

	//if checkIsLoginRequired(statusPage) {
	//	return errors.New("no content. please, (re-)login first")
	//}
	//println("---[DEBUG]---")
	//println(statusPage)
	//println("---[DEBUG]---")

}

func printDebugReportGs316EPx(string) {
	//var reqUrl string
	//reqUrl = fmt.Sprintf("http://%s/iss/specific/poe.html", drc.Address)
	//requestPage(args, drc.Address, reqUrl)
}
