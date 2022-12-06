package archive

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*archiveFileResource)(nil)

func NewArchiveFileResource() resource.Resource {
	return &archiveFileResource{}
}

type archiveFileResource struct{}

func (d *archiveFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:        `**NOTE**: This resource is deprecated, use data source instead.`,
		DeprecationMessage: `**NOTE**: This resource is deprecated, use data source instead.`,
		Blocks: map[string]schema.Block{
			"source": schema.SetNestedBlock{
				Description: "Specifies attributes of a single source file to include into the archive.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"content": schema.StringAttribute{
							Description: "Add this content to the archive with `filename` as the filename.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(
									fwpath.MatchRoot("source_file"),
									fwpath.MatchRoot("source_dir"),
									fwpath.MatchRoot("source_content"),
									fwpath.MatchRoot("source_content_filename"),
								),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"filename": schema.StringAttribute{
							Description: "Set this as the filename when declaring a `source`.",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The sha1 checksum hash of the output.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of archive to generate. NOTE: `zip` is supported.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_content": schema.StringAttribute{
				Description: "Add only this content to the archive with `source_content_filename` as the filename.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_content_filename": schema.StringAttribute{
				Description: "Set this as the filename when using `source_content`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_file": schema.StringAttribute{
				Description: "Package this file into the archive.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_dir"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_dir": schema.StringAttribute{
				Description: "Package entire contents of this directory into the archive.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"excludes": schema.SetAttribute{
				Description: "Specify files to ignore when reading the `source_dir`.",
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"output_path": schema.StringAttribute{
				Description: "The output of the archive file.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"output_size": schema.Int64Attribute{
				Description: "The byte size of the output archive file.",
				Computed:    true,
			},
			"output_sha": schema.StringAttribute{
				Description: "The SHA1 checksum of output archive file.",
				Computed:    true,
			},
			"output_base64sha256": schema.StringAttribute{
				Description: "The base64-encoded SHA256 checksum of output archive file.",
				Computed:    true,
			},
			"output_md5": schema.StringAttribute{
				Description: "The MD5 checksum of output archive file.",
				Computed:    true,
			},
			"output_file_mode": schema.StringAttribute{
				Description: "String that specifies the octal file mode for all archived files. For example: `\"0666\"`. " +
					"Setting this will ensure that cross platform usage of this module will not vary the modes of archived " +
					"files (and ultimately checksums) resulting in more deterministic behavior.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
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
