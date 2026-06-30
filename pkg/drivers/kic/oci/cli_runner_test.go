/*
Copyright 2020 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oci

import (
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"testing"

	"k8s.io/minikube/pkg/minikube/out"
	"k8s.io/minikube/pkg/minikube/tests"
)

func resetPodmanRootlessProbe() {
	podmanRootlessProbeOnce = sync.Once{}
	podmanRootlessDetected = false
	podmanRootlessDetectionDone = false
}

func TestRunCmdWarnSlowOnce(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	func TestPrefixCmdPodmanExternalHostSkipsSudo(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("linux-only podman sudo behavior")
		}
		resetPodmanRootlessProbe()
		t.Setenv("MINIKUBE_ROOTLESS", "false")
		t.Setenv("CONTAINER_HOST", "ssh://127.0.0.1/run/user/1000/podman/podman.sock")

		cmd := PrefixCmd(exec.Command(Podman, "version"))
		if cmd.Args[0] == "sudo" {
			t.Fatalf("expected external podman command without sudo, got %q", strings.Join(cmd.Args, " "))
		}
	}

	func TestPrefixCmdPodmanRootlessForcedSkipsSudo(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("linux-only podman sudo behavior")
		}
		resetPodmanRootlessProbe()
		t.Setenv("MINIKUBE_ROOTLESS", "true")
		t.Setenv("CONTAINER_HOST", "")

		cmd := PrefixCmd(exec.Command(Podman, "version"))
		if cmd.Args[0] == "sudo" {
			t.Fatalf("expected rootless podman command without sudo, got %q", strings.Join(cmd.Args, " "))
		}
	}
	f1 := tests.NewFakeFile()
	out.SetErrFile(f1)

	cmd := exec.Command("sleep", "3")
	_, err := runCmd(cmd, true)

	if err != nil {
		t.Errorf("runCmd has error: %v", err)
	}

	if !strings.Contains(f1.String(), "Executing \"sleep 3\" took an unusually long time") {
		t.Errorf("runCmd does not print the correct log, instead print :%v", f1.String())
	}

	f2 := tests.NewFakeFile()
	out.SetErrFile(f2)

	cmd = exec.Command("sleep", "3")
	_, err = runCmd(cmd, true)

	if err != nil {
		t.Errorf("runCmd has error: %v", err)
	}

	if strings.Contains(f2.String(), "Executing \"sleep 3\" took an unusually long time") {
		t.Errorf("runCmd does not print the correct log, instead print :%v", f2.String())
	}
}
