package logger

//why does this not work???????
// import (
// 	"net/url"
// )

//to be defined in UsageLoggers or something of the sort
const (
	urlByDefault = ""
)

/**
 * Initialize enabled logger using default url.
 */
func NewBaseLoggerAgent(_agent string) *BaseLogger {
	return &BaseLogger{agent: _agent, enableable: true}
}

/**
 * Initialize enabled/disabled logger using default url.
 */
func NewBaseLoggerAgentEnabled(_agent string, _enabled bool) *BaseLogger {
	return &BaseLogger{agent: _agent, enableable: true}
}

//

/**
 * Initialize enabled/disabled logger using url.
 */
func NewBaseLogger(_agent string, _url string, _enabled bool) (*BaseLogger, error) {
	baselogger := &BaseLogger{}
	baselogger.agent = _agent
	baselogger.host = "please implement host_lookup()"
	baselogger.version = "please implement version_lookup()"
	//baselogger.queue is nil by default

	//set options in priority order
	baselogger.enabled = _enabled
	//I believe comparing with an empty string should be the same as comparing with nil
	if _url == "" {
		baselogger.url = urlByDefault
		if baselogger.url == "" {
			baselogger.enabled = false
		} else {
			baselogger.url = _url
		}
	}

	//validate url when present
	if baselogger.url != "" {
		//I'm replacing the try catch statement with the error return
		//but I'll implement that later, this line is a dummy
		baselogger.url_parsed = baselogger.url + ";parsed"
		//the rest is throwing errors
	}
}

type BaseLogger struct {
	agent            string
	enableable       bool
	enabled          bool
	host             string
	queue            []string
	skip_compression bool
	skip_submission  bool
	submit_failures  int
	submit_successes int
	url              string
	//url parsed should be of type URL, but I can't get the package to import
	url_parsed string
	version    string
}
