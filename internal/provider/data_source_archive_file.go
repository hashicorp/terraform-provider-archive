// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
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

func (d *archiveFileDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.AtLeastOneOf(
			fwpath.MatchRoot("source"),
			fwpath.MatchRoot("source_content_filename"),
			fwpath.MatchRoot("source_file"),
			fwpath.MatchRoot("source_dir"),
		),
	}
}

func (d *archiveFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates an archive from content, a file, or directory of files. " +
			"The archive is built during the terraform plan, so you must persist the archive through to the terraform apply. " +
			"See the `archive_file` resource for an alternative if you cannot persist the file, " +
			"such as in a multi-phase CI or build server context.",
		Blocks: map[string]schema.Block{
			"source": schema.SetNestedBlock{
				Description: "Specifies attributes of a single source file to include into the archive. " +
					"One and only one of `source`, `source_content_filename` (with `source_content`), `source_file`, " +
					"or `source_dir` must be specified.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"content": schema.StringAttribute{
							Description: "Add this content to the archive with `filename` as the filename.",
							Required:    true,
						},
						"filename": schema.StringAttribute{
							Description: "Set this as the filename when declaring a `source`.",
							Required:    true,
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The sha1 checksum hash of the output.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of archive to generate. NOTE: `zip` and `tar.gz` is supported.",
				Required:    true,
			},
			"source_content": schema.StringAttribute{
				Description: "Add only this content to the archive with `source_content_filename` as the filename. " +
					"One and only one of `source`, `source_content_filename` (with `source_content`), `source_file`, " +
					"or `source_dir` must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
					),
				},
			},
			"source_content_filename": schema.StringAttribute{
				Description: "Set this as the filename when using `source_content`. " +
					"One and only one of `source`, `source_content_filename` (with `source_content`), `source_file`, " +
					"or `source_dir` must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_dir"),
					),
				},
			},
			"source_file": schema.StringAttribute{
				Description: "Package this file into the archive. " +
					"One and only one of `source`, `source_content_filename` (with `source_content`), `source_file`, " +
					"or `source_dir` must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_dir"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
			},
			"source_dir": schema.StringAttribute{
				Description: "Package entire contents of this directory into the archive. " +
					"One and only one of `source`, `source_content_filename` (with `source_content`), `source_file`, " +
					"or `source_dir` must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						fwpath.MatchRoot("source_file"),
						fwpath.MatchRoot("source_content"),
						fwpath.MatchRoot("source_content_filename"),
					),
				},
			},
			"excludes": schema.SetAttribute{
				Description: "Specify files/directories to ignore when reading the `source_dir`. " +
					"Supports glob file matching patterns including doublestar/globstar (`**`) patterns.",
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
			"exclude_symlink_directories": schema.BoolAttribute{
				Optional: true,
				Description: "Boolean flag indicating whether symbolically linked directories should be excluded during " +
					"the creation of the archive. Defaults to `false`.",
			},
			"output_path": schema.StringAttribute{
				Description: "The output of the archive file.",
				Required:    true,
			},
			"output_size": schema.Int64Attribute{
				Description: "The byte size of the output archive file.",
				Computed:    true,
			},
			"output_file_mode": schema.StringAttribute{
				Description: "String that specifies the octal file mode for all archived files. For example: `\"0666\"`. " +
					"Setting this will ensure that cross platform usage of this module will not vary the modes of archived " +
					"files (and ultimately checksums) resulting in more deterministic behavior.",
				Optional: true,
			},
			"output_md5": schema.StringAttribute{
				Description: "MD5 of output file",
				Computed:    true,
			},
			"output_sha": schema.StringAttribute{
				Description: "SHA1 checksum of output file",
				Computed:    true,
			},
			"output_sha256": schema.StringAttribute{
				Description: "SHA256 checksum of output file",
				Computed:    true,
			},
			"output_base64sha256": schema.StringAttribute{
				Description: "Base64 Encoded SHA256 checksum of output file",
				Computed:    true,
			},
			"output_sha512": schema.StringAttribute{
				Description: "SHA512 checksum of output file",
				Computed:    true,
			},
			"output_base64sha512": schema.StringAttribute{
				Description: "Base64 Encoded SHA512 checksum of output file",
				Computed:    true,
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

	switch {
	case !model.SourceDir.IsNull():
		excludeList := make([]string, len(model.Excludes.Elements()))

		if !model.Excludes.IsNull() {
			var elements []types.String
			model.Excludes.ElementsAs(ctx, &elements, false)

			for i, elem := range elements {
				excludeList[i] = elem.ValueString()
			}
		}

		opts := ArchiveDirOpts{
			Excludes: excludeList,
		}

		if !model.ExcludeSymlinkDirectories.IsNull() {
			opts.ExcludeSymlinkDirectories = model.ExcludeSymlinkDirectories.ValueBool()
		}

		if err := archiver.ArchiveDir(model.SourceDir.ValueString(), opts); err != nil {
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
	model.OutputSize = types.Int64Value(fi.Size())

	checksums, err := genFileChecksums(outputPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Hash generation error",
			fmt.Sprintf("error generating checksums: %s", err),
		)
	}
	model.OutputMd5 = types.StringValue(checksums.md5Hex)
	model.OutputSha = types.StringValue(checksums.sha1Hex)
	model.OutputSha256 = types.StringValue(checksums.sha256Hex)
	model.OutputBase64Sha256 = types.StringValue(checksums.sha256Base64)
	model.OutputSha512 = types.StringValue(checksums.sha512Hex)
	model.OutputBase64Sha512 = types.StringValue(checksums.sha512Base64)

	model.ID = types.StringValue(checksums.sha1Hex)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (d *archiveFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

type fileModel struct {
	ID                        types.String `tfsdk:"id"`
	Source                    types.Set    `tfsdk:"source"` // sourceModel
	Type                      types.String `tfsdk:"type"`
	SourceContent             types.String `tfsdk:"source_content"`
	SourceContentFilename     types.String `tfsdk:"source_content_filename"`
	SourceFile                types.String `tfsdk:"source_file"`
	SourceDir                 types.String `tfsdk:"source_dir"`
	Excludes                  types.Set    `tfsdk:"excludes"`
	ExcludeSymlinkDirectories types.Bool   `tfsdk:"exclude_symlink_directories"`
	OutputPath                types.String `tfsdk:"output_path"`
	OutputSize                types.Int64  `tfsdk:"output_size"`
	OutputFileMode            types.String `tfsdk:"output_file_mode"`
	OutputMd5                 types.String `tfsdk:"output_md5"`
	OutputSha                 types.String `tfsdk:"output_sha"`
	OutputSha256              types.String `tfsdk:"output_sha256"`
	OutputBase64Sha256        types.String `tfsdk:"output_base64sha256"`
	OutputSha512              types.String `tfsdk:"output_sha512"`
	OutputBase64Sha512        types.String `tfsdk:"output_base64sha512"`
}

type sourceModel struct {
	Content  types.String `tfsdk:"content"`
	Filename types.String `tfsdk:"filename"`
}

type fileChecksums struct {
	md5Hex       string
	sha1Hex      string
	sha256Hex    string
	sha256Base64 string
	sha512Hex    string
	sha512Base64 string
}

func genFileChecksums(filename string) (fileChecksums, error) {
	var checksums fileChecksums

	data, err := os.ReadFile(filename)
	if err != nil {
		return checksums, fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}

	md5Sum := md5.Sum(data)
	checksums.md5Hex = hex.EncodeToString(md5Sum[:])

	sha1Sum := sha1.Sum(data)
	checksums.sha1Hex = hex.EncodeToString(sha1Sum[:])

	sha256Sum := sha256.Sum256(data)
	checksums.sha256Hex = hex.EncodeToString(sha256Sum[:])
	checksums.sha256Base64 = base64.StdEncoding.EncodeToString(sha256Sum[:])

	sha512Sum := sha512.Sum512(data)
	checksums.sha512Hex = hex.EncodeToString(sha512Sum[:])
	checksums.sha512Base64 = base64.StdEncoding.EncodeToString(sha512Sum[:])

	return checksums, nil
}
