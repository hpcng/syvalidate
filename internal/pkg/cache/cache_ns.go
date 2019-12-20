// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// Version is the version of the cache implementation
	Version = "1.0.0"

	// CacheLocalationEnvDir is the name of the environment variable that can be used
	// to set the location of the cache for SyBlockFS
	CacheLocalationEnvDir = "SY_BLOCKFS_DIR"

	defaultCacheDirName     = ".syblockfs"
	defaultNamespaceDirName = "ns"
	defaultDataDirName      = "data"
)

func getBasedir() string {
	if os.Getenv(CacheLocalationEnvDir) != "" {
		return os.Getenv(CacheLocalationEnvDir)
	} else {
		return filepath.Join(os.Getenv("HOME"), defaultCacheDirName)
	}
}

func getNamespaceDir(basedir string) string {
	return filepath.Join(basedir, defaultNamespaceDirName)
}

// AddNamespaces add a list of namespaces to the local cache
// It is okay if the namespace is already in the cache.
func AddNamespaces(basedir string, namespaces []string) error {
	dir := getNamespaceDir(basedir)
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read %s: %s", dir, err)
	}

	for _, ns := range namespaces {
		for _, e := range entries {
			if ns != e.Name() {
				newNSDir := filepath.Join(dir, ns)
				err := os.MkdirAll(newNSDir, 0700)
				if err != nil {
					return fmt.Errorf("failed to create %s: %s", newNSDir, err)
				}
			}
		}
	}

	return nil
}

func LoadNamespaces(basedir string) ([]string, error) {
	dir := getNamespaceDir(basedir)
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %s", dir, err)
	}

	var namespaces []string
	for _, e := range entries {
		namespaces = append(namespaces, e.Name())
	}

	return namespaces, nil
}
