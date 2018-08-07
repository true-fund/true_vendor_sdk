package true_vendor_sdk

import (
	"net/http"
	"time"
	"bytes"
	"context"
	"encoding/json"
	"github.com/alexmay23/httpshared/shared"
	"github.com/alexmay23/httputils"
)

type VendorApp struct {
	key string
	secret string
}


func NewApplication(Key, Secret string)*VendorApp {
	return &VendorApp{key:Key, secret:Secret}
}


func (self *VendorApp)CreateEntity(description, url, language, id string)(shared.IDContainer, error){
	container := shared.IDContainer{}
	err := self.makeRequest("https://trueapi.alexmay23.com/vendor/entity", map[string]interface{}{"description":description,
	"language":language, "url": url, "client_uid": id}, &container)
	return container, err
}

func (self *VendorApp)makeRequest(url string, parameters map[string]interface{}, value interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	b, err := json.Marshal(parameters)
	if err != nil {
		return  err
	}
	body := bytes.NewReader(b)
	req, _ := http.NewRequest("POST", url, body)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("True-Auth-Vendor-Handshake", "SALAM")
	req.Header.Set("True-Auth-Vendor-Key", self.key)
	req.Header.Set("True-Auth-Vendor-Secret", self.secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return  err
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var serverErrors httputils.Errors
		err = decoder.Decode(&serverErrors)
		if err != nil {
			return err
		}
		return serverErrors
	} else{
		err = decoder.Decode(value)
		if err != nil {
			return err
		}
		return nil
	}
}


