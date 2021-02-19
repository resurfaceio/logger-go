package logger

import (
	"net/url"
)

/**
 * Initialize enabled logger using default url.
 */
func NewBaseLoggerAgent(_agent string) *BaseLogger {
	return NewBaseLogger(_agent, "UsageLoggers.urlByDefault()", true, nil)
}

/**
 * Initialize enabled/disabled logger using default url.
 */
func NewBaseLoggerAgentEnabled(_agent string, _enabled bool) *BaseLogger {
	return NewBaseLogger(_agent, UsageLoggers.urlByDefault(), _enabled, nil)
}

/**
 * Initialize enabled logger using url.
 */
func NewBaseLoggerAgentUrl(_agent string, _url string) *BaseLogger {
	return NewBaseLogger(_agent, _url, true, nil)
}

/**
 * Initialize enabled logger using queue.
 */
func NewBaseLoggerAgentQueue(_agent string, _queue []string) *BaseLogger {
	return NewBaseLogger(_agent, UsageLoggers.urlByDefault(), true, _queue)
}

/**
* Initialize enabled/disabled logger using queue.
 */
func NewBaseLoggerAgentQueueEnabled(_agent string, _queue []string, _enabled bool) *BaseLogger {
	return NewBaseLogger(_agent, UsageLoggers.urlByDefault(), _enabled, _queue)
}

//main constructor
func NewBaseLogger(_agent string, _url string, _enabled bool, _queue []string) *BaseLogger {
	baseLogger := &BaseLogger{}
	baseLogger.agent = _agent
	baseLogger.host = "please implement host_lookup()"
	baseLogger.version = "please implement version_lookup()"
	baseLogger.queue = _queue

	//set options in priority order
	baseLogger.enabled = _enabled
	//I believe comparing with an empty string should be the same as comparing with nil
	if _url == "" {
		baseLogger.url = "urlByDefault, in the java version this is defined in UsageLoggers"
		if baseLogger.url == "" {
			baseLogger.enabled = false
		} else {
			baseLogger.url = _url
		}
	}

	//validate url when present
	if baseLogger.url != "" {
		//I'm replacing the try catch statement with the error return
		//but I'll implement that later, this line is a dummy
		baseLogger.urlParsed, _ = url.Parse(baseLogger.url)
		//the rest is throwing errors
	}
	baseLogger.enableable = (baseLogger.url != "")

	return baseLogger
}

//these functions could be redundant, basically setters
func (obj BaseLogger) Enable() {
	obj.enabled = true
}

func (obj BaseLogger) Disable() {
	obj.enabled = false
}

/**
 * Returns host identifier for this logger.
 */
func (obj BaseLogger) hostLookup() string {
	dyno := "this is what System.getenv(dyno) will return"
	if dyno == "" {
		return dyno
	}
	//TODO: implement try/catch style error return
	//Implement whatever this host code is
	return "InetAddress.getLocalHost().getHostName()"
}

/**
 * Submits JSON message to intended destination.
 */
func (obj BaseLogger) Submit(msg string) {
	//woah congrats you submitted the message
	//TODO: implement submit func
}
func (obj BaseLogger) versionLookup() string { return "0.0.0.wehaventstartedityet" }

func (obj BaseLogger) Agent() string         { return obj.agent }
func (obj BaseLogger) Enableable() bool      { return obj.enableable }
func (obj BaseLogger) Enabled() bool         { return obj.enabled }
func (obj BaseLogger) Host() string          { return obj.host }
func (obj BaseLogger) Queue() []string       { return obj.queue }
func (obj BaseLogger) SkipCompression() bool { return obj.skipCompression }
func (obj BaseLogger) SkipSubmission() bool  { return obj.skipSubmission }
func (obj BaseLogger) SubmitFailures() int   { return obj.submitFailures }
func (obj BaseLogger) SubmitSuccesses() int  { return obj.submitSuccesses }
func (obj BaseLogger) Url() string           { return obj.url }
func (obj BaseLogger) UrlParsed() *url.URL   { return obj.urlParsed }
func (obj BaseLogger) Version() string       { return obj.version }

func (obj BaseLogger) SetSkipCompression(_skipCompression bool) {obj.skipCompression = _skipCompression}
func (obj BaseLogger) SetSkipSubmission(_skipSubmission bool) {obj.skipSubmission = _skipSubmission}

type BaseLogger struct {
	agent      string
	enableable bool
	enabled    bool
	host       string
	//easiest to implement a queue in go by using slices, need enqueue and dequque methods
	queue           []string
	skipCompression bool
	skipSubmission  bool
	submitFailures  int
	submitSuccesses int
	url             string
	urlParsed       *url.URL
	version         string
}
