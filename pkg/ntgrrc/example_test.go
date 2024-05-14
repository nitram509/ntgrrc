package ntgrrc

import "fmt"

func ExampleNtgrrcSession_DoLogin() {
	host := "127.0.0.1"                              // the IP address or host name
	host = host + fmt.Sprintf(":%d", mockServerPort) // only required in unit test context
	passw := "secret"
	session := NewSession()
	err := session.DoLogin(host, passw)
	if err != nil {
		println(err)
	}
}

func ExampleNtgrrcSession_DetectNetgearModel() {
	host := "127.0.0.1"                              // the IP address or host name
	host = host + fmt.Sprintf(":%d", mockServerPort) // only required in unit test context
	session := NewSession()                          // for detecting the model, you don't need a password
	model, err := session.DetectNetgearModel(host)
	if err != nil {
		println(err)
	}
	fmt.Print(model)

	// Output: GS316EP
}
