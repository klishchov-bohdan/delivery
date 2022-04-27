package helper

type TestCaseGetBearerString struct {
	Name         string
	BearerString string
	Expected     string
}

type TestCaseValidateToken struct {
	Name            string
	TokenString     string
	IsValidExpected bool
}
