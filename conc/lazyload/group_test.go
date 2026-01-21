package lazyload

import (
	"testing"
)

func TestGroup_LoadOrNew(t *testing.T) {
	t.Run("Load existing value", func(t *testing.T) {
		// Arrange
		group := Group[int]{
			New: func(key string) (int, error) {
				return 42, nil
			},
		}
		group.m.Store("testKey", 42)

		// Act
		obj, err, wasLoaded := group.LoadOrNew("testKey", nil)

		// Assert
		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}
		if obj != 42 {
			t.Errorf("Expected value 42, got: %v", obj)
		}
		if !wasLoaded {
			t.Error("Expected value to be loaded from cache")
		}
	})

	t.Run("Create new value", func(t *testing.T) {
		// Arrange
		group := Group[int]{
			New: func(key string) (int, error) {
				return 100, nil
			},
		}

		// Act
		obj, err, wasLoaded := group.LoadOrNew("newKey", nil)

		// Assert
		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}
		if obj != 100 {
			t.Errorf("Expected value 100, got: %v", obj)
		}
		if wasLoaded {
			t.Error("Expected value to be created, not loaded from cache")
		}
	})

	t.Run("Nil function error", func(t *testing.T) {
		// Arrange
		group := Group[int]{}

		// Act
		obj, err, wasLoaded := group.LoadOrNew("newKey", nil)

		// Assert
		if obj != 0 {
			t.Errorf("Expected value to be 0, got: %v", obj)
		}
		if err == nil || err != ErrNilFunction {
			t.Errorf("Expected error to be ErrNilFunction, got: %v", err)
		}
		if wasLoaded {
			t.Error("Expected value to not be loaded")
		}
	})
}
