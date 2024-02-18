package limit

import (
	"math"
	"testing"
	"time"
)

func TestTokenRateLimit_Consume(t *testing.T) {
	limiter := NewTokenRateLimit(10, 20, 5)

	// Ensure initial tokens are correct
	if limiter.currentTokens != 10 {
		t.Errorf("Expected initial tokens to be 10, got %f", limiter.currentTokens)
	}

	// Consume 5 tokens
	limiter.Consume(func(currentTokens float64) float64 {
		return 5
	})

	// Ensure tokens are consumed correctly
	if math.Round(limiter.currentTokens) != 5 {
		t.Errorf("Expected tokens to be 5, got %f", limiter.currentTokens)
	}

	// Consume 10 tokens
	limiter.Consume(func(currentTokens float64) float64 {
		return 10
	})

	// Ensure tokens are consumed correctly and not below 0
	if limiter.currentTokens >= 0 {
		t.Errorf("Expected tokens to be 0, got %f", limiter.currentTokens)
	}

	// Consume 5 tokens, but only 2 tokens are available
	limiter.Consume(func(currentTokens float64) float64 {
		return math.Min(currentTokens, 5)
	})

	// Ensure tokens are consumed correctly and not below 0
	if limiter.currentTokens != 0 {
		t.Errorf("Expected tokens to be 0, got %f", limiter.currentTokens)
	}
}

func TestTokenRateLimit_ConsumeWithWait(t *testing.T) {
	limiter := NewTokenRateLimit(5, 10, 2)

	startTime := time.Now()

	// Consume 2 tokens immediately
	limiter.ConsumeWithWait(func(currentTokens float64) float64 {
		return 2
	})

	// Ensure tokens are consumed correctly
	if math.Round(limiter.currentTokens) != 3 {
		t.Errorf("Expected tokens to be 3, got %f", limiter.currentTokens)
	}

	// Start a goroutine to consume 5 tokens after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		limiter.ConsumeWithWait(func(currentTokens float64) float64 {
			return 5
		})
	}()

	// Consume 3 tokens immediately, should block until the goroutine adds more tokens
	limiter.ConsumeWithWait(func(currentTokens float64) float64 {
		return 3
	})

	// Ensure tokens are consumed correctly and the total waiting time is around 2 seconds
	elapsedTime := time.Since(startTime)
	if elapsedTime < 2*time.Second {
		t.Errorf("Expected waiting time to be around 2 seconds, got %s", elapsedTime)
	}
	if limiter.currentTokens != 0 {
		t.Errorf("Expected tokens to be 0, got %f", limiter.currentTokens)
	}
}

func TestTokenRateLimit_RefillTokens(t *testing.T) {
	lastRefillTime := time.Now().Add(-5 * time.Second)
	fillRate := 2.5

	refilledTokens := calculateRefilledTokens(lastRefillTime, fillRate)

	// Ensure refilled tokens are calculated correctly
	expectedTokens := fillRate * 5
	if refilledTokens != expectedTokens {
		t.Errorf("Expected refilled tokens to be %f, got %f", expectedTokens, refilledTokens)
	}
}

func TestTokenRateLimit_InvalidParameters(t *testing.T) {
	// Initialize with invalid parameters
	limiter := NewTokenRateLimit(20, 10, -5)

	// Ensure initial tokens are set correctly
	if limiter.currentTokens != 20 {
		t.Errorf("Expected initial tokens to be 20, got %f", limiter.currentTokens)
	}

	// Consume tokens with invalid handler
	limiter.Consume(func(currentTokens float64) float64 {
		return -10
	})

	// Ensure tokens are not consumed
	if limiter.currentTokens != 20 {
		t.Errorf("Expected tokens to be 20, got %f", limiter.currentTokens)
	}
}

func TestTokenRateLimit_FillRateUpdate(t *testing.T) {
	limiter := NewTokenRateLimit(10, 20, 5)

	// Update fill rate to 10 tokens per second
	limiter.fillRate = 10

	// Consume 5 tokens
	limiter.Consume(func(currentTokens float64) float64 {
		return 5
	})

	// Ensure tokens are consumed correctly based on updated fill rate
	if limiter.currentTokens != 5 {
		t.Errorf("Expected tokens to be 5, got %f", limiter.currentTokens)
	}
}
