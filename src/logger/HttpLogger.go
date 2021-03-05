package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const loggerAgent string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct{
	baseLogger BaseLogger
	enabled bool
	queue []string
	skip_compression bool
	skip_submission bool
	rules string
	url string
}

 /**
 * Initialize logger using default url and default rules.
 */
func NewHttpLogger() *HttpLogger {
	baseLogger := NewBaseLoggerAgent(loggerAgent)
	return HttpLogger("", true, nil, "", baseLogger)
}

/**
* Initialize enabled/disabled logger using default url and default rules.
*/
func NewHttpLoggerEnabled(enabled bool) *HttpLogger {
	baseLogger := NewBaseLoggerAgentEnabled(loggerAgent, _enabled)
	return HttpLogger()
}

/**
* Initialize logger using specified url and default rules.
*/
func HttpLoggerUrl(url string) *HttpLogger {
	return
}

/**
* Initialize logger using specified url and specified rules.
*/
func HttpLoggerUrlRules(url string, rules string) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using specified url and default rules.
*/

func NewHttpLoggerUrlEnabled(url string, enabled bool) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using specified url and specified rules.
*/

func NewHttpLoggerUrlEnabledRules(url string, enabled bool, rules string) *HttpLogger{
	return
}

/**
* Initialize enabled logger using queue and default rules.
*/
func NewHttpLoggerQueue(queue []string) *HttpLogger {
	return
} 

/**
* Initialize enabled logger using queue and specified rules.
*/
func NewHttpLoggerQueueRules(queue []string, rules string) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using queue and default rules.
*/
func NewHttpLoggerQueueEnabled(queue []string, enabled bool) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using queue and specified rules.
*/
func NewHttpLoggerQueueEnabledRules(queue []string, enabled bool, rules string) *HttpLogger {
	return
}

// initialize HttpLogger 
func HttpLogger(_url string, _enabled bool, _queue []string, 
	_rules string, _baseLogger BaseLogger) (*HttpLogger, error) { 

	httpRules, err := NewHttpRules(_rules)

	if err != nil {
		return nil, err
	}

	exHttpLogger := &HttpLogger{
		baseLogger: _baseLogger,
		rules:	httpRules,
	}
	return exHttpLogger
}

//slightly confused on the implementation here
func initialize(ex_rules string){

	// specific rules we need to account for
	rules = new HttpRules(rules)

	// IS allow_http_url allowed in our golang implementation
	if((ex_url != nil) && (ex_url.HasPrefix("http:") && !ex_rules.allow_http_url)) {
		// WHERE DID enableable come from
		ex_enableable = false; // NOT DEFINED in struct yet!!!!!!!!!!!
		ex_enabled = false;
	}
	
}

// don't worry about yet
func HttpRules getRules(){
	return rules
}


func submitIfPassing(var details []string){
	details = rules.apply(details) // MMMMMMMMMMMMMMMMMMMM try again

	if details == nil {
		return
	}

	details = append(["host", host]) // where does host come from?

	submit(json.Unmarshal([]byte(details), &value))// mmmmmm 
}