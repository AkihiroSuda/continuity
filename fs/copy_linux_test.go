// +build linux

package fs

import (
	"io"
	"path/filepath"
	"crypto/rand"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/containerd/continuity/testutil"
)

func TestCopyReflinkWithXFS(t *testing.T) {
	testutil.RequiresRoot(t)
	mnt, err := ioutil.TempDir("", "containerd-"+t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(mnt)

	loopback, err := testutil.NewLoopback(1 << 30) // sparse file (max=1GB)
	if err != nil {
		t.Fatal(err)
	}
	mkfs := []string{"mkfs.xfs", "-m", "crc=0", "-n", "ftype=1"}
	if out, err := exec.Command(mkfs[0], append(mkfs[1:], loopback.Device)...).CombinedOutput(); err != nil {
		// not fatal
		t.Skipf("could not mkfs (%v) %s: %v (out: %q)", mkfs, loopback.Device, err, string(out))
	}
	loopbackSize, err := loopback.HardSize()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Loopback file size (after mkfs (%v)): %d", mkfs, loopbackSize)
	if out, err := exec.Command("mount", loopback.Device, mnt).CombinedOutput(); err != nil {
		// not fatal
		t.Skipf("could not mount %s: %v (out: %q)", loopback.Device, err, string(out))
	}
	unmounted := false
	defer func() {
		if !unmounted{
			testutil.Unmount(t, mnt)
		}
		loopback.Close()
	}()

	aPath := filepath.Join(mnt, "a")
	aSize := int64(100 << 20) // 100MB
	a, err := os.Create(aPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.CopyN(a, rand.Reader, aSize); err != nil {
		a.Close()
		t.Fatal(err)
	}
	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	bPath := filepath.Join(mnt, "b")
	if err := CopyFile(bPath, aPath); err != nil {
		t.Fatal(err)
	}
	testutil.Unmount(t, mnt)
	unmounted = true
	loopbackSize, err = loopback.HardSize()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Loopback file size (after copying a %d-byte file): %d", aSize, loopbackSize)
	allowedSize := int64(120 << 20) // 120MB
	if loopbackSize > allowedSize {
		t.Fatalf("expected <= %d, got %d", allowedSize, loopbackSize)
	}
}
