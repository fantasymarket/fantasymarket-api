package game

import (
	"fantasymarket/game/structs"
)

func loadEvents() (map[string]structs.EventSettings, error) {

	return map[string]structs.EventSettings{}, nil
	// eventSettings := []structs.EventSettings{}

	// eventData, err := ioutil.ReadFile("./game/events.yaml")
	// if err != nil {
	// 	return nil, err
	// }

	// if err := yaml.Unmarshal(eventData, &eventSettings); err != nil {
	// 	return nil, err
	// }

	// eventSettingsMap := map[string]structs.EventSettings{}
	// for _, event := range eventSettings {
	// 	eventSettingsMap[event.EventID] = event
	// }

	// return eventSettingsMap, nil
}
