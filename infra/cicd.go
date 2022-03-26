package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/iam"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

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
		Project: pulumi.String(ProjectID),
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
