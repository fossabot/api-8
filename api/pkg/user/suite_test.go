package user

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/pinterkode/pinterkode/api/pkg/utils/testhelper"
)

type UserTestSuite struct {
	testhelper.Suite
}

func TestUserSuite(t *testing.T) {
	s := &UserTestSuite{}
	s.NeedDB = true
	suite.Run(t, s)
}
