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
		//TODO: impement functionality of "true".equals(System.getenv("USAGE_LOGGERS_DISABLE"))"
		//in java bricked is final and both are static
		//should probably think about why that's the case and how to translate that into go
		//for now bricked is just set to false, this should change later
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
	//TODO: implement call functionality
	//String url = System.getProperty("USAGE_LOGGERS_URL");
	url := "dummy"
	if url == "" {
		//return url from env variables
		envUrl, _ := os.LookupEnv("USAGE_LOGGERS_URL")
		return envUrl
	} else {
		//this on the other hand should return the url from the system properties
		//not sure if this idea translates over to Go
		return url
	}

}
