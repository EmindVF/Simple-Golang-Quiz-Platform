package utility

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Details of an JWT.
type TokenDetails struct {
	Token     string
	UserId    int32
	ExpiresIn int64
}

func CreateToken(userId int32, ttl time.Duration, privateKey []byte) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: now.Add(ttl).Unix(),
		UserId:    userId,
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		panic(fmt.Errorf("error parsing private key to jwt: %v", err))
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = strconv.FormatInt(int64(userId), 10)
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims).SignedString(key)
	if err != nil {
		panic(fmt.Errorf("error creating jwt token: %v", err))
	}
	return td, nil
}

func ValidateToken(token string, publicKey []byte) (*TokenDetails, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		panic(fmt.Errorf("error parsing public key to jwt: %v", err))
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validation error: invalid token")
	}

	var userId int32 = 0
	sub := claims["sub"]
	if subStr, ok := sub.(string); ok {
		subInt, err := strconv.ParseInt(subStr, 10, 32)
		if err != nil {
			panic(fmt.Errorf("cannot parse jwt claim"))
		}
		userId = int32(subInt)
	} else {
		panic(fmt.Errorf("cannot parse jwt claim"))
	}

	return &TokenDetails{
		UserId: userId,
	}, nil
}
