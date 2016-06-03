package main

import (
	"github.com/coreos/fleet/schema"
	"errors"
)

const UnitStateErrorText = "Error retrieving unit states"
const SetUnitTargetStateErrorText = "Error setting unit target state"

type mockFleetApi struct {
}

func (m *mockFleetApi) UnitStates() ([]*schema.UnitState, error) {
	return []*schema.UnitState{&schema.UnitState{}}, nil
}

func (m *mockFleetApi) SetUnitTargetState(name, target string) error {
	return nil
}

type mockFleetApiError struct {
}

func (m *mockFleetApiError) UnitStates() ([]*schema.UnitState, error) {
	return []*schema.UnitState{&schema.UnitState{}}, errors.New(UnitStateErrorText)
}

func (m *mockFleetApiError) SetUnitTargetState(name, target string) error {
	return errors.New(SetUnitTargetStateErrorText)
}
