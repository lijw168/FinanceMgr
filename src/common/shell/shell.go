package shell

import (
	"os/exec"
)

func Shell(bash string) (string, error) {
	out, err := exec.Command("sh", "-c", bash).Output()

	return string(out[:]), err
}
