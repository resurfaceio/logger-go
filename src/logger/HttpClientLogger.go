package logger

import (

)

func NewHttpClientLogger() {

}

type HttpClientLogger struct {
	baseLogger BaseLogger
	HttpRules rules
}