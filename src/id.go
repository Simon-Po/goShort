package main
import (
	"strings"
	"errors"
	"github.com/google/uuid"
) 

type idGen interface {
	Generate(length int32) (string,error)
} 

type Uuid struct {

}



func (id *Uuid) Generate(length int32) (string, error) {
	uuidLong := strings.ReplaceAll(uuid.New().String(), "-", "")

	if length > 32 {
		return "", errors.New("length too long")
	}
	if length != 0 {
		return uuidLong[:length], nil
	}
	return uuidLong, nil
}
