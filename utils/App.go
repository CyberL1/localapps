package utils

import (
	"context"
	"encoding/json"
	"fmt"
	dbClient "localapps/db/client"
	"localapps/types"
)

func GetAppData(appId string) (*types.ApiAppResponse, error) {
	client, _ := dbClient.GetClient()
	app, err := client.GetAppByAppId(context.Background(), appId)

	if err != nil {
		return nil, fmt.Errorf("app \"%s\" not found", appId)
	}

	var appTyped types.ApiAppResponse
	appBytes, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal app: %v", err)
	}
	json.Unmarshal(appBytes, &appTyped)

	var partsMap map[string]string
	err = json.Unmarshal([]byte(app.Parts), &partsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal app parts: %v", err)
	}

	appTyped.Parts = partsMap
	return &appTyped, nil
}
