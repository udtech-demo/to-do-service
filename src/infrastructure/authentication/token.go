package authentication

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"todo-service/src/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/google/uuid"
)

var (
	path, getWdErr = os.Getwd()
	privateKey     = path + "/rsa_keys/private_key.pem"
	publicKey      = path + "/rsa_keys/public_key.pem"
)

type jwtConfigurator struct {
	logger          *zap.Logger
	accessTokenDur  time.Duration
	refreshTokenDur time.Duration
	signer          *rsa.PrivateKey
	verifier        *rsa.PublicKey
	signingMethod   *jwt.SigningMethodRSA
}

type JwtConfigurator interface {
	CreateTokenPair(ctx context.Context, userUUID uuid.UUID) (*models.SessionDetails, error)
	ValidateJwtToken(token string) (*jwt.Token, error)
}

func NewJwtConfigurator(logger *zap.Logger, privateKeyPath, publicKeyPath string) JwtConfigurator {
	// Check os.Getwd() err
	if getWdErr != nil {
		panic(getWdErr)
	}

	if privateKeyPath != "" {
		privateKey = privateKeyPath
	}

	if publicKeyPath != "" {
		publicKey = publicKeyPath
	}

	signBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		panic(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	verifyBytes, err := ioutil.ReadFile(publicKey)
	if err != nil {
		panic(err)
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}

	return &jwtConfigurator{
		accessTokenDur:  viper.GetDuration("jwt.at_lifetime"),
		refreshTokenDur: viper.GetDuration("jwt.rt_lifetime"),
		signer:          signKey,
		verifier:        verifyKey,
		signingMethod:   jwt.SigningMethodRS256,
		logger:          logger,
	}
}

func (jc *jwtConfigurator) CreateTokenPair(ctx context.Context, userUUID uuid.UUID) (*models.SessionDetails, error) {
	// 1. Resolve lifetime

	now := time.Now()
	acExp := now.Add(jc.accessTokenDur).Unix()
	rtExp := now.Add(jc.refreshTokenDur).Unix()

	// 2. Make access and refresh claims
	at := jwt.NewWithClaims(jc.signingMethod, &models.JwtCustomClaim{
		ID: userUUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: acExp,
		},
	})

	accessToken, err := at.SignedString(jc.signer)
	if err != nil {
		jc.logger.Sugar().Errorf("create: sign token: %s", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	rt := jwt.NewWithClaims(jc.signingMethod, &models.JwtCustomClaim{
		ID: userUUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: rtExp,
		},
	})

	refreshToken, err := rt.SignedString(jc.signer)
	if err != nil {
		jc.logger.Sugar().Errorf("create: sign token: %s", err)
		return nil, errors.New(http.StatusText(http.StatusInternalServerError))
	}

	tokens := models.SessionDetails{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AtExpires:    acExp,
		RtExpires:    rtExp,
	}

	return &tokens, nil
}

func (jc *jwtConfigurator) ValidateJwtToken(token string) (*jwt.Token, error) {
	// Verify and extract claims from a token:
	return jwt.ParseWithClaims(token, &models.JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return jc.verifier, nil
	})
}
