package contacts

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
)

const subscriptionTopicsUsage = `Subscription topic preferences as a JSON array, e.g. '[{"subscription_topic_id":7,"domain_id":4321,"status":"unsubscribed"}]'; topic and domain pairs left out keep their status`

var validTopicStatuses = []string{"subscribed", "unsubscribed"}

type subscriptionTopicInput struct {
	SubscriptionTopicID int64  `json:"subscription_topic_id"`
	DomainID            int64  `json:"domain_id"`
	Status              string `json:"status"`
}

func parseSubscriptionTopics(raw string) ([]subscriptionTopicInput, error) {
	if !strings.HasPrefix(strings.TrimSpace(raw), "[") {
		return nil, fmt.Errorf("invalid --subscription-topics: must be a JSON array")
	}

	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()

	topics := []subscriptionTopicInput{}
	if err := dec.Decode(&topics); err != nil {
		return nil, fmt.Errorf("invalid --subscription-topics JSON: %w", err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("invalid --subscription-topics JSON: unexpected data after the array")
	}

	for i, topic := range topics {
		if topic.SubscriptionTopicID <= 0 {
			return nil, fmt.Errorf("invalid --subscription-topics: entry %d: subscription_topic_id must be greater than 0", i+1)
		}
		if topic.DomainID <= 0 {
			return nil, fmt.Errorf("invalid --subscription-topics: entry %d: domain_id must be greater than 0", i+1)
		}
		if !slices.Contains(validTopicStatuses, topic.Status) {
			return nil, fmt.Errorf("invalid --subscription-topics: entry %d: status must be one of %s", i+1, strings.Join(validTopicStatuses, ", "))
		}
	}

	return topics, nil
}
