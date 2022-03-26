package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func deployFirestore(ctx *pulumi.Context, config *config.Config) error {
	_, err := appengine.NewApplication(ctx, formatName(ctx, "db"), &appengine.ApplicationArgs{
		Project:      pulumi.String(ProjectID),
		LocationId:   pulumi.String(config.Require("region")),
		DatabaseType: pulumi.String("CLOUD_FIRESTORE"),
	})
	if err != nil {
		return err
	}

	return nil
}
