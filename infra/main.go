package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	GitHubActionsTokenURL string = "https://token.actions.githubusercontent.com"
	ProjectID             string = "espressolog"
	ProjectName           string = "espresso-chronicle"
)

func formatName(ctx *pulumi.Context, thing string) string {
	return fmt.Sprintf("%s-%s-%s", ProjectName, ctx.Stack(), thing)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get stack specific configuration
		config := config.New(ctx, ProjectName)

		// Create Auth0 Resources
		clientID, err := deployAuth0Application(ctx, config)
		if err != nil {
			return err
		}

		clientIDArray := pulumi.ToStringArrayOutput([]pulumi.StringOutput{*clientID})

		err = deployAuth0SocialConnections(ctx, config, &clientIDArray)
		if err != nil {
			return err
		}

		// Create the deployment identity
		err = deployCIServiceAccount(ctx)
		if err != nil {
			return err
		}

		// Create the Firestore database
		err = deployFirestore(ctx, config)
		if err != nil {
			return err
		}

		// Create Docker artifact registry for API
		err = deployAPIArtifactRegistry(ctx, config)
		if err != nil {
			return err
		}

		return nil
	})
}
