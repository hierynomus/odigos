package config

import (
	"fmt"
	"net/url"

	odigosv1 "github.com/keyval-dev/odigos/api/odigos/v1alpha1"
	commonconf "github.com/keyval-dev/odigos/autoscaler/controllers/common"
	"github.com/keyval-dev/odigos/common"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	stackstateURLKey      = "STACKSTATE_URL"
	stackstateAPITOKENKey = "${STACKSTATE_API_TOKEN}"
)

var (
	ErrStackStateURLNotSpecified      = fmt.Errorf("StackState url not specified")
	ErrStackStateAPITOKENNotSpecified = fmt.Errorf("Api token not specified")
)

type StackState struct{}

func (n *StackState) DestType() common.DestinationType {
	return common.StackStateDestinationType
}

func (n *StackState) ModifyConfig(dest *odigosv1.Destination, currentConfig *commonconf.Config) {

	if !n.requiredVarsExists(dest) {
		log.Log.V(0).Info("StackState config is missing required variables")
		return
	}

	baseURL, err := parsetheDTurl(dest.Spec.Data[stackstateURLKey])
	if err != nil {
		log.Log.V(0).Info("StackState url is not a valid")
		return
	}

	currentConfig.Exporters["stackstate"] = commonconf.GenericMap{
		"endpoint": baseURL + "/api/v2/otlp",
		"headers": commonconf.GenericMap{
			"Authorization": "Api-Token ${STACKSTATE_API_TOKEN}",
		},
	}

	if isTracingEnabled(dest) {
		currentConfig.Service.Pipelines["traces/stackstate"] = commonconf.Pipeline{
			Receivers:  []string{"otlp"},
			Processors: []string{"batch"},
			Exporters:  []string{"stackstate"},
		}
	}

	if isMetricsEnabled(dest) {
		currentConfig.Service.Pipelines["metrics/stackstate"] = commonconf.Pipeline{
			Receivers:  []string{"otlp"},
			Processors: []string{"batch"},
			Exporters:  []string{"stackstate"},
		}
	}

	if isLoggingEnabled(dest) {
		currentConfig.Service.Pipelines["logs/stackstate"] = commonconf.Pipeline{
			Receivers:  []string{"otlp"},
			Processors: []string{"batch"},
			Exporters:  []string{"stackstate"},
		}
	}
}
func (g *StackState) requiredVarsExists(dest *odigosv1.Destination) bool {
	if _, ok := dest.Spec.Data[stackstateURLKey]; !ok {
		return false
	}
	return true
}

func parsetheStSurl(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		return parsetheStSurl(fmt.Sprintf("https://%s", rawURL))
	}

	return fmt.Sprintf("https://%s", u.Host), nil
}
