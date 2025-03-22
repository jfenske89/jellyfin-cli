package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type ModelList interface {
	ToJson() []byte

	ToMap() (map[string]any, error)

	JsonPath(string) gjson.Result

	Iterate()
}

type Model interface {
	ToJson() []byte

	ToMap() (map[string]any, error)

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

func (m *genericModel) ToMap() (map[string]any, error) {
	mapped := make(map[string]any)

	err := json.Unmarshal(m.jsonBytes, &mapped)

	return mapped, err
}

func (m *genericModel) JsonPath(path string) gjson.Result {
	// it's not real json path
	path = strings.TrimPrefix(path, "$")
	path = strings.TrimPrefix(path, ".")

	return gjson.GetBytes(m.jsonBytes, path)
}

func (m *genericModel) JsonPathAsDate(path string) time.Time {
	result := m.JsonPath(path)

	if result.Exists() {
		switch result.Type {
		case gjson.String:
			// 2024-02-18T16:31:11.9906841Z
			if timeResult, err := time.Parse("2006-01-02T15:04:05.9999999Z", result.String()); err == nil {
				return timeResult
			}

		case gjson.Number:
			return time.Unix(result.Int(), 0)
		}
	}

	return time.Time{}
}

func (m *genericModel) MarshalJSON() ([]byte, error) {
	return m.jsonBytes, nil
}

func (m *genericModel) UnmarshalJSON(data []byte) error {
	m.jsonBytes = data

	return nil
}

func buildModels[T any](jsonBytes []byte, constructor func([]byte) T) ([]T, error) {
	var models []T

	gjson.ParseBytes(jsonBytes).ForEach(func(key, value gjson.Result) bool {
		models = append(models, constructor([]byte(value.String())))

		return true
	})

	return models, nil
}

func buildModel[T any](jsonBytes []byte, constructor func([]byte) T) (T, error) {
	return constructor(jsonBytes), nil
}
