package a10

import (
	"a10bridge/a10/api"
	"a10bridge/a10/v2"
	"a10bridge/a10/v3"
	"a10bridge/config"
	"strconv"
)

//BuildClient builds a10 client
func BuildClient(a10Instance *config.A10Instance) (api.Client, api.A10Error) {
	var err api.A10Error
	var client api.Client

	switch a10Instance.APIVersion {
	case 2:
		client, err = v2.Connect(a10Instance)
		break
	case 3:
		client, err = v3.Connect(a10Instance)
		break
	default:
		err = buildError("Unsupported a10 api version " + strconv.Itoa(a10Instance.APIVersion))
	}
	return client, err
}
