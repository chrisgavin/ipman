package actions

import "fmt"

type DNSAction interface {
	ToString() string
	GetName() string
	GetType() string
}

type BaseDNSAction struct {
	DNSAction
	Name string
	Type string
}

type DNSCreateRecordAction struct {
	BaseDNSAction
	Data string
}

func (action *DNSCreateRecordAction) ToString() string {
	return fmt.Sprintf("+ [%s] %s = %s", action.GetType(), action.GetName(), action.Data)
}

func (action *BaseDNSAction) GetName() string {
	return action.Name
}

func (action *BaseDNSAction) GetType() string {
	return action.Type
}

type DNSDeleteRecordAction struct {
	BaseDNSAction
}

func (action *DNSDeleteRecordAction) ToString() string {
	return fmt.Sprintf("- [%s] %s", action.GetType(), action.GetName())
}

type DNSUpdateRecordAction struct {
	BaseDNSAction
	OldData string
	NewData string
}

func (action *DNSUpdateRecordAction) ToString() string {
	return fmt.Sprintf("* [%s] %s = %s -> %s", action.GetType(), action.GetName(), action.OldData, action.NewData)
}
