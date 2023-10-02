// © 2016-2023 Graylog, Inc.

package logger

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asaskevich/govalidator"
)

type baseLogger struct {
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
	bundleSize      int
	msgQueue        chan string
	submitQueue     chan strings.Builder
	wg              sync.WaitGroup
	stop            chan bool
}

// BaseLogger constructor
func newBaseLogger(_agent string, _url string, _enabled interface{}, _queue []string) *baseLogger {
	usageLoggers, _ := GetUsageLoggers()

	_enabled = (_enabled == nil) || (_enabled.(bool))
	if _queue == nil && _url == "" {
		_url = usageLoggers.UrlByDefault()
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
		if parsingError != nil || !isUrl {
			_url = ""
			_urlParsed = nil
			_enabled = false
		}
	}

	_enableable := _url != "" || _queue != nil

	config := usageLoggers.ConfigByDefault()

	constructedBaseLogger := &baseLogger{
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
		bundleSize:      config["BUNDLE_SIZE"],
		msgQueue:        make(chan string, config["MESSAGE_QUEUE_SIZE"]),
		submitQueue:     make(chan strings.Builder, config["BUNDLE_QUEUE_SIZE"]),
		stop:            make(chan bool, 1),
	}

	constructedBaseLogger.wg.Add(1)
	go constructedBaseLogger.dispatcher()

	return constructedBaseLogger
}

func (logger *baseLogger) Enable() {
	logger.enabled = logger.enableable
}

func (logger *baseLogger) Disable() {
	logger.enabled = false
}

func (logger *baseLogger) worker() {
	defer logger.wg.Done()
work:
	for {
		submission, open := <-logger.submitQueue
		if submission.Len() > 0 {
			bundle := submission.String()
			logger.submit(bundle)
		}
		if !open {
			break work
		}
	}
}

func (logger *baseLogger) dispatcher() {
	defer logger.wg.Done()
	buffer := strings.Builder{}
	created := time.Now()
	logger.wg.Add(1)
	go logger.worker()
dispatch:
	for {
		select {
		case msg := <-logger.msgQueue:
			if msg != "" {
				if buffer.Len() < logger.bundleSize {
					buffer.WriteString(msg + "\n")
				} else {
					buffer.WriteString(msg)
					logger.submitQueue <- buffer
					buffer = strings.Builder{}
					created = time.Now()
				}
			}
		case flush := <-logger.stop:
			if flush {
				select {
				case msg, open := <-logger.msgQueue:
					buffer.WriteString(msg + "\n")
					if !open {
						for msg := range logger.msgQueue {
							buffer.WriteString(msg + "\n")
						}
					}
				default:
				}
			}
			if buffer.Len() != 0 {
				logger.submitQueue <- buffer
			}
			close(logger.submitQueue)
			break dispatch
		default:
			if buffer.Len() != 0 && time.Since(created) > time.Second {
				logger.submitQueue <- buffer
				buffer = strings.Builder{}
				created = time.Now()
			}
		}
	}
}

func (logger *baseLogger) ndjsonHandler(msg string) {
	if msg == "" || logger.skipSubmission || !logger.Enabled() {
		//do nothing
	} else if logger.queue != nil {
		logger.queue = append(logger.queue, msg)
		atomic.AddInt64(&logger.submitSuccesses, 1)
		return
	} else {
		logger.msgQueue <- msg
	}
}

/**
 * Submits JSON message to intended destination.
 */
func (logger *baseLogger) submit(msg string) {

	var submitRequest *http.Request
	var reqError error

	if !logger.skipCompression { // Compression will not be skipped

		var body bytes.Buffer

		zWriter := zlib.NewWriter(&body)

		b, err := zWriter.Write([]byte(msg))
		if err != nil || b != len([]byte(msg)) {
			log.Println("error compressing log: ", err)
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}

		err = zWriter.Close()

		if err != nil {
			log.Println("error closing compression writer: ", err)
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}

		submitRequest, reqError = http.NewRequest("POST", logger.url, &body)

		if reqError != nil {
			fmt.Printf("Error creating submit request: %s", reqError.Error())
			log.Println("Error making submit request...")
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}

		submitRequest.Header.Set("Content-Encoding", "deflated")
		submitRequest.Header.Set("Content-Type", "application/ndjson; charset=UTF-8")
		submitRequest.Header.Set("User-Agent", "Resurface/"+logger.version+" ("+logger.agent+")")

	} else { // Compression will be skipped

		submitRequest, reqError = http.NewRequest("POST", logger.url, bytes.NewBuffer([]byte(msg)))

		if reqError != nil {
			fmt.Printf("Error creating submit request: %s", reqError.Error())
			atomic.AddInt64(&logger.submitFailures, 1)
			return
		}

		submitRequest.Header.Set("Content-Type", "application/ndjson; charset=UTF-8")
		submitRequest.Header.Set("User-Agent", "Resurface/"+logger.version+" ("+logger.agent+")")
	}

	submitResponse, err := httpLoggerClient.Do(submitRequest)

	if err != nil {
		atomic.AddInt64(&logger.submitFailures, 1)
		return
	}
	if submitResponse != nil && submitResponse.StatusCode == 204 {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(submitResponse.Body)
		_, err := io.ReadAll(submitResponse.Body)

		if err != nil {
			log.Println(err)
		}

		atomic.AddInt64(&logger.submitSuccesses, 1)
		return
	} else {
		if submitResponse == nil {
			log.Println("Response is nil")
		} else {
			log.Println("Response from fluke:", submitResponse.StatusCode)
		}
		log.Println("An unknown error occurred")
		atomic.AddInt64(&logger.submitFailures, 1)
		return
	}

}

func (logger *baseLogger) stopDispatcher() {
	logger.Disable()
	logger.stop <- true
	logger.wg.Wait()
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

func versionLookup() string {
	version := "3.3.1"
	return version
}

func (logger *baseLogger) Enabled() bool {
	usageLoggers, _ := GetUsageLoggers()
	return logger.enabled && usageLoggers.IsEnabled()
}

func (logger *baseLogger) Queue() []string {
	return logger.queue
}
