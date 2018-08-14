package main

import (
	"path"
	"path/filepath"
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"regexp"
	"reflect"
)

const (
	defaultDirectoryPermission = 0755
	templatesDir               = "templates"
	filterOut                  = "__filter_out"
)

type filterCmd struct {
	chartPath       string
	outputDir       string
	valuesFile      string
	overwriteValues bool
}

func (cmd *filterCmd) run() error {
	chart := cmd.chartPath

	// copy chart to output-dir if provided
	if cmd.outputDir != "" {
		cmd.outputDir = path.Join(cmd.outputDir, path.Base(cmd.chartPath))
		err := deepCopy(cmd.chartPath, cmd.outputDir)
		if err != nil {
			return err
		}
		chart = cmd.outputDir
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

	// 只先固定檢查在第二層的 key
	values, err = filter(values, func(exp string) error {
		r := regexp.MustCompile(exp)
		return deleteFilesIfMatch(templatesPath, r)
	})
	if err != nil {
		return err
	}

	if cmd.overwriteValues {
		b, err := yaml.Marshal(values)
		if err != nil {
			return err
		}
		out := path.Join(cmd.outputDir, path.Base(cmd.valuesFile))
		fmt.Printf("overwrote %s\n", out)
		ioutil.WriteFile(out, b, defaultDirectoryPermission)
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

func vals(valuesFile string) (yaml.MapSlice, error) {
	base := yaml.MapSlice{}
	bytes, err := ioutil.ReadFile(valuesFile)
	if err != nil {
		return yaml.MapSlice{}, err
	}
	if err := yaml.Unmarshal(bytes, &base); err != nil {
		return yaml.MapSlice{}, fmt.Errorf("failed to parse %s: %s", base, err)
	}
	return base, nil
}

func filter(slice yaml.MapSlice, consume func(regexp string) error) (items []yaml.MapItem, err error) {
	for _, v := range slice {
		nSlice, isNSlice := v.Value.(yaml.MapSlice)

		// If next level is not type of slice, then just set the key to that value
		if !isNSlice || len(nSlice) <= 0 {
			items = append(items, v)
		}

		var nItem []yaml.MapItem
		for _, nv := range nSlice {
			if nv.Key == filterOut {
				if nv.Value != nil {
					regexp, isString := nv.Value.(string)
					if !isString {
						return []yaml.MapItem{}, fmt.Errorf("value of %s must be string, but got %v", nv.Key, reflect.TypeOf(nv.Value))
					}
					err = consume(regexp)
					if err != nil {
						return
					}
					break
				} else {
					continue
				}
			}
			nItem = append(nItem, nv)
		}

		if len(nItem) > 0 {
			items = append(items, yaml.MapItem{Key: v.Key, Value: nItem})
		}
	}
	return
}
