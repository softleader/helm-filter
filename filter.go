package main

import (
	"github.com/otiai10/copy"
	"path"
	"path/filepath"
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"regexp"
	"github.com/kataras/iris/core/errors"
)

const (
	defaultDirectoryPermission = 0755
	templatesDir               = "templates"
	__filter_regexp            = "__filter_regexp"
)

type filterCmd struct {
	chartPath    string
	isolationDir string
	valuesFile   string
}

func (cmd *filterCmd) run() error {
	chart := cmd.chartPath

	// isolate chart path if provided
	if cmd.isolationDir != "" {
		err := deepCopy(cmd.chartPath, cmd.isolationDir)
		if err != nil {
			return err
		}
		chart = cmd.isolationDir
	}

	templatesPath := path.Join(chart, templatesDir)
	// verify templates path exists
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		return fmt.Errorf("templates '%s' does not exist", templatesPath)
	}

	values, err := vals(cmd.valuesFile)
	if err != nil {
		return err
	}

	// 只先 support 定義在第二層避免一直 loop 下去找
	for k, v := range values {
		switch vv := v.(type) {
		case map[interface{}]interface{}:
			for kk, vvv := range vv {
				if kk == __filter_regexp {
					switch exp := vvv.(type) {
					case string:
						r := regexp.MustCompile(exp)
						err := deleteFilesIfMatch(templatesPath, r)
						if err != nil {
							return err
						}
						delete(values, k)
						break
					default:
						return errors.New(fmt.Sprintf("value of %s must be string", kk))
					}
					break
				}
			}
		}
	}

	b, err := yaml.Marshal(values)
	if err != nil {
		return err
	}

	out := path.Join(cmd.isolationDir, path.Base(cmd.valuesFile))
	fmt.Printf("wrote %s\n", out)
	ioutil.WriteFile(out, b, defaultDirectoryPermission)

	return nil
}

func deepCopy(src, dst string) error {
	// 先 copy 到 tmp 防止新的目錄也在同一層會無限 loop 下去
	tmp, err := ioutil.TempDir(os.TempDir(), "helm-filter")
	defer os.RemoveAll(tmp)
	if err != nil {
		return err
	}
	err = copy.Copy(src, tmp)
	if err != nil {
		return err
	}
	err = ensureDirectoryExist(dst)
	if err != nil {
		return err
	}
	err = copy.Copy(tmp, dst)
	if err != nil {
		return err
	}
	return nil
}

func ensureDirectoryExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if dir, err = filepath.Abs(dir); err != nil {
			return err
		}
		return os.MkdirAll(dir, defaultDirectoryPermission)
	}
	return nil
}

func deleteFilesIfMatch(templatesPath string, r *regexp.Regexp) error {
	return filepath.Walk(templatesPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && r.MatchString(f.Name()) {
			fmt.Printf("filtering out '%s'", path)
			return os.Remove(path)

		}
		return err
	})
}

func vals(valuesFile string) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	bytes, err := ioutil.ReadFile(valuesFile)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if err := yaml.Unmarshal(bytes, &base); err != nil {
		return map[string]interface{}{}, fmt.Errorf("failed to parse %s: %s", base, err)
	}
	return base, nil
}
