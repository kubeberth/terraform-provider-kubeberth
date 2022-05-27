package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kubeberth/kubeberth-go"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = serverResourceType{}
var _ tfsdk.Resource = serverResource{}
var _ tfsdk.ResourceWithImportState = serverResource{}

type serverResourceType struct{}

func (t serverResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "running",
				Type: types.StringType,
				Required: true,
			},
			"running": {
				MarkdownDescription: "running",
				Type: types.BoolType,
				Required: true,
			},
			"cpu": {
				MarkdownDescription: "cpu",
				Type: types.StringType,
				Required: true,
			},
			"memory": {
				MarkdownDescription: "memory",
				Type: types.StringType,
				Required: true,
			},
			"mac_address": {
				MarkdownDescription: "mac_address",
				Type: types.StringType,
				Required: true,
			},
			"hostname": {
				MarkdownDescription: "hostname",
				Type: types.StringType,
				Required: true,
			},
			"disk": {
				MarkdownDescription: "disk",
				Type: types.StringType,
				Required: true,
			},
			"cloudinit": {
				MarkdownDescription: "cloudinit",
				Type: types.StringType,
				Required: true,
			},
		},
	}, nil
}

func (t serverResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return serverResource{
		provider: provider,
	}, diags
}

type serverResourceData struct {
	Name       string `tfsdk:"name"`
	Running    bool   `tfsdk:"running"`
	CPU        string `tfsdk:"cpu"`
	Memory     string `tfsdk:"memory"`
	MACAddress string `tfsdk:"mac_address"`
	HostName   string `tfsdk:"hostname"`
	Disk       string `tfsdk:"disk"`
	CloudInit  string `tfsdk:"cloudinit"`
}

type serverResource struct {
	provider provider
}

func createNewServer(data *serverResourceData) *kubeberth.Server {
	server := &kubeberth.Server{
		Name: data.Name,
		Running: strconv.FormatBool(data.Running),
		CPU: data.CPU,
		Memory: data.Memory,
		MACAddress: data.MACAddress,
		HostName: data.HostName,
		Disk: &kubeberth.AttachedDisk{
			Name: data.Disk,
		},
		CloudInit: &kubeberth.AttachedCloudInit{
			Name: data.CloudInit,
		},
	}

	return server
}

func (r serverResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data serverResourceData

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

	newServer := createNewServer(&data)
	createdServer, err := r.provider.client.CreateServer(ctx, newServer)
	tflog.Trace(ctx, fmt.Sprintf("server: %+v\n", createdServer))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create server, got error: %s", err))
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

func (r serverResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data serverResourceData

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

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r serverResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data serverResourceData

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

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r serverResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data serverResourceData

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
	ok, err := r.provider.client.DeleteServer(ctx, data.Name)
	tflog.Trace(ctx, fmt.Sprintf("server: %+v\n", ok))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}
}

func (r serverResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
