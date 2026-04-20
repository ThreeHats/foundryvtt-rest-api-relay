package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ParamSource indicates where to look for a parameter.
type ParamSource string

const (
	FromBody   ParamSource = "body"
	FromQuery  ParamSource = "query"
	FromParams ParamSource = "params"
)

// ParamType indicates the expected type for coercion.
type ParamType string

const (
	TypeString  ParamType = "string"
	TypeNumber  ParamType = "number"
	TypeBoolean ParamType = "boolean"
	TypeArray   ParamType = "array"
	TypeObject  ParamType = "object"
)

// ParamDef defines a parameter to be extracted from the request.
type ParamDef struct {
	Name        string
	From        []ParamSource // Sources to check in order
	Type        ParamType     // Expected type for coercion
	Required    bool
	Description string // Human-readable description for API docs
}

// Params is a map of extracted parameter values.
type Params map[string]interface{}

// GetString returns a string param or empty string.
func (p Params) GetString(key string) string {
	v, _ := p[key].(string)
	return v
}

// GetInt returns an int param or 0.
func (p Params) GetInt(key string) int {
	switch v := p[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		i, _ := strconv.Atoi(v)
		return i
	}
	return 0
}

// GetFloat returns a float64 param or 0.
func (p Params) GetFloat(key string) float64 {
	switch v := p[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
	return 0
}

// GetBool returns a boolean param or false.
func (p Params) GetBool(key string) bool {
	v, _ := p[key].(bool)
	return v
}

// Has returns true if the param exists and is non-nil.
func (p Params) Has(key string) bool {
	v, ok := p[key]
	return ok && v != nil
}

// ExtractParams extracts parameters from the request according to ParamDef specs.
// Body should be pre-parsed JSON (map[string]interface{}).
func ExtractParams(r *http.Request, body map[string]interface{}, defs []ParamDef) (Params, error) {
	params := make(Params)
	var query url.Values // parsed lazily on first FromQuery access

	for _, def := range defs {
		var value interface{}

		for _, source := range def.From {
			switch source {
			case FromQuery:
				if query == nil {
					query = r.URL.Query()
				}
				if v := query.Get(def.Name); v != "" {
					value = v
				}
			case FromBody:
				if body != nil {
					if v, ok := body[def.Name]; ok {
						value = v
					}
				}
			case FromParams:
				if v := chi.URLParam(r, def.Name); v != "" {
					value = v
				}
			}
			if value != nil {
				break
			}
		}

		if value == nil {
			continue
		}

		// Type coercion
		if def.Type != "" {
			coerced, err := coerceType(def.Name, value, def.Type)
			if err != nil {
				return nil, err
			}
			value = coerced
		}

		params[def.Name] = value
	}

	return params, nil
}

// ValidateRequired checks that all required params are present.
func ValidateRequired(params Params, defs []ParamDef) error {
	for _, def := range defs {
		if def.Required {
			if !params.Has(def.Name) {
				return fmt.Errorf("'%s' is required", def.Name)
			}
		}
	}
	return nil
}

func coerceType(name string, value interface{}, expected ParamType) (interface{}, error) {
	switch expected {
	case TypeNumber:
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("'%s' must be a valid number", name)
			}
			return f, nil
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return nil, fmt.Errorf("'%s' must be a valid number", name)
			}
			return f, nil
		default:
			return nil, fmt.Errorf("'%s' must be a valid number", name)
		}

	case TypeBoolean:
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			lower := strings.ToLower(v)
			if lower == "true" {
				return true, nil
			}
			if lower == "false" {
				return false, nil
			}
			return nil, fmt.Errorf("'%s' must be a valid boolean", name)
		default:
			return nil, fmt.Errorf("'%s' must be a valid boolean", name)
		}

	case TypeArray:
		switch v := value.(type) {
		case []interface{}:
			return v, nil
		case string:
			var arr []interface{}
			if err := json.Unmarshal([]byte(v), &arr); err != nil {
				return nil, fmt.Errorf("'%s' must be a valid array", name)
			}
			return arr, nil
		default:
			return nil, fmt.Errorf("'%s' must be an array", name)
		}

	case TypeObject:
		switch v := value.(type) {
		case map[string]interface{}:
			return v, nil
		default:
			return nil, fmt.Errorf("'%s' must be an object", name)
		}

	case TypeString:
		switch v := value.(type) {
		case string:
			return v, nil
		default:
			return nil, fmt.Errorf("'%s' must be a string", name)
		}
	}

	return value, nil
}
