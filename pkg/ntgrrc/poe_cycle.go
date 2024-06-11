package ntgrrc

import (
	"errors"
	"fmt"
	"net/url"
)

// PoeCyclePower does a power cycle per each given port, whereas
// the given slice contains port numbers (starting with 1) for which the power cycle will happen.
// Returns an ordered list of type PoePortSetting with the new settings of each provided port number
func (session *NtgrrcSession) PoeCyclePower(ports []int) ([]PoePortSetting, error) {
	err := ensureModelIs30x(session, session.address)
	if err != nil {
		return nil, err
	}

	poeExt := &poeExtValues{}

	settings, err := requestPoeConfiguration(session, session.address, poeExt)
	if err != nil {
		return nil, err
	}

	poeSettings := url.Values{
		"hash":   {poeExt.Hash},
		"ACTION": {"Reset"},
	}

	for _, switchPort := range ports {
		if switchPort > len(settings) || switchPort < 1 {
			return nil, errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(settings)))
		}
		poeSettings.Add(fmt.Sprintf("port%d", switchPort-1), "checked")
	}

	result, err := requestPoeSettingsUpdate(session, session.address, poeSettings.Encode())
	if result != "SUCCESS" {
		return nil, errors.New(result)
	}
	if err != nil {
		return nil, err
	}

	settings, err = requestPoeConfiguration(session, session.address, poeExt)
	if err != nil {
		return nil, err
	}

	changedPorts := collectChangedPoePortConfiguration(ports, settings)
	return changedPorts, nil
}
