package units

import (
	"errors"
	"fmt"
	"sync"
)

var (
	errGraphAlreadyHas = errors.New("units already exist with specified Signature")
)

// allUnits todo
var allUnits = unitsGraph{
	nodes:   map[UnitSignature]Units{},
	RWMutex: sync.RWMutex{},
}

// unitsGraph todo
type unitsGraph struct {
	nodes map[UnitSignature]Units
	sync.RWMutex
}

// contains todo
func (ug *unitsGraph) contains(id UnitSignature) bool { // TODO: consider making private
	_, exists := ug.Get(id)
	return exists
}

// tryAddNode todo
func (ug *unitsGraph) tryAddNode(toAdd Units) error { // TODO: consider making private
	sig := toAdd.Signature()
	if ug.contains(sig) {
		return errors.Join(errGraphAlreadyHas, fmt.Errorf(`duplicate Signature: %s`, sig))
	}

	ug.addNode(toAdd, sig)
	return nil
}

// addNode TODO
func (ug *unitsGraph) addNode(toAdd Units, sig UnitSignature) {
	ug.Lock()
	defer ug.Unlock()

	ug.nodes[sig] = toAdd
}

// removeNode TODO
func (ug *unitsGraph) removeNode(id UnitSignature) { // TODO: consider making private
	ug.Lock()
	defer ug.Unlock()
	delete(ug.nodes, id)
}

// Get TODO
func (ug *unitsGraph) Get(id UnitSignature) (unit Units, exists bool) {
	ug.RLock()
	defer ug.RUnlock()

	unit, exists = ug.nodes[id]
	return
}
