package oci

import (
	"testing"

	"k8s.io/minikube/pkg/minikube/constants"
)

func TestRoutableHostIPFromInsideExternalRootlessDaemon(t *testing.T) {
	orig := CachedDaemonInfo
	defer func() {
		CachedDaemonInfo = orig
	}()
	CachedDaemonInfo = func(_ string) (SysInfo, error) {
		return SysInfo{Rootless: true}, nil
	}

	t.Setenv(constants.PodmanContainerHostEnv, "ssh://192.0.2.11/run/user/1000/podman/podman.sock")

	ip, err := RoutableHostIPFromInside(Podman, "minikube", "minikube")
	if err != nil {
		t.Fatalf("RoutableHostIPFromInside() error = %v", err)
	}
	if got, want := ip.String(), "192.0.2.11"; got != want {
		t.Fatalf("RoutableHostIPFromInside() = %q, want %q", got, want)
	}
}
