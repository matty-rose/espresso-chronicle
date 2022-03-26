package main

import (
	"fmt"

	"github.com/pulumi/pulumi-auth0/sdk/v2/go/auth0"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func deployAuth0Application(ctx *pulumi.Context, config *config.Config) (*pulumi.StringOutput, error) {
	client, err := auth0.NewClient(ctx, formatName(ctx, fmt.Sprintf("auth0-%s", "app")), &auth0.ClientArgs{
		AppType:                 pulumi.String("spa"),
		IsFirstParty:            pulumi.Bool(true),
		Name:                    pulumi.String(ProjectName),
		TokenEndpointAuthMethod: pulumi.String("client_secret_post"),
	})
	if err != nil {
		return nil, err
	}

	return &client.ClientId, nil
}

func deployAuth0SocialConnections(ctx *pulumi.Context, config *config.Config, clients *pulumi.StringArrayOutput) error {
	for _, social := range []string{"google-oauth2"} {
		_, err := auth0.NewConnection(ctx, formatName(ctx, fmt.Sprintf("auth0-%s", social)), &auth0.ConnectionArgs{
			Strategy:       pulumi.String(social),
			EnabledClients: clients,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
