package ntgrrc

func ExampleNtgrrcSession_DoLogin() {
	host := "127.0.0.1" // the IP address or host name
	passw := "secret"
	session := NewSession()
	err := session.DoLogin(host, passw)
	if err != nil {
		panic(err)
	}
}

func ExampleNtgrrcSession_DetectNetgearModel() {
	host := "127.0.0.1" // the IP address or host name
	session := NewSession()
	model, err := session.DetectNetgearModel(host)
	if err != nil {
		panic(err)
	}
	println(model)
	// Output: GS308EPP
}
