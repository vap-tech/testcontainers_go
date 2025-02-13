package rates

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vap-tech/testcontainers_go/testhelpers"
	"log"
	"testing"
	"time"
)

type RatesRepoTestSuite struct {
	suite.Suite
	oracleContainer *testhelpers.OracleContainer
	repository      *Repository
	ctx             context.Context
}

func (suite *RatesRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	oraContainer, err := testhelpers.CreateOracleContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.oracleContainer = oraContainer
	repository, err := NewRepository(suite.oracleContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository
}

func (suite *RatesRepoTestSuite) TearDownSuite() {
	if err := suite.oracleContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating oracle container: %s", err)
	}
}

func (suite *RatesRepoTestSuite) TestGetRates() {
	t := suite.T()

	d := time.Date(2025, time.January, 26, 0, 0, 0, 0, time.UTC)
	rate, err := suite.repository.GetRates(d)
	assert.NoError(t, err)
	assert.NotNil(t, rate.day)
	assert.Equal(t, rate.value, 98.2636)

	d = time.Date(2025, time.January, 25, 0, 0, 0, 0, time.UTC)
	rate, err = suite.repository.GetRates(d)
	assert.NoError(t, err)
	assert.NotNil(t, rate.day)
	assert.Equal(t, rate.value, 99.0978)

}

func TestRatesRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RatesRepoTestSuite))
}
