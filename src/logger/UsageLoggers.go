package logger

import (
	"sync"
	//Library used for getting environment variables and other useful env things
	"os"
	//used for converting the string returned by lookupEnv
	"strconv"
)

//name change since helper uses the name "once"
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
		//lookup returns a false bool is it fails, along with a nil value.
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
func (uLogger UsageLoggers) Enable() {
	if !uLogger.bricked {
		uLogger.disabled = false
	}
}

/**
* Disable all usage loggers.
 */
func (uLogger UsageLoggers) Disable() {
	uLogger.disabled = true
}

/**
* Returns true if usage loggers are generally enabled.
 */
func (uLogger UsageLoggers) IsEnabled() bool {
	return !uLogger.disabled
}

/**
* Returns url to use by default.
 */
func (uLogger UsageLoggers) UrlByDefault() string {
	url := os.Getenv("USAGE_LOGGERS_URL")
	return url
}
