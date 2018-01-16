package v3

import (
	"a10bridge/a10/api"
	"a10bridge/args"
	"a10bridge/model"
)

type Client struct {
	sessionId string
}

func Connect(arguments args.Args) (api.Client, error) {
	var client api.Client
	return client, nil
}

func (client Client) Close() error {
	return nil
}

func (client Client) GetServer(serverName string) (*model.Node, error) {
	var server *model.Node
	return server, nil
}
