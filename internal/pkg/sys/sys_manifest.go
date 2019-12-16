// Copyright (c) 2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package sys

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sylabs/singularity-mpi/pkg/manifest"
	"github.com/sylabs/singularity-mpi/pkg/syexec"
)

func createManifestForCmdOutput(dir string, cmdBin string, args []string) error {
	var cmd syexec.SyCmd
	var err error
	var execRes syexec.Result

	// get output of uname -msron and save it
	cmd.BinPath, err = exec.LookPath(cmdBin)
	if err != nil {
		return fmt.Errorf("failed to execute uname: %s", err)
	}
	cmd.CmdArgs = args
	cmd.ManifestDir = dir
	cmd.ManifestName = cmdBin
	cmd.ManifestData = []string{execRes.Stdout}
	execRes = cmd.Run()
	if execRes.Err != nil {
		return fmt.Errorf("failed to execute command: %s", execRes.Err)
	}

	return nil
}

func createSystemManifest(dir string) error {
	err := createManifestForCmdOutput(dir, "uname", []string{"-msron"})
	if err != nil {
		return fmt.Errorf("unable to create manifest based on uname output: %s", err)
	}
	return nil
}

func createPCIManifest(dir string) error {
	err := createManifestForCmdOutput(dir, "lspci", nil)
	if err != nil {
		return fmt.Errorf("unable to create manifest based on uname output: %s", err)
	}
	return nil
}

func createCPUManifest(filepath string) error {
	filesToHash := []string{"/proc/cpuinfo"}
	hashData := manifest.HashFiles(filesToHash)
	err := manifest.Create(filepath, hashData)
	if err != nil {
		return fmt.Errorf("failed to create manifest for /proc/cpuinfo")
	}
	return nil
}

func createHardDriveManifest(filepath string) error {
	// read /proc/partitions and extract only the information we need
	d, err := ioutil.ReadFile("/proc/partitions")
	if err != nil {
		return fmt.Errorf("failed to read /proc/partitions")
	}
	content := string(d)
	lines := strings.Split(content, "\n")

	var data []string
	for _, line := range lines {
		if !strings.Contains(line, "loop") {
			data = append(data, line)
		}
	}

	err = manifest.Create(filepath, data)
	if err != nil {
		return fmt.Errorf("failed to create manifest for /proc/cpuinfo")
	}

	return nil
}

func createMemoryManifest(filepath string) error {
	// read /proc/meminfo and extract only the information we need
	d, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return fmt.Errorf("failed to read /proc/meminfo: %s", err)
	}
	content := string(d)
	lines := strings.Split(content, "\n")

	var data []string
	for _, line := range lines {
		if strings.Contains(line, "MemTotal") || strings.Contains(line, "SwapTotal") {
			data = append(data, line)
			if len(data) == 2 {
				break
			}
		}
	}

	err = manifest.Create(filepath, data)
	if err != nil {
		return fmt.Errorf("failed to create manifest for /proc/meminfo")
	}

	return nil
}

func CreatePlatformManifest(dir string) error {
	memManifestPath := filepath.Join(dir, "memory.MANIFEST")
	hddManifestPath := filepath.Join(dir, "hdd.MANIFEST")
	cpuManifestPath := filepath.Join(dir, "cpu.MANIFEST")

	err := createMemoryManifest(memManifestPath)
	if err != nil {
		return fmt.Errorf("failed to create memory manifest: %s", err)
	}
	err = createHardDriveManifest(hddManifestPath)
	if err != nil {
		return fmt.Errorf("failed to create hdd manifest: %s", err)
	}
	err = createCPUManifest(cpuManifestPath)
	if err != nil {
		return fmt.Errorf("failed to create CPU manifest: %s", err)
	}
	err = createPCIManifest(dir)
	if err != nil {
		return fmt.Errorf("failed to create PCI manifest: %s", err)
	}
	err = createSystemManifest(dir)
	if err != nil {
		return fmt.Errorf("failed to create system manifest: %s", err)
	}

	return nil
}
