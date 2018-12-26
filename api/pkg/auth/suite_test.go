package auth

import (
	"testing"

	"github.com/devlover-id/api/pkg/utils/testhelper"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	testhelper.Suite
}

func TestAuthSuite(t *testing.T) {
	s := &AuthTestSuite{}
	s.NeedDB = true
	suite.Run(t, s)
}
