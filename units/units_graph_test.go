package units

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestUnitsGraph(t *testing.T) {
	testUg := unitsGraph{
		nodes:   map[UnitSignature]Units{},
		RWMutex: sync.RWMutex{},
	}
	t.Run("Contains, addNode, and removeNode", func(t *testing.T) {
		aSig := UnitSignature("aSig")
		assert.False(t, testUg.contains(aSig), "graph should not initially contain an unknown signature")
		testUg.removeNode(aSig)
		assert.False(t, testUg.contains(aSig), "removing a nonexistent signature should be a no-op")
		testUg.addNode(nil, aSig)
		assert.True(t, testUg.contains(aSig), "should now contain node")
		testUg.removeNode(aSig)
		assert.False(t, testUg.contains(aSig), "graph should not initially contain an unknown signature") // TODO: this
	})

	t.Run("tryAddNode", func(t *testing.T) {
		// TODO: this
	})

	t.Run("Get", func(t *testing.T) {
		// TODO: this
	})
}
