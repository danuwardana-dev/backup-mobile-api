package middleware

import (
	"backend-mobile-api/app/config"
	"backend-mobile-api/helpers"
	"backend-mobile-api/internal/repository/redis"
	"backend-mobile-api/model/entity"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"time"
)

type customMiddleware struct {
	logger     *helpers.CustomLogger
	jwtConfig  *config.Jwt
	Redis      redis.Redis
	rootConfig *config.Root
}

func NewCustomMiddleware(
	jwtConfig *config.Jwt,
	logger *helpers.CustomLogger,
	Redis redis.Redis,
	rootConfig *config.Root,
) CustomMiddleware {
	return &customMiddleware{
		jwtConfig:  jwtConfig,
		logger:     logger,
		Redis:      Redis,
		rootConfig: rootConfig,
	}
}

type CustomMiddleware interface {
	CreateTokens(ctx context.Context, user *Claims) (TokenData, error)
	generateToken(ctx context.Context, user *Claims, key string, expiration time.Duration) (string, error)
	ParseJwtToken(context.Context, string, *rsa.PublicKey) (*jwt.Token, error)
	EncodePublicKeyRSA(ctx context.Context, strKey string) (*rsa.PublicKey, error)
	ClaimJWT(context.Context, *jwt.Token) (*map[string]interface{}, error)
	RefreshToken(ctx context.Context, stringToken string, user *entity.User) (bool, *TokenData, error)
	EncodePrivateKeyRSA(ctx context.Context, strKey string) (*rsa.PrivateKey, error)
	GeneratePublicKeyPem(ctx context.Context, privateKey *rsa.PrivateKey) []byte
	GenerateRsaKeyBioMetric(ctx context.Context) (string, error)

	AccessMiddleware(excludeUrl *ExcludeURLValidation, list *ListRouth) echo.MiddlewareFunc
	AccessLogger(ctx context.Context)
}

func (s *customMiddleware) EncodePublicKeyRSA(ctx context.Context, strKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(strKey))
	if block == nil {
		err := errors.New("failed decode key")
		s.logger.ErrorLogger(ctx, "encodeRSA.pem.Decode", err)
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		s.logger.ErrorLogger(ctx, "encodeRSA.parsePublicKey", err)
		return nil, err
	}
	rsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		err := errors.New("Failed to cast to RSA public key")
		s.logger.ErrorLogger(ctx, "encodeRSA.publickey.parsePublicKey", err)
		return nil, err
	}
	return rsaPubKey, nil
}

func (svc *customMiddleware) EncodePrivateKeyRSA(ctx context.Context, strKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(strKey))
	if block == nil {
		err := errors.New("failed decode key")
		svc.logger.ErrorLogger(ctx, "EncodePrivateKeyRSA.pem", err)
		return nil, err
	}
	rsaPriv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			svc.logger.ErrorLogger(ctx, "EncodePrivateKeyRSA.x509.Parse", err)
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
		// Konversi ke *rsa.PrivateKey
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			err = fmt.Errorf("not an RSA private key")
			svc.logger.ErrorLogger(ctx, "EncodePrivateKeyRSA.x509.Parse", err)
			return nil, err
		}
		return rsaKey, nil
	}

	return rsaPriv, nil

}
func (svc *customMiddleware) GeneratePublicKeyPem(ctx context.Context, privateKey *rsa.PrivateKey) []byte {
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		svc.logger.ErrorLogger(ctx, "GeneratePublicKeyPem", err)
		return nil
	}
	publicPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return publicPem
}
