package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/appengine"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/iam"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	GitHubActionsTokenURL string = "https://token.actions.githubusercontent.com"
	ProjectId             string = "espressolog"
	ProjectName           string = "espresso-chronicle"
)

func formatName(ctx *pulumi.Context, thing string) string {
	return fmt.Sprintf("%s-%s-%s", ProjectName, ctx.Stack(), thing)
}

func deployCIServiceAccount(ctx *pulumi.Context) error {
	sa, err := serviceaccount.NewAccount(ctx, formatName(ctx, "cicd"), &serviceaccount.AccountArgs{
		AccountId:   pulumi.String(formatName(ctx, "cicd")),
		DisplayName: pulumi.String("CICD Service Account"),
	})
	if err != nil {
		return err
	}

	saEmail := sa.Email.ApplyT(func(dnsName string) string {
		return "serviceAccount:" + dnsName
	}).(pulumi.StringOutput)

	_, err = projects.NewIAMMember(ctx, "cicd-iam", &projects.IAMMemberArgs{
		Member:  saEmail,
		Project: pulumi.String(ProjectId),
		Role:    pulumi.String("roles/owner"),
	})
	if err != nil {
		return err
	}

	gha := pulumi.String("cicd")

	pool, err := iam.NewWorkloadIdentityPool(ctx, "github-actions", &iam.WorkloadIdentityPoolArgs{
		WorkloadIdentityPoolId: gha,
	})
	if err != nil {
		return err
	}

	_, err = iam.NewWorkloadIdentityPoolProvider(ctx, "github-actions", &iam.WorkloadIdentityPoolProviderArgs{
		WorkloadIdentityPoolId:         pool.WorkloadIdentityPoolId,
		WorkloadIdentityPoolProviderId: gha,
		AttributeMapping: pulumi.StringMap{
			"google.subject":       pulumi.String("assertion.sub"),
			"attribute.actor":      pulumi.String("assertion.actor"),
			"attribute.aud":        pulumi.String("assertion.aud"),
			"attribute.repository": pulumi.String("assertion.repository"),
		},
		Oidc: &iam.WorkloadIdentityPoolProviderOidcArgs{
			IssuerUri: pulumi.String(GitHubActionsTokenURL),
		},
	})
	if err != nil {
		return err
	}

	member := pool.Name.ApplyT(func(name string) string {
		return "principalSet://iam.googleapis.com/" + name + "/attribute.repository/matty-rose/espresso-chronicle"
	}).(pulumi.StringOutput)

	_, err = serviceaccount.NewIAMMember(ctx, "github-actions", &serviceaccount.IAMMemberArgs{
		Member:           member,
		Role:             pulumi.String("roles/iam.workloadIdentityUser"),
		ServiceAccountId: sa.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func deployFirestore(ctx *pulumi.Context, config *config.Config) error {
	_, err := appengine.NewApplication(ctx, formatName(ctx, "db"), &appengine.ApplicationArgs{
		Project:      pulumi.String(ProjectId),
		LocationId:   pulumi.String(config.Require("region")),
		DatabaseType: pulumi.String("CLOUD_FIRESTORE"),
	})
	if err != nil {
		return err
	}

	return nil
}

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

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get stack specific configuration
		config := config.New(ctx, ProjectName)

		// Create the deployment identity
		err := deployCIServiceAccount(ctx)
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
