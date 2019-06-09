package helm

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kostkobv/mannequin/pkg"
)

var defaultBinPath = "helm"
var checkVer = regexp.MustCompile(`SemVer:"(v\d+.\d+.\d+)"`)

func CheckInstalled(binpath string) (string, error) {
	if binpath == "" {
		binpath = defaultBinPath
	}

	cmd := exec.Command(binpath, "version", "--client")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	res := checkVer.FindStringSubmatch(string(out))
	if len(res) < 1 {
		return "", errors.New("helm is not installed")
	}

	return res[1], nil
}

func Deploy(w io.Writer, vars pkg.VarStorer, lc LConfig) error {
	if err := lc.Validate(); err != nil {
		return err
	}

	if lc.BinaryPath == "" {
		lc.BinaryPath = defaultBinPath
	}
	if lc.Namespace == "" {
		lc.Namespace = lc.ReleaseName
	}

	args := []string{lc.BinaryPath, "upgrade", "--install", "--namespace", vars.Replace(lc.Namespace)}
	if lc.ValuesPath != "" {
		args = append(args, "--values", vars.Replace(lc.ValuesPath))
	}
	if lc.Set != nil {
		for k, v := range lc.Set {
			args = append(args, "--set", vars.Replace(k)+"="+vars.Replace(v))
		}
	}
	if lc.Flags != nil {
		for k, v := range lc.Flags {
			args = append(args, vars.Replace(k))
			if v != "" {
				args = append(args, vars.Replace(v))
			}
		}
	}
	args = append(args, lc.ReleaseName, lc.ChartPath)

	cmd := exec.Command(args[0], args[1:]...)

	out, err := cmd.Output()
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			return err
		}

		return fmt.Errorf("failed to run deployment: error code %d: %s", ee.ExitCode(), strings.TrimSpace(string(ee.Stderr)))
	}

	fmt.Fprintln(w, "Deploying:")
	fmt.Fprintln(w, "----------------------------------------------------")
	fmt.Fprintf(w, "%s\n", out)
	fmt.Fprintln(w, "----------------------------------------------------")
	fmt.Fprintln(w, "Successfully deployed!")

	return nil
}
