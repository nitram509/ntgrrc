package ntgrrc

import "fmt"

func ExampleNtgrrcSession_DoLogin() {
	host := "127.0.0.1"                              // the IP address or host name
	host = host + fmt.Sprintf(":%d", mockServerPort) // only required in unit test context
	passw := "secret"
	session := NewSession()
	err := session.DoLogin(host, passw)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
}

func ExampleNtgrrcSession_DetectNetgearModel() {
	host := "127.0.0.1"                              // the IP address or host name
	host = host + fmt.Sprintf(":%d", mockServerPort) // only required in unit test context
	session := NewSession()                          // for detecting the model, you don't need a password
	model, err := session.DetectNetgearModel(host)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(model)

	// Output: GS316EP
}

func ExampleNtgrrcSession_GetPoePortStatus() {
	host := "127.0.0.1"                              // the IP address or host name
	host = host + fmt.Sprintf(":%d", mockServerPort) // only required in unit test context
	passw := "secret"
	session := NewSession()
	err := session.DoLogin(host, passw)
	if err != nil {
		fmt.Println(err)
	}
	status, err := session.GetPoePortStatus()
	if err != nil {
		fmt.Println(err)
	}

	// print out all port names
	for _, portStatus := range status {
		fmt.Println(portStatus.PortIndex, portStatus.PortName)
	}

	// convenient helper method to print all status
	PrettyPrintPoePortStatus(MarkdownFormat, status) // print all status items

	// Output:
}
