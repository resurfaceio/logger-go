// © 2016-2023 Graylog, Inc.

package logger

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// sync.Once for UsageLoggers
var onceUsageLoggers sync.Once

type UsageLoggers struct {
	bricked  bool
	disabled bool
}

var usageLoggers *UsageLoggers

func getEnvVar(key string) string {
	// Override existing (or not) env vars with values from .env (e.g. test/dev envs)
	_ = godotenv.Overload()
	// If there's no .env file, the value will be loaded from the existing (or not) env vars for the current process
	value := os.Getenv(key)
	return value
}

func GetUsageLoggers() (*UsageLoggers, error) {
	//var envError error
	var parseError error
	onceUsageLoggers.Do(func() {
		//lookup returns a false bool if it fails, along with a nil value.
		//We can ignore this because parsebool will throw an error anyway if this fails
		envLookup, _ := os.LookupEnv("USAGE_LOGGERS_DISABLE")
		if envLookup == "" { // added to avoid failing test, I'm still not sure the best way to test the env variable lookups in internal testing.
			envLookup = "false"
		}
		_bricked, err := strconv.ParseBool(envLookup)
		parseError = err
		usageLoggers = &UsageLoggers{
			bricked:  _bricked,
			disabled: _bricked,
		}
	})
	return usageLoggers, parseError
}

/**
* Enable all usage loggers, except those explicitly disabled.
 */
func (uLogger *UsageLoggers) Enable() {
	if !uLogger.bricked {
		uLogger.disabled = false
	}
}

/**
* Disable all usage loggers.
 */
func (uLogger *UsageLoggers) Disable() {
	uLogger.disabled = true
}

/**
* Returns true if usage loggers are generally enabled.
 */
func (uLogger *UsageLoggers) IsEnabled() bool {
	return !uLogger.disabled
}

/**
* Returns url to use by default.
 */
func (uLogger *UsageLoggers) UrlByDefault() string {
	url := getEnvVar("USAGE_LOGGERS_URL")
	if url == "" {
		log.Println("USAGE_LOGGERS_URL env var not set or does not exist; logger disabled")
	}
	return url
}

/**
* Returns additional configuration for the base logger
 */
func (uLogger *UsageLoggers) ConfigByDefault() map[string]int {
	config := map[string]int{
		"BODY_LIMIT":         1024 * 1024,     // HTTP payload bytes
		"BUNDLE_SIZE":        5 * 1024 * 1024, // NDJSON bundle size, in bytes
		"MESSAGE_QUEUE_SIZE": 1000,            // Number of buffered messages
		"BUNDLE_QUEUE_SIZE":  128,             // Number of buffered NDJSON bundles
	}
	for key := range config {
		value := getEnvVar("USAGE_LOGGERS_" + key)
		if parsedValue, err := strconv.Atoi(value); err == nil {
			config[key] = parsedValue
		}
	}
	return config
}
