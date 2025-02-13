package customer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vap-tech/testcontainers_go/testhelpers"
	"log"
	"testing"
)

type CustomerRepoTestSuite struct {
	suite.Suite
	oracleContainer *testhelpers.OracleContainer
	repository      *Repository
	ctx             context.Context
}

func (suite *CustomerRepoTestSuite) SetupSuite() {
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

func (suite *CustomerRepoTestSuite) TearDownSuite() {
	if err := suite.oracleContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating oracle container: %s", err)
	}
}

func (suite *CustomerRepoTestSuite) TestCreateCustomer() {
	t := suite.T()

	customer, err := suite.repository.CreateCustomer(Customer{
		Id:    5,
		Name:  "Vitaliy",
		Email: "v.petrenko@gmail.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, customer.Id)
	assert.Equal(t, customer.Name, "Vitaliy")
	assert.Equal(t, customer.Email, "v.petrenko@gmail.com")
}

func (suite *CustomerRepoTestSuite) TestGetCustomerByEmail() {
	t := suite.T()

	customer, err := suite.repository.GetCustomerByEmail("a.petrenko@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "Anna", customer.Name)
	assert.Equal(t, "a.petrenko@gmail.com", customer.Email)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepoTestSuite))
}
