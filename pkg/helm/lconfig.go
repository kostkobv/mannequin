package helm

import (
	"errors"
)

// LConfig of Helm values.
/*
   - helm upgrade
       --install
       --namespace "$NAMESPACE_PRODUCTION"
       --values helm/production.yml
       --set image.repository="$IMAGE_LOCATION"
       --set image.tag="$VERSION"
       --set googleCloud.credentials="${GCLOUD_BASE64_PRODUCTION}"
       --set vault.username="${VAULT_PRODUCTION_USERNAME}"
       --set vault.password="${VAULT_PRODUCTION_PASSWORD}"
       --set internalAuth.key="${INTERNALAUTH_PRODUCTION_KEY}"
       --set internalAuth.value="${INTERNALAUTH_PRODUCTION_VALUE}"
       --version "1.0.0"
       $NAME
       $HELM_CHART
*/
type LConfig struct {
	BinaryPath  string            `yaml:"binary_path,omitempty,flow"`
	ValuesPath  string            `yaml:"values,omitempty,flow"`
	Namespace   string            `yaml:"namespace,omitempty,flow"`
	Set         map[string]string `yaml:"set,omitempty,flow"`
	Flags       map[string]string `yaml:"flags,omitempty,flow"`
	ReleaseName string            `yaml:"release_name,omitempty,flow"`
	ChartPath   string            `yaml:"chart,flow"`
}

// New is a constructor for LConfig.
func New(chart, name string) (LConfig, error) {
	switch {
	case chart == "":
		return LConfig{}, errors.New("chart path is required")
	case name == "":
		return LConfig{}, errors.New("release name is required")
	}

	return LConfig{
		BinaryPath:  defaultBinPath,
		ReleaseName: name,
		ChartPath:   chart,
	}, nil
}

// Validate the LConfig.
func (lc *LConfig) Validate() error {
	switch {
	case lc.ChartPath == "":
		return errors.New("chart path is required")
	case lc.ReleaseName == "":
		return errors.New("release name is required")
	}

	return nil
}
