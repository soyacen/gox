package hashring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		hr := New()
		assert.Empty(t, hr.Nodes())
	})

	t.Run("multiple nodes", func(t *testing.T) {
		hr := New("a", "b", "c")
		assert.Equal(t, []string{"a", "b", "c"}, hr.Nodes())
	})

	t.Run("duplicate nodes deduped by map", func(t *testing.T) {
		hr := New("a", "a", "b")
		// map dedupes but we only have 2 unique keys
		assert.Equal(t, []string{"a", "b"}, hr.Nodes())
	})
}

func TestNewWeighted(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		weights := map[string]int{"a": 1, "b": 2}
		hr := NewWeighted(weights)
		assert.Equal(t, []string{"a", "b"}, hr.Nodes())
	})

	t.Run("zero weight treated as one", func(t *testing.T) {
		weights := map[string]int{"a": 0}
		hr := NewWeighted(weights)
		assert.Equal(t, []string{"a"}, hr.Nodes())
	})

	t.Run("negative weight treated as one", func(t *testing.T) {
		weights := map[string]int{"a": -1}
		hr := NewWeighted(weights)
		assert.Equal(t, []string{"a"}, hr.Nodes())
	})
}

func TestGet(t *testing.T) {
	t.Run("empty ring", func(t *testing.T) {
		hr := New()
		node, ok := hr.Get("key")
		assert.False(t, ok)
		assert.Empty(t, node)
	})

	t.Run("single node", func(t *testing.T) {
		hr := New("only")
		node, ok := hr.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "only", node)
	})

	t.Run("multiple nodes", func(t *testing.T) {
		hr := New("a", "b", "c")
		node, ok := hr.Get("some-key")
		assert.True(t, ok)
		assert.Contains(t, []string{"a", "b", "c"}, node)
	})

	t.Run("consistent result for same key", func(t *testing.T) {
		hr := New("a", "b", "c", "d")
		n1, _ := hr.Get("user:123")
		n2, _ := hr.Get("user:123")
		assert.Equal(t, n1, n2)
	})
}

func TestGetN(t *testing.T) {
	t.Run("empty ring", func(t *testing.T) {
		hr := New()
		assert.Nil(t, hr.GetN("key", 2))
	})

	t.Run("n zero", func(t *testing.T) {
		hr := New("a", "b")
		assert.Nil(t, hr.GetN("key", 0))
	})

	t.Run("n negative", func(t *testing.T) {
		hr := New("a", "b")
		assert.Nil(t, hr.GetN("key", -1))
	})

	t.Run("n greater than node count", func(t *testing.T) {
		hr := New("a", "b")
		nodes := hr.GetN("key", 5)
		assert.Len(t, nodes, 2)
		assert.Equal(t, []string{"a", "b"}, sorted(nodes))
	})

	t.Run("returns distinct nodes", func(t *testing.T) {
		hr := New("a", "b", "c", "d")
		nodes := hr.GetN("key", 3)
		assert.Len(t, nodes, 3)
		// All returned nodes should be distinct
		seen := make(map[string]bool)
		for _, n := range nodes {
			assert.False(t, seen[n], "node %s duplicated", n)
			seen[n] = true
		}
	})

	t.Run("consistent ordering for same key", func(t *testing.T) {
		hr := New("a", "b", "c", "d")
		n1 := hr.GetN("user:456", 3)
		n2 := hr.GetN("user:456", 3)
		assert.Equal(t, n1, n2)
	})
}

func TestAdd(t *testing.T) {
	t.Run("add new node", func(t *testing.T) {
		hr := New("a", "b")
		hr.Add("c", 1)
		assert.Equal(t, []string{"a", "b", "c"}, hr.Nodes())
	})

	t.Run("add with default weight", func(t *testing.T) {
		hr := New("a")
		hr.Add("b", 0)
		// weight 0 should be treated as 1
		assert.Equal(t, []string{"a", "b"}, hr.Nodes())
	})

	t.Run("update existing node weight", func(t *testing.T) {
		hr := NewWeighted(map[string]int{"a": 1, "b": 1})
		// Verify both have equal presence by sampling many keys
		before := countDistribution(hr, 1000)
		assert.InDelta(t, before["a"], before["b"], 200)

		hr.Add("b", 5)
		after := countDistribution(hr, 1000)
		// b should now get significantly more keys
		assert.Greater(t, after["b"], after["a"])
	})

	t.Run("chaining", func(t *testing.T) {
		hr := New("a").Add("b", 1).Add("c", 1)
		assert.Equal(t, []string{"a", "b", "c"}, hr.Nodes())
	})
}

func TestRemove(t *testing.T) {
	t.Run("remove existing", func(t *testing.T) {
		hr := New("a", "b", "c")
		ok := hr.Remove("b")
		assert.True(t, ok)
		assert.Equal(t, []string{"a", "c"}, hr.Nodes())
	})

	t.Run("remove non-existing", func(t *testing.T) {
		hr := New("a", "b")
		ok := hr.Remove("z")
		assert.False(t, ok)
		assert.Equal(t, []string{"a", "b"}, hr.Nodes())
	})

	t.Run("keys remap after remove", func(t *testing.T) {
		hr := New("a", "b", "c")
		node, _ := hr.Get("my-key")
		assert.Contains(t, []string{"a", "b", "c"}, node)

		hr.Remove(node)
		newNode, ok := hr.Get("my-key")
		assert.True(t, ok)
		assert.NotEqual(t, node, newNode)
	})
}

func TestDistribution(t *testing.T) {
	nodes := []string{"node1", "node2", "node3", "node4", "node5"}
	hr := New(nodes...)

	// Count how many keys each node gets
	counts := countDistribution(hr, 10000)

	// Each node should get roughly 20%
	for _, node := range nodes {
		count := counts[node]
		assert.InDelta(t, 2000, count, 500,
			"node %s got %d keys, expected ~2000", node, count)
	}
}

func TestWeightedDistribution(t *testing.T) {
	// node1 weight 1, node2 weight 3
	weights := map[string]int{"light": 1, "heavy": 3}
	hr := NewWeighted(weights)

	counts := countDistribution(hr, 10000)

	// heavy should get roughly 3x more keys than light
	lightCount := counts["light"]
	heavyCount := counts["heavy"]

	assert.Greater(t, heavyCount, lightCount)
	assert.InDelta(t, float64(heavyCount)/float64(lightCount), 3.0, 0.5,
		"expected heavy/light ratio ~3, got %.2f", float64(heavyCount)/float64(lightCount))
}

func TestNodes(t *testing.T) {
	hr := New("c", "a", "b")
	assert.Equal(t, []string{"a", "b", "c"}, hr.Nodes())
}

// countDistribution samples many keys and counts which node each maps to.
func countDistribution(hr *HashRing, samples int) map[string]int {
	counts := make(map[string]int)
	for i := 0; i < samples; i++ {
		key := fmt.Sprintf("key-%d", i)
		node, _ := hr.Get(key)
		counts[node]++
	}
	return counts
}

// sorted returns a sorted copy of the slice.
func sorted(ss []string) []string {
	out := make([]string, len(ss))
	copy(out, ss)
	for i := 0; i < len(out)-1; i++ {
		for j := i + 1; j < len(out); j++ {
			if out[i] > out[j] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

func BenchmarkGet(b *testing.B) {
	hr := New("node1", "node2", "node3", "node4", "node5")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hr.Get(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkGetN(b *testing.B) {
	hr := New("node1", "node2", "node3", "node4", "node5")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		hr.GetN(fmt.Sprintf("key-%d", i), 3)
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hr := New("node1", "node2", "node3")
		hr.Add("node4", 1)
	}
}
