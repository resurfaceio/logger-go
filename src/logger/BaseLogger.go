package logger

import (
	"bytes"
	"compress/flate"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"

	"github.com/asaskevich/govalidator"
)

// BaseLogger constructor
func NewBaseLogger(_agent string, _url string, _enabled interface{}, _queue []string) *BaseLogger {
	usageLoggers, _ := GetUsageLoggers()

	_enabled = (_enabled == nil) || (_enabled.(bool) == true)
	if _queue == nil && _url == "" {
		_url = usageLoggers.UrlByDefault()
		// WIP 03.26.21
		if _url == "" {
			_enabled = false
		}
	}

	var _urlParsed *url.URL
	var parsingError error
	//validate url when present
	if _url != "" {
		_urlParsed, parsingError = url.ParseRequestURI(_url)
		isUrl := govalidator.IsURL(_url)
		// fmt.Println("parsed url: " + _urlParsed.String())
		// fmt.Println("is Url:" + strconv.FormatBool(isUrl))
		if parsingError != nil || !isUrl {
			_url = ""
			_urlParsed = nil
			_enabled = false
		}
	}

	_enableable := (_url != "" || _queue != nil)

	constructedBaseLogger := &BaseLogger{
		agent:           _agent,
		enableable:      _enableable,
		enabled:         _enabled.(bool),
		host:            hostLookup(),
		queue:           _queue,
		skipCompression: false,
		skipSubmission:  false,
		submitFailures:  0,
		submitSuccesses: 0,
		url:             _url,
		urlParsed:       _urlParsed,
		version:         versionLookup(),
	}
	return constructedBaseLogger
}

func (logger *BaseLogger) Enable() {
	logger.enabled = logger.enableable
}

func (logger *BaseLogger) Disable() {
	logger.enabled = false
}

/**
 * Submits JSON message to intended destination.
 */
func (logger *BaseLogger) Submit(msg string) {
	//woah congrats you submitted the message
	if msg == "" || logger.SkipSubmission() || !logger.Enabled() {
		//do nothing
	} else if logger.queue != nil {
		logger.queue = append(logger.queue, msg)
		atomic.AddInt64(&logger.submitSuccesses, 1)
		return
	} else {
		// not 100% sure this works (needs testing) and should add some error handling
		submitRequest, err := http.NewRequest("POST", logger.url, bytes.NewBuffer([]byte(msg)))
		if err != nil {
			fmt.Printf("Error creating submit request: %s", err.Error())
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}

		submitRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")
		submitRequest.Header.Set("User-Agent", "Resurface/"+logger.version+" ("+logger.agent+")")

		if !logger.skipCompression {
			submitRequest.Header.Set("Content-Encoding", "deflated")

			var b bytes.Buffer
			w, err := flate.NewWriter(&b, 0)
			if err != nil {
				fmt.Printf("Error applying deflate compression to message: %s\n", err.Error())
				atomic.AddInt64(&logger.submitFailures, 1)
				return
			}

			_, err = w.Write([]byte(msg))
			if err != nil {
				fmt.Printf("Error writing message to io.Writer: %s\n", err.Error())
				atomic.AddInt64(&logger.submitFailures, 1)
				return
			}
			w.Close()

			err = submitRequest.Write(w)
			if err != nil {
				fmt.Printf("Error writing message to request: %s\n", err.Error())
				atomic.AddInt64(&logger.submitFailures, 1)
				return
			}
		}

		submitResponse, err := http.DefaultClient.Do(submitRequest)
		if err != nil {
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}

		if submitResponse.StatusCode == 204 {
			atomic.AddInt64(&logger.submitSuccesses, 1)
			return
		} else {
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}
	}
}

/**
 * Returns host identifier.
 * These are utility functions that can be static if this wasn't Go
 */
func hostLookup() string {
	dyno, dynoIsPresent := os.LookupEnv("DYNO")
	if dynoIsPresent && dyno != "" {
		return dyno
	}

	hostName, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return hostName
}

func versionLookup() string { return "1.0.0" }

func (logger *BaseLogger) Agent() string    { return logger.agent }
func (logger *BaseLogger) Enableable() bool { return logger.enableable }
func (logger *BaseLogger) Enabled() bool {
	usageLoggers, _ := GetUsageLoggers()
	return logger.enabled && usageLoggers.IsEnabled()
}
func (logger *BaseLogger) Host() string           { return logger.host }
func (logger *BaseLogger) Queue() []string        { return logger.queue }
func (logger *BaseLogger) SkipCompression() bool  { return logger.skipCompression }
func (logger *BaseLogger) SkipSubmission() bool   { return logger.skipSubmission }
func (logger *BaseLogger) SubmitFailures() int64  { return logger.submitFailures }
func (logger *BaseLogger) SubmitSuccesses() int64 { return logger.submitSuccesses }
func (logger *BaseLogger) Url() string            { return logger.url }
func (logger *BaseLogger) UrlParsed() *url.URL    { return logger.urlParsed }
func (logger *BaseLogger) Version() string        { return logger.version }

func (logger *BaseLogger) SetSkipCompression(_skipCompression bool) {
	logger.skipCompression = _skipCompression
}
func (logger *BaseLogger) SetSkipSubmission(_skipSubmission bool) {
	logger.skipSubmission = _skipSubmission
}

type BaseLogger struct {
	agent           string
	enableable      bool
	enabled         bool
	host            string
	queue           []string
	skipCompression bool
	skipSubmission  bool
	submitFailures  int64
	submitSuccesses int64
	url             string
	urlParsed       *url.URL
	version         string
}
