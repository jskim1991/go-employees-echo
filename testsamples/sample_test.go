package testsamples

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SampleTestSuite struct {
	suite.Suite
	testConstant string
}

func (suite *SampleTestSuite) SetupTest() {
	fmt.Println("SetupTest")
	suite.testConstant = "hello"
}

func (suite *SampleTestSuite) TearDownTest() {
	fmt.Println("TearDownTest")
}

func (suite *SampleTestSuite) SetupSubTest() {
	fmt.Println("SetupSubTest")
}

func (suite *SampleTestSuite) TearDownSubTest() {
	fmt.Println("TearDownSubTest")
}

func TestSampleTestSuite(t *testing.T) {
	suite.Run(t, new(SampleTestSuite))
}

func (suite *SampleTestSuite) TestSuiteOne() {
	suite.Run("testConstant should equal string hello", func() {
		// ... test code
		assert.Equal(suite.T(), "hello", suite.testConstant)
	})

	suite.Run("true should be true", func() {
		// ... test code
		assert.Equal(suite.T(), true, true)
	})

	suite.Run("1 should be 1", func() {
		// ... test code
		assert.Equal(suite.T(), 1, 1)
	})
}
