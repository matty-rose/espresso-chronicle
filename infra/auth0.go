package main

import (
	"fmt"

	"github.com/pulumi/pulumi-auth0/sdk/v2/go/auth0"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func deployAuth0SocialConnections(ctx *pulumi.Context, config *config.Config) error {
	for _, social := range []string{"google-oauth2"} {
		_, err := auth0.NewConnection(ctx, formatName(ctx, fmt.Sprintf("auth0-%s", social)), &auth0.ConnectionArgs{
			Strategy: pulumi.String(social),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
