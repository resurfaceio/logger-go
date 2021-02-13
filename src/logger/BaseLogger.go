package logger

//this doesn't resolve, don't know why
//not super important
// import (
// 	"net/url"
// )

/**
 * Initialize enabled logger using default url.
 */
func NewBaseLoggerAgent(_agent string) *BaseLogger {
	return NewBaseLogger(_agent, "", true)
}

/**
 * Initialize enabled/disabled logger using default url.
 */
func NewBaseLoggerEnabledUrl(_agent string, _enabled bool) *BaseLogger {
	return NewBaseLogger(_agent, "", _enabled)
}

/**
 * Initialize enabled logger using url.
 */
func NewBaseLoggerUrl(_agent string, _url string) *BaseLogger {
	return NewBaseLogger(_agent, _url, true)
}

//

/**
 * Initialize enabled/disabled logger using url.
 */

//use two return values to handle errors, ie (*baselogger, error)
func NewBaseLogger(_agent string, _url string, _enabled bool) *BaseLogger {
	baselogger := &BaseLogger{}
	baselogger.agent = _agent
	baselogger.host = "please implement host_lookup()"
	baselogger.version = "please implement version_lookup()"
	//baselogger.queue is nil by default

	//set options in priority order
	baselogger.enabled = _enabled
	//I believe comparing with an empty string should be the same as comparing with nil
	if _url == "" {
		baselogger.url = "urlByDefault, in the java version this is defined in UsageLoggers"
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
	baselogger.enableable = (baselogger.url != "")

	return baselogger
}

//TODO: constructor using queue

//TODO: constructor using queue and enabled/diabled

//TODO: Enable/disable?

//Getters and Setters

func (obj BaseLogger) getAgent() string         { return obj.agent }
func (obj BaseLogger) getHost() string          { return obj.host }
func (obj BaseLogger) getQueue() []string       { return obj.queue }
func (obj BaseLogger) getSkipCompression() bool { return obj.skip_compression }
func (obj BaseLogger) getSkipSubmission() bool  { return obj.skip_submission }
func (obj BaseLogger) getUrl() string           { return obj.url }
func (obj BaseLogger) getVersion() string       { return obj.version }
func (obj BaseLogger) getEnableable() bool      { return obj.enableable }
func (obj BaseLogger) getEnableabled() bool     { return obj.enabled }
func (obj BaseLogger) getSubmitFailues() int    { return obj.submit_failures }
func (obj BaseLogger) getSubmitSuccesses() int  { return obj.submit_successes }

func (obj BaseLogger) setSkipCompression(_skip_compresssion bool) {
	obj.skip_compression = _skip_compresssion
}
func (obj BaseLogger) setSkipSubmission(_skip_submission bool) {
	obj.skip_submission = _skip_submission
}

//End Getters and Setters

/**
 * Returns host identifier for this logger.
 */
func host_lookup() string {
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
func submit(msg string) {
	//woah congrats you submitted the message
	//TODO: implement submit func
}

//the java version uses a getter function to access this
//it isn't used anywhere else
//could use this const version instead
//const version string = "0.0.0.wehaventstartedityet"
func version_lookup() string { return "0.0.0.wehaventstartedityet" }

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
