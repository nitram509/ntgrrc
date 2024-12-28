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
		fmt.Println("Warning, prior error: " + err.Error())
		printDebugNotLoggedIn(args, drc.Address, err)
	}
	printDebugLoggedIn(args, model, drc.Address)
	return nil
}

func printDebugNotLoggedIn(args *GlobalOptions, host string, err error) {
	fmt.Println("---[DEBUG: not logged in]---")
	fmt.Println(fmt.Sprintf("Not logged in error: %s", err))
	fmt.Println("Please try to login and run `debug-report` command again, in order to detect the model and get even more debug information")
	reqUrls := []string{
		fmt.Sprintf("http://%s/", host),
		fmt.Sprintf("http://%s/login.cgi", host),
		fmt.Sprintf("http://%s/wmi/login", host),
		fmt.Sprintf("http://%s/redirect.html", host),
	}
	for _, reqUrl := range reqUrls {
		body, err := doUnauthenticatedHttpRequestAndReadResponse(args, "GET", reqUrl, "")
		fmt.Println(fmt.Sprintf("---[RESPONSE: %s]---", reqUrl))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		} else {
			fmt.Println(body)
		}
		fmt.Println("---[/RESPONSE]---")
	}
	fmt.Println("---[/DEBUG]---")
}

func printDebugLoggedIn(args *GlobalOptions, model NetgearModel, host string) {
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
		fmt.Println(fmt.Sprintf("---[DEBUG: model '%s']---", model))
		for _, reqUrl := range reqUrls {
			body, err := doHttpRequestAndReadResponse(args, "GET", host, reqUrl, "")
			fmt.Println(fmt.Sprintf("---[RESPONSE: %s]---", reqUrl))
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
			} else if checkIsLoginRequired(body) {
				fmt.Println("WARN: it seems the session token expired, please re-login")
			} else {
				fmt.Println(body)
			}
			fmt.Println("---[/RESPONSE]---")
		}
		fmt.Println("---[/DEBUG]---")
	}
}
