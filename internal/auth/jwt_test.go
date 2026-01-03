package auth

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenGenerator(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()

	assert.NotNil(t, tg)
	assert.NotNil(t, tg.cache)
	assert.Equal(t, MaxCacheSize, tg.maxSize)
	assert.Equal(t, CacheTTLSeconds*time.Second, tg.cacheTTL)
	assert.Equal(t, APITokenTTLSeconds*time.Second, tg.tokenTTL)
	assert.False(t, tg.disableCache)
}

func TestNewTokenGeneratorWithConfig(t *testing.T) {
	t.Parallel()

	maxSize := 5
	cacheTTL := 60 * time.Second
	tokenTTL := 90 * time.Second

	tg := NewTokenGeneratorWithConfig(maxSize, cacheTTL, tokenTTL)

	assert.NotNil(t, tg)
	assert.Equal(t, maxSize, tg.maxSize)
	assert.Equal(t, cacheTTL, tg.cacheTTL)
	assert.Equal(t, tokenTTL, tg.tokenTTL)
}

func TestTokenGenerator_GenerateToken_ValidAPIKey(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "test-api-key.test-secret"

	token, err := tg.GenerateToken(apiKey)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token structure
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3, "JWT token should have 3 parts")
}

func TestTokenGenerator_GenerateToken_InvalidAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		apiKey string
		errMsg error
	}{
		{
			name:   "empty API key",
			apiKey: "",
			errMsg: ErrEmptyAPIKey,
		},
		{
			name:   "no separator",
			apiKey: "invalidkey",
			errMsg: ErrInvalidAPIKey,
		},
		{
			name:   "multiple separators",
			apiKey: "key.secret.extra",
			errMsg: ErrInvalidAPIKey,
		},
		{
			name:   "empty key part",
			apiKey: ".secret",
			errMsg: ErrInvalidAPIKey,
		},
		{
			name:   "empty secret part",
			apiKey: "key.",
			errMsg: ErrInvalidAPIKey,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tg := NewTokenGenerator()
			token, err := tg.GenerateToken(tt.apiKey)

			assert.Error(t, err)
			assert.Empty(t, token)
			assert.ErrorIs(t, err, tt.errMsg)
		})
	}
}

func TestTokenGenerator_VerifyToken(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "12345678.abcdefg"

	// Generate token
	token, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token
	claims, err := VerifyToken(token, "abcdefg")
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Check claims
	assert.Equal(t, "12345678", claims.APIKey)
	assert.NotZero(t, claims.Timestamp)
	assert.NotNil(t, claims.ExpiresAt)

	// Check expiration is in the future
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))

	// Check token validity period (should be ~210 seconds)
	expectedExpiration := time.UnixMilli(claims.Timestamp).Add(APITokenTTLSeconds * time.Second)
	assert.WithinDuration(t, expectedExpiration, claims.ExpiresAt.Time, time.Second)
}

func TestTokenGenerator_VerifyToken_CustomHeaders(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "test-key.test-secret"

	// Generate token
	tokenString, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)

	// Parse token without verification to check headers
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, &Claims{})
	require.NoError(t, err)

	// Check custom headers
	assert.Equal(t, "HS256", token.Header["alg"])
	assert.Equal(t, "SIGN", token.Header["sign_type"])
}

func TestTokenGenerator_Caching(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "cache-test.secret123"

	// Generate token first time
	token1, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token1)

	// Generate token second time (should be cached)
	token2, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Tokens should be identical (from cache)
	assert.Equal(t, token1, token2)

	// Verify cache size
	assert.Equal(t, 1, tg.GetCacheSize())
}

func TestTokenGenerator_CacheExpiration(t *testing.T) {
	// Create generator with short cache TTL for testing
	tg := NewTokenGeneratorWithConfig(10, 100*time.Millisecond, 500*time.Millisecond)
	apiKey := "expire-test.secret456"

	// Generate token
	token1, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token1)

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Generate token again (should be new token)
	token2, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Tokens should be different (cache expired)
	assert.NotEqual(t, token1, token2)
}

func TestTokenGenerator_DisableCache(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "disable-cache.secret789"

	// Disable cache
	tg.DisableCache()

	// Generate token first time
	token1, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token1)

	// Sleep briefly to ensure different timestamp
	time.Sleep(10 * time.Millisecond)

	// Generate token second time (should NOT be cached)
	token2, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Tokens should be different (caching disabled, different timestamps)
	assert.NotEqual(t, token1, token2)

	// Cache should be empty
	assert.Equal(t, 0, tg.GetCacheSize())
}

func TestTokenGenerator_EnableCache(t *testing.T) {
	tg := NewTokenGenerator()
	apiKey := "enable-cache.secret-abc"

	// Disable cache first
	tg.DisableCache()

	// Generate token (not cached)
	token1, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token1)
	assert.Equal(t, 0, tg.GetCacheSize())

	// Enable cache
	tg.EnableCache()

	// Generate token (should be cached)
	token2, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.NotEmpty(t, token2)
	assert.Equal(t, 1, tg.GetCacheSize())

	// Generate again (should return cached)
	token3, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)
	assert.Equal(t, token2, token3)
}

func TestTokenGenerator_CacheEviction(t *testing.T) {
	// Create generator with small cache size
	tg := NewTokenGeneratorWithConfig(3, 60*time.Second, 90*time.Second)

	// Generate tokens for different API keys
	for i := 1; i <= 4; i++ {
		apiKey := "key" + string(rune('0'+i)) + ".secret"
		_, err := tg.GenerateToken(apiKey)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Cache size should not exceed max size
	assert.Equal(t, 3, tg.GetCacheSize())
}

func TestTokenGenerator_ClearCache(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()

	// Generate multiple tokens
	for i := 1; i <= 5; i++ {
		apiKey := "clear-key" + string(rune('0'+i)) + ".secret"
		_, err := tg.GenerateToken(apiKey)
		require.NoError(t, err)
	}

	assert.Equal(t, 5, tg.GetCacheSize())

	// Clear cache
	tg.ClearCache()

	assert.Equal(t, 0, tg.GetCacheSize())
}

func TestTokenGenerator_ClearExpiredTokens(t *testing.T) {
	// Create generator with short cache TTL
	tg := NewTokenGeneratorWithConfig(10, 100*time.Millisecond, 500*time.Millisecond)

	// Generate some tokens
	for i := 1; i <= 3; i++ {
		apiKey := "expire-key" + string(rune('0'+i)) + ".secret"
		_, err := tg.GenerateToken(apiKey)
		require.NoError(t, err)
	}

	assert.Equal(t, 3, tg.GetCacheSize())

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Clear expired tokens
	tg.ClearExpiredTokens()

	assert.Equal(t, 0, tg.GetCacheSize())
}

func TestTokenGenerator_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "concurrent.secret-test"

	var wg sync.WaitGroup
	numGoroutines := 10

	// Generate tokens concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := tg.GenerateToken(apiKey)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()

	// Should have exactly 1 cached token (all goroutines used same key)
	assert.Equal(t, 1, tg.GetCacheSize())
}

func TestTokenGenerator_MultipleConcurrentKeys(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()

	var wg sync.WaitGroup
	numGoroutines := 20

	// Generate tokens for different keys concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// Use modulo to have some key collisions
			apiKey := "key" + string(rune('0'+(index%5))) + ".secret"
			_, err := tg.GenerateToken(apiKey)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Should have 5 cached tokens (keys 0-4)
	assert.Equal(t, 5, tg.GetCacheSize())
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		token  string
		secret string
	}{
		{
			name:   "malformed token",
			token:  "invalid.token.format",
			secret: "secret",
		},
		{
			name:   "empty token",
			token:  "",
			secret: "secret",
		},
		{
			name:   "wrong secret",
			token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			secret: "wrong-secret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			claims, err := VerifyToken(tt.token, tt.secret)

			assert.Error(t, err)
			assert.Nil(t, claims)
		})
	}
}

func TestTokenGenerator_TokenTimestamp(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "timestamp-test.secret-xyz"

	// Get current time truncated to milliseconds for fair comparison
	beforeGeneration := time.Now().Truncate(time.Millisecond)

	// Generate token
	token, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)

	afterGeneration := time.Now().Truncate(time.Millisecond)

	// Verify token
	claims, err := VerifyToken(token, "secret-xyz")
	require.NoError(t, err)

	// Timestamp should be within a reasonable range of generation time
	tokenTime := time.UnixMilli(claims.Timestamp)

	// Token time should be between before and after generation (with tolerance)
	assert.True(t, !tokenTime.Before(beforeGeneration),
		"token timestamp should not be before generation started")
	assert.True(t, !tokenTime.After(afterGeneration.Add(time.Millisecond)),
		"token timestamp should not be significantly after generation ended")
}

func TestClaims_ExpirationInMilliseconds(t *testing.T) {
	t.Parallel()

	tg := NewTokenGenerator()
	apiKey := "exp-test.secret-def"

	// Generate token
	token, err := tg.GenerateToken(apiKey)
	require.NoError(t, err)

	// Verify token
	claims, err := VerifyToken(token, "secret-def")
	require.NoError(t, err)

	// Calculate expected expiration
	expectedExp := time.UnixMilli(claims.Timestamp).Add(APITokenTTLSeconds * time.Second)

	// Verify expiration is correct (within 1 second tolerance)
	assert.WithinDuration(t, expectedExp, claims.ExpiresAt.Time, time.Second)

	// Verify token is valid for approximately 210 seconds
	duration := claims.ExpiresAt.Time.Sub(time.UnixMilli(claims.Timestamp))
	assert.InDelta(t, APITokenTTLSeconds, duration.Seconds(), 1.0)
}
