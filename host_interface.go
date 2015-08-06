package zabbix

import (
	"github.com/AlekSi/reflector"
)

type (
	InterfaceType int
)

const (
	Agent InterfaceType = 1
	SNMP  InterfaceType = 2
	IPMI  InterfaceType = 3
	JMX   InterfaceType = 4
)

// https://www.zabbix.com/documentation/2.0/manual/appendix/api/hostinterface/definitions
type HostInterface struct {
	DNS    string        `json:"dns"`
	IP     string        `json:"ip"`
	Main   int           `json:"main"`
	Port   string        `json:"port"`
	Type   InterfaceType `json:"type"`
	UseIP  int           `json:"useip"`
	HostId string        `json:"hostid"`
	Id     string        `json:"interfaceid"`
}

type HostInterfaces []HostInterface

// Wrapper for host.get: https://www.zabbix.com/documentation/2.0/manual/appendix/api/hostinterface/get
func (api *API) HostsInterfaceGet(params Params) (res HostInterfaces, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("hostinterface.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets host interface by host Id.
func (api *API) HostsInterfacesGetByIds(ids []string) (res HostInterfaces, err error) {
	return api.HostsInterfaceGet(Params{"hostids": ids})
}
