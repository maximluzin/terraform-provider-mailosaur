package provider

import (
	"context"
	"fmt"
	"github.com/mailosaur/mailosaur-go"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = mailosaurServerResourceType{}
var _ tfsdk.Resource = mailosaurServerResource{}

type mailosaurServerResourceType struct{}

func (t mailosaurServerResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Virtual server",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "Name of virtual server",
				Optional:            false,
				Required:            true,
				Type:                types.StringType,
			},
			"id": {
				Computed:            true,
				MarkdownDescription: "Id of virtual server",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"email": {
				Computed:            true,
				MarkdownDescription: "Email address for the server",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"password": {
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for use with SMTP and POP3",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (t mailosaurServerResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return mailosaurServerResource{
		provider: provider,
	}, diags
}

type mailosaurServerResourceData struct {
	Name     types.String `tfsdk:"name"`
	Id       types.String `tfsdk:"id"`
	Password types.String `tfsdk:"password"`
	Email    types.String `tfsdk:"email"`
}

type mailosaurServerResource struct {
	provider provider
}

func (r mailosaurServerResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data mailosaurServerResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	options := mailosaur.ServerCreateOptions{
		Name: data.Name.Value,
	}

	server, err := r.provider.client.Servers.Create(options)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create server, got error: %s", err))
		return
	}
	data.Name = types.String{Value: server.Name}
	data.Id = types.String{Value: server.Id}

	password, err := r.provider.client.Servers.GetPassword(server.Id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read password of server, got error: %s", err))

	}
	data.Password = types.String{Value: password}

	email := r.provider.client.Servers.GenerateEmailAddress(server.Id)
	data.Email = types.String{Value: email}

	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r mailosaurServerResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data mailosaurServerResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.provider.client.Servers.Get(data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server, got error: %s", err))
		return
	}

	if data.Password.Null {
		password, err := r.provider.client.Servers.GetPassword(server.Id)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read password of server, got error: %s", err))

		}
		data.Password = types.String{Value: password}
	}
	if data.Email.Null {
		email := r.provider.client.Servers.GenerateEmailAddress(server.Id)
		data.Email = types.String{Value: email}
	}
	data.Name = types.String{Value: server.Name}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r mailosaurServerResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data mailosaurServerResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	retrievedServer, err := r.provider.client.Servers.Get(data.Id.Value)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server, got error: %s", err))
		return
	}

	retrievedServer.Name = data.Name.Value

	_, err = r.provider.client.Servers.Update(retrievedServer.Id, retrievedServer)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server, got error: %s", err))
		return
	}
	if data.Password.Null {
		password, err := r.provider.client.Servers.GetPassword(retrievedServer.Id)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read password of server, got error: %s", err))

		}
		data.Password = types.String{Value: password}
	}
	if data.Email.Null {
		email := r.provider.client.Servers.GenerateEmailAddress(retrievedServer.Id)
		data.Email = types.String{Value: email}
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r mailosaurServerResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data mailosaurServerResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.client.Servers.Delete(data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
		return
	}
}
