package zabbix

import (
	"fmt"
	"github.com/AlekSi/reflector"
)

type (
	TriggerType  int
        FlagsType    int
        PriorityType int
        ValueFlagsType int
)

const (
	SingleEvents    TriggerType = 0
	MultipleEvents  TriggerType = 1

        PlainTrigger    FlagsType = 0
        DiscoveredTrigger FlagsType = 4

        NotClassified   PriorityType = 0
        Information     PriorityType = 1
        Warning         PriorityType = 2
        Average         PriorityType = 3
        High            PriorityType = 4
        Disaster        PriorityType = 5

        Enabled         StatusType = 0
        Disabled        StatusType = 1

        Ok              ValueType = 0
        Problem         ValueType = 1

	Updated     ValueFlagsType = 0
	Unknown       ValueFlagsType = 1
)

// https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/definitions
type Trigger struct {
	TriggerId      string    `json:"trigerid,omitempty"`
	Description string    `json:"description"`
	Expression string    `json:"expression"`
	Comments string    `json:"comments"`
	Error       string       `json:"error"`
        Flags       int         `json:"flags"`
        Priority        int   `json:"priority"`
        Status          int   `json:"status"`
        TemplateId  string    `json:"templateid"`
        Type            TriggerType `json:"type"`
        Url         string    `json:"url"`
        Value           int   `json:"value"`
        ValueFlags      int     `json:"value_flags"`
}

type TriggerPrototype struct {
	TriggerId      string    `json:"trigerid,omitempty"`
	Description string    `json:"description"`
	Expression string    `json:"expression"`
	Comments string    `json:"comments"`
	Error       string       `json:"error"`
        Flags       int         `json:"flags"`
        Priority        int   `json:"priority"`
        Status          int   `json:"status"`
        Type            TriggerType `json:"type"`
}

type Triggers []Trigger
type TriggerPrototypes []TriggerPrototype

// Converts slice to map by key. Panics if there are duplicate keys.
func (triggers Triggers) ById() (res map[string]Trigger) {
	res = make(map[string]Trigger, len(triggers))
	for _, i := range triggers {
		_, present := res[i.TriggerId]
		if present {
			panic(fmt.Errorf("Duplicate key %s", i.TriggerId))
		}
		res[i.TriggerId] = i
	}
	return
}

// Wrapper for item.get https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/get
func (api *API) TriggerGet(params Params) (res Triggers, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("trigger.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Wrapper for trigger.create: https://www.zabbix.com/documentation/2.0/manual/appendix/api/trigger/create
func (api *API) TriggersCreate(triggers TriggerPrototypes) (err error) {
	response, err := api.CallWithError("trigger.create", triggers)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	triggerids := result["triggerids"].([]interface{})
	for i, id := range triggerids {
		triggers[i].TriggerId = id.(string)
	}
	return
}

// Wrapper for trigger.delete: https://www.zabbix.com/documentation/2.0/manual/appendix/api/trigger/delete
// Cleans TriggerId in all Triggers elements if call succeed.
func (api *API) TriggersDelete(triggers Triggers) (err error) {
	ids := make([]string, len(triggers))
	for i, trigger := range triggers {
		ids[i] = trigger.TriggerId
	}

	err = api.TriggersDeleteByIds(ids)
	if err == nil {
		for i := range triggers {
			triggers[i].TriggerId = ""
		}
	}
	return
}

// Wrapper for trigger.delete: https://www.zabbix.com/documentation/2.0/manual/appendix/api/trigger/delete
func (api *API) TriggersDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("trigger.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	triggerids1, ok := result["triggerids"].([]interface{})
	l := len(triggerids1)
	if !ok {
		// some versions actually return map there
		triggerids2 := result["triggerids"].(map[string]interface{})
		l = len(triggerids2)
	}
	if len(ids) != l {
		err = &ExpectedMore{len(ids), l}
	}
	return
}
