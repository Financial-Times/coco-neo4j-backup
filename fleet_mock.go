package main

import (
	"github.com/coreos/fleet/schema"
	"errors"
)

// UnitStateErrorText expected error text when retrieving unit states.
const UnitStateErrorText = "Error retrieving unit states"

// SetUnitTargetStateErrorText expected error text when setting unit state.
const SetUnitTargetStateErrorText = "Error setting unit target state"

type mockFleetAPI struct {
}

func (m *mockFleetAPI) UnitStates() ([]*schema.UnitState, error) {
	return []*schema.UnitState{&schema.UnitState{}}, nil
}

func (m *mockFleetAPI) SetUnitTargetState(name, target string) error {
	return nil
}

type mockFleetAPIError struct {
}

func (m *mockFleetAPIError) UnitStates() ([]*schema.UnitState, error) {
	return []*schema.UnitState{&schema.UnitState{}}, errors.New(UnitStateErrorText)
}

func (m *mockFleetAPIError) SetUnitTargetState(name, target string) error {
	return errors.New(SetUnitTargetStateErrorText)
}
