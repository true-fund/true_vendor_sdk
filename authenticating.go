package true_vendor_sdk

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/alexmay23/httputils"
)

const (
	EntityIDKey  = "entity_id"
	AssistantKey = "assistant_id"
)

var (
	failedToken       httputils.Error = httputils.UndefinedKeyError("TRUE_VENDOR_SDK_FAILED_TOKEN_ERROR", "failed token error")
	accessDeniedError httputils.Error = httputils.UndefinedKeyError("TRUE_VENDOR_SDK_ACCESS_DENIED_ERROR", "access denied error")
)

type Authenticating interface {
	GenerateTokenFor(assistantID, entityId string) (string, error)
	ValidateTokenFor(token string, assistantId, entityId string) error
}

type DefaultAuthenticating struct {
	secret string
}

func (self *DefaultAuthenticating) ValidateTokenFor(token, assistantId, entityId string) error {
	data, err := self.getDataFromToken(token)
	if err != nil {
		return failedToken.AsServerError(400)
	}
	if data[AssistantKey] == assistantId && data[EntityIDKey] == entityId {
		return nil
	}
	return accessDeniedError.AsServerError(403)
}

func NewDefaultAuthenticating(secret string) *DefaultAuthenticating {
	return &DefaultAuthenticating{secret: secret}
}

func (self *DefaultAuthenticating) GenerateTokenFor(assistantID, entityId string) (string, error) {
	return self.generateToken(map[string]string{"assistant": assistantID, "entity": entityId})
}

func (self *DefaultAuthenticating) getDataFromToken(token string) (map[string]interface{}, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(self.secret), nil
	})
	if err != nil {
		return nil, failedToken.AsServerError(400)
	}
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if ok && tokenObj.Valid {
		return claims, nil
	}
	return nil, failedToken.AsServerError(400)
}

func (self *DefaultAuthenticating) generateToken(value interface{}) (string, error) {
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

	token, err := tokenObj.SignedString([]byte(self.secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
