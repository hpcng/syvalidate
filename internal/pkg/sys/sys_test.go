// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package sys

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gvallee/go_util/pkg/util"
)

func TestCreatePlatformManifest(t *testing.T) {
	expectedFiles := []string{"memory.MANIFEST", "hdd.MANIFEST", "cpu.MANIFEST", "lspci.MANIFEST", "uname.MANIFEST"}
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(dir)

	err = CreatePlatformManifest(dir)
	if err != nil {
		t.Fatalf("failed to create the platform's manifests: %s", err)
	}

	for _, file := range expectedFiles {
		path := filepath.Join(dir, file)
		if !util.FileExists(path) {
			t.Fatalf("%s does not exist", path)
		}
	}
}
