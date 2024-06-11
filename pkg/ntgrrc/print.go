package ntgrrc

import (
	"fmt"
	"github.com/nitram509/ntgrrc/pkg/ntgrrc/mapping"
)

type PrintFormat string

const (
	MarkdownFormat PrintFormat = "md"
	JsonFormat     PrintFormat = "json"
)

// PrettyPrintPoePortStatus prints Markdown or JSON information of the PoePortStatus list items
func PrettyPrintPoePortStatus(format PrintFormat, statuses []PoePortStatus) {
	var header = []string{"Port ID", "Port Name", "Status", "PortPwr class", "Voltage (V)", "Current (mA)", "PortPwr (W)", "Temp. (°C)", "Error status"}
	var content [][]string
	for _, status := range statuses {
		var row []string
		row = append(row, fmt.Sprintf("%d", status.PortIndex))
		row = append(row, status.PortName)
		row = append(row, status.PoePortStatus)
		row = append(row, status.PoePowerClass)
		row = append(row, fmt.Sprintf("%d", status.VoltageInVolt))
		row = append(row, fmt.Sprintf("%d", status.CurrentInMilliAmps))
		row = append(row, fmt.Sprintf("%.2f", status.PowerInWatt))
		row = append(row, fmt.Sprintf("%d", status.TemperatureInCelsius))
		row = append(row, status.ErrorStatus)
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("poe_status", header, content)
	default:
		panic("not implemented format: " + format)
	}
}

func PrettyPrintPoePortSettings(format PrintFormat, settings []PoePortSetting) {
	var header = []string{"Port ID", "Port Name", "Port Power", "Mode", "Priority", "Limit Type", "Limit (W)", "Type", "Longer Detection Time"}
	var content [][]string
	for _, setting := range settings {
		var row []string
		row = append(row, fmt.Sprintf("%d", setting.PortIndex))
		row = append(row, setting.PortName)
		row = append(row, asTextPortPower(setting.PortPwr))
		row = append(row, mapping.BidiMapLookup(setting.PwrMode, mapping.PwrModeMap))
		row = append(row, mapping.BidiMapLookup(setting.PortPrio, mapping.PortPrioMap))
		row = append(row, mapping.BidiMapLookup(setting.LimitType, mapping.LimitTypeMap))
		row = append(row, setting.PwrLimit)
		row = append(row, mapping.BidiMapLookup(setting.DetecType, mapping.DetecTypeMap))
		row = append(row, mapping.BidiMapLookup(setting.LongerDetect, mapping.LongerDetectMap))
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("poe_settings", header, content)
	default:
		panic("not implemented format: " + format)
	}
}

func asTextPortPower(portPwr bool) string {
	if portPwr {
		return "enabled"
	}
	return "disabled"
}

func PrettyPrintPortSettings(format PrintFormat, settings []Port) {
	var header = []string{"Port ID", "Port Name", "Speed", "Ingress Limit", "Egress Limit", "Flow Control"}
	var content [][]string
	for _, setting := range settings {
		var row []string
		row = append(row, fmt.Sprintf("%d", setting.Index))
		row = append(row, setting.Name)
		setting.Speed = mapping.BidiMapLookup(setting.Speed, mapping.PortSpeedMap)
		row = append(row, setting.Speed)
		setting.IngressRateLimit = mapping.BidiMapLookup(setting.IngressRateLimit, mapping.PortRateLimitMap)
		row = append(row, setting.IngressRateLimit)
		setting.EgressRateLimit = mapping.BidiMapLookup(setting.EgressRateLimit, mapping.PortRateLimitMap)
		row = append(row, setting.EgressRateLimit)
		setting.FlowControl = mapping.BidiMapLookup(setting.FlowControl, mapping.PortFlowControlMap)
		row = append(row, setting.FlowControl)
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("port_settings", header, content)
	default:
		panic("not implemented format: " + format)
	}
}
