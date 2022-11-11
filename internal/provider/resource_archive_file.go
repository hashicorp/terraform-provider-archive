package archive

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*archiveFileResource)(nil)

func NewArchiveFileResource() resource.Resource {
	return &archiveFileResource{}
}

type archiveFileResource struct{}

func (d *archiveFileResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:        `**NOTE**: This resource is deprecated, use data source instead.`,
		DeprecationMessage: `**NOTE**: This resource is deprecated, use data source instead.`,
		Blocks: map[string]tfsdk.Block{
			"source": {
				Description: "Specifies attributes of a single source file to include into the archive.",
				Attributes: map[string]tfsdk.Attribute{
					"content": {
						Description: "Add this content to the archive with `filename` as the filename.",
						Type:        types.StringType,
						Required:    true,
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(
								fwpath.MatchRoot("source_file"),
								fwpath.MatchRoot("source_dir"),
								fwpath.MatchRoot("source_content"),
								fwpath.MatchRoot("source_content_filename"),
							),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.RequiresReplace(),
						},
					},
					"filename": {
						Description: "Set this as the filename when declaring a `source`.",
						Type:        types.StringType,
						Required:    true,
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.RequiresReplace(),
						},
					},
				},
				NestingMode: tfsdk.BlockNestingModeSet,
			},
		},
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "The sha1 checksum hash of the output.",
				Type:        types.StringType,
				Computed:    true,
			},
			"type": {
				Description: "The type of archive to generate. NOTE: `zip` is supported.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"source_content": {
				Description: "Add only this content to the archive with `source_content_filename` as the filename.",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
					),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"source_content_filename": {
				Description: "Set this as the filename when using `source_content`.",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
					),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"source_file": {
				Description: "Package this file into the archive.",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(
						fwpath.MatchRoot("source_dir"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"source_dir": {
				Description: "Package entire contents of this directory into the archive.",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"excludes": {
				Description: "Specify files to ignore when reading the `source_dir`.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"output_path": {
				Description: "The output of the archive file.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"output_size": {
				Description: "The byte size of the output archive file.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"output_sha": {
				Description: "The SHA1 checksum of output archive file.",
				Type:        types.StringType,
				Computed:    true,
			},
			"output_base64sha256": {
				Description: "The base64-encoded SHA256 checksum of output archive file.",
				Type:        types.StringType,
				Computed:    true,
			},
			"output_md5": {
				Description: "The MD5 checksum of output archive file.",
				Type:        types.StringType,
				Computed:    true,
			},
			"output_file_mode": {
				Description: "String that specifies the octal file mode for all archived files. For example: `\"0666\"`. " +
					"Setting this will ensure that cross platform usage of this module will not vary the modes of archived " +
					"files (and ultimately checksums) resulting in more deterministic behavior.",
				Type:     types.StringType,
				Optional: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (d *archiveFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model fileModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateModel(ctx, &model)...)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (d *archiveFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model fileModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateModel(ctx, &model)...)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func updateModel(ctx context.Context, model *fileModel) diag.Diagnostics {
	var diags diag.Diagnostics
	outputPath := model.OutputPath.ValueString()

	outputDirectory := path.Dir(outputPath)
	if outputDirectory != "" {
		if _, err := os.Stat(outputDirectory); err != nil {
			if err := os.MkdirAll(outputDirectory, 0755); err != nil {
				diags.AddError(
					"Output path error",
					fmt.Sprintf("error creating output path: %s", err),
				)
				return diags
			}
		}
	}

	if err := archive(ctx, *model); err != nil {
		diags.AddError(
			"Archive creation error",
			fmt.Sprintf("error creating archive: %s", err),
		)
		return diags
	}

	// Generate archived file stats
	fi, err := os.Stat(outputPath)
	if err != nil {
		diags.AddError(
			"Archive output error",
			fmt.Sprintf("error reading output: %s", err),
		)
		return diags
	}

	sha1, base64sha256, md5, err := genFileShas(outputPath)
	if err != nil {
		diags.AddError(
			"Hash generation error",
			fmt.Sprintf("error generating hashed: %s", err),
		)
		return diags
	}

	model.OutputSha = types.StringValue(sha1)
	model.OutputBase64Sha256 = types.StringValue(base64sha256)
	model.OutputMd5 = types.StringValue(md5)
	model.OutputSize = types.Int64Value(fi.Size())

	model.ID = types.StringValue(sha1)

	return diags
}

func (d *archiveFileResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (d *archiveFileResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (d *archiveFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}
