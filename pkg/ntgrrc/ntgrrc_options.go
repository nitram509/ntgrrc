package ntgrrc

// NtgrrcSession configuration options for all commands
type NtgrrcSession struct {
	PrintVerbose bool
	TokenDir     string
	model        NetgearModel
	token        string
	address      string
}

func NewSession() *NtgrrcSession {
	return &NtgrrcSession{}
}
