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
	"reflect"
)

const (
	defaultDirectoryPermission = 0755
	templatesDir               = "templates"
	filter                     = "__filter"
)

type filterCmd struct {
	chartPath       string
	isolationDir    string
	valuesFile      string
	overwriteValues bool
}

func (cmd *filterCmd) run() error {
	chart := cmd.chartPath

	// isolate chart path if provided
	if cmd.isolationDir != "" {
		cmd.isolationDir = path.Join(cmd.isolationDir, path.Base(cmd.chartPath))
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
				if kk == filter {
					switch exp := vvv.(type) {
					case string:
						r := regexp.MustCompile(exp)
						err := deleteFilesIfMatch(templatesPath, r)
						if err != nil {
							return err
						}
						delete(values, k)
					case nil:
						delete(vv, kk)
					default:
						return errors.New(fmt.Sprintf("value of %s must be string, but got %v", kk, reflect.TypeOf(exp)))
					}
					break
				}
			}
		}
	}

	if cmd.overwriteValues {
		b, err := yaml.Marshal(values)
		if err != nil {
			return err
		}
		out := path.Join(cmd.isolationDir, path.Base(cmd.valuesFile))
		fmt.Printf("overwrote %s\n", out)
		ioutil.WriteFile(out, b, defaultDirectoryPermission)
	}

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
			fmt.Printf("filtering out '%s'\n", path)
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
