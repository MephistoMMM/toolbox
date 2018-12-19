package udf

import (
	"fmt"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"gopkg.in/yaml.v2"
	"strings"
)

//AsInt converts source into int
func AsInt(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.ToInt(source)
}

//AsInt converts source into int
func AsString(source interface{}, state data.Map) (interface{}, error) {
	if toolbox.IsSlice(source) || toolbox.IsMap(source) || toolbox.IsStruct(source) {
		text, err := toolbox.AsJSONText(source)
		if err == nil {
			return text, nil
		}
	}
	return toolbox.AsString(source), nil
}

//AsFloat converts source into float64
func AsFloat(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.AsFloat(source), nil
}

//AsBool converts source into bool
func AsBool(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.AsBoolean(source), nil
}

//AsMap converts source into map
func AsMap(source interface{}, state data.Map) (interface{}, error) {
	if source == nil || toolbox.IsMap(source) {
		return source, nil
	}
	source = convertToTextIfNeeded(source)
	if text, ok := source.(string); ok {
		text = strings.TrimSpace(text)
		aMap := map[string]interface{}{}
		if strings.HasPrefix(text, "{") || strings.HasSuffix(text, "}") {
			if err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(text)).Decode(&aMap); err != nil {
				return nil, err
			}
		}
		if err := yaml.NewDecoder(strings.NewReader(toolbox.AsString(source))).Decode(&aMap); err != nil {
			return nil, err
		}
		return toolbox.NormalizeKVPairs(aMap)
	}
	return toolbox.ToMap(source)
}

//AsCollection converts source into a slice
func AsCollection(source interface{}, state data.Map) (interface{}, error) {
	if source == nil || toolbox.IsSlice(source) {
		return source, nil
	}
	source = convertToTextIfNeeded(source)
	if text, ok := source.(string); ok {
		text = strings.TrimSpace(text)
		if strings.HasPrefix(text, "[") || strings.HasSuffix(text, "[") {
			aSlice := []interface{}{}
			if err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(text)).Decode(&aSlice); err != nil {
				return nil, err
			}
		}
		var aSlice interface{}
		if err := yaml.NewDecoder(strings.NewReader(toolbox.AsString(source))).Decode(&aSlice); err != nil {
			return nil, err
		}
		return toolbox.NormalizeKVPairs(aSlice)
	}
	return nil, fmt.Errorf("unable convert to slice, unsupported type: %T", source)
}

//AsData converts source into map or slice
func AsData(source interface{}, state data.Map) (interface{}, error) {
	if source == nil || toolbox.IsMap(source) || toolbox.IsSlice(source) {
		return source, nil
	}
	var aData interface{}
	source = convertToTextIfNeeded(source)
	if text, ok := source.(string); ok {
		text = strings.TrimSpace(text)
		if strings.HasPrefix(text, "[") || strings.HasSuffix(text, "[") || strings.HasPrefix(text, "{") || strings.HasSuffix(text, "}") {
			if err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(text)).Decode(&aData); err != nil {
				return nil, err
			}
		}
		if err := yaml.NewDecoder(strings.NewReader(toolbox.AsString(source))).Decode(&aData); err != nil {
			return nil, err
		}
		return toolbox.NormalizeKVPairs(aData)
	}
	return source, nil
}

func convertToTextIfNeeded(data interface{}) interface{} {
	if data == nil {
		return data
	}
	if bs, ok := data.([]byte); ok {
		return string(bs)
	}
	return data
}