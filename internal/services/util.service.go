package services

import (
	"encoding/json"
	"net"
	"net/netip"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hubkudev/sentinel/internal/repositories"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/oschwald/geoip2-golang"
	"golang.org/x/crypto/bcrypt"
)

type UtilService interface {
	GenerateRandomID(length int) string
	ValidateInput(payload any) string
	GenerateHash(password string) string
	CompareHash(password string, hash string) bool
	ParseIP(str string) *netip.Addr
	ParseTimestamp(str string) time.Time
	LookupIP(ipStr string) *geoip2.City
	LookupAlphabetFromIdx(index int) string
	ByteToJSON(v []byte) interface{}
}

type UtilServiceImpl struct {
	Validate *validator.Validate
	IPRepo   repositories.IPDBRepo
}

func InitUtilService(
	validate *validator.Validate,
	ipRepo repositories.IPDBRepo,
) UtilServiceImpl {
	return UtilServiceImpl{
		Validate: validate,
		IPRepo:   ipRepo,
	}
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

			if err.Tag() == "max" {
				errMessage = err.Field() + " field is too long"
				break
			}

			if err.Tag() == "email" {
				errMessage = err.Field() + " field is invalid"
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

			if err.Tag() == "eqfield" && err.Field() == "ConfirmPassword" {
				errMessage = "Password & confirm password do not match"
				break
			}

			if err.Tag() == "min" && err.Field() == "Password" {
				errMessage = "Minimum length of a password is 8 characters"
				break
			}

			if err.Tag() == "password" {
				errMessage = "Password is too weak; include at least 1 uppercase letter, 1 symbol, and 1 number."
				break
			}

			if err.Tag() == "min" {
				errMessage = err.Field() + " field is too short. Make it at least 3 characters."
				break
			}

			// raw error which is not covered above
			errMessage = "Error on field " + err.StructField()
		}
	}

	return errMessage
}

func (s *UtilServiceImpl) GenerateHash(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return ""
	}
	return string(hashed)
}

func (s *UtilServiceImpl) CompareHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
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
	record, err := s.IPRepo.GetIP(ip)
	if err != nil {
		return nil
	}
	return record
}

func (s *UtilServiceImpl) LookupAlphabetFromIdx(index int) string {
	result := ""

	for index >= 0 {
		result = string(rune('A' + (index % 26)))
		index = index/26 - 1
	}

	return result
}

func (s *UtilServiceImpl) ByteToJSON(v []byte) interface{} {
	var result interface{}
	if err := json.Unmarshal(v, &result); err != nil {
		return nil
	}
	return result
}
