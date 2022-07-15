package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kubeberth/kubeberth-go"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
//var _ tfsdk.ResourceType = serverResourceType{}
//var _ tfsdk.Resource = serverResource{}
//var _ tfsdk.ResourceWithImportState = serverResource{}

type serverResourceType struct{}

func (t serverResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "name",
				Type:                types.StringType,
				Required:            true,
			},
			"running": {
				MarkdownDescription: "running",
				Type:                types.BoolType,
				Optional:            true,
			},
			"cpu": {
				MarkdownDescription: "cpu",
				Type:                types.Int64Type,
				Required:            true,
			},
			"memory": {
				MarkdownDescription: "memory",
				Type:                types.StringType,
				Required:            true,
			},
			"mac_address": {
				MarkdownDescription: "mac_address",
				Type:                types.StringType,
				Optional:            true,
			},
			"hostname": {
				MarkdownDescription: "hostname",
				Type:                types.StringType,
				Required:            true,
			},
			"hosting": {
				MarkdownDescription: "hosting",
				Type:                types.StringType,
				Optional:            true,
			},
			"disks": {
				MarkdownDescription: "disks",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"isoimage": {
				MarkdownDescription: "isoimage",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
				}),
			},
			"cloudinit": {
				MarkdownDescription: "cloudinit",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
				}),
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

type diskData struct {
	Name types.String `tfsdk:"name"`
}

type isoimageData struct {
	Name types.String `tfsdk:"name"`
}

type cloudinitData struct {
	Name types.String `tfsdk:"name"`
}

type serverResourceData struct {
	Name       types.String   `tfsdk:"name"`
	Running    types.Bool     `tfsdk:"running"`
	CPU        types.Int64    `tfsdk:"cpu"`
	Memory     types.String   `tfsdk:"memory"`
	MACAddress types.String   `tfsdk:"mac_address"`
	Hostname   types.String   `tfsdk:"hostname"`
	Hosting    types.String   `tfsdk:"hosting"`
	Disks      []diskData     `tfsdk:"disks"`
	ISOImage   *isoimageData  `tfsdk:"isoimage"`
	CloudInit  *cloudinitData `tfsdk:"cloudinit"`
}

type serverResource struct {
	provider provider
}

func newRequestServer(data *serverResourceData) *kubeberth.RequestServer {
	cpu    := resource.MustParse(strconv.FormatInt(data.CPU.Value, 10))
	memory := resource.MustParse(data.Memory.Value)
	disks  := []kubeberth.AttachedDisk{}
	for _, disk := range data.Disks {
		disks = append(disks, kubeberth.AttachedDisk{Name: disk.Name.Value})
	}

	server := &kubeberth.RequestServer{
		Name:       data.Name.Value,
		Running:    data.Running.Value,
		CPU:        &cpu,
		Memory:     &memory,
		MACAddress: data.MACAddress.Value,
		Hostname:   data.Hostname.Value,
		Hosting:    data.Hosting.Value,
		Disks:      disks,
	}

	if data.ISOImage != nil {
		if !data.ISOImage.Name.Null {
			server.ISOImage = &kubeberth.AttachedISOImage{Name: data.ISOImage.Name.Value}
		}
	}

	if data.CloudInit != nil {
		if !data.CloudInit.Name.Null {
			server.CloudInit = &kubeberth.AttachedCloudInit{Name: data.CloudInit.Name.Value}
		}
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

	requestServer := newRequestServer(&data)
	responseServer, err := r.provider.client.CreateServer(ctx, requestServer)
	tflog.Trace(ctx, fmt.Sprintf("server: %+v\n", responseServer))
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

	responseServer, err := r.provider.client.GetServer(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("server: %+v\n", responseServer))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a resource")

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

	requestServer := newRequestServer(&data)
	responseServer, err := r.provider.client.UpdateServer(ctx, data.Name.Value, requestServer)
	tflog.Trace(ctx, fmt.Sprintf("server: %+v\n", responseServer))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a resource")

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

	ok, err := r.provider.client.DeleteServer(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("server: %+v\n", ok))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}
}

func (r serverResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
