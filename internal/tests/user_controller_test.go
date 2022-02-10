package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/controller"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"github.com/klishchov-bohdan/delivery/internal/services/mocks"
	"github.com/klishchov-bohdan/delivery/internal/tests/helper"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type UserControllerSuite struct {
	suite.Suite
	cfg *config.Config
}

func (suite *UserControllerSuite) SetupSuite() {
	suite.cfg = &config.Config{
		AccessSecret:         "access_secret",
		RefreshSecret:        "refresh_secret",
		AccessTokenLifeTime:  1,
		RefreshTokenLifeTime: 2,
	}

}

func (suite *UserControllerSuite) SetupTest() {

}

func (suite *UserControllerSuite) TearDownSuite() {

}

func (suite *UserControllerSuite) TearDownTest() {

}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}

func (suite *UserControllerSuite) TestGetProfile() {
	user, _ := models.CreateUser("test", "test@test.com", "someTestPassword123")
	userJSON, _ := json.Marshal(user)
	validAccessToken, _ := token.GenerateToken(user.ID, suite.cfg.AccessTokenLifeTime, suite.cfg.AccessSecret)
	testCases := []helper.TestCaseGetProfile{
		{
			Name:         "Valid access bearer string",
			BearerString: "Bearer " + validAccessToken,
			Method:       "GET",
			MockBehavior: func(s *mocks.MockUserService, uid uuid.UUID) {
				s.EXPECT().GetUserByID(uid).Return(user, nil)
			},
			ExpectedStatusCode:   200,
			ExpectedResponseBody: string(userJSON),
		},
		{
			Name:         "Not existing user",
			BearerString: "Bearer " + validAccessToken,
			Method:       "GET",
			MockBehavior: func(s *mocks.MockUserService, uid uuid.UUID) {
				s.EXPECT().GetUserByID(uid).Return(nil, errors.New("user does not exists"))
			},
			ExpectedStatusCode:   401,
			ExpectedResponseBody: "",
		},
		{
			Name:                 "Method not allowed",
			BearerString:         "Bearer " + validAccessToken,
			Method:               "POST",
			MockBehavior:         func(s *mocks.MockUserService, uid uuid.UUID) {},
			ExpectedStatusCode:   405,
			ExpectedResponseBody: "",
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.Name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			us := mocks.NewMockUserService(c)
			testCase.MockBehavior(us, user.ID)
			service := &services.Manager{User: us}
			uc := controller.NewUserController(service, suite.cfg)
			// test server
			r := chi.NewRouter()
			r.Get("/profile", uc.GetProfile)
			// test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.Method, "/profile", bytes.NewBufferString(""))
			req.Header.Set("Authorization", testCase.BearerString)
			// perform http
			r.ServeHTTP(w, req)
			// assert
			assert.Equal(t, testCase.ExpectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), testCase.ExpectedResponseBody)

		})
	}
}
