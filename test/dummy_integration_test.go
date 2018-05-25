package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

// The full suite of integration tests:
type fullSuite struct {
	suite.Suite
}

// What to do to set it up:
func (suite *fullSuite) SetupSuite() {
}

func (suite *fullSuite) TearDownSuite() {
}

func (suite *fullSuite) SetupTest() {
}

func (suite *fullSuite) TearDownTest() {
}

func (suite *fullSuite) TestIntegration() {
	assert.Equal(suite.T(), "PASS", integration(), "Cmd+Q for more information.")
}

func TestFullSuite(t *testing.T) {
	fullTests := new(fullSuite)
	suite.Run(t, fullTests)
}

//////////////////////////////////////////////////////////////////////

// The DB only suite of integration tests:
type dbOnlySuite struct {
	suite.Suite
}

func (suite *dbOnlySuite) TestDBSetup() {
	assert.Equal(suite.T(), "PASS", dbSetup(), "Cmd+Q for more information.")
}

func TestDBOnlySuite(t *testing.T) {
	dbOnlyTests := new(dbOnlySuite)
	suite.Run(t, dbOnlyTests)
}
