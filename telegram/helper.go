// This file is part of Botgoram
// Botgoram is free software: see LICENSE.txt for more details.

package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func optStr(ret url.Values, key, val string) {
	if val != "" {
		ret.Set(key, val)
	}
}

func optInt(ret url.Values, key string, val int) {
	if val != 0 {
		ret.Set(key, strconv.Itoa(val))
	}
}

func optBool(ret url.Values, key string, val bool) {
	if val {
		ret.Set(key, "true")
	}
}

func optFloat(ret url.Values, key string, val float64) {
	if val != 0.0 {
		ret.Set(key, strconv.FormatFloat(val, 'f', -1, 64))
	}
}

func optJSON(ret url.Values, key string, val interface{}) {
	if data, err := json.Marshal(val); err == nil {
		ret.Set(key, string(data))
	}
}

func mapStr(ret map[string]interface{}, key, val string) {
	if val != "" {
		ret[key] = val
	}
}

func mapInt(ret map[string]interface{}, key string, val int) {
	if val != 0 {
		ret[key] = val
	}
}

func mapBool(ret map[string]interface{}, key string, val bool) {
	if val {
		ret[key] = true
	}
}

func mapFloat(ret map[string]interface{}, key string, val float64) {
	if val != 0.0 {
		ret[key] = val
	}
}

func mapJSON(ret map[string]interface{}, key string, val interface{}) {
	if val != nil {
		ret[key] = val
	}
}
