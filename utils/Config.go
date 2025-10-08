package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"localapps/constants"
	dbClient "localapps/db/client"
	"os"
	"path/filepath"

	db "localapps/db/generated"
	"localapps/types"
	"reflect"
)

var ServerConfig types.ServerConfig
var CliConfig types.CliConfig

func UpdateServerConfigCache() error {
	err := validateServerConfig()
	if err != nil {
		return err
	}

	client, _ := dbClient.GetClient()
	config, err := client.GetConfig(context.Background())
	if err != nil {
		return err
	}

	configMap := make(map[string]string)
	for _, c := range config {
		configMap[c.Key] = c.Value.String
	}

	configType := reflect.TypeOf(ServerConfig)
	for i := range configType.NumField() {
		field := configType.Field(i)

		if _, ok := configMap[field.Name]; ok {
			fieldValue := reflect.ValueOf(&ServerConfig).Elem().FieldByName(field.Name)
			json.Unmarshal([]byte(configMap[field.Name]), fieldValue.Addr().Interface())
		}
	}
	return nil
}

func validateServerConfig() error {
	client, _ := dbClient.GetClient()
	config, err := client.GetConfig(context.Background())
	if err != nil {
		return err
	}

	configStruct := reflect.TypeOf(types.ServerConfig{})
	var missingKeys []string

	configMap := make(map[string]string)
	for _, c := range config {
		configMap[c.Key] = c.Value.String
	}

	for i := range configStruct.NumField() {
		field := configStruct.Field(i)

		if _, ok := configMap[field.Name]; !ok {
			missingKeys = append(missingKeys, field.Name)
		}
	}

	for _, k := range missingKeys {
		field, _ := reflect.TypeOf(types.ServerConfig{}).FieldByName(k)
		defaultValue := field.Tag.Get("default")

		client.SetConfigKey(context.Background(), db.SetConfigKeyParams{Key: k, Value: sql.NullString{String: defaultValue, Valid: true}})
	}
	return nil
}

func UpdateCliConfigCache() error {
	err := validateCliConfig()
	if err != nil {
		return err
	}

	configFile, err := os.Open(filepath.Join(constants.LocalappsDir, "cli-config.json"))
	if err != nil {
		return fmt.Errorf("cannot find cli config file: %s", err)
	}

	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	var config types.CliConfig
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("failed to decode cli config: %s", err)
	}

	CliConfig = config
	return nil
}

func validateCliConfig() error {
	_, err := os.Open(filepath.Join(constants.LocalappsDir, "cli-config.json"))
	if err != nil {
		os.WriteFile(filepath.Join(constants.LocalappsDir, "cli-config.json"), []byte("{\"server\":{\"url\":\"http://localhost:8080\"}}"), 0644)
	}
	return nil
}

func SaveCliConfig() error {
	data, err := json.Marshal(CliConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal CliConfig: %s", err)
	}
	os.WriteFile(filepath.Join(constants.LocalappsDir, "cli-config.json"), data, 0644)
	return nil
}
