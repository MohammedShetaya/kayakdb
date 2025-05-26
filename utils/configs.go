package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// LoadConfigurations loads configurations from a json file into the passed object
// , then override the json values with any value that is set in the environment variables based on the object definition
// NOTE: this implementation assumes that anything can be configuration in json, but some of the configs can be set using env variables
// NOTE: the env vars can only set primitive values string, int, float, etc.
func LoadConfigurations(configObject any) (any, error) {

	if reflect.TypeOf(configObject).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("passed config object should be a pointer")
	}

	err := readConfigFromJson("raft.json", configObject)
	if err != nil {
		return nil, err
	}

	err = readConfigFromEnvironmentVars(configObject)
	if err != nil {
		return nil, err
	}

	return configObject, nil
}

func readConfigFromJson(fileName string, configObject any) error {

	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	// Decode the file into json
	err = json.Unmarshal(file, configObject)

	if err != nil {
		return fmt.Errorf("error decodign json file: %w", err)
	}

	return nil
}

func readConfigFromEnvironmentVars(configObject any) error {
	v := reflect.ValueOf(configObject).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envKey := fieldType.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envVal, exists := os.LookupEnv(envKey)
		if !exists {
			continue
		}

		if !field.CanSet() {
			continue
		}

		switch fieldType.Type.Kind() {
		case reflect.String:
			field.SetString(envVal)

		case reflect.Bool:
			val, err := strconv.ParseBool(envVal)
			if err != nil {
				return fmt.Errorf("unable to parse %s as bool: %w", envKey, err)
			}
			field.SetBool(val)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(envVal, 10, fieldType.Type.Bits())
			if err != nil {
				return fmt.Errorf("unable to parse %s as int: %w", envKey, err)
			}
			field.SetInt(val)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(envVal, 10, fieldType.Type.Bits())
			if err != nil {
				return fmt.Errorf("unable to parse %s as uint: %w", envKey, err)
			}
			field.SetUint(val)

		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(envVal, fieldType.Type.Bits())
			if err != nil {
				return fmt.Errorf("unable to parse %s as float: %w", envKey, err)
			}
			field.SetFloat(val)

		default:
			return fmt.Errorf("unsupported kind %s for field %s", fieldType.Type.Kind(), fieldType.Name)
		}
	}

	return nil
}
