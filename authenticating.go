package true_vendor_sdk

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/wolvesdev/gohttplib"
)

// Keys
const (
	EntityKey    = "entity"
	AssistantKey = "assistant"
)

var (
	failedToken       = gohttplib.NewError("TRUE_VENDOR_SDK_FAILED_TOKEN_ERROR", "failed token error", "", nil)
	accessDeniedError = gohttplib.NewError("TRUE_VENDOR_SDK_ACCESS_DENIED_ERROR", "access denied error", "", nil)
)

// Authenticating interface
type Authenticating interface {
	GenerateTokenFor(assistantID, entityID string) (string, error)
	ValidateTokenFor(token string, assistantID, entityID string) error
}

// DefaultAuthenticating is implementation of Authenticating
type DefaultAuthenticating struct {
	secret string
}

// ValidateTokenFor is fn to validate token and  assistantID with entityID
func (auth *DefaultAuthenticating) ValidateTokenFor(token, assistantID, entityID string) error {
	data, err := auth.getDataFromToken(token)
	if err != nil {
		return failedToken.AsServerError(400)
	}
	if data[AssistantKey].(string) == assistantID && data[EntityKey].(string) == entityID {
		return nil
	}
	return accessDeniedError.AsServerError(403)
}

// NewDefaultAuthenticating is init
func NewDefaultAuthenticating(secret string) *DefaultAuthenticating {
	return &DefaultAuthenticating{secret: secret}
}

// GenerateTokenFor assistantID and entityID
func (auth *DefaultAuthenticating) GenerateTokenFor(assistantID, entityID string) (string, error) {
	return auth.generateToken(map[string]string{AssistantKey: assistantID, EntityKey: entityID})
}

func (auth *DefaultAuthenticating) getDataFromToken(token string) (map[string]interface{}, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.secret), nil
	})
	if err != nil {
		return nil, failedToken.AsServerError(400)
	}
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if ok && tokenObj.Valid {
		return claims["value"].(map[string]interface{}), nil
	}
	return nil, failedToken.AsServerError(400)
}

func (auth *DefaultAuthenticating) generateToken(value map[string]string) (string, error) {
	claims := struct {
		Value interface{} `json:"value"`
		jwt.StandardClaims
	}{
		value,
		jwt.StandardClaims{
			Issuer: "place",
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenObj.SignedString([]byte(auth.secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
