package tools

import (
	"github.com/pioz/faker"
	"github.com/rs/zerolog/log"
)

func ManageTestError(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}
}

func FakerBuild(v interface{}) {
	err := faker.Build(v)
	ManageTestError(err)
}

func BoolPointer(b bool) *bool {
	return &b
}

type testError struct {
	msg string
}

func NewTestError(msg string) error {
	return testError{msg: msg}
}
func (e testError) Error() string {
	return e.msg
}
