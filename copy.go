package main

import (
	"io/ioutil"
	"os"
	"github.com/otiai10/copy"
)

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
