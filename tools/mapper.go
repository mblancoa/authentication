package tools

import (
	"github.com/devfeel/mapper"
	"log"
	"reflect"
)

// Mapper mappers and sets value from struct fromObj to toObj
func Mapper(fromObj, toObj interface{}) {
	err := mapper.Mapper(fromObj, toObj)
	if err != nil {
		from := reflect.TypeOf(fromObj).Name()
		to := reflect.TypeOf(toObj).Name()
		log.Fatalf("Error Mappering from %s to %s", from, to)
	}
}
