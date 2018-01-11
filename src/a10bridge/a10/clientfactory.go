package a10

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v2"
	"a10bridge/args"
	"strconv"
)

//BuildClient builds a10 client
func BuildClient(arguments args.Args) (api.Client, api.A10Error) {
	var err api.A10Error
	var client api.Client

	switch arguments.A10ApiVersion {
	case 2:
		client, err = v2.Connect(arguments)
		break
		/*
			case 3:
				client, err = v3.Connect(arguments)
				break
		*/
	default:
		buildError("Unsupported a10 api version " + strconv.Itoa(arguments.A10ApiVersion))
	}
	return client, err
}
