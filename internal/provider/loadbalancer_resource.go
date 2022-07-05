package provider

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kubeberth/kubeberth-go"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
//var _ tfsdk.ResourceType = loadbalancerResourceType{}
//var _ tfsdk.Resource = loadbalancerResource{}
//var _ tfsdk.ResourceWithImportState = loadbalancerResource{}

type loadbalancerResourceType struct{}

func (t loadbalancerResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server..
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "name",
				Type: types.StringType,
				Required: true,
			},
			"backends": {
				MarkdownDescription: "backends",
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"server": {
						Type: types.StringType,
						Required: true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"ports": {
				MarkdownDescription: "ports",
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type: types.StringType,
						Required: true,
					},
					"protocol": {
						Type: types.StringType,
						Required: true,
					},
					"port": {
						Type: types.Int64Type,
						Required: true,
					},
					"target_port": {
						Type: types.Int64Type,
						Required: true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (t loadbalancerResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return loadbalancerResource{
		provider: provider,
	}, diags
}

type destinationData struct {
	Server types.String `tfsdk:"server"`
}

type portData struct {
	Name       types.String `tfsdk:"name"`
	Protocol   types.String `tfsdk:"protocol"`
	Port       types.Int64  `tfsdk:"port"`
	TargetPort types.Int64  `tfsdk:"target_port"`
}

type loadbalancerResourceData struct {
	Name       types.String      `tfsdk:"name"`
	Backends   []destinationData `tfsdk:"backends"`
	Ports      []portData        `tfsdk:"ports"`
}

type loadbalancerResource struct {
	provider provider
}

func newRequestLoadBalancer(data *loadbalancerResourceData) *kubeberth.RequestLoadBalancer {
	

	backends := []kubeberth.Destination{}
	for _, destination := range data.Backends {
		backends = append(backends, kubeberth.Destination{Server: destination.Server.Value})
	}

	ports := []kubeberth.Port{}
	for _, port := range data.Ports {
		ports = append(ports, kubeberth.Port{
			Name:       port.Name.Value,
			Protocol:   (corev1.Protocol)(port.Protocol.Value),
			Port:       (int32)(port.Port.Value),
			TargetPort: intstr.FromInt((int)(port.TargetPort.Value)),
		})
	}

	loadbalancer := &kubeberth.RequestLoadBalancer{
		Name: data.Name.Value,
		Backends: backends,
		Ports: ports,
	}

	return loadbalancer
}

func (r loadbalancerResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data loadbalancerResourceData

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

	requestLoadBalancer := newRequestLoadBalancer(&data)
	responseLoadBalancer, err := r.provider.client.CreateLoadBalancer(ctx, requestLoadBalancer)
	tflog.Trace(ctx, fmt.Sprintf("loadbalancer: %+v\n", responseLoadBalancer))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create loadbalancer, got error: %s", err))
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

func (r loadbalancerResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data loadbalancerResourceData

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

	responseLoadBalancer, err := r.provider.client.GetLoadBalancer(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("loadbalancer: %+v\n", responseLoadBalancer))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read loadbalancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "read a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r loadbalancerResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data loadbalancerResourceData

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

	requestLoadBalancer := newRequestLoadBalancer(&data)
	responseLoadBalancer, err := r.provider.client.UpdateLoadBalancer(ctx, data.Name.Value, requestLoadBalancer)
	tflog.Trace(ctx, fmt.Sprintf("loadbalancer: %+v\n", responseLoadBalancer))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update loadbalancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated a resource")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r loadbalancerResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data loadbalancerResourceData

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

	ok, err := r.provider.client.DeleteLoadBalancer(ctx, data.Name.Value)
	tflog.Trace(ctx, fmt.Sprintf("loadbalancer: %+v\n", ok))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete loadbalancer, got error: %s", err))
		return
	}
}

func (r loadbalancerResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
