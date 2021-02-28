package logger

import (
	"net/url"
	"sync/atomic"
)

/**
 * Initialize enabled logger using default url.
 */
func NewBaseLoggerAgent(_agent string) *BaseLogger {
	return NewBaseLogger(_agent, UsageLoggers.urlByDefault(), true, nil)
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

	//I believe comparing with an empty string should be the same as comparing with nil
	if _url == "" {
		_url = UsageLoggers.urlByDefault()
		if _url == "" {
			_enabled = false
		}
	}

	var _urlParsed *url.URL
	var parsingError error
	//validate url when present
	if _url != "" {
		_urlParsed, parsingError = url.Parse(_url)
		if parsingError != nil {
			_url = ""
			_urlParsed = nil
			_enabled = false
		}
	}
	_enableable := (_url != "")

	constructedBaseLogger := &BaseLogger{
		agent:           _agent,
		enableable:      _enableable,
		enabled:         _enabled,
		host:            hostLookup(),
		queue:           _queue,
		skipCompression: false,
		skipSubmission:  false,
		url:             _url,
		urlParsed:       _urlParsed,
		version:         versionLookup(),
	}
	return constructedBaseLogger
}

func (obj BaseLogger) Enable() {
	obj.enabled = true
}

func (obj BaseLogger) Disable() {
	obj.enabled = false
}

/**
 * Submits JSON message to intended destination.
 */
func (obj BaseLogger) Submit(msg string) {
	//woah congrats you submitted the message
	//TODO: implement submit func
	if msg == "" || obj.SkipSubmission() || !obj.Enabled() {
		//do nothing
	} else if obj.queue != nil {
		queue.add(msg)
		obj.submitSuccesses.increment()
	} else {
		// HttpURLConnection url_connection = (HttpURLConnection) this.url_parsed.openConnection();
		// url_connection.setConnectTimeout(5000);
		// url_connection.setReadTimeout(1000);
		// url_connection.setRequestMethod("POST");
		// url_connection.setRequestProperty("Content-Type", "application/json; charset=UTF-8");
		// url_connection.setRequestProperty("User-Agent", "Resurface/" + version + " (" + agent + ")");
		// url_connection.setDoOutput(true);
		if obj.SkipCompression() {
			url_connection.setRequestProperty("Content-Encoding", "deflated")
		} else {

		}

	}
}

/**
 * Returns host identifier.
 * These are utility functions that can be static if this wasn't Go
 */
func hostLookup() string {
	dyno := "this is what System.getenv(dyno) will return"
	if dyno == "" {
		return dyno
	}
	//TODO: implement try/catch style error return
	//Implement whatever this host code is
	return InetAddress.getLocalHost().getHostName()
}
func versionLookup() string { return "0.0.0.wehaventstartedityet" }

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

func (obj BaseLogger) SetSkipCompression(_skipCompression bool) {
	obj.skipCompression = _skipCompression
}
func (obj BaseLogger) SetSkipSubmission(_skipSubmission bool) { obj.skipSubmission = _skipSubmission }

type BaseLogger struct {
	agent      string
	enableable bool
	enabled    bool
	host       string
	//easiest to implement a queue in go by using slices, need enqueue and dequque methods
	queue           []string
	skipCompression bool
	skipSubmission  bool
	submitFailures  atomic.Value
	submitSuccesses atomic.Value
	url             string
	urlParsed       *url.URL
	version         string
}
