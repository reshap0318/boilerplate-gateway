package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/reshap0318/serv-uam/internal/helpers"
)

// JWKSAuth validates a Bearer token locally by fetching the api-gateway's
// public keys from GATEWAY_JWKS_URL (the same JWKS the gateway exposes for
// its own token verification) and checking the RS256 signature + expiry.
//
// Not wired into any route yet — GatewayAuth (trusting the gateway's
// identity headers) is what's actually used today. This is kept ready for
// a future need: verifying a caller's token directly in this service,
// e.g. a service-to-service call that doesn't go through the gateway.
func JWKSAuth() gin.HandlerFunc {
	cache := newJWKSCache(helpers.GetEnv("GATEWAY_JWKS_URL", ""), 10*time.Minute)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims := &helpers.JWTClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			kid, _ := token.Header["kid"].(string)
			return cache.getKey(kid)
		})
		if err != nil || !token.Valid {
			helpers.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

// jwksCache fetches and caches the gateway's public keys so tokens can be
// verified locally without calling the gateway on every request.
type jwksCache struct {
	url string
	ttl time.Duration

	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
}

func newJWKSCache(url string, ttl time.Duration) *jwksCache {
	return &jwksCache{url: url, ttl: ttl, keys: make(map[string]*rsa.PublicKey)}
}

func (j *jwksCache) getKey(kid string) (*rsa.PublicKey, error) {
	j.mu.RLock()
	key, ok := j.keys[kid]
	fresh := time.Since(j.fetchedAt) < j.ttl
	j.mu.RUnlock()

	if ok && fresh {
		return key, nil
	}

	if err := j.refresh(); err != nil {
		return nil, err
	}

	j.mu.RLock()
	defer j.mu.RUnlock()
	key, ok = j.keys[kid]
	if !ok {
		return nil, fmt.Errorf("unknown key id: %s", kid)
	}
	return key, nil
}

func (j *jwksCache) refresh() error {
	if j.url == "" {
		return fmt.Errorf("GATEWAY_JWKS_URL is not configured")
	}

	resp, err := http.Get(j.url)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var body jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	keys := make(map[string]*rsa.PublicKey, len(body.Keys))
	for _, k := range body.Keys {
		pubKey, err := jwkToRSAPublicKey(k)
		if err != nil {
			continue
		}
		keys[k.Kid] = pubKey
	}

	j.mu.Lock()
	j.keys = keys
	j.fetchedAt = time.Now()
	j.mu.Unlock()

	return nil
}

func jwkToRSAPublicKey(k jwk) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, err
	}

	e := 0
	for _, b := range eBytes {
		e = e<<8 | int(b)
	}

	return &rsa.PublicKey{N: new(big.Int).SetBytes(nBytes), E: e}, nil
}
