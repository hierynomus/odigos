package backend

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/keyval-dev/odigos/common"
	"github.com/spf13/cobra"
)

type StackState struct{}

func (d *StackState) Name() common.DestinationType {
	return common.StackStateDestinationType
}

func (d *StackState) ParseFlags(cmd *cobra.Command, selectedSignals []common.ObservabilitySignal) (*ObservabilityArgs, error) {
	apiKey := cmd.Flag("api-key").Value.String()
	if apiKey == "" {
		return nil, fmt.Errorf("API key required for StackState backend, please specify --api-key")
	}

	targetUrl := cmd.Flag("url").Value.String()
	if targetUrl == "" {
		return nil, fmt.Errorf("URL required for StackState backend, please specify --url")
	}

	_, err := url.Parse(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid url specified: %s", err)
	}

	if !strings.Contains(targetUrl, "datadoghq.com") {
		return nil, fmt.Errorf("%s is not a valid datadog url", targetUrl)
	}

	return &ObservabilityArgs{
		Data: map[string]string{
			"DATADOG_SITE": targetUrl,
		},
		Secret: map[string]string{
			"DATADOG_API_KEY": apiKey,
		},
	}, nil
}

func (d *StackState) SupportedSignals() []common.ObservabilitySignal {
	return []common.ObservabilitySignal{
		common.TracesObservabilitySignal,
		common.MetricsObservabilitySignal,
		common.LogsObservabilitySignal,
	}
}
