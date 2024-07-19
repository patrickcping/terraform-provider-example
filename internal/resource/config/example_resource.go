package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &exampleResource{}
)

// ExampleResource is a helper function to simplify the provider implementation.
func ExampleResource() resource.Resource {
	return &exampleResource{}
}

// exampleResource is the resource implementation.
type exampleResource struct {
}

type exampleResourceModel struct {
	Id              types.String `tfsdk:"id"`
	StringVal       types.String `tfsdk:"string_val"`
	SingleNestedVal types.Object `tfsdk:"single_nested_val"`
}

type exampleSingleNestedValResourceModel struct {
	StringVal        types.String `tfsdk:"string_val"`
	StringValStatic  types.String `tfsdk:"string_val_static"`
	StringValDynamic types.String `tfsdk:"string_val_dynamic"`
}

// GetSchema defines the schema for the resource.
func (r *exampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Example resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed id",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"string_val": schema.StringAttribute{
				Description: "Optional string attribute",
				Optional:    true,
			},
			"single_nested_val": schema.SingleNestedAttribute{
				Description: "Optional single nested object",
				Optional:    true,

				Attributes: map[string]schema.Attribute{
					"string_val": schema.StringAttribute{
						Description: "Required string attribute",
						Required:    true,
					},

					"string_val_static": schema.StringAttribute{
						Description: "Computed single nested attribute object, using \"UseStateForUnknown()\" plan modifier",
						Computed:    true,

						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},

					"string_val_dynamic": schema.StringAttribute{
						Description: "Computed single nested attribute object, without using \"UseStateForUnknown()\" plan modifier",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *exampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state exampleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stub API
	data := plan
	data.SetId()

	resp.Diagnostics.Append(state.toState(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *exampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, state exampleResourceModel

	// Stub API
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data = state

	resp.Diagnostics.Append(state.toState(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state exampleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stub API
	data := plan
	data.SetId()

	resp.Diagnostics.Append(state.toState(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// No backend so no logic needed
func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *exampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *exampleResourceModel) SetId() {
	m.Id = types.StringValue("id")
}

func (m *exampleResourceModel) toState(ctx context.Context, data exampleResourceModel) (diags diag.Diagnostics) {
	m.Id = data.Id
	m.StringVal = data.StringVal

	singleNestedObjectTyp := map[string]attr.Type{
		"string_val":         types.StringType,
		"string_val_static":  types.StringType,
		"string_val_dynamic": types.StringType,
	}

	var d diag.Diagnostics

	singleNestedPlan := exampleSingleNestedValResourceModel{
		StringVal:        types.StringValue("required string"),
		StringValStatic:  types.StringValue("service static string"),
		StringValDynamic: types.StringValue("service dynamic string"),
	}

	if !data.SingleNestedVal.IsNull() {
		d = data.SingleNestedVal.As(context.Background(), &singleNestedPlan, basetypes.ObjectAsOptions{})
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		singleNestedPlan.StringValStatic = types.StringValue("service static string")
		singleNestedPlan.StringValDynamic = types.StringValue("service dynamic string, updated because this block is defined")
	}

	m.SingleNestedVal, d = types.ObjectValueFrom(ctx, singleNestedObjectTyp, singleNestedPlan)
	diags = append(diags, d...)

	return
}
