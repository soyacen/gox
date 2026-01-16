package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/soyacen/gox/backoff"
)

func TestCall(t *testing.T) {
	maxAttempts := 3
	ctx := context.Background()
	method := func(attemptTime int) error {
		fmt.Println(attemptTime)
		if attemptTime < maxAttempts {
			return errors.New("mock error")
		}
		return nil
	}
	backoffFunc := backoff.Constant(time.Second)
	err := Call(ctx, uint(maxAttempts), backoffFunc, method)
	assert.Nil(t, err)
}


func TestRetry(t *testing.T) {
	err := MaxAttempts(3).Backoff(backoff.Constant(time.Second)).Exec(context.Background(), func(ctx context.Context, attempts uint) error {
		fmt.Println(attempts)
		if attempts < 3 {
			return errors.New("mock error")
		}
		return nil
	})
	assert.Nil(t, err)
}
