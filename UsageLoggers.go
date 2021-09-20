// © 2016-2021 Resurface Labs Inc.

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

func GetUsageLoggers() (*UsageLoggers, error) {
	//var envError error
	var parseError error
	onceUsageLoggers.Do(func() {
		//lookup returns a false bool if it fails, along with a nil value.
		//We can ignore this because parsebool will throw an error anyway if this fails
		envLookup, _ := os.LookupEnv("USAGE_LOGGERS_DISABLE")
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
	url := ""
	err := godotenv.Load()
	if err != nil {
		log.Println("env file not loaded, logger disabled: ", err)
	} else {
		url, exists := os.LookupEnv("USAGE_LOGGERS_URL")

		if url == "" || !exists {
			log.Println("USAGE_LOGGERS_URL env var not set or does not exist; logger dissabled")
		}
	}
	return url
}
