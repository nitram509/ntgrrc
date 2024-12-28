package main

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type PortSettingKey string

const (
	Index            PortSettingKey = "Index"
	Name             PortSettingKey = "Name"
	Speed            PortSettingKey = "Speed"
	IngressRateLimit PortSettingKey = "IngressRateLimit"
	EgressRateLimit  PortSettingKey = "EgressRateLimit"
	FlowControl      PortSettingKey = "FlowControl"
)

type PortSetting struct {
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
	FlowControl      string  `optional:"" help:"enable/disable flow control on port ['Off', 'On']" short:"c"`
}

func (portSet *PortSetCommand) Run(args *GlobalOptions) error {
	model := args.model
	if len(model) == 0 {
		var err error
		model, err = detectNetgearModel(args, portSet.Address)
		if err != nil {
			return err
		}
	}
	if isModel30x(model) {
		return portSet.runPortSetGs30xEPx(args)
	}
	if isModel316(model) {
		return portSet.runPortSetGs316EPx(args)
	}
	panic(fmt.Sprintf("model '%s' not supported", model))
}

func (portSet *PortSetCommand) runPortSetGs30xEPx(args *GlobalOptions) error {
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

		ingressRateLimit, err := comparePortSettings(IngressRateLimit, portSetting.IngressRateLimit, portSet.IngressRateLimit)
		if err != nil {
			return err
		}

		egressRateLimit, err := comparePortSettings(EgressRateLimit, portSetting.EgressRateLimit, portSet.EgressRateLimit)
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
			"IngressRate":  {ingressRateLimit},
			"EgressRate":   {egressRateLimit},
			"priority":     {"0"},
		}

		requestUrl := fmt.Sprintf("http://%s/port_status.cgi", portSet.Address)
		result, err := postPage(args, portSet.Address, requestUrl, portUpdateValues.Encode())
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
	prettyPrintPortSettings(args.model, args.OutputFormat, changedPorts)

	return err
}

func (portSet *PortSetCommand) runPortSetGs316EPx(args *GlobalOptions) (err error) {
	_, token, err := readTokenAndModel2GlobalOptions(args, portSet.Address)
	if err != nil {
		return err
	}

	currentSettings, _, err := requestPortSettings(args, portSet.Address)
	if err != nil {
		return err
	}

	for _, portId := range portSet.Ports {
		const gs316MaxPorts = 16
		if portId < 1 || portId > gs316MaxPorts {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", portId, gs316MaxPorts))
		}

		currentSetting := currentSettings[portId-1]

		newSetting, err := createPortSettingUpdatePayloadGs316ep(portSet, currentSetting, token, strconv.Itoa(portId))
		if err != nil {
			return err
		}

		requestUrl := fmt.Sprintf("http://%s/iss/specific/dashboard.html", portSet.Address)
		result, err := postPage(args, portSet.Address, requestUrl, newSetting.Encode())
		if err != nil {
			return err
		}

		if result != "SUCCESS" {
			return errors.New(result)
		}
	}

	updatedSettings, _, err := requestPortSettings(args, portSet.Address)
	if err != nil {
		return err
	}

	updatedSettings = filter(updatedSettings, func(status PortSetting) bool {
		return slices.Contains(portSet.Ports, int(status.Index))
	})
	prettyPrintPortSettings(args.model, args.OutputFormat, updatedSettings)

	return err
}

func createPortSettingUpdatePayloadGs316ep(portSet *PortSetCommand, currentSetting PortSetting, token string, portId string) (url.Values, error) {
	// If the port name was not set by the user, set it to the existing name (otherwise an empty port name is always considered to be the
	// "new" value which blanks the port name on the setting next update)
	if portSet.Name == nil {
		portSet.Name = &currentSetting.Name
	}

	newSetting := url.Values{
		"Gambit":    {token},
		"TYPE":      {"portInfo"},
		"PORT_NO":   {portId},
		"PORT_NAME": {*portSet.Name},
		// default values, for all requests (not entirely sure about the meaning)
		"COLOR1G":    {"NOTSET"},
		"COLOR100M":  {"NOTSET"},
		"FREQUENCY":  {"-1"},
		"BRIGHTNESS": {"undefined"},
		"STATUS":     {"0"},
	}

	if portSet.IngressRateLimit != "" {
		newVal := bidiMapLookup(portSet.IngressRateLimit, portRateLimitMap)
		if newVal == unknown {
			return nil, errors.New(fmt.Sprintf("port ingres setting '%s' could not be set. Accepted values are: %s", portSet.IngressRateLimit, valuesAsString(portRateLimitMap)))
		}
		newSetting.Add("INGRESS", newVal)
	} else {
		newSetting.Add("INGRESS", "NOTSET")
	}

	if portSet.EgressRateLimit != "" {
		newVal := bidiMapLookup(portSet.EgressRateLimit, portRateLimitMap)
		if newVal == unknown {
			return nil, errors.New(fmt.Sprintf("port egress setting '%s' could not be set. Accepted values are: %s", portSet.EgressRateLimit, valuesAsString(portRateLimitMap)))
		}
		newSetting.Add("EGRESS", newVal)
	} else {
		newSetting.Add("EGRESS", "NOTSET")
	}

	if portSet.FlowControl != "" {
		flowControlValue := "4"
		if strings.ToLower(portSet.FlowControl) == "off" {
			flowControlValue = "1"
		}
		newSetting.Add("FLOW_CONTROL", flowControlValue)
	} else {
		newSetting.Add("FLOW_CONTROL", "NOTSET")
	}

	if portSet.Speed != "" {
		switch portSet.Speed {
		case portSpeedAuto:
			newSetting.Add("PORT_CTRL_MODE", "1")
		case portSpeedDisable:
			newSetting.Add("PORT_CTRL_MODE", "3")
		case portSpeed10Mhalf:
			newSetting.Add("PORT_CTRL_MODE", "2")
			newSetting.Add("PORT_CTRL_SPEED", "1")
			newSetting.Add("PORT_CTRL_DUPLEX", "2")
		case portSpeed10Mfull:
			newSetting.Add("PORT_CTRL_MODE", "2")
			newSetting.Add("PORT_CTRL_SPEED", "1")
			newSetting.Add("PORT_CTRL_DUPLEX", "1")
		case portSpeed100Nhalf:
			newSetting.Add("PORT_CTRL_MODE", "2")
			newSetting.Add("PORT_CTRL_SPEED", "2")
			newSetting.Add("PORT_CTRL_DUPLEX", "2")
		case portSpeed100Mfull:
			newSetting.Add("PORT_CTRL_MODE", "2")
			newSetting.Add("PORT_CTRL_SPEED", "2")
			newSetting.Add("PORT_CTRL_DUPLEX", "1")
		default:
			return nil, errors.New(fmt.Sprintf("port speed setting '%s' could not be set. Accepted values are: %s", portSet.Speed, valuesAsString(portSpeedMap)))
		}
	} else {
		newSetting.Add("PORT_CTRL_MODE", "NOTSET")
		newSetting.Add("PORT_CTRL_SPEED", "NOTSET")
		newSetting.Add("PORT_CTRL_DUPLEX", "NOTSET")
	}
	return newSetting, nil
}

func collectChangedPortConfiguration(ports []int, settings []PortSetting) (changedPorts []PortSetting) {
	for _, configuredPort := range ports {
		for _, portSetting := range settings {
			if int(portSetting.Index) == configuredPort {
				changedPorts = append(changedPorts, portSetting)
			}
		}
	}

	return changedPorts
}

func comparePortSettings(name PortSettingKey, defaultValue string, newValue string) (string, error) {
	if len(newValue) == 0 && name != Name {
		return defaultValue, nil
	}

	switch name {
	case Name:
		if defaultValue != newValue {
			if len(newValue) <= 16 {
				return newValue, nil
			} else {
				return defaultValue, errors.New("port name could not be set. PortSetting name must be 16 characters or less")
			}
		}
		return defaultValue, nil
	case Speed:
		speed := bidiMapLookup(newValue, portSpeedMap)
		if speed == unknown {
			return speed, errors.New("port speed could not be set. Accepted values are: " + valuesAsString(portSpeedMap))
		}
		return speed, nil
	case IngressRateLimit:
		inRateLimit := bidiMapLookup(newValue, portRateLimitMap)
		if inRateLimit == unknown {
			return inRateLimit, errors.New("ingress rate limit could not be set. Accepted values are: " + valuesAsString(portRateLimitMap))
		}
		return inRateLimit, nil
	case EgressRateLimit:
		outRateLimit := bidiMapLookup(newValue, portRateLimitMap)
		if outRateLimit == unknown {
			return outRateLimit, errors.New("egress rate limit could not be set. Accepted values are: " + valuesAsString(portRateLimitMap))
		}
		return outRateLimit, nil
	case FlowControl:
		flowControl := bidiMapLookup(newValue, portFlowControlMap)
		if flowControl == unknown {
			return flowControl, errors.New("flow control could not be set. Accepted values are: " + valuesAsString(portFlowControlMap))
		}
		return flowControl, nil
	default:
		return defaultValue, errors.New("could not find port setting")
	}

}
