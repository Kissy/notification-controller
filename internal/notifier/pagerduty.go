/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package notifier

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/fluxcd/pkg/runtime/events"
)

// PagerDuty holds the proxy URL
type PagerDuty struct {
	ProxyURL string
	Token    string
}

// NewPagerDuty returns a PagerDuty object
func NewPagerDuty(proxyURL string, token string) (*PagerDuty, error) {
	return &PagerDuty{
		ProxyURL: proxyURL,
		Token: token,
	}, nil
}

// Post PagerDuty event
func (p *PagerDuty) Post(event events.Event) error {
	eventPayload := &pagerduty.V2Payload{
		Summary:   event.Message,
		Source:    event.ReportingInstance,
		Severity:  event.Severity,
		Timestamp: event.Timestamp.String(),
		Component: event.ReportingController,
		Group:     "flux-system", // TODO get namespace
		Class:     event.InvolvedObject.Kind,
		Details:   event.Metadata,
	}

	pagerdutyEvent := pagerduty.V2Event{
		RoutingKey: p.Token,
		Action:     "trigger",
		DedupKey:   "", // TODO based on details
		Client:     "", // TODO
		ClientURL:  "", // TODO
		Payload:    eventPayload,
	}

	if event.Severity == events.EventSeverityInfo {
		pagerdutyEvent.Action = "resolve"
	}

	_, err := pagerduty.ManageEvent(pagerdutyEvent)
	if err != nil {
		return fmt.Errorf("postMessage failed: %w", err)
	}
	return nil
}
