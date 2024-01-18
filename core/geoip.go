package core

import (
	"net"
	"os"
	"path"

	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

const CITY_MMDB_NAME = "GeoLite2-City.mmdb"
const ASN_MMDB_NAME = "GeoLite2-ASN.mmdb"

type GeoipResolver struct {
	citydb *geoip2.Reader
	asndb  *geoip2.Reader
}

func NewGeoipResolver() *GeoipResolver {
	resolver := &GeoipResolver{}
	var err error

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	citydb_name := path.Join(wd, CITY_MMDB_NAME)
	fs, _ := os.Stat(citydb_name)
	if fs == nil {
		return &GeoipResolver{}
	}

	resolver.citydb, err = geoip2.Open(citydb_name)
	if err != nil {
		panic(err)
	}

	asndb_name := path.Join(wd, ASN_MMDB_NAME)
	fs, _ = os.Stat(asndb_name)
	if fs == nil {
		return &GeoipResolver{}
	}

	resolver.asndb, err = geoip2.Open(asndb_name)
	if err != nil {
		panic(err)
	}

	logrus.Info("geoip database setup successfully")
	return resolver
}

type GeoipResult struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	ISOCountry string `json:"country_iso"`
	ASN        int    `json:"asn"`
}

func (r *GeoipResolver) Resolve(ipstr string) *GeoipResult {
	if r.citydb == nil || r.asndb == nil {
		return nil
	}

	ip := net.ParseIP(ipstr)
	result := &GeoipResult{}

	city, err := r.citydb.City(ip)
	if err != nil {
		return nil
	}

	result.City = city.City.Names["en"]
	result.Country = city.Country.Names["en"]
	result.ISOCountry = city.Country.IsoCode

	asn, err := r.asndb.ASN(ip)
	if err != nil {
		return nil
	}

	result.ASN = int(asn.AutonomousSystemNumber)

	return result
}
