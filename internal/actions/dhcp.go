package actions

import (
	"fmt"
)

type DHCPAction interface {
	ToString() string
	GetName() string
	GetProviderState() interface{}
}

type BaseDHCPAction struct {
	DHCPAction
	ProviderState interface{}
	Name          string
}

func (action *BaseDHCPAction) GetName() string {
	return action.Name
}

func (action *BaseDHCPAction) GetProviderState() interface{} {
	return action.ProviderState
}

type DHCPCreateReservationAction struct {
	BaseDHCPAction
	MAC     string
	Address string
}

func (action *DHCPCreateReservationAction) ToString() string {
	return fmt.Sprintf("+ %s [%s] = %s", action.GetName(), action.MAC, action.Address)
}

type DHCPDeleteReservationAction struct {
	BaseDHCPAction
}

func (action *DHCPDeleteReservationAction) ToString() string {
	if action.GetName() == "" {
		return fmt.Sprintf("- unnamed (%s)", action.ProviderState)
	}
	return fmt.Sprintf("- %s", action.GetName())
}

type DHCPUpdateReservationAction struct {
	BaseDHCPAction
	OldMAC     string
	NewMAC     string
	OldAddress string
	NewAddress string
}
