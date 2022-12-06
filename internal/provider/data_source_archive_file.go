package archive

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*archiveFileDataSource)(nil)

func NewArchiveFileDataSource() datasource.DataSource {
	return &archiveFileDataSource{}
}

type archiveFileDataSource struct{}

func (d *archiveFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates an archive from content, a file, or directory of files.",
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
						},
						"filename": schema.StringAttribute{
							Description: "Set this as the filename when declaring a `source`.",
							Required:    true,
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
			},
			"output_path": schema.StringAttribute{
				Description: "The output of the archive file.",
				Required:    true,
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
			},
		},
	}
}

func archive(ctx context.Context, model fileModel) error {
	archiveType := model.Type.ValueString()
	outputPath := model.OutputPath.ValueString()

	archiver := getArchiver(archiveType, outputPath)
	if archiver == nil {
		return fmt.Errorf("archive type not supported: %s", archiveType)
	}

	outputFileMode := model.OutputFileMode.ValueString()
	if outputFileMode != "" {
		archiver.SetOutputFileMode(outputFileMode)
	}

	switch true {
	case !model.SourceDir.IsNull():
		excludeList := make([]string, len(model.Excludes.Elements()))

		if !model.Excludes.IsNull() {
			var elements []types.String
			model.Excludes.ElementsAs(ctx, &elements, false)

			for i, elem := range elements {
				excludeList[i] = elem.ValueString()
			}
		}

		if err := archiver.ArchiveDir(model.SourceDir.ValueString(), excludeList); err != nil {
			return fmt.Errorf("error archiving directory: %s", err)
		}
	case !model.SourceFile.IsNull():
		if err := archiver.ArchiveFile(model.SourceFile.ValueString()); err != nil {
			return fmt.Errorf("error archiving file: %s", err)
		}
	case !model.SourceContentFilename.IsNull():
		content := model.SourceContent.ValueString()

		if err := archiver.ArchiveContent([]byte(content), model.SourceContentFilename.ValueString()); err != nil {
			return fmt.Errorf("error archiving content: %s", err)
		}
	case !model.Source.IsNull():
		content := make(map[string][]byte)

		var elements []sourceModel
		model.Source.ElementsAs(ctx, &elements, false)

		for _, elem := range elements {
			content[elem.Filename.ValueString()] = []byte(elem.Content.ValueString())
		}

		if err := archiver.ArchiveMultiple(content); err != nil {
			return fmt.Errorf("error archiving content: %s", err)
		}
	}

	return nil
}

func (d *archiveFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model fileModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	outputPath := model.OutputPath.ValueString()

	outputDirectory := path.Dir(outputPath)
	if outputDirectory != "" {
		if _, err := os.Stat(outputDirectory); err != nil {
			if err := os.MkdirAll(outputDirectory, 0755); err != nil {
				resp.Diagnostics.AddError(
					"Output path error",
					fmt.Sprintf("error creating output path: %s", err),
				)
				return
			}
		}
	}

	if err := archive(ctx, model); err != nil {
		resp.Diagnostics.AddError(
			"Archive creation error",
			fmt.Sprintf("error creating archive: %s", err),
		)
		return
	}

	// Generate archived file stats
	fi, err := os.Stat(outputPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Archive output error",
			fmt.Sprintf("error reading output: %s", err),
		)
		return
	}

	sha1, base64sha256, md5, err := genFileShas(outputPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Hash generation error",
			fmt.Sprintf("error generating hashed: %s", err),
		)
	}

	model.OutputSha = types.StringValue(sha1)
	model.OutputBase64Sha256 = types.StringValue(base64sha256)
	model.OutputMd5 = types.StringValue(md5)
	model.OutputSize = types.Int64Value(fi.Size())

	model.ID = types.StringValue(sha1)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (d *archiveFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

type fileModel struct {
	ID                    types.String `tfsdk:"id"`
	Source                types.Set    `tfsdk:"source"` // sourceModel
	Type                  types.String `tfsdk:"type"`
	SourceContent         types.String `tfsdk:"source_content"`
	SourceContentFilename types.String `tfsdk:"source_content_filename"`
	SourceFile            types.String `tfsdk:"source_file"`
	SourceDir             types.String `tfsdk:"source_dir"`
	Excludes              types.Set    `tfsdk:"excludes"`
	OutputPath            types.String `tfsdk:"output_path"`
	OutputSize            types.Int64  `tfsdk:"output_size"`
	OutputSha             types.String `tfsdk:"output_sha"`
	OutputBase64Sha256    types.String `tfsdk:"output_base64sha256"`
	OutputMd5             types.String `tfsdk:"output_md5"`
	OutputFileMode        types.String `tfsdk:"output_file_mode"`
}

type sourceModel struct {
	Content  types.String `tfsdk:"content"`
	Filename types.String `tfsdk:"filename"`
}

func genFileShas(filename string) (string, string, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", "", "", fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}
	h := sha1.New()
	h.Write(data)
	sha1 := hex.EncodeToString(h.Sum(nil))

	h256 := sha256.New()
	h256.Write(data)
	shaSum := h256.Sum(nil)
	sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])

	md5 := md5.New()
	md5.Write(data)
	md5Sum := hex.EncodeToString(md5.Sum(nil))

	return sha1, sha256base64, md5Sum, nil
}
