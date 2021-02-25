package logger

func NewUsageLoggers() *UsageLoggers {
	//TODO: impement functionality of "true".equals(System.getenv("USAGE_LOGGERS_DISABLE"))"
	//in java bricked is final and both are static
	//should probably think about why that's the case and how to translate that into go
	//for now bricked is just set to false, this should change later
	_bricked := false
	constructedUsageLogger := &UsageLoggers{
		bricked:  _bricked,
		disabled: _bricked,
	}
	return constructedUsageLogger
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
func urlByDefault() string {
	//String url = System.getProperty("USAGE_LOGGERS_URL");
	url := ""
	if url == "" {
		//return System.getenv("USAGE_LOGGERS_URL")
		return "dummy"
	} else {
		return url
	}

}

type UsageLoggers struct {
	bricked  bool
	disabled bool
}
