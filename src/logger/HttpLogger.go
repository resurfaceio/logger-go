package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var AGENT string = "HttpLogger.go"

//base HttpLogger definition
type HttpLogger struct{
	enabled bool
	queue []string
	skip_compression bool
	skip_submission bool
	rules string
	url string
	BaseLogger // HMMMM?
}

// initialize HttpLogger either function?? parameters to pass
func NewHttpLogger(ex_url string, ex_enabled bool, ex_queue []string, 
				   ex_rules string) *HttpLogger { 

exHttpLogger := &HttpLogger{
	enabled: 			ex_enabled,
	queue: 				ex_queue,
	skip_compression:	false,
	skip_submission:	false,
	rules:				ex_rules,
	url:				ex_url
	}
	return exHttpLogger
}

// Some of them say initializing a logger using default url and rules but don't
// actually pass those as parameters

 /**
 * Initialize logger using default url and default rules.
 */
func NewNewHttpLogger() *HttpLogger {
	initialize(nil)
	return NewHttpLogger("", true, nil, "")

}

/**
* Initialize enabled/disabled logger using default url and default rules.
*/
func NewHttpLoggerEnabled(ex_enabled bool) *HttpLogger {
	return NewHttpLogger()
}

/**
* Initialize logger using specified url and default rules.
*/
func NewHttpLoggerUrl(ex_url string) *HttpLogger {
	return
}

/**
* Initialize logger using specified url and specified rules.
*/
func NewHttpLoggerUrlRules(ex_url string, ex_rules string) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using specified url and default rules.
*/

func NewHttpLoggerUrlEnabled(ex_url string, ex_enabled bool) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using specified url and specified rules.
*/

func NewHttpLoggerUrlEnabledRules(ex_url string, ex_enabled bool, ex_rules string) *HttpLogger{
	return
}

/**
* Initialize enabled logger using queue and default rules.
*/
func NewHttpLoggerQueue(ex_queue []string) *HttpLogger {
	return
} 

/**
* Initialize enabled logger using queue and specified rules.
*/
func NewHttpLoggerQueueRules(ex_queue []string, ex_rules string) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using queue and default rules.
*/
func NewHttpLoggerQueueEnabled(ex_queue []string, ex_enabled bool) *HttpLogger {
	return
}

/**
* Initialize enabled/disabled logger using queue and specified rules.
*/
func NewHttpLoggerQueueEnabledRules(ex_queue []string, ex_enabled bool, ex_rules string) *HttpLogger {
	return
}

/**
* Initialize a new logger.
*/

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