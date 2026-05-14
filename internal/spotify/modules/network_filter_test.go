package modules

import (
	"testing"
)

// TestNetworkFilterModuleImplementsInterface verifies that NetworkFilterModule
// properly implements the Module interface.
func TestNetworkFilterModuleImplementsInterface(t *testing.T) {
	var _ Module = (*NetworkFilterModule)(nil)
}

// TestNewNetworkFilterModule verifies the constructor creates a properly initialized module.
func TestNewNetworkFilterModule(t *testing.T) {
	module := NewNetworkFilterModule(true)
	
	if module.Name() != "network_filter" {
		t.Errorf("Expected name 'network_filter', got '%s'", module.Name())
	}
	
	if !module.Enabled() {
		t.Error("Expected module to be enabled")
	}
	
	blockList := module.GetBlockList()
	if len(blockList) == 0 {
		t.Error("Expected non-empty default block list")
	}
}

// TestNetworkFilterModuleInterfaceMethods verifies all Module interface methods work.
func TestNetworkFilterModuleInterfaceMethods(t *testing.T) {
	module := NewNetworkFilterModule(false)
	
	// Test Name()
	if module.Name() != "network_filter" {
		t.Errorf("Name() = %s, want 'network_filter'", module.Name())
	}
	
	// Test Enabled() and SetEnabled()
	if module.Enabled() {
		t.Error("Enabled() = true, want false")
	}
	module.SetEnabled(true)
	if !module.Enabled() {
		t.Error("After SetEnabled(true), Enabled() = false, want true")
	}
	
	// Test CSS() - should return empty string
	if css := module.CSS(); css != "" {
		t.Errorf("CSS() = %q, want empty string", css)
	}
	
	// Test JS() - should return empty string (for now)
	if js := module.JS(); js != "" {
		t.Errorf("JS() = %q, want empty string", js)
	}
	
	// Test Selectors() - should return empty slice
	if selectors := module.Selectors(); len(selectors) != 0 {
		t.Errorf("Selectors() = %v, want empty slice", selectors)
	}
}

// TestBlockListManagement verifies block list operations work correctly.
func TestBlockListManagement(t *testing.T) {
	module := NewNetworkFilterModule(true)
	
	// Test initial block list
	initialList := module.GetBlockList()
	if len(initialList) == 0 {
		t.Fatal("Expected non-empty initial block list")
	}
	
	// Test AddPattern
	testPattern := "test-ad-domain.com"
	module.AddPattern(testPattern)
	updatedList := module.GetBlockList()
	if len(updatedList) != len(initialList)+1 {
		t.Errorf("After AddPattern, expected %d patterns, got %d", len(initialList)+1, len(updatedList))
	}
	
	found := false
	for _, p := range updatedList {
		if p == testPattern {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Pattern %q not found in block list after AddPattern", testPattern)
	}
	
	// Test RemovePattern
	module.RemovePattern(testPattern)
	finalList := module.GetBlockList()
	if len(finalList) != len(initialList) {
		t.Errorf("After RemovePattern, expected %d patterns, got %d", len(initialList), len(finalList))
	}
	
	for _, p := range finalList {
		if p == testPattern {
			t.Errorf("Pattern %q still found in block list after RemovePattern", testPattern)
		}
	}
	
	// Test SetBlockList
	newList := []string{"pattern1.com", "pattern2.com"}
	module.SetBlockList(newList)
	currentList := module.GetBlockList()
	if len(currentList) != 2 {
		t.Errorf("After SetBlockList, expected 2 patterns, got %d", len(currentList))
	}
}
