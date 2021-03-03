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

func GetUsageLoggers() *UsageLoggers {
	onceUsageLoggers.Do(func() {
		envLookup, _ := os.LookupEnv("USAGE_LOGGERS_DISABLE")
		_bricked, err := strconv.ParseBool(envLookup)
		if err != nil {
			//help????
		}
		usageLoggers = &UsageLoggers{
			bricked:  _bricked,
			disabled: _bricked,
		}
	})
	return usageLoggers
}

/**
* Enable all usage loggers, except those explicitly disabled.
 */
func (obj UsageLoggers) Enable() {
	obj.disabled = !obj.bricked
}

/**
* Disable all usage loggers.
 */
func (obj UsageLoggers) Disable() {
	obj.disabled = true
}

/**
* Returns true if usage loggers are generally enabled.
 */
func (obj UsageLoggers) IsEnabled() bool {
	return !obj.disabled
}

/**
* Returns url to use by default.
 */
func (obj UsageLoggers) UrlByDefault() string {
	url, _ := os.LookupEnv("USAGE_LOGGERS_URL")
	return url
}
