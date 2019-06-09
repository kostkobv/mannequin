package kubectl

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var checkVer = regexp.MustCompile(`GitVersion:"(v\d+.\d+.\d+)"`)

func CheckInstalled() (string, error) {
	cmd := exec.Command("kubectl", "version", "--client")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	res := checkVer.FindStringSubmatch(string(out))
	if len(res) < 1 {
		return "", errors.New("kubectl is not installed")
	}

	return res[1], nil
}

func CheckContext(expected string) error {
	cmd := exec.Command("kubectl", "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	k8sCtx := strings.TrimSpace(string(out))
	if k8sCtx != expected {
		return fmt.Errorf("unexpected context %s: expected %s", k8sCtx, expected)
	}

	return nil
}

func UseContext(expected string) error {
	cmd := exec.Command("kubectl", "config", "use-context", expected)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	k8sCtx := strings.TrimSpace(string(out))
	if k8sCtx != `Switched to context "`+expected+`".` {
		return fmt.Errorf("couldn't set context %s: %s", expected, k8sCtx)
	}

	return nil
}

func CheckAndUseContext(expected string) error {
	if err := CheckContext(expected); err != nil {
		if err := UseContext(expected); err != nil {
			return err
		}
	}

	return nil
}
