package mannequin

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	ver           = "v0.0.1"
	defaultK8SCtx = "minikube"
)

// Mnqn represents working application data layer.
type Mnqn struct {
	K8SContext string
	Version    string
	Config     *Config
	LocalVars  LocalVars
	w          io.Writer
}

// New is a constructor for Mnqn.
// Initiates new clean config if none is provided.
func New(w io.Writer, cfg *Config) (Mnqn, error) {
	if w == nil {
		return Mnqn{}, errors.New("writer is required")
	}
	if cfg == nil {
		c, err := NewConfig(ver)
		if err != nil {
			return Mnqn{}, fmt.Errorf("couldn't create new configuration: %s", err)
		}
		cfg = &c
	}
	if err := cfg.Validate(); err != nil {
		return Mnqn{}, fmt.Errorf("configuration is not valid: %s", err)
	}

	return Mnqn{
		Version:    ver,
		Config:     cfg,
		LocalVars:  map[string]string{},
		K8SContext: defaultK8SCtx,
		w:          w,
	}, nil
}

// Write impl.
func (m Mnqn) Write(p []byte) (n int, err error) {
	return m.w.Write(p)
}

type LocalVars map[string]string

// RegisterVar to the execution context.
func (lv *LocalVars) Register(name, val string) error {
	switch {
	case name == "" && val == "":
		return errors.New("attempt to register variable without name and value")
	case name == "":
		return fmt.Errorf("name for variable with value \"%s\" is not provided", val)
	}

	if lv == nil {
		return errors.New("local variables are not defined")
	}

	(*lv)[name] = val

	return nil
}

// Var by name.
func (lv *LocalVars) Var(name string) (string, error) {
	switch {
	case name == "":
		return "", errors.New("name is required")
	case lv == nil:
		return "", errors.New("local variables are not defined")
	}

	val, ok := (*lv)[name]
	if !ok {
		return "", errors.New("variable is not set")
	}

	return val, nil
}

// Replace the variable placeholder with it's value (if registered).
// Expected variable template is `$VARIABLE_NAME` where VARIABLE_NAME is the name of the
// variable with which it's registered.
// Returns provided string if there is nothing to replace.
func (lv *LocalVars) Replace(value string) string {
	if lv == nil {
		return value
	}

	for k, v := range *lv {
		value = strings.ReplaceAll(value, "$"+k, v)
	}

	return value
}
