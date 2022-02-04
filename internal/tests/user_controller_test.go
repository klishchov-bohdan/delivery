package tests

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/services/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserControllerSuite struct {
	suite.Suite
	mockUserService *mocks.MockUserService
}

func (suite *UserControllerSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	suite.mockUserService = mocks.NewMockUserService(ctrl)

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
	id, _ := uuid.FromBytes([]byte("d246a1b8-85b9-11ec-a5c3-7c10c940acb3"))
	suite.mockUserService.EXPECT().DeleteUser(id).Return(id, nil)
}
