package provider

import (
	"context"
	"fmt"

	"github.com/kubeberth/kubeberth-go"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
//var _ tfsdk.ResourceType = cloudinitResourceType{}
//var _ tfsdk.Resource = cloudinitResource{}
//var _ tfsdk.ResourceWithImportState = cloudinitResource{}

type cloudinitResourceType struct{}

func (t cloudinitResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "name",
				Type: types.StringType,
				Required: true,
			},
			"user_data": {
				MarkdownDescription: "user_data",
				Type: types.StringType,
				Optional: true,
			},
			"network_data": {
				MarkdownDescription: "network_data",
				Type: types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (t cloudinitResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return cloudinitResource{
		provider: provider,
	}, diags
}

type cloudinitResourceData struct {
	Name        types.String `tfsdk:"name"`
	UserData    types.String `tfsdk:"user_data"`
	NetworkData types.String `tfsdk:"network_data"`
}

type cloudinitResource struct {
	provider provider
}

func createNewCloudInit(data *cloudinitResourceData) *kubeberth.CloudInit {
	cloudinit := &kubeberth.CloudInit{
		Name: data.Name.Value,
	}

	if !data.UserData.Null {
		cloudinit.UserData = data.UserData.Value
	}
	if !data.NetworkData.Null {
		cloudinit.NetworkData = data.NetworkData.Value
	}

	return cloudinit
}

func (r cloudinitResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data cloudinitResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.CreateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	newCloudInit := createNewCloudInit(&data)
	createdCloudInit, err := r.provider.client.CreateCloudInit(ctx, newCloudInit)
	tflog.Trace(ctx, fmt.Sprintf("cloudinit: %+v\n", createdCloudInit))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create cloudinit, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	//data.Id = types.String{Value: "example-id"}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information

	tflog.Trace(ctx, "created a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cloudinitResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data cloudinitResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.ReadExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	cloudinit, err := r.provider.client.GetCloudInit(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("cloudinit: %+v\n", cloudinit))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read cloudinit, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cloudinitResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data cloudinitResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.UpdateExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	newCloudInit := createNewCloudInit(&data)
	updatedCloudInit, err := r.provider.client.UpdateCloudInit(ctx, data.Name.Value, newCloudInit)
	tflog.Trace(ctx, fmt.Sprintf("cloudinit: %+v\n", updatedCloudInit))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update cloudinit, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r cloudinitResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data cloudinitResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// example, err := d.provider.client.DeleteExample(...)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }

	ok, err := r.provider.client.DeleteCloudInit(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("cloudinit: %+v\n", ok))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete cloudinit, got error: %s", err))
		return
	}
}

func (r cloudinitResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
