package unsafesql

import (
	"testing"
)

func TestAnd(t *testing.T) {
	conditions := []string{"age > 18", "status = 'active'"}
	combiner := And()
	result := combiner.Combine(conditions)
	expected := "(age > 18 AND status = 'active')"
	if result != expected {
		t.Errorf("And() = %v, want %v", result, expected)
	}

	// Test empty conditions
	result = combiner.Combine([]string{})
	if result != "" {
		t.Errorf("And() with empty conditions = %v, want \"\"", result)
	}
}

func TestOr(t *testing.T) {
	conditions := []string{"status = 'active'", "status = 'pending'"}
	combiner := Or()
	result := combiner.Combine(conditions)
	expected := "(status = 'active' OR status = 'pending')"
	if result != expected {
		t.Errorf("Or() = %v, want %v", result, expected)
	}

	// Test empty conditions
	result = combiner.Combine([]string{})
	if result != "" {
		t.Errorf("Or() with empty conditions = %v, want \"\"", result)
	}
}

func TestWhere(t *testing.T) {
	conditions := []string{"age > 18", "status = 'active'"}
	result := Where(conditions, And())
	expected := "WHERE (age > 18 AND status = 'active')"
	if result != expected {
		t.Errorf("Where() = %v, want %v", result, expected)
	}

	// Test empty conditions
	result = Where([]string{}, And())
	if result != "" {
		t.Errorf("Where() with empty conditions = %v, want \"\"", result)
	}
}

func TestMustWhere(t *testing.T) {
	conditions := []string{"age > 18", "status = 'active'"}
	result := MustWhere(conditions, And())
	expected := "WHERE (age > 18 AND status = 'active')"
	if result != expected {
		t.Errorf("MustWhere() = %v, want %v", result, expected)
	}

	// Test panic on empty conditions
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustWhere() with empty conditions should panic")
		}
	}()
	MustWhere([]string{}, And())
}

func TestHaving(t *testing.T) {
	conditions := []string{"SUM(amount) > 1000", "COUNT(*) > 10"}
	result := Having(conditions, And())
	expected := "HAVING (SUM(amount) > 1000 AND COUNT(*) > 10)"
	if result != expected {
		t.Errorf("Having() = %v, want %v", result, expected)
	}

	// Test empty conditions
	result = Having([]string{}, And())
	if result != "" {
		t.Errorf("Having() with empty conditions = %v, want \"\"", result)
	}
}

func TestMustHaving(t *testing.T) {
	conditions := []string{"SUM(amount) > 1000", "COUNT(*) > 10"}
	result := MustHaving(conditions, And())
	expected := "HAVING (SUM(amount) > 1000 AND COUNT(*) > 10)"
	if result != expected {
		t.Errorf("MustHaving() = %v, want %v", result, expected)
	}

	// Test panic on empty conditions
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustHaving() with empty conditions should panic")
		}
	}()
	MustHaving([]string{}, And())
}

func TestEq(t *testing.T) {
	conditions := []string{}
	result := Eq(conditions, "name", "'john'")
	expected := []string{"(name = 'john')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Eq() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Eq([]string{}, "", "'john'")
	if len(result) != 0 {
		t.Errorf("Eq() with empty field should return empty slice")
	}

	result = Eq([]string{}, "name", "")
	if len(result) != 0 {
		t.Errorf("Eq() with empty value should return empty slice")
	}
}

func TestMustEq(t *testing.T) {
	conditions := []string{}
	result := MustEq(conditions, "name", "'john'")
	expected := []string{"(name = 'john')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustEq() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustEq() with empty field should panic")
		}
	}()
	MustEq([]string{}, "", "'john'")
}

func TestNe(t *testing.T) {
	conditions := []string{}
	result := Ne(conditions, "status", "'inactive'")
	expected := []string{"(status <> 'inactive')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Ne() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Ne([]string{}, "", "'inactive'")
	if len(result) != 0 {
		t.Errorf("Ne() with empty field should return empty slice")
	}
}

func TestMustNe(t *testing.T) {
	conditions := []string{}
	result := MustNe(conditions, "status", "'inactive'")
	expected := []string{"(status <> 'inactive')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNe() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNe() with empty field should panic")
		}
	}()
	MustNe([]string{}, "", "'inactive'")
}

func TestGt(t *testing.T) {
	conditions := []string{}
	result := Gt(conditions, "age", "18")
	expected := []string{"(age > 18)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Gt() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Gt([]string{}, "", "18")
	if len(result) != 0 {
		t.Errorf("Gt() with empty field should return empty slice")
	}
}

func TestMustGt(t *testing.T) {
	conditions := []string{}
	result := MustGt(conditions, "age", "18")
	expected := []string{"(age > 18)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustGt() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGt() with empty field should panic")
		}
	}()
	MustGt([]string{}, "", "18")
}

func TestLt(t *testing.T) {
	conditions := []string{}
	result := Lt(conditions, "age", "65")
	expected := []string{"(age < 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Lt() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Lt([]string{}, "", "65")
	if len(result) != 0 {
		t.Errorf("Lt() with empty field should return empty slice")
	}
}

func TestMustLt(t *testing.T) {
	conditions := []string{}
	result := MustLt(conditions, "age", "65")
	expected := []string{"(age < 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLt() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLt() with empty field should panic")
		}
	}()
	MustLt([]string{}, "", "65")
}

func TestGe(t *testing.T) {
	conditions := []string{}
	result := Ge(conditions, "age", "18")
	expected := []string{"(age >= 18)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Ge() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Ge([]string{}, "", "18")
	if len(result) != 0 {
		t.Errorf("Ge() with empty field should return empty slice")
	}
}

func TestMustGe(t *testing.T) {
	conditions := []string{}
	result := MustGe(conditions, "age", "18")
	expected := []string{"(age >= 18)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustGe() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGe() with empty field should panic")
		}
	}()
	MustGe([]string{}, "", "18")
}

func TestLe(t *testing.T) {
	conditions := []string{}
	result := Le(conditions, "age", "65")
	expected := []string{"(age <= 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Le() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Le([]string{}, "", "65")
	if len(result) != 0 {
		t.Errorf("Le() with empty field should return empty slice")
	}
}

func TestMustLe(t *testing.T) {
	conditions := []string{}
	result := MustLe(conditions, "age", "65")
	expected := []string{"(age <= 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLe() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLe() with empty field should panic")
		}
	}()
	MustLe([]string{}, "", "65")
}

func TestBetween(t *testing.T) {
	conditions := []string{}
	result := Between(conditions, "age", "18", "65")
	expected := []string{"(age BETWEEN 18 AND 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Between() = %v, want %v", result, expected)
	}

	// Test empty field or values
	result = Between([]string{}, "", "18", "65")
	if len(result) != 0 {
		t.Errorf("Between() with empty field should return empty slice")
	}
}

func TestMustBetween(t *testing.T) {
	conditions := []string{}
	result := MustBetween(conditions, "age", "18", "65")
	expected := []string{"(age BETWEEN 18 AND 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustBetween() = %v, want %v", result, expected)
	}

	// Test panic on empty field or values
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBetween() with empty field should panic")
		}
	}()
	MustBetween([]string{}, "", "18", "65")
}

func TestNotBetween(t *testing.T) {
	conditions := []string{}
	result := NotBetween(conditions, "age", "18", "65")
	expected := []string{"(age NOT BETWEEN 18 AND 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("NotBetween() = %v, want %v", result, expected)
	}
}

func TestMustNotBetween(t *testing.T) {
	conditions := []string{}
	result := MustNotBetween(conditions, "age", "18", "65")
	expected := []string{"(age NOT BETWEEN 18 AND 65)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNotBetween() = %v, want %v", result, expected)
	}

	// Test panic on empty field or values
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNotBetween() with empty field should panic")
		}
	}()
	MustNotBetween([]string{}, "", "18", "65")
}

func TestLike(t *testing.T) {
	conditions := []string{}
	result := Like(conditions, "name", "'%john%'")
	expected := []string{"(name LIKE '%john%')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Like() = %v, want %v", result, expected)
	}

	// Test empty field or value
	result = Like([]string{}, "", "'%john%'")
	if len(result) != 0 {
		t.Errorf("Like() with empty field should return empty slice")
	}
}

func TestMustLike(t *testing.T) {
	conditions := []string{}
	result := MustLike(conditions, "name", "'%john%'")
	expected := []string{"(name LIKE '%john%')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLike() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLike() with empty field should panic")
		}
	}()
	MustLike([]string{}, "", "'%john%'")
}

func TestNotLike(t *testing.T) {
	conditions := []string{}
	result := NotLike(conditions, "name", "'%john%'")
	expected := []string{"(name NOT LIKE '%john%')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("NotLike() = %v, want %v", result, expected)
	}
}

func TestMustNotLike(t *testing.T) {
	conditions := []string{}
	result := MustNotLike(conditions, "name", "'%john%'")
	expected := []string{"(name NOT LIKE '%john%')"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNotLike() = %v, want %v", result, expected)
	}

	// Test panic on empty field or value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNotLike() with empty field should panic")
		}
	}()
	MustNotLike([]string{}, "", "'%john%'")
}

func TestIn(t *testing.T) {
	conditions := []string{}
	result := In(conditions, "id", []string{"1", "2", "3"})
	expected := []string{"(id IN (1, 2, 3))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("In() = %v, want %v", result, expected)
	}

	// Test empty field or values
	result = In([]string{}, "", []string{"1", "2", "3"})
	if len(result) != 0 {
		t.Errorf("In() with empty field should return empty slice")
	}

	result = In([]string{}, "id", []string{})
	if len(result) != 0 {
		t.Errorf("In() with empty values should return empty slice")
	}
}

func TestMustIn(t *testing.T) {
	conditions := []string{}
	result := MustIn(conditions, "id", []string{"1", "2", "3"})
	expected := []string{"(id IN (1, 2, 3))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustIn() = %v, want %v", result, expected)
	}

	// Test panic on empty field or values
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustIn() with empty field should panic")
		}
	}()
	MustIn([]string{}, "", []string{"1", "2", "3"})
}

func TestNotIn(t *testing.T) {
	conditions := []string{}
	result := NotIn(conditions, "id", []string{"1", "2", "3"})
	expected := []string{"(id NOT IN (1, 2, 3))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("NotIn() = %v, want %v", result, expected)
	}
}

func TestMustNotIn(t *testing.T) {
	conditions := []string{}
	result := MustNotIn(conditions, "id", []string{"1", "2", "3"})
	expected := []string{"(id NOT IN (1, 2, 3))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNotIn() = %v, want %v", result, expected)
	}

	// Test panic on empty field or values
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNotIn() with empty field should panic")
		}
	}()
	MustNotIn([]string{}, "", []string{"1", "2", "3"})
}

func TestIsNull(t *testing.T) {
	conditions := []string{}
	result := IsNull(conditions, "email")
	expected := []string{"(email IS NULL)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("IsNull() = %v, want %v", result, expected)
	}
}

func TestMustIsNull(t *testing.T) {
	conditions := []string{}
	result := MustIsNull(conditions, "email")
	expected := []string{"(email IS NULL)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustIsNull() = %v, want %v", result, expected)
	}

	// Test panic on empty field
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustIsNull() with empty field should panic")
		}
	}()
	MustIsNull([]string{}, "")
}

func TestIsNotNull(t *testing.T) {
	conditions := []string{}
	result := IsNotNull(conditions, "email")
	expected := []string{"(email IS NOT NULL)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("IsNotNull() = %v, want %v", result, expected)
	}
}

func TestMustIsNotNull(t *testing.T) {
	conditions := []string{}
	result := MustIsNotNull(conditions, "email")
	expected := []string{"(email IS NOT NULL)"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustIsNotNull() = %v, want %v", result, expected)
	}

	// Test panic on empty field
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustIsNotNull() with empty field should panic")
		}
	}()
	MustIsNotNull([]string{}, "")
}

func TestExists(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT 1 FROM orders WHERE user_id = users.id"
	result := Exists(conditions, subquery)
	expected := []string{"(EXISTS (SELECT 1 FROM orders WHERE user_id = users.id))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Exists() = %v, want %v", result, expected)
	}
}

func TestMustExists(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT 1 FROM orders WHERE user_id = users.id"
	result := MustExists(conditions, subquery)
	expected := []string{"(EXISTS (SELECT 1 FROM orders WHERE user_id = users.id))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustExists() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustExists() with empty subquery should panic")
		}
	}()
	MustExists([]string{}, "")
}

func TestNotExists(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT 1 FROM orders WHERE user_id = users.id"
	result := NotExists(conditions, subquery)
	expected := []string{"(NOT EXISTS (SELECT 1 FROM orders WHERE user_id = users.id))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("NotExists() = %v, want %v", result, expected)
	}
}

func TestMustNotExists(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT 1 FROM orders WHERE user_id = users.id"
	result := MustNotExists(conditions, subquery)
	expected := []string{"(NOT EXISTS (SELECT 1 FROM orders WHERE user_id = users.id))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNotExists() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNotExists() with empty subquery should panic")
		}
	}()
	MustNotExists([]string{}, "")
}

func TestEqAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := EqAll(conditions, "amount", subquery)
	expected := []string{"(amount = ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("EqAll() = %v, want %v", result, expected)
	}
}

func TestMustEqAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustEqAll(conditions, "amount", subquery)
	expected := []string{"(amount = ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustEqAll() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustEqAll() with empty subquery should panic")
		}
	}()
	MustEqAll([]string{}, "amount", "")
}

func TestEqAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := EqAny(conditions, "amount", subquery)
	expected := []string{"(amount = ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("EqAny() = %v, want %v", result, expected)
	}
}

func TestMustEqAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustEqAny(conditions, "amount", subquery)
	expected := []string{"(amount = ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustEqAny() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustEqAny() with empty subquery should panic")
		}
	}()
	MustEqAny([]string{}, "amount", "")
}

func TestNeAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := NeAll(conditions, "amount", subquery)
	expected := []string{"(amount <> ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("NeAll() = %v, want %v", result, expected)
	}
}

func TestMustNeAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustNeAll(conditions, "amount", subquery)
	expected := []string{"(amount <> ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNeAll() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNeAll() with empty subquery should panic")
		}
	}()
	MustNeAll([]string{}, "amount", "")
}

func TestNeAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := NeAny(conditions, "amount", subquery)
	expected := []string{"(amount <> ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("NeAny() = %v, want %v", result, expected)
	}
}

func TestMustNeAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustNeAny(conditions, "amount", subquery)
	expected := []string{"(amount <> ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustNeAny() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNeAny() with empty subquery should panic")
		}
	}()
	MustNeAny([]string{}, "amount", "")
}

func TestGtAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := GtAll(conditions, "amount", subquery)
	expected := []string{"(amount > ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("GtAll() = %v, want %v", result, expected)
	}
}

func TestMustGtAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustGtAll(conditions, "amount", subquery)
	expected := []string{"(amount > ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustGtAll() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGtAll() with empty subquery should panic")
		}
	}()
	MustGtAll([]string{}, "amount", "")
}

func TestGtAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := GtAny(conditions, "amount", subquery)
	expected := []string{"(amount > ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("GtAny() = %v, want %v", result, expected)
	}
}

func TestMustGtAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustGtAny(conditions, "amount", subquery)
	expected := []string{"(amount > ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustGtAny() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGtAny() with empty subquery should panic")
		}
	}()
	MustGtAny([]string{}, "amount", "")
}

func TestLtAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := LtAll(conditions, "amount", subquery)
	expected := []string{"(amount < ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("LtAll() = %v, want %v", result, expected)
	}
}

func TestMustLtAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustLtAll(conditions, "amount", subquery)
	expected := []string{"(amount < ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLtAll() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLtAll() with empty subquery should panic")
		}
	}()
	MustLtAll([]string{}, "amount", "")
}

func TestLtAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := LtAny(conditions, "amount", subquery)
	expected := []string{"(amount < ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("LtAny() = %v, want %v", result, expected)
	}
}

func TestMustLtAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustLtAny(conditions, "amount", subquery)
	expected := []string{"(amount < ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLtAny() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLtAny() with empty subquery should panic")
		}
	}()
	MustLtAny([]string{}, "amount", "")
}

func TestGeAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := GeAll(conditions, "amount", subquery)
	expected := []string{"(amount >= ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("GeAll() = %v, want %v", result, expected)
	}
}

func TestMustGeAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustGeAll(conditions, "amount", subquery)
	expected := []string{"(amount >= ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustGeAll() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGeAll() with empty subquery should panic")
		}
	}()
	MustGeAll([]string{}, "amount", "")
}

func TestGeAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := GeAny(conditions, "amount", subquery)
	expected := []string{"(amount >= ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("GeAny() = %v, want %v", result, expected)
	}
}

func TestMustGeAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustGeAny(conditions, "amount", subquery)
	expected := []string{"(amount >= ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustGeAny() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGeAny() with empty subquery should panic")
		}
	}()
	MustGeAny([]string{}, "amount", "")
}

func TestLeAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := LeAll(conditions, "amount", subquery)
	expected := []string{"(amount <= ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("LeAll() = %v, want %v", result, expected)
	}
}

func TestMustLeAll(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustLeAll(conditions, "amount", subquery)
	expected := []string{"(amount <= ALL (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLeAll() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLeAll() with empty subquery should panic")
		}
	}()
	MustLeAll([]string{}, "amount", "")
}

func TestLeAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := LeAny(conditions, "amount", subquery)
	expected := []string{"(amount <= ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("LeAny() = %v, want %v", result, expected)
	}
}

func TestMustLeAny(t *testing.T) {
	conditions := []string{}
	subquery := "SELECT amount FROM orders WHERE user_id = 1"
	result := MustLeAny(conditions, "amount", subquery)
	expected := []string{"(amount <= ANY (SELECT amount FROM orders WHERE user_id = 1))"}
	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("MustLeAny() = %v, want %v", result, expected)
	}

	// Test panic on empty subquery
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLeAny() with empty subquery should panic")
		}
	}()
	MustLeAny([]string{}, "amount", "")
}

func TestSelect(t *testing.T) {
	// Test with fields
	result := Select("id", "name", "email")
	expected := "SELECT id, name, email"
	if result != expected {
		t.Errorf("Select() = %v, want %v", result, expected)
	}

	// Test without fields (default to *)
	result = Select()
	expected = "SELECT *"
	if result != expected {
		t.Errorf("Select() without fields = %v, want %v", result, expected)
	}
}

func TestFrom(t *testing.T) {
	result := From("users")
	expected := "FROM users"
	if result != expected {
		t.Errorf("From() = %v, want %v", result, expected)
	}

	// Test panic on empty table
	defer func() {
		if r := recover(); r == nil {
			t.Error("From() with empty table should panic")
		}
	}()
	From("")
}

func TestLeftJoin(t *testing.T) {
	result := LeftJoin("orders", "users.id = orders.user_id")
	expected := "LEFT JOIN orders ON users.id = orders.user_id"
	if result != expected {
		t.Errorf("LeftJoin() = %v, want %v", result, expected)
	}

	// Test with multiple conditions
	result = LeftJoin("orders", "users.id = orders.user_id", "orders.status = 'active'")
	expected = "LEFT JOIN orders ON users.id = orders.user_id AND orders.status = 'active'"
	if result != expected {
		t.Errorf("LeftJoin() with multiple conditions = %v, want %v", result, expected)
	}

	// Test panic on empty table or conditions
	defer func() {
		if r := recover(); r == nil {
			t.Error("LeftJoin() with empty table should panic")
		}
	}()
	LeftJoin("", "users.id = orders.user_id")
}

func TestRightJoin(t *testing.T) {
	result := RightJoin("orders", "users.id = orders.user_id")
	expected := "RIGHT JOIN orders ON users.id = orders.user_id"
	if result != expected {
		t.Errorf("RightJoin() = %v, want %v", result, expected)
	}

	// Test with multiple conditions
	result = RightJoin("orders", "users.id = orders.user_id", "orders.status = 'active'")
	expected = "RIGHT JOIN orders ON users.id = orders.user_id AND orders.status = 'active'"
	if result != expected {
		t.Errorf("RightJoin() with multiple conditions = %v, want %v", result, expected)
	}

	// Test panic on empty table or conditions
	defer func() {
		if r := recover(); r == nil {
			t.Error("RightJoin() with empty table should panic")
		}
	}()
	RightJoin("", "users.id = orders.user_id")
}

func TestGroupBy(t *testing.T) {
	// Test with fields
	result := GroupBy([]string{"department", "status"})
	expected := "GROUP BY department, status"
	if result != expected {
		t.Errorf("GroupBy() = %v, want %v", result, expected)
	}

	// Test with empty fields
	result = GroupBy([]string{})
	if result != "" {
		t.Errorf("GroupBy() with empty fields = %v, want \"\"", result)
	}
}

func TestMustGroupBy(t *testing.T) {
	result := MustGroupBy([]string{"department", "status"})
	expected := "GROUP BY department, status"
	if result != expected {
		t.Errorf("MustGroupBy() = %v, want %v", result, expected)
	}

	// Test panic on empty fields
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGroupBy() with empty fields should panic")
		}
	}()
	MustGroupBy([]string{})
}

func TestOrderBy(t *testing.T) {
	// Test with fields
	result := OrderBy([]string{"created_at DESC", "name ASC"})
	expected := "ORDER BY created_at DESC, name ASC"
	if result != expected {
		t.Errorf("OrderBy() = %v, want %v", result, expected)
	}

	// Test with empty fields
	result = OrderBy([]string{})
	if result != "" {
		t.Errorf("OrderBy() with empty fields = %v, want \"\"", result)
	}
}

func TestMustOrderBy(t *testing.T) {
	result := MustOrderBy([]string{"created_at DESC", "name ASC"})
	expected := "ORDER BY created_at DESC, name ASC"
	if result != expected {
		t.Errorf("MustOrderBy() = %v, want %v", result, expected)
	}

	// Test panic on empty fields
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustOrderBy() with empty fields should panic")
		}
	}()
	MustOrderBy([]string{})
}

func TestLimit(t *testing.T) {
	// Test with positive number
	result := Limit(10)
	expected := "LIMIT 10"
	if result != expected {
		t.Errorf("Limit() = %v, want %v", result, expected)
	}

	// Test with negative number
	result = Limit(-1)
	if result != "" {
		t.Errorf("Limit() with negative number = %v, want \"\"", result)
	}
}

func TestMustLimit(t *testing.T) {
	result := MustLimit(10)
	expected := "LIMIT 10"
	if result != expected {
		t.Errorf("MustLimit() = %v, want %v", result, expected)
	}

	// Test panic on negative number
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLimit() with negative number should panic")
		}
	}()
	MustLimit(-1)
}

func TestOffset(t *testing.T) {
	// Test with positive number
	result := Offset(20)
	expected := "OFFSET 20"
	if result != expected {
		t.Errorf("Offset() = %v, want %v", result, expected)
	}

	// Test with negative number
	result = Offset(-1)
	if result != "" {
		t.Errorf("Offset() with negative number = %v, want \"\"", result)
	}
}

func TestMustOffset(t *testing.T) {
	result := MustOffset(20)
	expected := "OFFSET 20"
	if result != expected {
		t.Errorf("MustOffset() = %v, want %v", result, expected)
	}

	// Test panic on negative number
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustOffset() with negative number should panic")
		}
	}()
	MustOffset(-1)
}
