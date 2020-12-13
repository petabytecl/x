package viperx

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/petabytecl/x/stringslice"
	"github.com/petabytecl/x/stringsx"
)

// GetFloat64 returns a float64 value from viper config or the fallback value.
func GetFloat64(key string, fallback float64) float64 {
	v := viper.GetFloat64(key)
	if v == 0 {
		return fallback
	}

	return v
}

// GetInt returns an int value from viper config or the fallback value.
func GetInt(key string, fallback int) int {
	v := viper.GetInt(key)
	if v == 0 {
		return fallback
	}

	return v
}

// GetDuration returns a duration from viper config or the fallback value.
func GetDuration(key string, fallback time.Duration) time.Duration {
	v := viper.GetDuration(key)
	if v == 0 {
		return fallback
	}

	return v
}

// GetString returns a string from viper config or the fallback value.
func GetString(key string, fallback string) string {
	v := viper.GetString(key)
	if len(v) == 0 {
		return fallback
	}

	return v
}

// GetBool returns a bool from viper config or false.
func GetBool(key string, fallback bool) bool {
	if !viper.IsSet(key) {
		return fallback
	}

	return viper.GetBool(key)
}

// GetStringSlice returns a string slice from viper config or the fallback value.
func GetStringSlice(key string, fallback []string) []string {
	v := viper.GetStringSlice(key)
	r := make([]string, 0, len(v))
	for _, s := range v {
		if len(s) == 0 {
			continue
		}

		if strings.Contains(s, ",") {
			r = append(r, stringslice.TrimSpaceEmptyFilter(stringsx.Splitx(s, ","))...)
		} else {
			r = append(r, s)
		}
	}

	if len(r) == 0 {
		return fallback
	}

	return r
}

// GetStringMapConfig returns a string map using all settings which will lookup env vars
func GetStringMapConfig(paths ...string) map[string]interface{} {
	node := viper.AllSettings()

	for _, path := range paths {
		subNode, ok := node[path].(map[string]interface{})
		if !ok {
			return make(map[string]interface{})
		}

		node = subNode
	}

	return node
}
