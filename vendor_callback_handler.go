package true_vendor_sdk

import (
	"github.com/techpro-studio/gohttplib"
	"github.com/techpro-studio/gohttplib/validator"
	"net/http"
)

// OKJSON  is default response for ok
var OKJSON = map[string]interface{}{"ok": 1}

// VendorCallbackHandler callback handler for vendor api
type VendorCallbackHandler struct {
	useCase VendorUseCase
}

// NewVendorCallbackHandler is init for VendorCallbackHandler
func NewVendorCallbackHandler(useCase VendorUseCase) *VendorCallbackHandler {
	return &VendorCallbackHandler{useCase: useCase}
}

// PostRegistrable is a Router with register Post handlers
type PostRegistrable interface {
	Post(route string, handler http.Handler)
}

// RegisterAPI in  router
func (handler *VendorCallbackHandler) RegisterAPI(router PostRegistrable) {
	router.Post("/vendor/reserve", http.HandlerFunc(handler.ReserveVerificationHandler))
	router.Post("/vendor/reserve/callback", http.HandlerFunc(handler.ReserveCallbackHandler))
	router.Post("/vendor/cancel", http.HandlerFunc(handler.CancelReserveHandler))
	router.Post("/vendor/resolve", http.HandlerFunc(handler.ResolveVerification))
	router.Post("/vendor/resolve/callback", http.HandlerFunc(handler.ResolveCallbackHandler))
}

func (handler *VendorCallbackHandler) writeMapOrError(w http.ResponseWriter, m map[string]interface{}, err error) {
	if err != nil {
		if serverError, ok := err.(gohttplib.ServerError); ok {
			serverError.Write(w)
			return
		}
		description := err.Error()
		code := "VENDOR_ERROR"
		gohttplib.NewServerError(400, "undefined", description, code, nil).Write(w)
	} else {
		gohttplib.WriteJson(w, m, 200)
	}
}

// ReserveVerificationHandler implementation
func (self *VendorCallbackHandler) ReserveVerificationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := validator.GetValidatedBody(r, reserveVerificationValidatorMap())
	if err != nil {
		err.(gohttplib.ServerError).Write(w)
	}
	token, err := self.useCase.ReservationVerification(r.Context(), body["id"].(string), body["user_id"].(string))
	self.writeMapOrError(w, map[string]interface{}{"token": token}, err)
}

// ReserveCallbackHandler implementation
func (handler *VendorCallbackHandler) ReserveCallbackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := validator.GetValidatedBody(r, idValidatorMap())
	if err != nil {
		err.(gohttplib.ServerError).Write(w)
	}
	err = handler.useCase.ReservedCallback(r.Context(), body["id"].(string))
	handler.writeMapOrError(w, OKJSON, err)
}

// ResolveCallbackHandler implementation
func (handler *VendorCallbackHandler) ResolveCallbackHandler(w http.ResponseWriter, r *http.Request) {
	body, err := validator.GetValidatedBody(r, idValidatorMap())
	if err != nil {
		err.(gohttplib.ServerError).Write(w)
	}
	err = handler.useCase.ResolvingCallback(r.Context(), body["id"].(string))
	handler.writeMapOrError(w, OKJSON, err)
}

// CancelReserveHandler implementation
func (handler *VendorCallbackHandler) CancelReserveHandler(w http.ResponseWriter, r *http.Request) {
	body, err := validator.GetValidatedBody(r, idValidatorMap())
	if err != nil {
		err.(gohttplib.ServerError).Write(w)
	}
	err = handler.useCase.CancelReservationCallback(r.Context(), body["id"].(string))
	handler.writeMapOrError(w, OKJSON, err)
}

// ResolveVerification implementation
func (handler *VendorCallbackHandler) ResolveVerification(w http.ResponseWriter, r *http.Request) {
	body, err := validator.GetValidatedBody(r, resolveVerificationValidatorMap())
	if err != nil {
		err.(gohttplib.ServerError).Write(w)
	}
	token, err := handler.useCase.ResolvingVerification(r.Context(), body["id"].(string), body["token"].(string), body["verdict"].(string),
		body["note"].(string))
	handler.writeMapOrError(w, map[string]interface{}{"token": token}, err)
}
