# Go Comment Standards

**Project**: gox
**Standard Package**: `backoff/`
**Language**: English only

---

## Function Comments

```go
// FunctionName describes what the function does.
// Additional details if needed (multiple lines allowed).
//
// Parameters:
//   - param1: Description of parameter 1
//   - param2: Description of parameter 2
//
// Returns:
//   - ReturnType: Description of return value
//   - error: Error description if applicable
func FunctionName(param1 Type, param2 Type) (ReturnType, error) {
```

### Rules
- Start with the function name followed by a verb
- Blank line before `Parameters:` section
- Use bullet list with `-` for each parameter
- Use bullet list with `-` for each return value
- Include `Parameters:` even if no params (use `// Parameters: //   - None`)
- Include `Returns:` even if no returns (use `// Returns: //   - None`)

---

## Type Comments

```go
// TypeName is a ... (description of what the type represents).
// Additional details about usage or purpose.
type TypeName struct {
```

### Rules
- Start with the type name followed by "is a"
- Describe what the type represents
- Add usage notes if non-obvious

---

## Method Comments

```go
// MethodName performs ... (what the method does).
//
// Parameters:
//   - param: Description
//
// Returns:
//   - *ReceiverType: The receiver for method chaining (if applicable)
func (r *ReceiverType) MethodName(param Type) *ReceiverType {
```

### Rules
- Same as function comments
- Start with the method name
- Document receiver mutations if any

---

## Package Comments

```go
// Package packagex provides ... (one-line summary).
//
// Additional details about the package's purpose and usage.
package packagex
```

### Rules
- Must start with "Package packagex" (matching package name)
- One-line summary is required
- Additional paragraphs allowed after blank line

---

## Interface Comments

```go
// InterfaceName is a ... (description).
//
// Methods:
//   - Method1: description
//   - Method2: description
type InterfaceName interface {
```

### Rules
- Start with the interface name followed by "is a"
- Document each method if not self-explanatory

---

## Error Variable Comments

```go
// ErrSomething is returned when ... (condition).
var ErrSomething = errors.New("packagex: something failed")
```

### Rules
- Start with the variable name
- Describe when the error is returned
- Use package prefix in error message

---

## Constant Comments

```go
// ConstantName represents ... (meaning and usage).
const ConstantName = value
```

### Rules
- Start with the constant name
- Describe what the constant represents

---

## Generic Type/Function Comments

```go
// FunctionName does something with type T.
//
// Type Parameters:
//   - T: Description of the type parameter
//
// Parameters:
//   - param: Description
//
// Returns:
//   - T: Description
func FunctionName[T any](param T) T {
```

### Rules
- Add `Type Parameters:` section before `Parameters:`
- Document each type parameter

---

## Examples

### Simple Function

```go
// Constant returns a backoff function that waits for a fixed period of time between calls.
// This is useful when you want consistent retry intervals regardless of the attempt number.
// Returns 0 when attempt is 0.
//
// Parameters:
//   - delta: The fixed duration to wait between each retry attempt
//
// Returns:
//   - Func: A backoff function that always returns the same duration (or 0 for attempt 0)
func Constant(delta time.Duration) Func {
```

### Method with Receiver

```go
// Submit submits an object to the Group.
//
// Parameters:
//   - obj: the object to submit
//
// Returns:
//   - error: nil on success, ErrClosed if the Group is closed
func (g *Group[Obj]) Submit(obj Obj) error {
```

### Interface

```go
// Cipher is the interface for AES encryption/decryption with different modes.
//
// Methods:
//   - Encrypt: encrypts plaintext using the configured mode
//   - Decrypt: decrypts ciphertext using the configured mode
type Cipher interface {
```

---

## Checklist

Before submitting code:

- [ ] All exported functions have doc comments
- [ ] All exported types have doc comments
- [ ] All exported methods have doc comments
- [ ] All exported variables/constants have doc comments
- [ ] Comments are in English
- [ ] Comments follow the format above
- [ ] Parameters are documented
- [ ] Return values are documented
- [ ] Chinese comments have been translated
