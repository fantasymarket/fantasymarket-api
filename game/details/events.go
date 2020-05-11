package details

import (
	"gopkg.in/yaml.v2"
)

// LoadEventDetails loads all event details from the events.yaml
func LoadEventDetails() (map[string]EventDetails, error) {

	//return map[string]EventDetails{}, nil
	eventDetails := []EventDetails{}

	eventData, err := EventsYamlBytes()
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(eventData, &eventDetails); err != nil {
		return nil, err
	}

	eventDetailsMap := map[string]EventDetails{}
	for _, event := range eventDetails {
		eventDetailsMap[event.EventID] = event
	}

	return eventDetailsMap, nil
}
