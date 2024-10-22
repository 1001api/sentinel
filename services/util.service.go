package services

import (
	"net"
	"net/netip"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/oschwald/geoip2-golang"
)

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(fl.Field().String())
}

type UtilService interface {
	GenerateRandomID(length int) string
	ValidateInput(payload any) string
	ParseIP(str string) *netip.Addr
	ParseTimestamp(str string) time.Time
	LookupIP(ipStr string) *geoip2.City
}

type UtilServiceImpl struct {
	Validate *validator.Validate
	IPReader *geoip2.Reader
}

func (s *UtilServiceImpl) GenerateRandomID(length int) string {
	// generate random string id using nanoid package
	id, _ := gonanoid.New(length)
	return id
}

func (s *UtilServiceImpl) ValidateInput(payload any) string {
	if payload == nil {
		return "Invalid Payload"
	}

	// save error messages here
	var errMessage string

	errors := s.Validate.Struct(payload)
	if errors != nil {
		// loop through all possible errors,
		// then give appropriate message based on
		// defined error tag, StructField, etc
		for _, err := range errors.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				errMessage = err.Field() + " field is required"
				break
			}

			if err.Tag() == "timestamp" {
				errMessage = err.Field() + " field is invalid, please use ISO8601 date format"
				break
			}

			if err.Tag() == "uuid" {
				errMessage = err.Field() + " field is not a valid ID"
				break
			}

			if err.Tag() == "ip_addr" {
				errMessage = err.Field() + " field is not a valid IP address"
				break
			}

			// raw error which is not covered above
			errMessage = "Error on field " + err.StructField()
		}
	}

	return errMessage
}

func (s *UtilServiceImpl) ParseIP(str string) *netip.Addr {
	if str == "" {
		return nil
	}
	ip, err := netip.ParseAddr(str)
	if err != nil {
		return nil
	}
	return &ip
}

func (s *UtilServiceImpl) ParseTimestamp(str string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, str)
	return parsedTime
}

func (s *UtilServiceImpl) LookupIP(ipStr string) *geoip2.City {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	record, err := s.IPReader.City(ip)
	if err != nil {
		return nil
	}
	return record
}
