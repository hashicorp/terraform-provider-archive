package archive

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/karrick/godirwalk"
	"github.com/mholt/archiver"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Archive type is now determined from the output_path's file extenstion.",
			},
			"source": {
				Type:     schema.TypeSet,
				Optional: true,
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
				ForceNew: true,
			},
			"output_size": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
			"output_sha": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA1 checksum of output file",
			},
			"output_base64sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Base64 Encoded SHA256 checksum of output file",
			},
			"output_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "MD5 of output file",
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
	d.Set("output_sha", sha1)
	d.Set("output_base64sha256", base64sha256)
	d.Set("output_md5", md5)

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

func checkMatch(fileName string, excludes []string) (value bool) {
	for _, exclude := range excludes {
		if exclude == "" {
			continue
		}

		if exclude == fileName {
			return true
		}
	}
	return false
}

func getFileList(dirName string, excludes []string) ([]string, error) {
	var files []string

	err := godirwalk.Walk(dirName, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			relname, err := filepath.Rel(dirName, osPathname)
			if err != nil {
				return nil
			}

			if checkMatch(relname, []string{".", ".."}) {
				return nil
			}

			shouldExclude := checkMatch(relname, excludes)
			fullName := filepath.FromSlash(fmt.Sprintf("%s/%s", dirName, relname))

			if de.IsDir() {
				if shouldExclude {
					return filepath.SkipDir
				}
			} else if shouldExclude {
				return nil
			}

			files = append(files, fullName)
			return nil
		},
	})

	return files, err
}

func archive(d *schema.ResourceData) error {
	outputPath := d.Get("output_path").(string)
	var filesToArchive []string

	compressor := archiver.MatchingFormat(outputPath)
	if compressor == nil {
		return fmt.Errorf("cannot compress unsupported file type: %s", outputPath)
	}

	var err error

	if dir, ok := d.GetOk("source_dir"); ok {
		var excludeList []string
		if excludes, ok := d.GetOk("excludes"); ok {
			excludeList = expandStringList(excludes.(*schema.Set).List())
		}

		filesToArchive, err = getFileList(dir.(string), excludeList)
		if err != nil {
			return fmt.Errorf("could not walk dir: %s", dir.(string))
		}
	} else if file, ok := d.GetOk("source_file"); ok {
		filesToArchive = append(filesToArchive, file.(string))
	} else if fileName, ok := d.GetOk("source_content_filename"); ok {
		return fmt.Errorf("source_content not supported for %s", fileName)
	} else if v, ok := d.GetOk("source"); ok {
		vL := v.(*schema.Set).List()
		content := make(map[string][]byte)
		for _, v := range vL {
			src := v.(map[string]interface{})
			content[src["filename"].(string)] = []byte(src["content"].(string))
		}

		keys := reflect.ValueOf(content).MapKeys()
		return fmt.Errorf("cannot compress %d source blocks", len(keys))
	} else {
		return fmt.Errorf("one of 'source_dir', 'source_file', 'source_content_filename' must be specified")
	}

	check := compressor.Make(outputPath, filesToArchive)
	if check != nil {
		return fmt.Errorf("could not archive to %s: %s", outputPath, check)
	}

	return nil
}

func genFileShas(filename string) (string, string, string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", "", fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}
	h := sha1.New()
	h.Write([]byte(data))
	sha1 := hex.EncodeToString(h.Sum(nil))

	h256 := sha256.New()
	h256.Write([]byte(data))
	shaSum := h256.Sum(nil)
	sha256base64 := base64.StdEncoding.EncodeToString(shaSum[:])

	md5 := md5.New()
	md5.Write([]byte(data))
	md5Sum := hex.EncodeToString(md5.Sum(nil))

	return sha1, sha256base64, md5Sum, nil
}
