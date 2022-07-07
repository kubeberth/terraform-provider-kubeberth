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
//var _ tfsdk.ResourceType = diskResourceType{}
//var _ tfsdk.Resource = diskResource{}
//var _ tfsdk.ResourceWithImportState = diskResource{}

type diskResourceType struct{}

func (t diskResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"source": {
				MarkdownDescription: "source",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"archive": {
						Type:     types.StringType,
						Optional: true,
					},
					"disk": {
						Type:     types.StringType,
						Optional: true,
					},
				}),
			},
		},
	}, nil
}

func (t diskResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return diskResource{
		provider: provider,
	}, diags
}

type sourceData struct {
	Archive types.String `tfsdk:"archive"`
	Disk    types.String `tfsdk:"disk"`
}

type diskResourceData struct {
	Name   types.String `tfsdk:"name"`
	Size   types.String `tfsdk:"size"`
	Source *sourceData  `tfsdk:"source"`
}

type diskResource struct {
	provider provider
}

func newRequestDisk(data *diskResourceData) *kubeberth.RequestDisk {
	requestDisk := &kubeberth.RequestDisk{
		Name:   data.Name.Value,
		Size:   data.Size.Value,
		Source: &kubeberth.AttachedSource{},
	}

	if data.Source != nil {
		if !data.Source.Archive.Null {
			requestDisk.Source.Archive = &kubeberth.AttachedArchive{
				Name: data.Source.Archive.Value,
			}
		}
		if !data.Source.Disk.Null {
			requestDisk.Source.Disk = &kubeberth.AttachedDisk{
				Name: data.Source.Disk.Value,
			}
		}
	}

	return requestDisk
}

func (r diskResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data diskResourceData

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

	requestDisk := newRequestDisk(&data)
	responseDisk, err := r.provider.client.CreateDisk(ctx, requestDisk)
	tflog.Trace(ctx, fmt.Sprintf("disk: %+v\n", responseDisk))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create disk, got error: %s", err))
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

func (r diskResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data diskResourceData

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

	responseDisk, err := r.provider.client.GetDisk(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("disk: %+v\n", responseDisk))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read disk, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r diskResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data diskResourceData

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

	requestDisk := newRequestDisk(&data)
	responseDisk, err := r.provider.client.UpdateDisk(ctx, data.Name.Value, requestDisk)
	tflog.Trace(ctx, fmt.Sprintf("disk: %+v\n", responseDisk))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update disk, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r diskResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data diskResourceData

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

	ok, err := r.provider.client.DeleteDisk(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("disk: %+v\n", ok))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete disk, got error: %s", err))
		return
	}
}

func (r diskResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
