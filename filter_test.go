package main

import (
	"testing"
	"path"
	"os"
)

func TestRun(t *testing.T) {
	chart := "mychart"
	tmp := "/Users/Matt/tmp"
	//helm := exec.Command("sh", "-c", "helm create "+chart)
	//helm.Dir = tmp
	//err := helm.Run()
	//if err != nil {
	//	t.Error(err)
	//}
	cmd := filterCmd{
		chartPath: path.Join(tmp, chart),
	}
	cmd.isolationDir = path.Join(cmd.chartPath, "filtered")
	cmd.valuesFile = path.Join(cmd.chartPath, "values-filter.yaml")

	os.RemoveAll(cmd.isolationDir)

	err := cmd.run()
	if err != nil {
		t.Error(err)
	}
}
