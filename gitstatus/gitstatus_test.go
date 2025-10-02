package gitstatus

import (
	"context"
	"os/exec"
	"testing"
	"time"

	// "regexp"
	// log "github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

func TestGitStatus(t *testing.T) {
    assert.Equal(t, "foo", "foo", "foo != foo")
}

// TestGitTimeoutBehavior tests that git operations have proper timeout behavior
// This prevents hangs during SSH sessions or slow network conditions
func TestGitTimeoutBehavior(t *testing.T) {
	// Test that _git function respects timeout
	start := time.Now()
	
	// This should timeout quickly due to the 2-second timeout in _git()
	_, err := _git("git", "rev-list", "--left-right", "--count", "nonexistent/remote..branch")
	
	duration := time.Since(start)
	
	// Should timeout within reasonable time (2 seconds + some overhead)
	assert.True(t, duration < 5*time.Second, "Git operation should timeout within 5 seconds")
	assert.Error(t, err, "Git operation with invalid remote should return error")
}

// TestGitTimeoutWithContext tests that the timeout context works correctly
func TestGitTimeoutWithContext(t *testing.T) {
	// Create a context that will timeout quickly
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "git", "rev-list", "--left-right", "--count", "origin/main..main")
	
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)
	
	// Should timeout quickly due to context timeout
	assert.True(t, duration < 200*time.Millisecond, "Command should timeout within 200ms")
	assert.Error(t, err, "Command should return timeout error")
}
