// Package auth provides JWT token generation for Z.ai SDK authentication.
package auth

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// CacheTTLSeconds is the duration for which tokens are cached.
	CacheTTLSeconds = 3 * 60 // 180 seconds (3 minutes)

	// APITokenTTLSeconds is the duration for which tokens are valid.
	// This is 30 seconds longer than cache TTL to ensure fresh tokens.
	APITokenTTLSeconds = CacheTTLSeconds + 30 // 210 seconds (3.5 minutes)

	// MaxCacheSize is the maximum number of tokens to cache.
	MaxCacheSize = 10
)

var (
	// ErrInvalidAPIKey is returned when the API key format is invalid.
	ErrInvalidAPIKey = errors.New("invalid API key format: expected 'key.secret'")

	// ErrEmptyAPIKey is returned when the API key is empty.
	ErrEmptyAPIKey = errors.New("API key cannot be empty")
)

// Claims represents the JWT claims for Z.ai authentication.
type Claims struct {
	APIKey    string `json:"api_key"`
	Timestamp int64  `json:"timestamp"`
	jwt.RegisteredClaims
}

// TokenCache represents a cached token with its creation time.
type TokenCache struct {
	Token     string
	CreatedAt time.Time
}

// TokenGenerator generates and caches JWT tokens for authentication.
type TokenGenerator struct {
	cache     map[string]*TokenCache
	cacheMu   sync.RWMutex
	maxSize   int
	cacheTTL  time.Duration
	tokenTTL  time.Duration
	disableCache bool
}

// NewTokenGenerator creates a new token generator with default settings.
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{
		cache:     make(map[string]*TokenCache),
		maxSize:   MaxCacheSize,
		cacheTTL:  CacheTTLSeconds * time.Second,
		tokenTTL:  APITokenTTLSeconds * time.Second,
		disableCache: false,
	}
}

// NewTokenGeneratorWithConfig creates a new token generator with custom configuration.
func NewTokenGeneratorWithConfig(maxSize int, cacheTTL, tokenTTL time.Duration) *TokenGenerator {
	return &TokenGenerator{
		cache:     make(map[string]*TokenCache),
		maxSize:   maxSize,
		cacheTTL:  cacheTTL,
		tokenTTL:  tokenTTL,
		disableCache: false,
	}
}

// DisableCache disables token caching.
func (tg *TokenGenerator) DisableCache() {
	tg.disableCache = true
}

// EnableCache enables token caching.
func (tg *TokenGenerator) EnableCache() {
	tg.disableCache = false
}

// GenerateToken generates a JWT token from the API key.
// The API key should be in the format "apikey.secret".
// If caching is enabled, it returns a cached token if available and not expired.
func (tg *TokenGenerator) GenerateToken(apiKey string) (string, error) {
	if apiKey == "" {
		return "", ErrEmptyAPIKey
	}

	// Check cache first if caching is enabled
	if !tg.disableCache {
		if cachedToken := tg.getCachedToken(apiKey); cachedToken != "" {
			return cachedToken, nil
		}
	}

	// Generate new token
	token, err := tg.generateToken(apiKey)
	if err != nil {
		return "", err
	}

	// Cache the token if caching is enabled
	if !tg.disableCache {
		tg.cacheToken(apiKey, token)
	}

	return token, nil
}

// generateToken creates a new JWT token.
func (tg *TokenGenerator) generateToken(apiKey string) (string, error) {
	// Split API key into key and secret
	parts := strings.Split(apiKey, ".")
	if len(parts) != 2 {
		return "", ErrInvalidAPIKey
	}

	key := parts[0]
	secret := parts[1]

	if key == "" || secret == "" {
		return "", ErrInvalidAPIKey
	}

	// Get current time in milliseconds
	now := time.Now()
	timestampMs := now.UnixMilli()
	expirationMs := timestampMs + (tg.tokenTTL.Milliseconds())

	// Create claims
	claims := &Claims{
		APIKey:    key,
		Timestamp: timestampMs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.UnixMilli(expirationMs)),
		},
	}

	// Create token with custom headers
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"

	// Sign the token
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// getCachedToken retrieves a cached token if it exists and is not expired.
func (tg *TokenGenerator) getCachedToken(apiKey string) string {
	tg.cacheMu.RLock()
	defer tg.cacheMu.RUnlock()

	cached, exists := tg.cache[apiKey]
	if !exists {
		return ""
	}

	// Check if cache entry is expired
	if time.Since(cached.CreatedAt) > tg.cacheTTL {
		return ""
	}

	return cached.Token
}

// cacheToken stores a token in the cache.
func (tg *TokenGenerator) cacheToken(apiKey, token string) {
	tg.cacheMu.Lock()
	defer tg.cacheMu.Unlock()

	// Evict old entries if cache is full
	if len(tg.cache) >= tg.maxSize {
		tg.evictOldest()
	}

	tg.cache[apiKey] = &TokenCache{
		Token:     token,
		CreatedAt: time.Now(),
	}
}

// evictOldest removes the oldest cache entry.
func (tg *TokenGenerator) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range tg.cache {
		if oldestKey == "" || cached.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.CreatedAt
		}
	}

	if oldestKey != "" {
		delete(tg.cache, oldestKey)
	}
}

// ClearCache removes all cached tokens.
func (tg *TokenGenerator) ClearCache() {
	tg.cacheMu.Lock()
	defer tg.cacheMu.Unlock()

	tg.cache = make(map[string]*TokenCache)
}

// ClearExpiredTokens removes expired tokens from the cache.
func (tg *TokenGenerator) ClearExpiredTokens() {
	tg.cacheMu.Lock()
	defer tg.cacheMu.Unlock()

	now := time.Now()
	for key, cached := range tg.cache {
		if now.Sub(cached.CreatedAt) > tg.cacheTTL {
			delete(tg.cache, key)
		}
	}
}

// GetCacheSize returns the current number of cached tokens.
func (tg *TokenGenerator) GetCacheSize() int {
	tg.cacheMu.RLock()
	defer tg.cacheMu.RUnlock()

	return len(tg.cache)
}

// VerifyToken verifies a JWT token and returns its claims.
// This is primarily used for testing.
func VerifyToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
