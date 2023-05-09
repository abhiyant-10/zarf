// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package test provides e2e tests for Zarf.
package test

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/defenseunicorns/zarf/src/pkg/utils/exec"
)

// ZarfE2ETest Struct holding common fields most of the tests will utilize.
type ZarfE2ETest struct {
	ZarfBinPath     string
	Arch            string
	ApplianceMode   bool
	RunClusterTests bool
}

// GetCLIName looks at the OS and CPU architecture to determine which Zarf binary needs to be run.
func GetCLIName() string {
	var binaryName string
	if runtime.GOOS == "linux" {
		binaryName = "zarf"
	} else if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			binaryName = "zarf-mac-apple"
		} else {
			binaryName = "zarf-mac-intel"
		}
	} else if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			binaryName = "zarf.exe"
		}
	}
	return binaryName
}

// Setup performs actions prior to each test.
func (e2e *ZarfE2ETest) Setup(t *testing.T) {
	t.Log("Test setup")
	// Output list of allocated cluster resources
	if runtime.GOOS != "windows" {
		_ = exec.CmdWithPrint("sh", "-c", "kubectl describe nodes |grep -A 99 Non\\-terminated")
	} else {
		t.Log("Skipping kubectl describe nodes on Windows")
	}
}

// SetupWithCluster performs actions for each test that requires a K8s cluster.
func (e2e *ZarfE2ETest) SetupWithCluster(t *testing.T) {
	if !e2e.RunClusterTests {
		t.Skip("")
	}
	e2e.Setup(t)
}

// Teardown performs actions prior to tearing down each test.
func (e2e *ZarfE2ETest) Teardown(t *testing.T) {
	t.Log("Test teardown")
}

// ExecZarfCommand executes a Zarf command.
func (e2e *ZarfE2ETest) ExecZarfCommand(commandString ...string) (string, string, error) {
	return exec.CmdWithContext(context.TODO(), exec.PrintCfg(), e2e.ZarfBinPath, commandString...)
}

// CleanFiles removes files and directories that have been created during the test.
func (e2e *ZarfE2ETest) CleanFiles(files ...string) {
	for _, file := range files {
		_ = os.RemoveAll(file)
	}
}

// GetMismatchedArch determines what architecture our tests are running on,
// and returns the opposite architecture.
func (e2e *ZarfE2ETest) GetMismatchedArch() string {
	switch e2e.Arch {
	case "arm64":
		return "amd64"
	default:
		return"arm64"
	}
}