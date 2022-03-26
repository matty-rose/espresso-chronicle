package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func deployAPIArtifactRegistry(ctx *pulumi.Context, config *config.Config) error {
	_, err := artifactregistry.NewRepository(ctx, "api", &artifactregistry.RepositoryArgs{
		Location:     pulumi.String(config.Require("region")),
		RepositoryId: pulumi.String("api"),
		Format:       pulumi.String("DOCKER"),
	})
	if err != nil {
		return err
	}

	return nil
}
