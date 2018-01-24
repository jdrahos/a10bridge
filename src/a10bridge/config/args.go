package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Args arguments
type Args struct {
	A10Pwd    *string
	A10Config *string
	Debug     *bool
}

func buildArguments() (*Args, error) {
	args := Args{
		A10Config: addStringFlag("a10-config", "path to a10 config yaml file"),
		A10Pwd:    addStringFlag("a10-pwd", "a10 password"),
		Debug:     addBoolFlag("debug", "run in debug mode"),
	}

	flag.Parse()

	if *args.Debug {
		args.printArgs()
	}

	return &args, args.validate()
}

//Validate validates if all required parameters were passed in
func (toValidate Args) validate() error {
	if len(strings.TrimSpace(*toValidate.A10Config)) == 0 {
		return errors.New("a10-config parameter is required")
	}

	if _, err := os.Stat(*toValidate.A10Config); os.IsNotExist(err) {
		return errors.New("a10-config parameter points to notexistent file")
	}

	return nil
}

//printArgs prints calculated arguments
func (inp Args) printArgs() {
	fmt.Println("Using following argument values:")
	fmt.Println("a10-config:", inp.A10Config)
	fmt.Println("a10-pwd:", inp.A10Pwd)
	fmt.Println()
}

//addStringFlag adds string flag using uppercased flagname to get the default value from environment variables
func addStringFlag(flagName, description string) *string {
	return flag.String(flagName, getEnv(flagName), description)
}

//addBoolFlag adds boolean flag using upercased flagname to get the default value from environment variables
func addBoolFlag(flagName, description string) *bool {
	boolDefault := false
	stringDefault := getEnv(flagName)

	if len(stringDefault) != 0 {
		boolDefault, _ = strconv.ParseBool(stringDefault)
	}

	return flag.Bool(flagName, boolDefault, description)
}

//getEnv get the flag value from environment variables
func getEnv(flagName string) string {
	upperFlagName := strings.ToUpper(flagName)
	normalizedEnvName := strings.Replace(upperFlagName, "-", "", -1)
	return strings.TrimSpace(os.Getenv(normalizedEnvName))
}
