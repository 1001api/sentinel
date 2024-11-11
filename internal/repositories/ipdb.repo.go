package repositories

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

type IPDBRepo interface {
	GetIP(ip net.IP) (*geoip2.City, error)
}

type IPDBRepoImpl struct {
	Reader *geoip2.Reader
}

func InitIPDBRepo(reader *geoip2.Reader) IPDBRepoImpl {
	return IPDBRepoImpl{
		Reader: reader,
	}
}

func (r *IPDBRepoImpl) GetIP(ip net.IP) (*geoip2.City, error) {
	return r.Reader.City(ip)
}
