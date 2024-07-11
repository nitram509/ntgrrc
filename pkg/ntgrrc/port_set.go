package ntgrrc

import (
	"errors"
	"fmt"
	"github.com/nitram509/ntgrrc/pkg/ntgrrc/mapping"
	"net/url"
	"strings"
)

const (
	indexSetting            setting = "Index"
	nameSetting             setting = "Name"
	speedSetting            setting = "Speed"
	ingressRateLimitSetting setting = "IngressRateLimit"
	egressRateLimitSetting  setting = "EgressRateLimit"
	flowControlSetting      setting = "FlowControl"
)

// Port used to show port settings
type Port struct {
	Index            int8 // port number (starting with 1)
	Name             string
	Speed            string
	IngressRateLimit string
	EgressRateLimit  string
	FlowControl      string
}

// SetPortSettingsRequest used to change port settings
type SetPortSettingsRequest struct {
	Ports            []int   // port number (starting with 1), use multiple times for setting multiple ports at once
	Name             *string // sets the name of a port, 1-16 character limit; use nil, to keep the current name
	Speed            string  // set the speed and duplex of the port ['100M full', '100M half', '10M full', '10M half', 'Auto', 'Disable']
	IngressRateLimit string  // set an incoming rate limit for the port ['1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit']
	EgressRateLimit  string  // set an outgoing rate limit for the port ['1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit']
	FlowControl      string  // enable/disable flow control on port ['Off', 'On']
}

func (session *NtgrrcSession) SetPortSettings(portSet SetPortSettingsRequest) ([]Port, error) {
	currentSettings, hash, err := requestPortSettings(session, session.address)
	if err != nil {
		return nil, err
	}

	err = ensureModelIs30x(session, session.address)
	if err != nil {
		return nil, err
	}

	for _, switchPort := range portSet.Ports {

		if switchPort > len(currentSettings) || switchPort < 1 {
			return nil, errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(currentSettings)))
		}

		portSetting := currentSettings[switchPort-1]

		// If the port name was not set by the user, set it to the existing name (otherwise an empty port name is always considered to be the
		// "new" value which blanks the port name on the setting next update)
		if portSet.Name == nil {
			portSet.Name = &portSetting.Name
		}

		name, err := comparePortSettings(nameSetting, portSetting.Name, *portSet.Name)
		if err != nil {
			return nil, err
		}

		speed, err := comparePortSettings(speedSetting, portSetting.Speed, portSet.Speed)
		if err != nil {
			return nil, err
		}

		inRateLimit, err := comparePortSettings(ingressRateLimitSetting, portSetting.IngressRateLimit, portSet.IngressRateLimit)
		if err != nil {
			return nil, err
		}

		outRateLimit, err := comparePortSettings(egressRateLimitSetting, portSetting.EgressRateLimit, portSet.EgressRateLimit)
		if err != nil {
			return nil, err
		}

		flowControl, err := comparePortSettings(flowControlSetting, portSetting.FlowControl, portSet.FlowControl)
		if err != nil {
			return nil, err
		}

		portUpdateValues := url.Values{
			"hash": {hash},
			fmt.Sprintf("%s%d", "port", portSetting.Index): {"checked"},
			"SPEED":        {speed},
			"FLOW_CONTROL": {flowControl},
			"DESCRIPTION":  {name},
			"IngressRate":  {inRateLimit},
			"EgressRate":   {outRateLimit},
			"priority":     {"0"},
		}

		result, err := requestPortSettingsUpdate(session, session.address, portUpdateValues.Encode())
		if err != nil {
			return nil, err
		}

		if result != "SUCCESS" {
			return nil, errors.New(result)
		}
	}

	currentSettings, _, err = requestPortSettings(session, session.address)
	if err != nil {
		return nil, err
	}

	changedPorts := collectChangedPortConfiguration(portSet.Ports, currentSettings)

	return changedPorts, err
}

func collectChangedPortConfiguration(ports []int, settings []Port) (changedPorts []Port) {
	for _, configuredPort := range ports {
		for _, portSetting := range settings {
			if int(portSetting.Index) == configuredPort {
				changedPorts = append(changedPorts, portSetting)
			}
		}
	}

	return changedPorts
}

func comparePortSettings(name setting, defaultValue string, newValue string) (string, error) {
	if len(newValue) == 0 && name != nameSetting {
		return defaultValue, nil
	}

	switch name {
	case nameSetting:
		if defaultValue != newValue {
			if len(newValue) <= 16 {
				return newValue, nil
			} else {
				return defaultValue, errors.New("port name could not be set. Port name must be 16 characters or less")
			}
		}
		return defaultValue, nil
	case speedSetting:
		speed := mapping.BidiMapLookup(newValue, mapping.PortSpeedMap)
		if speed == "unknown" {
			return speed, errors.New("port speed could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.PortSpeedMap))
		}
		return speed, nil
	case ingressRateLimitSetting:
		inRateLimit := mapping.BidiMapLookup(newValue, mapping.PortRateLimitMap)
		if inRateLimit == "unknown" {
			return inRateLimit, errors.New("ingress rate limit could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.PortRateLimitMap))
		}
		return inRateLimit, nil
	case egressRateLimitSetting:
		outRateLimit := mapping.BidiMapLookup(newValue, mapping.PortRateLimitMap)
		if outRateLimit == "unknown" {
			return outRateLimit, errors.New("egress rate limit could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.PortRateLimitMap))
		}
		return outRateLimit, nil
	case flowControlSetting:
		flowControl := mapping.BidiMapLookup(newValue, mapping.PortFlowControlMap)
		if flowControl == "unknown" {
			return flowControl, errors.New("flow control could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.PortFlowControlMap))
		}
		return flowControl, nil
	default:
		return defaultValue, errors.New("could not find port setting")
	}

}

func requestPortSettings(args *NtgrrcSession, host string) (portSettings []Port, hash string, err error) {
	model, _, err := readTokenAndModel2GlobalOptions(args, host)
	if err != nil {
		return portSettings, hash, err
	}

	var requestUrl string
	if isModel30x(model) {
		requestUrl = fmt.Sprintf("http://%s/dashboard.cgi", host)
	} else if isModel316(model) {
		requestUrl = fmt.Sprintf("http://%s/iss/specific/dashboard.html", host)
	} else {
		panic("model not supported")
	}

	dashboardData, err := requestPage(args, host, requestUrl)
	if err != nil {
		return portSettings, hash, err
	}

	if checkIsLoginRequired(dashboardData) {
		return portSettings, hash, errors.New("no content. please, (re-)login first")
	}

	hash, err = findHashInHtml(model, strings.NewReader(dashboardData))
	if err != nil {
		return portSettings, hash, err
	}

	portSettings, err = findPortSettingsInHtml(model, strings.NewReader(dashboardData))

	if err != nil {
		return portSettings, hash, err
	}

	return portSettings, hash, err

}

func requestPortSettingsUpdate(args *NtgrrcSession, host string, data string) (string, error) {
	requestUrl := fmt.Sprintf("http://%s/port_status.cgi", host)
	return postPage(args, host, requestUrl, data)
}
