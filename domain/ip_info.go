package domain

import (
	"encoding/json"
	"net"
)

type IpInfo struct {
	Ip        net.IP  `json:"ip"`
	Continent string  `db:"continent"  json:"continent"`
	Country   string  `db:"country"    json:"country"`
	StateProv string  `db:"state_prov" json:"state_prov"`
	City      string  `db:"city"       json:"city"`
	Latitude  float64 `db:"latitude"   json:"latitude"`
	Longitude float64 `db:"longitude"  json:"longitude"`
}

func (i *IpInfo) Bytes() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")

	return b
}

func (i *IpInfo) String() string {
	return string(i.Bytes())
}
