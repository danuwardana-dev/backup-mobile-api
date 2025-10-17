package middleware

import (
	"backend-mobile-api/model/dto"
	"backend-mobile-api/model/entity"
	"backend-mobile-api/model/enum"
	"backend-mobile-api/model/enum/pkgErr"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"strings"
)

type Claims struct {
	Uuid      string `json:"uuid"` //unix user
	Username  string `json:"username"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
}

type TokenData struct {
	AccessToken    string    `json:"access_token"`
	ExpiredAccess  time.Time `json:"expired_access"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiredRefresh time.Time `json:"expired_refresh"`
}

func (svc *customMiddleware) ValidateAuthorization(ctx context.Context, req *dto.ContextValue, excUrl *ExcludeURLValidation) (*dto.BaseResponse, error) {
	var roles []enum.RolesEnum

	//skip validation auth by list path
	for _, value := range excUrl.Authorization.ExcludeURL {
		if value == req.RequestPath {
			return nil, nil
		}
	}
	// vlaidate roles access
	roles = excUrl.Authorization.AccessByRole[req.HeaderPath]

	for {
		tokenString := req.HeaderAuthorization
		if tokenString == "" {
			err := errors.New("token is empty")
			svc.logger.ErrorLogger(ctx, "AuthV2.tokenString.empty", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
				Message:    pkgErr.UNAUTHORIZED_MSG,
			}, err
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(tokenString, prefix) {
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
				Message:    pkgErr.UNAUTHORIZED_MSG,
			}, errors.New("token is invalid")
		}
		tokenString = strings.TrimPrefix(tokenString, prefix)

		if redisData, err := svc.Redis.GetBlaclistJwt(ctx, tokenString); redisData == "active" {
			if err != nil {
				svc.logger.ErrorLogger(ctx, "AuthV2.redisData.active", err)
				return &dto.BaseResponse{
					StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
					Message:    pkgErr.SERVER_BUSY,
					Error:      err.Error(),
				}, err
			}
			err = errors.New("token is blacklisted")
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
				Message:    pkgErr.UNAUTHORIZED_MSG,
			}, err
		}
		rsaData, err := svc.EncodePublicKeyRSA(ctx, svc.jwtConfig.PublicKey)
		if err != nil {
			svc.logger.ErrorLogger(ctx, "AuthV2.rsaData.EncodePublicKeyRSA", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.UNDEFINED_ERROR_CODE,
				Message:    pkgErr.SERVER_BUSY,
				Error:      err.Error(),
			}, err
		}
		token, err := svc.ParseJwtToken(ctx, tokenString, rsaData)

		if err != nil {
			svc.logger.ErrorLogger(ctx, "AuthV2.token.ParseJwtToken", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
				Message:    pkgErr.UNAUTHORIZED_MSG,
				Error:      err.Error(),
			}, err
		}

		claimData, err := svc.ClaimJWT(ctx, token)
		if err != nil {
			svc.logger.ErrorLogger(ctx, "AuthV2.claimData.ClaimJWT", err)
			return &dto.BaseResponse{
				StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
				Message:    pkgErr.UNAUTHORIZED_MSG,
			}, err
		}
		if len(roles) > 0 {
			isUnAuthorized := false
			for _, rolesEnum := range roles {
				if strings.EqualFold(fmt.Sprint(rolesEnum), fmt.Sprint((*claimData)["role"])) {
					isUnAuthorized = true
				}
			}
			if !isUnAuthorized {
				err = errors.New("role access is not authorized")
				svc.logger.ErrorLogger(ctx, "AuthV2.isUnAuthorized.error", err)
				return &dto.BaseResponse{
					StatusCode: pkgErr.AUTH_UNAUTHORIZED_CODE,
					Message:    pkgErr.UNAUTHORIZED_MSG,
				}, err
			}
		}

		req.AuthEmail = fmt.Sprint((*claimData)["username"])
		req.AuthUUID = fmt.Sprint((*claimData)["uuid"])
		req.AuthRole = fmt.Sprint((*claimData)["role"])

		return nil, nil
	}
}

func (s *customMiddleware) RefreshToken(ctx context.Context, stringToken string, user *entity.User) (bool, *TokenData, error) {
	rsa, err := s.EncodePublicKeyRSA(ctx, s.jwtConfig.RefreshPublicKey)
	if err != nil {
		return false, nil, err
	}
	token, err := s.ParseJwtToken(ctx, stringToken, rsa)
	if err != nil {
		return false, nil, err
	}
	mapClaimData, err := s.ClaimJWT(ctx, token)
	if err != nil {
		return false, nil, err
	}
	MapRefresh := *mapClaimData
	claimData := Claims{
		Uuid:     fmt.Sprint(MapRefresh["uuid"]),
		Username: fmt.Sprint(MapRefresh["username"]),
		Role:     fmt.Sprint(MapRefresh["role"]),
	}
	if claimData.Uuid != user.UUID {
		return false, nil, nil
	}
	tokenData, err := s.CreateTokens(ctx, &Claims{
		Uuid:     fmt.Sprint(MapRefresh["uuid"]),
		Username: fmt.Sprint(MapRefresh["username"]),
		Role:     fmt.Sprint(MapRefresh["role"]),
	})
	if err != nil {
		return false, nil, err
	}
	return true, &tokenData, nil
}

func (s *customMiddleware) generateToken(ctx context.Context, user *Claims, key string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,

		"uuid": user.Uuid,
		"exp":  time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		s.logger.ErrorLogger(ctx, "generateToken.pem.Decode", errors.New("public key decode error"))
		return "", errors.New("customMiddleware block is nil")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		s.logger.ErrorLogger(context.Background(), "generateToken.customMiddleware.ParseECP", err)
		return "", err
	}
	return token.SignedString(privateKey)
}

func (s *customMiddleware) CreateTokens(ctx context.Context, user *Claims) (TokenData, error) {
	curentTime := time.Now()
	accessToken, err := s.generateToken(ctx, user, s.jwtConfig.SecreteKey, s.jwtConfig.Expiration)
	if err != nil {
		return TokenData{}, err
	}
	refreshToken, err := s.generateToken(ctx, user, s.jwtConfig.RefreshSecreteKey, s.jwtConfig.RefreshExpiration)
	if err != nil {
		return TokenData{}, err
	}
	return TokenData{
		AccessToken:    accessToken,
		ExpiredAccess:  curentTime.Add(s.jwtConfig.Expiration),
		RefreshToken:   refreshToken,
		ExpiredRefresh: curentTime.Add(s.jwtConfig.RefreshExpiration),
	}, nil
}

func (s *customMiddleware) ParseJwtToken(ctx context.Context, strJwt string, rsaPublicKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(strJwt, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			err := fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			s.logger.ErrorLogger(ctx, "ParseJwtToken.token.Method", err)
			return nil, err
		}
		return rsaPublicKey, nil
	})
	if err != nil {
		s.logger.ErrorLogger(ctx, "ParseJwtToken.customMiddleware.Parse", err)
		return nil, err
	}
	if !token.Valid {
		err = fmt.Errorf("invalid token")
		s.logger.ErrorLogger(ctx, "ParseJwtToken.token.Method", err)
		return nil, err
	}
	return token, nil
}

func (s *customMiddleware) ClaimJWT(ctx context.Context, token *jwt.Token) (*map[string]interface{}, error) {
	var claimMap = make(map[string]interface{})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for key, value := range claims {
			claimMap[key] = value
		}
		return &claimMap, nil
	}
	err := errors.New("invalid token")
	s.logger.ErrorLogger(ctx, "ClaimJWT.token.Claims.(customMiddleware.MapClaims)", err)
	return nil, err
}
func (svc *customMiddleware) GenerateRsaKeyBioMetric(ctx context.Context) (string, error) {
	privateKey, err := svc.EncodePrivateKeyRSA(ctx, svc.rootConfig.App.BiometricPrivateKey)
	if err != nil {
		svc.logger.ErrorLogger(ctx, "GenerateRsaKeyBioMetric", err)
		return "", err
	}
	byteKey := svc.GeneratePublicKeyPem(ctx, privateKey)
	endcodeData := base64.StdEncoding.EncodeToString(byteKey)
	return endcodeData, nil

}
