package configs

import (
	"log"

	"github.com/oschwald/geoip2-golang"
)

func InitIPDBCon() *geoip2.Reader {
	ipdb, err := geoip2.Open("internal/ipdb/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal("Error opening ipdb", err.Error())
	}
	return ipdb
}
