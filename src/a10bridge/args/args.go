package args

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
	APIURL, APIKey, APICert, A10URL, A10User, A10Pwd string
	A10ApiVersion                                    int
	Debug                                            bool
}

//Build builds arguments
func Build() (Args, error) {
	apiURL := addStringFlag("api-url", "url of kubernetes api endpoint")
	apiKey := addStringFlag("api-key", "path to key file")
	apiCert := addStringFlag("api-cert", "path to cert file")
	a10URL := addStringFlag("a10-url", "url of a10 api endpoint")
	a10User := addStringFlag("a10-usr", "a10 user")
	a10Pwd := addStringFlag("a10-pwd", "a10 password")
	a10ApiVersion := addIntFlag("a10-api-ver", "a10 api version, supported versions are 2 or 3")
	debug := addBoolFlag("debug", "run in debug mode")

	flag.Parse()

	args := Args{
		APIURL:        *apiURL,
		APIKey:        *apiKey,
		APICert:       *apiCert,
		A10URL:        *a10URL,
		A10User:       *a10User,
		A10Pwd:        *a10Pwd,
		A10ApiVersion: *a10ApiVersion,
		Debug:         *debug,
	}

	if args.Debug {
		args.printArgs()
	}

	return args, args.validate()
}

//Validate validates if all required parameters were passed in
func (inp Args) validate() error {
	if len(strings.TrimSpace(inp.A10URL)) == 0 {
		return errors.New("a10-url parameter is required")
	}

	if inp.A10ApiVersion != 2 && inp.A10ApiVersion != 3 {
		return errors.New("Unsupported value " + strconv.Itoa(inp.A10ApiVersion) + " for a10-api-ver parameter")
	}

	if len(strings.TrimSpace(inp.A10User)) == 0 {
		return errors.New("a10-usr parameter is required")
	}

	if len(strings.TrimSpace(inp.A10Pwd)) == 0 {
		return errors.New("a10-pwd parameter is required")
	}

	return nil
}

//printArgs prints calculated arguments
func (inp Args) printArgs() {
	fmt.Println("Using following argument values:")
	fmt.Println("api-url:", inp.APIURL)
	fmt.Println("api-key:", inp.APIKey)
	fmt.Println("api-cert:", inp.APICert)
	fmt.Println("a10-url:", inp.A10URL)
	fmt.Println("a10-usr:", inp.A10User)
	fmt.Println("a10-pwd:", inp.A10Pwd)
	fmt.Println("a10-api-ver:", inp.A10ApiVersion)
	fmt.Println()
}

//addStringFlag adds string flag using uppercased flagname to get the default value from environment variables
func addStringFlag(flagName, description string) *string {
	return flag.String(flagName, getEnv(flagName), description)
}

//addIntFlag adds int flag using uppercased flagname to get the default value from environment variables. If the
func addIntFlag(flagName, description string) *int {
	intDefault := 0
	stringDefault := getEnv(flagName)

	if len(stringDefault) != 0 {
		intDefault, _ = strconv.Atoi(stringDefault)
	}

	return flag.Int(flagName, intDefault, description)
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
