package services

import (
	"errors"
	"regexp"
)

var (
	ERROR_INVALID_ZIPCODE = errors.New("invalid zipcode")
)

type ZipcodeServiceInterface interface {
	CheckZipcode(zipcode string) error
}

type ZipcodeService struct{}

func NewZipcodeService() *ZipcodeService {
	return &ZipcodeService{}
}

func (z *ZipcodeService) isValidZipcode(zipcode string) bool {
	pattern := `^\d{5}-?\d{3}$`
	match, _ := regexp.MatchString(pattern, zipcode)
	return match
}

func (z *ZipcodeService) CheckZipcode(zipcode string) error {
	if !z.isValidZipcode(zipcode) {
		return ERROR_INVALID_ZIPCODE
	}

	return nil
}
