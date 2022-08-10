package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mailosaur/mailosaur-go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}

type provider struct {
	client     *mailosaur.MailosaurClient
	configured bool
	version    string
}

type providerData struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var apiKey string
	if data.ApiKey.Null {
		apiKey = os.Getenv("MAILOSAUR_API_KEY")
	} else {
		apiKey = data.ApiKey.Value
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Unable to find MAILOSAUR_API_KEY",
			"Api Key cannot be an empty string",
		)
	}

	p.client = mailosaur.New(apiKey)
	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"mailosaur_server": mailosaurServerResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Provider for resource management in Mailosaur",
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "Api key for access to Mailosaur. You can use env var MAILOSAUR_API_KEY for set this attribute",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
