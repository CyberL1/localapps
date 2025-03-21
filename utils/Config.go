package utils

import (
	"encoding/json"
	dbClient "localapps/db/client"

	db "localapps/db/generated"
	"localapps/types"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var CachedConfig types.Config

func UpdateConfigCache() error {
	err := ValidateConfig()
	if err != nil {
		return err
	}

	client, _ := dbClient.GetClient()
	config, err := client.GetConfig(dbClient.Ctx)
	if err != nil {
		return err
	}

	configMap := make(map[string]string)
	for _, c := range config {
		configMap[c.Key] = c.Value.String
	}

	configType := reflect.TypeOf(CachedConfig)
	for i := range configType.NumField() {
		field := configType.Field(i)
		fieldName := strings.ToLower(field.Name)

		if _, ok := configMap[strings.ToLower(fieldName)]; ok {
			fieldValue := reflect.ValueOf(&CachedConfig).Elem().FieldByName(field.Name)
			json.Unmarshal([]byte(configMap[fieldName]), fieldValue.Addr().Interface())
		}
	}
	return nil
}

func ValidateConfig() error {
	client, _ := dbClient.GetClient()
	config, err := client.GetConfig(dbClient.Ctx)
	if err != nil {
		return err
	}

	configStruct := reflect.TypeOf(types.Config{})
	var missingKeys []string

	configMap := make(map[string]string)
	for _, c := range config {
		configMap[c.Key] = c.Value.String
	}

	for i := range configStruct.NumField() {
		field := configStruct.Field(i)

		if _, ok := configMap[strings.ToLower(field.Name)]; !ok {
			missingKeys = append(missingKeys, strings.ToLower(field.Name))
		}
	}

	for _, k := range missingKeys {
		field, _ := reflect.TypeOf(types.Config{}).FieldByName(cases.Title(language.English).String(k))
		defaultValue := field.Tag.Get("default")

		client.SetConfigKey(dbClient.Ctx, db.SetConfigKeyParams{Key: k, Value: pgtype.Text{String: defaultValue, Valid: true}})
	}
	return nil
}
