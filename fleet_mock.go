package main

import "github.com/coreos/fleet/schema"

type mockFleetAPI struct {
}

func (m *mockFleetAPI) UnitStates() ([]*schema.UnitState, error) {
	return []*schema.UnitState{&schema.UnitState{}}, nil
}

func (m *mockFleetAPI) SetUnitTargetState(name, target string) error {
	return nil
}
