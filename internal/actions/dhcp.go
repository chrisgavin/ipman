package actions

import "fmt"

type DHCPAction interface {
	ToString() string
	GetName() string
}

type BaseDHCPAction struct {
	DHCPAction
	Name string
}

func (action *BaseDHCPAction) GetName() string {
	return action.Name
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
	MAC     string
	Address string
}

func (action *DHCPDeleteReservationAction) ToString() string {
	return fmt.Sprintf("- %s [%s] = %s", action.GetName(), action.MAC, action.Address)
}
