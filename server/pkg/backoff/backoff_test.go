package backoff

import (
	"testing"
	"time"
)

func TestBasicBackoff(t *testing.T) {
	min := 50 * time.Millisecond
	max := 5 * time.Second
	backoff := New(
		Min(min),
		Max(max),
	)
	for retry := 10; retry > 0; retry-- {
		d := backoff.Duration()
		if d < min {
			t.Fatalf("min should be >= %v", min)
		}
		if d > max {
			t.Fatalf("max should be <= %v", max)
		}
		t.Logf("retry: %d, duration: %v\n", retry, d)
	}
}

func TestReverseMinMaxBackoff(t *testing.T) {
	// reverse min and max
	min := 5 * time.Second
	max := 50 * time.Millisecond

	backoff := New(
		Min(min),
		Max(max),
	)
	for retry := 10; retry > 0; retry-- {
		d := backoff.Duration()
		if d < max {
			t.Fatalf("min should be >= %v", max)
		}
		if d > min {
			t.Fatalf("max should be <= %v", min)
		}
		t.Logf("retry: %d, duration: %v\n", retry, d)
	}
}

func TestCustomFactorBackoff(t *testing.T) {
	min := 50 * time.Millisecond
	max := 5 * time.Second
	backoff := New(
		Min(min),
		Max(max),
		Factor(3),
	)
	for retry := 10; retry > 0; retry-- {
		d := backoff.Duration()
		if d < min {
			t.Fatalf("min should be >= %v", min)
		}
		if d > max {
			t.Fatalf("max should be <= %v", max)
		}
		t.Logf("retry: %d, duration: %v\n", retry, d)
	}
}
