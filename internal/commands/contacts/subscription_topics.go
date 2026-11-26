package contacts

import (
	"fmt"
	"slices"
	"strings"

	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
)

const subscriptionTopicsUsage = `Subscription topic preferences as a JSON array, e.g. '[{"subscription_topic_id":7,"domain_id":4321,"status":"unsubscribed"}]'; topic and domain pairs left out keep their status`

var validTopicStatuses = []string{"subscribed", "unsubscribed"}

type subscriptionTopicInput struct {
	SubscriptionTopicID int64  `json:"subscription_topic_id"`
	DomainID            int64  `json:"domain_id"`
	Status              string `json:"status"`
}

func parseSubscriptionTopics(raw string) ([]subscriptionTopicInput, error) {
	topics, err := cmdutil.ParseJSONArrayFlag[subscriptionTopicInput]("subscription-topics", raw)
	if err != nil {
		return nil, err
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
