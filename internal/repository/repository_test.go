package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lotproject/user-service/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositoryTestSuite struct {
	suite.Suite
	db  *sqlx.DB
	cfg *config.Config
}

func Test_RepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupSuite() {
}

func (suite *RepositoryTestSuite) SetupTest() {
}

func (suite *RepositoryTestSuite) TearDownTest() {
}

func (suite *RepositoryTestSuite) Test_InitRepositories() {
	reps := InitRepositories(nil, nil)
	assert.IsType(suite.T(), (*Repositories)(nil), reps)
}
