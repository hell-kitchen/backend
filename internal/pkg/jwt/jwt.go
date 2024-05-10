package jwt

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/config"
	"go.uber.org/zap"
	"os"
	"time"
)

type Provider struct {
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	accessLifetime  int
	refreshLifetime int
}

func NewProvider(cfg *config.JWT, log *zap.Logger) (*Provider, error) {
	privateKeyRaw, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading private key file")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("error while reading private key file")
	}

	publicKeyRaw, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error while reading public key file")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyRaw)
	if err != nil {
		return nil, fmt.Errorf("error while reading public key file")
	}

	log.Info("public key", zap.ByteString("publicKey", publicKeyRaw))

	provider := &Provider{
		publicKey:       publicKey,
		privateKey:      privateKey,
		accessLifetime:  cfg.AccessTokenLifetime,
		refreshLifetime: cfg.RefreshTokenLifetime,
	}

	return provider, nil
}

func (provider *Provider) readKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return provider.publicKey, nil
}

func (provider *Provider) GetDataFromToken(token string) (uuid.UUID, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, provider.readKeyFunc)
	if err != nil {
		return uuid.Nil, err
	}

	if !parsed.Valid {
		return uuid.Nil, fmt.Errorf("invalid token: not valid")
	}

	if claims, ok := parsed.Claims.(*jwt.StandardClaims); ok {
		return uuid.Parse(claims.Issuer)
	}
	return uuid.Nil, fmt.Errorf("invalid token: cannot parse claims")
}

func (provider *Provider) CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error) {
	now := time.Now()

	var add time.Duration
	if isAccess {
		add = time.Duration(provider.accessLifetime) * time.Minute
	} else {
		add = time.Duration(provider.refreshLifetime) * time.Minute
	}

	claims := &jwt.StandardClaims{
		Issuer:    userID.String(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(add).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(provider.privateKey)
}
