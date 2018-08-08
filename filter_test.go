package main

import (
	"testing"
	"path"
	"os"
	"os/exec"
	"io/ioutil"
)

func TestRun(t *testing.T) {
	chart := "mychart"
	base, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(base)
	helm := exec.Command("sh", "-c", "helm create "+chart)
	helm.Dir = base
	err = helm.Run()
	if err != nil {
		t.Error(err)
	}

	cmd := filterCmd{
		chartPath: path.Join(base, chart),
	}
	cmd.outputDir = path.Join(base, "tmp")

	filter := path.Join(cmd.chartPath, "filter.yaml")
	ioutil.WriteFile(filter, []byte(`
ingress:
  __filter_out: ingress*
`), defaultDirectoryPermission)

	cmd.valuesFile = path.Join(cmd.chartPath, "values.yaml")
	cmd.overwriteValues = true
	os.RemoveAll(cmd.outputDir)

	err = cmd.run()
	if err != nil {
		t.Error(err)
	}
}
