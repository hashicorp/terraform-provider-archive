package archive

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-archive/internal/hashcode"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"filename": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_file", "source_dir"},
			},
			"source_content_filename": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_file", "source_dir"},
			},
			"source_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_content", "source_content_filename", "source_dir"},
			},
			"source_dir": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_content", "source_content_filename", "source_file"},
			},
			"excludes": {
				Type:          schema.TypeSet,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_content", "source_content_filename", "source_file"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_size": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
			"output_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "MD5 of output file",
			},
			"output_sha": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA1 checksum of output file",
			},
			"output_sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA256 checksum of output file",
			},
			"output_base64sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Base64 Encoded SHA256 checksum of output file",
			},
			"output_sha512": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA512 checksum of output file",
			},
			"output_base64sha512": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Base64 Encoded SHA512 checksum of output file",
			},
			"output_file_mode": {
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

	checksums, err := genFileChecksums(outputPath)
	if err != nil {
		return fmt.Errorf("could not generate file checksums: %s", err)
	}

	d.Set("output_md5", checksums.md5Hex)
	d.Set("output_sha", checksums.sha1Hex)
	d.Set("output_sha256", checksums.sha256Hex)
	d.Set("output_base64sha256", checksums.sha256Base64)
	d.Set("output_sha512", checksums.sha512Hex)
	d.Set("output_base64sha512", checksums.sha512Base64)
	d.Set("output_size", fi.Size())
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

type fileChecksums struct {
	md5Hex       string
	sha1Hex      string
	sha256Hex    string
	sha256Base64 string
	sha512Hex    string
	sha512Base64 string
}

func genFileChecksums(filename string) (fileChecksums, error) {
	checksums := fileChecksums{}

	data, err := ioutil.ReadFile(filename)
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
