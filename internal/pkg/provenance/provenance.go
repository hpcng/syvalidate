// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package provenance

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sylabs/singularity-mpi/pkg/manifest"
	"github.com/sylabs/syvalidate/internal/pkg/sys"
	"github.com/sylabs/syvalidate/pkg/syblockchainfs"
)

type Provenance struct {
	fs syblockchainfs.SyBlockchainFS

	cacheDir string

	manifests []string // series of manifests (abs path)
}

func Init() (Provenance, error) {
	var p Provenance
	var err error

	// We create a temporary directory that we use as cache for manifests
	// that the provenance package is creating
	p.cacheDir, err = ioutil.TempDir("", "provenance-cache-")
	if err != nil {
		return p, fmt.Errorf("failed to create provenance cache: %s", err)
	}

	// Init the blockchain file system
	fsInfo := syblockchainfs.Info{
		Connected: false, // todo: do not hardcode that, must be a configuration thingy
	}
	p.fs, err = syblockchainfs.Init(&fsInfo)
	if err != nil {
		return p, fmt.Errorf("unable to initialize file system: %s", err)
	}

	// Upon initialization, we ALWAYS capture the platform's manifest
	err = sys.CreatePlatformManifest(p.cacheDir)
	if err != nil {
		return p, fmt.Errorf("failed to create platform's manifests: %s", err)
	}

	// Create and add a manifest for the tool itself
	toolBin := os.Args[0]
	data := manifest.HashFiles([]string{toolBin})
	toolManifestPath := filepath.Join(p.cacheDir, "tool.MANIFEST")
	err = manifest.Create(toolManifestPath, data)
	if err != nil {
		return p, fmt.Errorf("failed to create manifest %s: %s", toolManifestPath, err)
	}

	return p, nil
}

func (p *Provenance) Add(manifestPath string) error {
	p.manifests = append(p.manifests, manifestPath)
	return nil
}

func (p *Provenance) Commit() error {
	// Save all manifest to our BockchainFS, this makes all the data (including output) persistent

	// Calculate the final hashes for the manifests themselves that did not have one yet

	// Generate stamp and publish it to create the blockchain
	myip := "toto"
	stamp := p.fs.CreateStamp(myip)
	for _, m := range p.manifests {
		err := stamp.AddManifest(m)
		if err != nil {
			return fmt.Errorf("failed to add manifest %s to stamp: %s", m, err)
		}
	}
	err := p.fs.CreateBlock(stamp)
	if err != nil {
		return fmt.Errorf("failed to create block: %s", err)
	}

	// Cleanup
	err = os.RemoveAll(p.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to delete provenance cache %s: %s", p.cacheDir, err)
	}

	return nil
}
