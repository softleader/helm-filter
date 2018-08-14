package main

import (
	"testing"
	"path"
	"os"
	"os/exec"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
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

func TestFilter(t *testing.T) {
	var slice []yaml.MapItem
	err := yaml.Unmarshal([]byte(`
replicaCount: 1
service:
 __filter_out: 
 type: ClusterIP
 port: 80
ingress:
 __filter_out: ingress*
 enabled: false
 annotations: {}
 path: /
 hosts:
 - chart-example.local
 tls: []
resources: {}
nodeSelector: {}
tolerations: []
affinity: {}`), &slice)
	if err != nil {
		t.Error(err)
	}

	marshalPrint(slice)

	filtered, err := filter(slice, func(regexp string) error {
		fmt.Println("~~~ found", regexp)
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	marshalPrint(filtered)
}

func marshalPrint(i interface{}) {
	b, e := yaml.Marshal(i)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(b))
}
