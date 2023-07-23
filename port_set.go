package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type PortSetting string

const (
	Index            Setting = "Index"
	Name             Setting = "Name"
	Speed            Setting = "Speed"
	IngressRateLimit Setting = "IngressRateLimit"
	EgressRateLimit  Setting = "EgressRateLimit"
	FlowControl      Setting = "FlowControl"
)

type Port struct {
	Index            int8
	Name             string
	Speed            string
	IngressRateLimit string
	EgressRateLimit  string
	FlowControl      string
}

type PortSetCommand struct {
	Address          string  `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports            []int   `required:"" help:"port number (starting with 1), use multiple times for setting multiple ports at once" short:"p" name:"port"`
	Name             *string `optional:"" help:"sets the name of a port, 1-16 character limit" short:"n"`
	Speed            string  `optional:"" help:"set the speed and duplex of the port ['100M full', '100M half', '10M full', '10M half', 'Auto', 'Disable']" short:"s"`
	IngressRateLimit string  `optional:"" help:"set an incoming rate limit for the port ['1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit']" short:"i"`
	EgressRateLimit  string  `optional:"" help:"set an outgoing rate limit for the port ['1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit']" short:"o"`
	FlowControl      string  `optional:"" help:"enable/disable flow control on port ['Off', 'On']"`
}

func (portSet *PortSetCommand) Run(args *GlobalOptions) error {
	settings, hash, err := requestPortSettings(args, portSet.Address)
	if err != nil {
		return err
	}

	for _, switchPort := range portSet.Ports {

		if switchPort > len(settings) || switchPort < 1 {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(settings)))
		}

		portSetting := settings[switchPort-1]

		// If the port name was not set by the user, set it to the existing name (otherwise an empty port name is always considered to be the
		// "new" value which blanks the port name on the setting next update)
		if portSet.Name == nil {
			portSet.Name = &portSetting.Name
		}

		name, err := comparePortSettings(Name, portSetting.Name, *portSet.Name)
		if err != nil {
			return err
		}

		speed, err := comparePortSettings(Speed, portSetting.Speed, portSet.Speed)
		if err != nil {
			return err
		}

		inRateLimit, err := comparePortSettings(IngressRateLimit, portSetting.IngressRateLimit, portSet.IngressRateLimit)
		if err != nil {
			return err
		}

		outRateLimit, err := comparePortSettings(EgressRateLimit, portSetting.EgressRateLimit, portSet.EgressRateLimit)
		if err != nil {
			return err
		}

		flowControl, err := comparePortSettings(FlowControl, portSetting.FlowControl, portSet.FlowControl)
		if err != nil {
			return err
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

		result, err := requestPortSettingsUpdate(args, portSet.Address, portUpdateValues.Encode())
		if err != nil {
			return err
		}

		if result != "SUCCESS" {
			return errors.New(result)
		}
	}

	settings, _, err = requestPortSettings(args, portSet.Address)
	if err != nil {
		return err
	}

	changedPorts := collectChangedPortConfiguration(portSet.Ports, settings)
	prettyPrintPortSettings(args.OutputFormat, changedPorts)

	return err
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

func comparePortSettings(name Setting, defaultValue string, newValue string) (string, error) {
	if len(newValue) == 0 && name != Name {
		return defaultValue, nil
	}

	switch name {
	case Name:
		if defaultValue != newValue {
			if len(newValue) <= 16 {
				return newValue, nil
			} else {
				return defaultValue, errors.New("port name could not be set. Port name must be 16 characters or less")
			}
		}
		return defaultValue, nil
	case Speed:
		speed := bidiMapLookup(newValue, portSpeedMap)
		if speed == "unknown" {
			return speed, errors.New("port speed could not be set. Accepted values are: " + valuesAsString(portSpeedMap))
		}
		return speed, nil
	case IngressRateLimit:
		inRateLimit := bidiMapLookup(newValue, portRateLimitMap)
		if inRateLimit == "unknown" {
			return inRateLimit, errors.New("ingress rate limit could not be set. Accepted values are: " + valuesAsString(portRateLimitMap))
		}
		return inRateLimit, nil
	case EgressRateLimit:
		outRateLimit := bidiMapLookup(newValue, portRateLimitMap)
		if outRateLimit == "unknown" {
			return outRateLimit, errors.New("egress rate limit could not be set. Accepted values are: " + valuesAsString(portRateLimitMap))
		}
		return outRateLimit, nil
	case FlowControl:
		flowControl := bidiMapLookup(newValue, portFlowControlMap)
		if flowControl == "unknown" {
			return flowControl, errors.New("flow control could not be set. Accepted values are: " + valuesAsString(portFlowControlMap))
		}
		return flowControl, nil
	default:
		return defaultValue, errors.New("could not find port setting")
	}

}

func requestPortSettings(args *GlobalOptions, host string) (portSettings []Port, hash string, err error) {

	url := fmt.Sprintf("http://%s/dashboard.cgi", host)

	portData, err := requestPage(args, host, url)
	if err != nil {
		return portSettings, hash, err
	}

	if checkIsLoginRequired(portData) {
		return portSettings, hash, errors.New("no content. please, (re-)login first")
	}

	hash, err = findHashInHtml(strings.NewReader(portData))
	if err != nil {
		return portSettings, hash, err
	}

	portSettings, err = findPortSettingsInHtml(strings.NewReader(portData))
	if err != nil {
		return portSettings, hash, err
	}

	return portSettings, hash, err

}

func requestPortSettingsUpdate(args *GlobalOptions, host string, data string) (string, error) {
	url := fmt.Sprintf("http://%s/port_status.cgi", host)
	return postPage(args, host, url, data)
}
