package helper

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/services/mocks"
)

type TestCaseGetProfile struct {
	Name                 string
	BearerString         string
	Method               string
	MockBehavior         func(s *mocks.MockUserService, uid uuid.UUID)
	ExpectedStatusCode   int
	ExpectedResponseBody string
}
