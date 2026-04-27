package slogx

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	attr := Bool("enabled", true)
	assert.Equal(t, "enabled", attr.Key)
	assert.Equal(t, slog.BoolValue(true), attr.Value)
}

func TestInt(t *testing.T) {
	attr := Int("count", 42)
	assert.Equal(t, "count", attr.Key)
	assert.Equal(t, slog.Int64Value(42), attr.Value)
}

func TestUint(t *testing.T) {
	attr := Uint("id", uint(100))
	assert.Equal(t, "id", attr.Key)
	assert.Equal(t, slog.Uint64Value(100), attr.Value)
}

func TestDuration(t *testing.T) {
	d := 5 * time.Second
	attr := Duration("elapsed", d)
	assert.Equal(t, "elapsed", attr.Key)
	assert.Equal(t, slog.DurationValue(d), attr.Value)
}

func TestFloat(t *testing.T) {
	attr := Float("score", 3.14)
	assert.Equal(t, "score", attr.Key)
	assert.Equal(t, slog.Float64Value(3.14), attr.Value)
}

func TestString(t *testing.T) {
	attr := String("name", "test")
	assert.Equal(t, "name", attr.Key)
	assert.Equal(t, slog.StringValue("test"), attr.Value)
}

func TestTime(t *testing.T) {
	now := time.Now()
	attr := Time("created", now)
	assert.Equal(t, "created", attr.Key)
	assert.Equal(t, slog.TimeValue(now), attr.Value)
}

func TestJson(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		attr := Json("data", map[string]int{"a": 1})
		assert.Equal(t, "data", attr.Key)
		assert.Equal(t, slog.StringValue(`{"a":1}`), attr.Value)
	})

	t.Run("with nil", func(t *testing.T) {
		attr := Json("data", nil)
		assert.Equal(t, "data", attr.Key)
		assert.Equal(t, slog.StringValue("<nil>"), attr.Value)
	})
}

func TestError(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		attr := Error("err", assert.AnError)
		assert.Equal(t, "err", attr.Key)
		assert.Equal(t, slog.StringValue(assert.AnError.Error()), attr.Value)
	})

	t.Run("with nil", func(t *testing.T) {
		attr := Error("err", nil)
		assert.Equal(t, "err", attr.Key)
		assert.Equal(t, slog.StringValue("<nil>"), attr.Value)
	})
}

func TestStruct(t *testing.T) {
	t.Run("with struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		p := Person{Name: "Alice", Age: 30}
		attr := Struct("person", p)
		assert.Equal(t, "person", attr.Key)
		assert.Equal(t, slog.KindString, attr.Value.Kind())
		assert.Contains(t, attr.Value.String(), "Name")
		assert.Contains(t, attr.Value.String(), "Alice")
		assert.Contains(t, attr.Value.String(), "Age")
		assert.Contains(t, attr.Value.String(), "30")
	})

	t.Run("with pointer", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}
		cfg := &Config{Host: "localhost", Port: 8080}
		attr := Struct("config", cfg)
		assert.Equal(t, "config", attr.Key)
		assert.Equal(t, slog.KindString, attr.Value.Kind())
		assert.Contains(t, attr.Value.String(), "Host")
		assert.Contains(t, attr.Value.String(), "localhost")
		assert.Contains(t, attr.Value.String(), "Port")
		assert.Contains(t, attr.Value.String(), "8080")
	})

	t.Run("with nil", func(t *testing.T) {
		attr := Struct("value", nil)
		assert.Equal(t, "value", attr.Key)
		assert.Equal(t, slog.StringValue("<nil>"), attr.Value)
	})

	t.Run("with nested struct", func(t *testing.T) {
		type Address struct {
			City  string
			State string
		}
		type Person struct {
			Name    string
			Address Address
		}
		p := Person{Name: "Bob", Address: Address{City: "NYC", State: "NY"}}
		attr := Struct("person", p)
		assert.Equal(t, "person", attr.Key)
		assert.Equal(t, slog.KindString, attr.Value.Kind())
		assert.Contains(t, attr.Value.String(), "Bob")
		assert.Contains(t, attr.Value.String(), "NYC")
		assert.Contains(t, attr.Value.String(), "NY")
	})
}
