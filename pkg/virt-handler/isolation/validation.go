package isolation

import (
	"encoding/json"
	"fmt"
	"os/exec"

	containerdisk "kubevirt.io/kubevirt/pkg/container-disk"
)

const (
	QEMUIMGPath = "/usr/bin/qemu-img"
)

func GetImageInfo(imagePath string, context IsolationResult) (*containerdisk.DiskInfo, error) {

	out, err := exec.Command(
		"/usr/bin/virt-chroot", "--user", "qemu", "--memory", "1000", "--cpu", "10", "--mount", context.MountNamespace(), "exec", "--",
		QEMUIMGPath, "info", imagePath, "--output", "json",
	).Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			if len(e.Stderr) > 0 {
				return nil, fmt.Errorf("failed to invoke qemu-img: %v: '%v'", err, string(e.Stderr))
			}
		}
		return nil, fmt.Errorf("failed to invoke qemu-img: %v", err)
	}

	info := &containerdisk.DiskInfo{}
	err = json.Unmarshal(out, info)
	if err != nil {
		return nil, fmt.Errorf("failed to parse disk info: %v", err)
	}
	return info, err
}
