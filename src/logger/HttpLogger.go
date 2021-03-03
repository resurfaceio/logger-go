package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HttpLogger struct{

}

func initialize(var rules string){

	this.rules = new HttpRules(rules)

	
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