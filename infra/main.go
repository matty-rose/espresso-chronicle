package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	ProjectId   string = "espressolog"
	ProjectName string = "espresso-chronicle"
)

func formatName(ctx *pulumi.Context, thing string) string {
	return fmt.Sprintf("%s-%s-%s", ProjectName, ctx.Stack(), thing)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get stack specific configuration
		config := config.New(ctx, ProjectName)

		// Create the Firestore database
		_, err := appengine.NewApplication(ctx, formatName(ctx, "db"), &appengine.ApplicationArgs{
			Project:      pulumi.String(ProjectId),
			LocationId:   pulumi.String(config.Require("region")),
			DatabaseType: pulumi.String("CLOUD_FIRESTORE"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
