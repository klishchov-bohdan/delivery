package tests

import (
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/tests/helper"
	"github.com/klishchov-bohdan/delivery/internal/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TokensTestSuite struct {
	suite.Suite
	cfg *config.Config
}

func (ts *TokensTestSuite) SetupSuite() {
	ts.cfg = &config.Config{
		AccessSecret:         "access_secret",
		RefreshSecret:        "refresh_secret",
		AccessTokenLifeTime:  1,
		RefreshTokenLifeTime: 2,
	}
}

func TestTokenService(t *testing.T) {
	suite.Run(t, new(TokensTestSuite))
}

func (ts *TokensTestSuite) TestGetTokenFromBearerString() {
	testCases := []helper.TestCaseGetBearerString{
		{
			Name:         "Get token successful",
			BearerString: "Bearer djvnrpenvdskjvnawpefo.aewoivnfpaeoghaierfurea.aurehfpaurfiaur",
			Expected:     "djvnrpenvdskjvnawpefo.aewoivnfpaeoghaierfurea.aurehfpaurfiaur",
		},
		{
			Name:         "Get empty token from incorrect string",
			BearerString: "Bearerdjvnrpenvdskjvnawpefo.aewoivnfpaeoghaierfurea.aurehfpaurfiaur",
			Expected:     "",
		},
	}

	for _, testCase := range testCases {
		ts.T().Run(testCase.Name, func(t *testing.T) {
			actual := token.GetTokenFromBearerString(testCase.BearerString)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}

func (ts *TokensTestSuite) TestValidateAccessToken() {
	// preparation
	userID := uuid.New()
	accessString, _ := token.GenerateToken(userID, ts.cfg.AccessTokenLifeTime, ts.cfg.AccessSecret)
	accessStringWithRefreshSecret, _ := token.GenerateToken(userID, ts.cfg.AccessTokenLifeTime, ts.cfg.RefreshSecret)
	accessExpiredString, _ := token.GenerateToken(userID, -1, ts.cfg.AccessSecret)
	testCases := []helper.TestCaseValidateToken{
		{
			Name:            "Valid access token string",
			TokenString:     accessString,
			IsValidExpected: true,
		},
		{
			Name:            "Invalid access token string",
			TokenString:     accessString + "hello_world",
			IsValidExpected: false,
		},
		{
			Name:            "Valid access token signed with refresh secret",
			TokenString:     accessStringWithRefreshSecret,
			IsValidExpected: false,
		},
		{
			Name:            "Expired access token",
			TokenString:     accessExpiredString,
			IsValidExpected: false,
		},
	}

	for _, testCase := range testCases {
		ts.T().Run(testCase.Name, func(t *testing.T) {
			isValid, _ := token.ValidateToken(testCase.TokenString, ts.cfg.AccessSecret)
			claims, _ := token.GetClaims(testCase.TokenString, ts.cfg.AccessSecret)
			if testCase.IsValidExpected {
				assert.True(t, isValid)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.ID)
			} else {
				assert.False(t, isValid)
				assert.Nil(t, claims)
			}
		})
	}
}

func (ts *TokensTestSuite) TestValidateRefreshToken() {
	// preparation
	userID := uuid.New()
	refreshString, _ := token.GenerateToken(userID, ts.cfg.RefreshTokenLifeTime, ts.cfg.RefreshSecret)
	refreshStringWithAccessSecret, _ := token.GenerateToken(userID, ts.cfg.RefreshTokenLifeTime, ts.cfg.AccessSecret)
	refreshExpiredString, _ := token.GenerateToken(userID, -1, ts.cfg.RefreshSecret)
	testCases := []helper.TestCaseValidateToken{
		{
			Name:            "Valid refresh token string",
			TokenString:     refreshString,
			IsValidExpected: true,
		},
		{
			Name:            "Invalid refresh token string",
			TokenString:     refreshString + "hello_world",
			IsValidExpected: false,
		},
		{
			Name:            "Valid refresh token signed with access secret",
			TokenString:     refreshStringWithAccessSecret,
			IsValidExpected: false,
		},
		{
			Name:            "Expired refresh token",
			TokenString:     refreshExpiredString,
			IsValidExpected: false,
		},
	}

	for _, testCase := range testCases {
		ts.T().Run(testCase.Name, func(t *testing.T) {
			isValid, _ := token.ValidateToken(testCase.TokenString, ts.cfg.RefreshSecret)
			claims, _ := token.GetClaims(testCase.TokenString, ts.cfg.RefreshSecret)
			if testCase.IsValidExpected {
				assert.True(t, isValid)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.ID)
			} else {
				assert.False(t, isValid)
				assert.Nil(t, claims)
			}
		})
	}
}
