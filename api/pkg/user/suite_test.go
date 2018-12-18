package user

import (
	"testing"

	"github.com/devlover-id/api/pkg/utils/testhelper"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	testhelper.Suite
}

func TestUserSuite(t *testing.T) {
	s := &UserTestSuite{}
	s.NeedDB = true
	suite.Run(t, s)
}
