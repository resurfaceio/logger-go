package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HttpLogger struct{
	agent string
	enabled bool
	queue []string
	skipCompression bool
	skipSubmission bool
	rules string
	url string


}

func initialize(var rules string){

	rules = new HttpRules(rules)

	
}

func HttpRules getRules(){
	return rules
}

func submitIfPassing(var details []string){
	details = rules.apply(details)

	if details == nil {
		return
	}

	details = append(details, "")

	submit()
}