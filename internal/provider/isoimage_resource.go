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
//var _ tfsdk.ResourceType = isoimageResourceType{}
//var _ tfsdk.Resource = isoimageResource{}
//var _ tfsdk.ResourceWithImportState = isoimageResource{}

type isoimageResourceType struct{}

func (t isoimageResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "name",
				Type:                types.StringType,
				Required:            true,
			},
			"size": {
				MarkdownDescription: "size",
				Type:                types.StringType,
				Required:            true,
			},
			"repository": {
				MarkdownDescription: "repository",
				Type:                types.StringType,
				Required:            true,
			},
		},
	}, nil
}

func (t isoimageResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return isoimageResource{
		provider: provider,
	}, diags
}

type isoimageResourceData struct {
	Name       types.String `tfsdk:"name"`
	Size       types.String `tfsdk:"size"`
	Repository types.String `tfsdk:"repository"`
}

type isoimageResource struct {
	provider provider
}

func createNewISOImage(data *isoimageResourceData) *kubeberth.RequestISOImage {
	isoimage := &kubeberth.RequestISOImage{
		Name:       data.Name.Value,
		Size:       data.Size.Value,
		Repository: data.Repository.Value,
	}

	return isoimage
}

func (r isoimageResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data isoimageResourceData

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

	newISOImage := createNewISOImage(&data)
	createdISOImage, err := r.provider.client.CreateISOImage(ctx, newISOImage)
	tflog.Trace(ctx, fmt.Sprintf("isoimage: %+v\n", createdISOImage))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create isoimage, got error: %s", err))
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

func (r isoimageResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data isoimageResourceData

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

	isoimage, err := r.provider.client.GetISOImage(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("isoimage: %+v\n", isoimage))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read isoimage, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r isoimageResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data isoimageResourceData

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

	newISOImage := createNewISOImage(&data)
	updatedISOImage, err := r.provider.client.UpdateISOImage(ctx, data.Name.Value, newISOImage)
	tflog.Trace(ctx, fmt.Sprintf("isoimage: %+v\n", updatedISOImage))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update isoimage, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r isoimageResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data isoimageResourceData

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

	ok, err := r.provider.client.DeleteISOImage(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("isoimage: %+v\n", ok))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete isoimage, got error: %s", err))
		return
	}
}

func (r isoimageResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
