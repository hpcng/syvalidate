// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package hashcash

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sylabs/syvalidate/internal/pkg/hash"
)

const (
	dummyIP = "192.1.23.2"
)

func createDummyManifest(t *testing.T, name string) (string, func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}

	path := filepath.Join(dir, name+".MANIFEST")
	data := []byte("dummy content")
	err = ioutil.WriteFile(path, data, 0777)
	if err != nil {
		t.Fatalf("unable to create %s: %s", path, err)
	}

	cleanup := func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("unable to delete %s: %s", dir, err)
		}
	}

	return path, cleanup
}

func TestAddManifest(t *testing.T) {
	path1, cleanup1 := createDummyManifest(t, "dummy1")
	defer cleanup1()

	stamp := Create(dummyIP)
	if stamp.ext != "" {
		t.Fatal("ext is not empty right after creation of new stamp")
	}

	path2, cleanup2 := createDummyManifest(t, "dummy2")
	defer cleanup2()

	err := stamp.AddManifest(path1)
	if err != nil {
		t.Fatalf("failed to add manifest %s: %s", path1, err)
	}

	err = stamp.AddManifest(path2)
	if err != nil {
		t.Fatalf("failed to add manifest %s: %s", path2, err)
	}

	hash1 := hash.HashFile(path1)
	hash2 := hash.HashFile(path2)

	expectedExt := "dummy1=" + hash1 + ";dummy2=" + hash2
	if stamp.ext != expectedExt {
		t.Fatalf("ext does not match expectation: %s vs. %s", stamp.ext, expectedExt)
	}
}

func TestCreate(t *testing.T) {
	stamp := Create(dummyIP)

	if stamp.version != Version {
		t.Fatalf("version does not match expectation: %s vs. %s", stamp.version, Version)
	}

	if stamp.bits != 20 {
		t.Fatalf("bits is %d instead of 20", stamp.bits)
	}

	if stamp.resource != dummyIP {
		t.Fatalf("resource does not match: %s vs %s", stamp.resource, dummyIP)
	}

	if stamp.ext != "" {
		t.Fatal("ext is not empty right after creation of new stamp")
	}

	if stamp.rand == "" {
		t.Fatal("rand is undefined")
	}

	if stamp.counter != "" {
		t.Fatalf("counter is not empty")
	}
}

func TestSerializeParse(t *testing.T) {
	stamp1 := Create(dummyIP)
	str := stamp1.Serialize()
	t.Logf("stamp: %s\n", str)
	stamp2, err := Parse(str)
	if err != nil {
		t.Fatalf("failed to parse stamp: %s", err)
	}

	if stamp1.version != stamp2.version {
		t.Fatalf("version mismatch: %s vs %s", stamp1.version, stamp2.version)
	}

	if stamp1.bits != stamp2.bits {
		t.Fatalf("bits mismatch: %d vs %d", stamp1.bits, stamp2.bits)
	}

	if stamp1.resource != stamp2.resource {
		t.Fatalf("resource mismatch: %s vs %s", stamp1.resource, stamp2.resource)
	}

	if stamp1.ext != stamp2.ext {
		t.Fatalf("ext mismatch: %s vs %s", stamp1.ext, stamp2.ext)
	}

	if stamp1.rand != stamp2.rand {
		t.Fatalf("rand mismatch: %s vs %s", stamp1.rand, stamp2.rand)
	}

	if stamp1.counter != stamp2.counter {
		t.Fatalf("counter mismatch: %s vs %s", stamp1.counter, stamp2.counter)
	}
}

/*
func TestIsValid(t *testing.T) {

}

func TestValidHashCash(t *testing.T) {

}
*/
