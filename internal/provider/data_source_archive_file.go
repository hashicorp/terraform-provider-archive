package archive

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-archive/internal/hashcode"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Description: "Generates an archive from content, a file, or directory of files.",

		Schema: map[string]*schema.Schema{
			"type": {
				Description: "The type of archive to generate. NOTE: `zip` is supported.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"source": {
				Description: "Specifies attributes of a single source file to include into the archive.",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Description: "Add this content to the archive with `filename` as the filename.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"filename": {
							Description: "Set this as the filename when declaring a `source`.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
					},
				},
				ConflictsWith: []string{"source_file", "source_dir", "source_content", "source_content_filename"},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["filename"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["content"].(string)))
					return hashcode.String(buf.String())
				},
			},
			"source_content": {
				Description:   "Add only this content to the archive with `source_content_filename` as the filename.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_file", "source_dir"},
			},
			"source_content_filename": {
				Description:   "Set this as the filename when using `source_content`.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_file", "source_dir"},
			},
			"source_file": {
				Description:   "Package this file into the archive.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_content", "source_content_filename", "source_dir"},
			},
			"source_dir": {
				Description:   "Package entire contents of this directory into the archive.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_content", "source_content_filename", "source_file"},
			},
			"excludes": {
				Description:   "Specify files to ignore when reading the `source_dir`.",
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_content", "source_content_filename", "source_file"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_path": {
				Description: "The output of the archive file.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"output_size": {
				Description: "The size of the output archive file.",
				Type:        schema.TypeInt,
				Computed:    true,
				ForceNew:    true,
			},
			"output_sha": {
				Description: "The SHA1 checksum of output archive file.",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
			},
			"output_base64sha256": {
				Description: "The base64-encoded SHA256 checksum of output archive file.",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
			},
			"output_md5": {
				Description: "The MD5 checksum of output archive file.",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
			},
			"output_file_mode": {
				Description: "String that specifies the octal file mode for all archived files. For example: `\"0666\"`. " +
					"Setting this will ensure that cross platform usage of this module will not vary the modes of archived " +
					"files (and ultimately checksums) resulting in more deterministic behavior.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
		},
	}
}

func dataSourceFileRead(d *schema.ResourceData, meta interface{}) error {
	outputPath := d.Get("output_path").(string)

	outputDirectory := path.Dir(outputPath)
	if outputDirectory != "" {
		if _, err := os.Stat(outputDirectory); err != nil {
			if err := os.MkdirAll(outputDirectory, 0755); err != nil {
				return err
			}
		}
	}

	if err := archive(d); err != nil {
		return err
	}

	// Generate archived file stats
	fi, err := os.Stat(outputPath)
	if err != nil {
		return err
	}

	sha1, base64sha256, md5, err := genFileShas(outputPath)
	if err != nil {
		return fmt.Errorf("could not generate file checksum sha256: %s", err)
	}

	err = d.Set("output_sha", sha1)
	if err != nil {
		return fmt.Errorf("read file error output_sha: %s", err)
	}
	err = d.Set("output_base64sha256", base64sha256)
	if err != nil {
		return fmt.Errorf("read file error output_base64sha256: %s", err)
	}
	err = d.Set("output_md5", md5)
	if err != nil {
		return fmt.Errorf("read file error output_md5: %s", err)
	}
	err = d.Set("output_size", fi.Size())
	if err != nil {
		return fmt.Errorf("read file error output_size: %s", err)
	}
	d.SetId(d.Get("output_sha").(string))

	return nil
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, len(configured))
	for i, v := range configured {
		vs[i] = v.(string)
	}
	return vs
}

func archive(d *schema.ResourceData) error {
	archiveType := d.Get("type").(string)
	outputPath := d.Get("output_path").(string)

	archiver := getArchiver(archiveType, outputPath)
	if archiver == nil {
		return fmt.Errorf("archive type not supported: %s", archiveType)
	}

	outputFileMode := d.Get("output_file_mode").(string)
	if outputFileMode != "" {
		archiver.SetOutputFileMode(outputFileMode)
	}

	if dir, ok := d.GetOk("source_dir"); ok {
		if excludes, ok := d.GetOk("excludes"); ok {
			excludeList := expandStringList(excludes.(*schema.Set).List())

			if err := archiver.ArchiveDir(dir.(string), excludeList); err != nil {
				return fmt.Errorf("error archiving directory: %s", err)
			}
		} else {
			if err := archiver.ArchiveDir(dir.(string), []string{""}); err != nil {
				return fmt.Errorf("error archiving directory: %s", err)
			}
		}
	} else if file, ok := d.GetOk("source_file"); ok {
		if err := archiver.ArchiveFile(file.(string)); err != nil {
			return fmt.Errorf("error archiving file: %s", err)
		}
	} else if filename, ok := d.GetOk("source_content_filename"); ok {
		content := d.Get("source_content").(string)
		if err := archiver.ArchiveContent([]byte(content), filename.(string)); err != nil {
			return fmt.Errorf("error archiving content: %s", err)
		}
	} else if v, ok := d.GetOk("source"); ok {
		vL := v.(*schema.Set).List()
		content := make(map[string][]byte)
		for _, v := range vL {
			src := v.(map[string]interface{})
			content[src["filename"].(string)] = []byte(src["content"].(string))
		}
		if err := archiver.ArchiveMultiple(content); err != nil {
			return fmt.Errorf("error archiving content: %s", err)
		}
	} else {
		return fmt.Errorf("one of 'source_dir', 'source_file', 'source_content_filename' must be specified")
	}
	return nil
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
