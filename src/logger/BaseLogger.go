package logger

//this doesn't resolve, don't know why
//not super important right now
// import (
// 	"net/url"
// )

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
	return NewBaseLogger(_agent, "UsageLoggers.urlByDefault()", _enabled, nil)
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
	return NewBaseLogger(_agent, "UsageLoggers.urlByDefault()", true, _queue)
}

/**
* Initialize enabled/disabled logger using queue.
 */
func NewBaseLoggerAgentQueueEnabled(_agent string, _queue []string, _enabled bool) *BaseLogger {
	return NewBaseLogger(_agent, "UsageLoggers.urlByDefault()", _enabled, _queue)
}

//main constructor
func NewBaseLogger(_agent string, _url string, _enabled bool, _queue []string) *BaseLogger {
	baseLogger := &BaseLogger{}
	baseLogger.Agent = _agent
	baseLogger.Host = "please implement host_lookup()"
	baseLogger.Version = "please implement version_lookup()"
	baseLogger.Queue = _queue

	//set options in priority order
	baseLogger.Enabled = _enabled
	//I believe comparing with an empty string should be the same as comparing with nil
	if _url == "" {
		baseLogger.Url = "urlByDefault, in the java version this is defined in UsageLoggers"
		if baseLogger.Url == "" {
			baseLogger.Enabled = false
		} else {
			baseLogger.Url = _url
		}
	}

	//validate url when present
	if baseLogger.Url != "" {
		//I'm replacing the try catch statement with the error return
		//but I'll implement that later, this line is a dummy
		baseLogger.UrlParsed = baseLogger.Url + ";parsed"
		//the rest is throwing errors
	}
	baseLogger.Enableable = (baseLogger.Url != "")

	return baseLogger
}

//these functions could be redundant, basically setters
//in the java version these use generics, asking rob if there's any functionality needed
func (obj BaseLogger) Enable() {
	obj.Enabled = true
}

func (obj BaseLogger) Disable() {
	obj.Enabled = false
}

//End Getters and Setters

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

type BaseLogger struct {
	Agent      string
	Enableable bool
	Enabled    bool
	Host       string
	//easiest to implement a queue in go by using slices, need enqueue and dequque methods
	Queue           []string
	SkipCompression bool
	SkipSubmission  bool
	SubmitFailures  int
	SubmitSuccesses int
	Url             string
	//url parsed should be of type URL, but I can't get the package to import
	UrlParsed string
	Version   string
}
