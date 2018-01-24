package config_test

import (
	"a10bridge/config"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestContext(t *testing.T) {
	tests := new(TestSuite)
	suite.Run(t, tests)
}

func (suite *TestSuite) TestBuildConfig_missingRequiredFlags() {
	original := os.Args
	defer func() { os.Args = original }()
	os.Args = original[0:1]

	flag.CommandLine = flag.NewFlagSet(``, flag.PanicOnError)
	_, err := config.BuildConfig()
	suite.Assert().NotNil(err)
}

func (suite *TestSuite) TestBuildConfig_withRequiredFlags() {
	original := os.Args
	defer func() { os.Args = original }()

	expectedPassword := "file_pwd"

	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config=testdata/config1.yaml")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	conf, err := config.BuildConfig()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(conf)

	suite.Assert().Equal(2, len(conf.A10Instances))

	suite.Assert().Equal("lga-lb01", conf.A10Instances[0].Name)
	suite.Assert().Equal("lga-lb02", conf.A10Instances[1].Name)

	suite.Assert().Equal("https://lga-lb01", conf.A10Instances[0].APIUrl)
	suite.Assert().Equal("https://lga-lb02", conf.A10Instances[1].APIUrl)

	suite.Assert().Equal("dingo", conf.A10Instances[0].UserName)
	suite.Assert().Equal("dongo", conf.A10Instances[1].UserName)

	suite.Assert().Equal(2, conf.A10Instances[0].APIVersion)
	suite.Assert().Equal(3, conf.A10Instances[1].APIVersion)

	suite.Assert().Equal(expectedPassword, conf.A10Instances[0].Password)
	suite.Assert().Equal(expectedPassword, conf.A10Instances[1].Password)
}

func (suite *TestSuite) TestBuildConfig_withRequiredFlagsFromEnvironment() {
	os.Setenv("A10CONFIG", "testdata/config1.yaml")
	defer os.Unsetenv("A10CONFIG")

	original := os.Args
	defer func() { os.Args = original }()

	expectedPassword := "file_pwd"
	os.Args = original[0:1]
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	conf, err := config.BuildConfig()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(conf)

	suite.Assert().Equal(2, len(conf.A10Instances))

	suite.Assert().Equal("lga-lb01", conf.A10Instances[0].Name)
	suite.Assert().Equal("lga-lb02", conf.A10Instances[1].Name)

	suite.Assert().Equal("https://lga-lb01", conf.A10Instances[0].APIUrl)
	suite.Assert().Equal("https://lga-lb02", conf.A10Instances[1].APIUrl)

	suite.Assert().Equal("dingo", conf.A10Instances[0].UserName)
	suite.Assert().Equal("dongo", conf.A10Instances[1].UserName)

	suite.Assert().Equal(2, conf.A10Instances[0].APIVersion)
	suite.Assert().Equal(3, conf.A10Instances[1].APIVersion)

	suite.Assert().Equal(expectedPassword, conf.A10Instances[0].Password)
	suite.Assert().Equal(expectedPassword, conf.A10Instances[1].Password)
}

func (suite *TestSuite) TestBuildConfig_passwordFromCli() {
	original := os.Args
	defer func() { os.Args = original }()

	expectedPassword := "cli_pwd"

	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-pwd="+expectedPassword)
	os.Args = append(os.Args, "-a10-config=testdata/config2.yaml")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	conf, err := config.BuildConfig()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(conf)

	suite.Assert().Equal(expectedPassword, conf.A10Instances[0].Password)
	suite.Assert().Equal(expectedPassword, conf.A10Instances[1].Password)
}

func (suite *TestSuite) TestBuildConfig_nameFromUrl() {
	original := os.Args
	defer func() { os.Args = original }()

	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-pwd=blah")
	os.Args = append(os.Args, "-a10-config=testdata/config4.yaml")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	conf, err := config.BuildConfig()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(conf)

	suite.Assert().Equal(conf.A10Instances[0].APIUrl, conf.A10Instances[0].Name)
	suite.Assert().Equal(conf.A10Instances[1].APIUrl, conf.A10Instances[1].Name)
}

func (suite *TestSuite) TestBuildConfig_debugMode() {
	original := os.Args
	defer func() { os.Args = original }()

	os.Args = original[0:1]
	os.Args = append(os.Args, "-debug")
	os.Args = append(os.Args, "-a10-config=testdata/config1.yaml")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	conf, err := config.BuildConfig()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(conf)
	suite.Assert().True(*conf.Arguments.Debug)
}

func (suite *TestSuite) TestBuildConfig_debugModeFromEnvironment() {
	os.Setenv("DEBUG", "true")
	defer os.Unsetenv("A10CONFIG")

	original := os.Args
	defer func() { os.Args = original }()

	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config=testdata/config1.yaml")
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	conf, err := config.BuildConfig()

	suite.Assert().Nil(err)
	suite.Assert().NotNil(conf)
	suite.Assert().True(*conf.Arguments.Debug)
}

func (suite *TestSuite) TestBuildConfig_notExistentConfigFile() {
	original := os.Args
	defer func() { os.Args = original }()
	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config=testdata/i.dont.exist")

	flag.CommandLine = flag.NewFlagSet(``, flag.PanicOnError)
	_, err := config.BuildConfig()
	suite.Assert().NotNil(err)
}

func (suite *TestSuite) TestBuildConfig_invalidConfigFile() {
	original := os.Args
	defer func() { os.Args = original }()
	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config=testdata/config3.yaml")

	flag.CommandLine = flag.NewFlagSet(``, flag.PanicOnError)
	_, err := config.BuildConfig()
	suite.Assert().NotNil(err)
}

func (suite *TestSuite) TestBuildConfig_unreadableConfigFile() {
	err := os.Chmod("testdata/config1.yaml", 0000)
	if err != nil {
		suite.T().Skip("Unable to prepare unreadable file for test")
	}
	defer os.Chmod("testdata/config1.yaml", 0660)

	original := os.Args
	defer func() { os.Args = original }()
	os.Args = original[0:1]
	os.Args = append(os.Args, "-a10-config=testdata/config1.yaml")

	flag.CommandLine = flag.NewFlagSet(``, flag.PanicOnError)
	_, err = config.BuildConfig()
	suite.Assert().NotNil(err)
}
