// © 2016-2022 Resurface Labs Inc.

package logger

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"

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
	msgQueue        chan string
	submitMutex     sync.Mutex
	poisonChan      chan string
	wg              sync.WaitGroup
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

	_enableable := (_url != "" || _queue != nil)

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
		msgQueue:        make(chan string, 10000),
		submitMutex:     sync.Mutex{},
		poisonChan:      make(chan string, 1),
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

func (logger *baseLogger) worker(buffer strings.Builder) {
	defer logger.wg.Done()
	logger.submitMutex.Lock()
	bundle := buffer.String()
	logger.submit(bundle)
	logger.submitMutex.Unlock()
}

func (logger *baseLogger) dispatcher() {
	defer logger.wg.Done()
	// Threshold that determines when NDJSON bundles are sent to Resurface
	thresh := 100
	buffer := strings.Builder{}
	var msg, poisonpill string
dispatch:
	for {
		select {
		case msg = <-logger.msgQueue:
			if buffer.Len() != thresh {
				buffer.WriteString(msg + "\n")
			} else {
				buffer.WriteString(msg)
				logger.wg.Add(1)
				go logger.worker(buffer)
				buffer = strings.Builder{}
			}
		case poisonpill = <-logger.poisonChan:
		default:
			if buffer.Len() != 0 {
				logger.wg.Add(1)
				go logger.worker(buffer)
				buffer = strings.Builder{}
			}
			if poisonpill == "POISON PILL" {
				break dispatch
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
		defer submitResponse.Body.Close()
		_, err := ioutil.ReadAll(submitResponse.Body)

		if err != nil {
			log.Println(err)
		}

		atomic.AddInt64(&logger.submitSuccesses, 1)
		return
	} else {
		if submitResponse == nil {
			log.Println("Response is nil")
		}
		log.Println("An unknown error occurred")
		atomic.AddInt64(&logger.submitFailures, 1)
		return
	}

}

func (logger *baseLogger) stopDispatcher() {
	logger.poisonChan <- "POISON PILL"
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
	version := "3.0.0"
	return version
}

func (logger *baseLogger) Enabled() bool {
	usageLoggers, _ := GetUsageLoggers()
	return logger.enabled && usageLoggers.IsEnabled()
}

func (logger *baseLogger) Queue() []string {
	return logger.queue
}
