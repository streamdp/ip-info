package domain

import (
	"encoding/json"
	"net"
)

type IpInfo struct {
	Ip        net.IP  `json:"ip"`
	Continent string  `json:"continent" db:"continent"`
	Country   string  `json:"country" db:"country"`
	StateProv string  `json:"state_prov" db:"state_prov"`
	City      string  `json:"city" db:"city"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

func (i *IpInfo) Bytes() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

func (i *IpInfo) String() string {
	return string(i.Bytes())
}
