package true_vendor_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/wolvesdev/gohttplib"
)

// VendorApp is an vendor applications structure
type VendorApp struct {
	key    string
	secret string
}

// TrueService is a struct to connect to true-fund vendor API
type TrueService interface {
	CreateEntity(description, url, language, id string, location *ObjectLocation) (gohttplib.IDContainer, error)
	DeleteEntity(id string) (gohttplib.IDContainer, error)
}

//ObjectLocation is an object location struct
type ObjectLocation struct {
	Latitude  float64
	Longitude float64
}

// NewApplication is init for Vendor app
func NewApplication(Key, Secret string) *VendorApp {
	return &VendorApp{key: Key, secret: Secret}
}

// DeleteEntity function deletes entity on true fund server. id is an client id of you service  item.
func (app *VendorApp) DeleteEntity(id string) (gohttplib.IDContainer, error) {
	container := gohttplib.IDContainer{}
	//
	trueURL := "https://api.true-fund.com/vendor/entity/delete"
	err := app.makeRequest(trueURL, map[string]interface{}{"id": id}, &container)
	return container, err
}

// CreateEntity creates entity on true fund server
func (app *VendorApp) CreateEntity(description, url, language, id string, location *ObjectLocation) (gohttplib.IDContainer, error) {
	container := gohttplib.IDContainer{}
	params := map[string]interface{}{"description": description,
		"language": language, "url": url, "client_uid": id}
	if location != nil {
		params["latitude"] = location.Latitude
		params["longitude"] = location.Longitude
	}
	trueURL := "https://api.true-fund.com/vendor/entity"
	err := app.makeRequest(trueURL, params, &container)
	return container, err
}

func (app *VendorApp) makeRequest(url string, parameters map[string]interface{}, value interface{}) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	b, err := json.Marshal(parameters)
	if err != nil {
		cancelFn()
		return err
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("True-Auth-Vendor-Handshake", "SALAM")
	req.Header.Set("True-Auth-Vendor-Key", app.key)
	req.Header.Set("True-Auth-Vendor-Secret", app.secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancelFn()
		return err
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var serverErrors gohttplib.Errors
		err = decoder.Decode(&serverErrors)
		if err != nil {
			return err
		}
		return gohttplib.ServerError{StatusCode: resp.StatusCode, Errors: serverErrors}
	}
	err = decoder.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
