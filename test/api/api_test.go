package api_test

import (
	"testing"

	"github.com/constant-money/constant-web-api/config"
	"github.com/constant-money/constant-web-api/daos"
	"github.com/stretchr/testify/suite"

	_ "github.com/go-sql-driver/mysql"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type APITestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *APITestSuite) SetupTest() {
	// load config
	conf := config.GetConfig()

	// init daos
	if err := daos.Init(conf); err != nil {
		panic(err)
	}

}

func (suite *APITestSuite) TestAPIInitSuccessfully() {

}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
