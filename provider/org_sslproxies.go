package provider

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	SSLPROXIES_URL = "http://www.sslproxies.org/"
)

type Org_sslproxies struct {
	client *http.Client
}

func (p *Org_sslproxies) SetClient(client *http.Client) {
	p.client = client
}

func (p *Org_sslproxies) Load() ([]*ProxyItem, error) {
	b, err := httpGet(SSLPROXIES_URL, p.client)
	if err != nil {
		return nil, errors.New("Failed to read stream")
	}

	startBytes := []byte("<tbody>")
	endBytes := []byte("</tbody>")

	tbodyStart := bytes.Index(b, startBytes)
	tbodyEnd := bytes.Index(b, endBytes)
	if tbodyEnd <= tbodyStart {
		return nil, errors.New("Failed to parse stream")
	}

	bytes := b[tbodyStart : tbodyEnd+len(endBytes)]
	tbl := Tbody{}
	err = xml.Unmarshal(bytes, &tbl)
	if err != nil {
		return nil, err
	}

	ret := make([]*ProxyItem, len(tbl.Tr))
	cnt := 0
	for _, tr := range tbl.Tr {
		item := p.convert(&tr)
		if item != nil {
			ret[cnt] = item
			cnt++
		}
	}

	return ret, nil
}

func (p *Org_sslproxies) convert(tr *Tr) *ProxyItem {
	if len(tr.Td) < 4 {
		return nil
	}

	port, err := strconv.Atoi(tr.Td[1])
	if err != nil || port == 0 {
		return nil
	}

	var t int
	if strings.Contains(tr.Td[6], "HTTPS") {
		t = 3
	} else if strings.Contains(tr.Td[6], "HTTP") {
		t = 1
	}

	var a int
	if tr.Td[4] == "anonymous" {
		a = 1
	} else {
		a = 0
	}

	return &ProxyItem{
		Host:      tr.Td[0],
		Port:      port,
		Type:      t,
		Anonymous: a,
	}
}
