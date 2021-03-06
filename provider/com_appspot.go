package provider

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
)

const (
	APPSPOT_URL = "http://%s.appspot.com/proxy.json?cache=%s"
)

type Com_appspot struct {
	Appid  string
	Cache bool
	client *http.Client
}

func (p *Com_appspot) Load() ([]*ProxyItem, error) {
	return p.load(nil)
}

func (p *Com_appspot) SetClient(client *http.Client) {
	p.client = client
}

func (p *Com_appspot) load(config interface{}) ([]*ProxyItem, error) {
	url := fmt.Sprintf(APPSPOT_URL, p.Appid, strconv.FormatBool(p.Cache))
	b, err := httpGet(url, p.client)
	if err != nil {
		return nil, err
	}

	ret := make([]*ProxyItem, 0)
	err = json.Unmarshal(b, &ret)

	return ret, err
}
