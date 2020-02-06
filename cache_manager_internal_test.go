package cache

import (
	"testing"
	"time"
)

func TestDistribution(t *testing.T) {
	cm, _ := NewManager(2, 6)

	err := cm.Put(0, 0, 300, 1)
	if err != nil {
		t.Error(err.Error())
	}
	err = cm.Put(2, 0, 300, 1)
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
	if cm.caches[0].Len() != 2 {
		t.Errorf("sharding error - cache should contain 2 elements, but contains %d", cm.caches[0].Len())
	}
	if cm.caches[1].Len() != 0 {
		t.Errorf("sharding error - cache should contain 0 elements, but contains %d", cm.caches[0].Len())
	}
}

func TestPut(t *testing.T) {
	cm, _ := NewManager(2, 6)

	err := cm.Put(0, 0, 300, 1)
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(100 * time.Millisecond)
	if cm.caches[0].Len() != 1 {
		t.Errorf("sharding error - cache should contain 1 element, but contains %d", cm.caches[0].Len())
	}
}
