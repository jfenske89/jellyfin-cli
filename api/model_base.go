package api

import (
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"
)

type ModelList interface {
	ToJson() []byte

	ToMap() (map[string]interface{}, error)

	JsonPath(string) gjson.Result

	Iterate()
}

type Model interface {
	ToJson() []byte

	ToMap() (map[string]interface{}, error)

	JsonPath(string) gjson.Result
}

type genericModel struct {
	jsonBytes []byte
}

func NewModel(jsonBytes []byte) Model {
	return &genericModel{
		jsonBytes: jsonBytes,
	}
}

func newGenericModel(jsonBytes []byte) *genericModel {
	return &genericModel{
		jsonBytes: jsonBytes,
	}
}

func (m *genericModel) ToJson() []byte {
	return m.jsonBytes
}

func (m *genericModel) ToMap() (map[string]interface{}, error) {
	mapped := make(map[string]interface{})

	err := json.Unmarshal(m.jsonBytes, &mapped)

	return mapped, err
}

func (m *genericModel) JsonPath(path string) gjson.Result {
	// it's not real json path
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")

	return gjson.GetBytes(m.jsonBytes, path)
}

func buildModels[T any](jsonBytes []byte, constructor func([]byte) T) ([]T, error) {
	var models []T

	gjson.ParseBytes(jsonBytes).ForEach(func(key, value gjson.Result) bool {
		models = append(models, constructor([]byte(value.String())))

		return true
	})

	return models, nil
}
