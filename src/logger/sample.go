package logger

import "testing"

func TestDatabaseSubmission(t *testing.T) {

	options := Options{
		Url:   "https://lastn-resurfaceio.herokuapp.com/messages",
		Rules: "include debug",
	}

	client, _ := NewNetHttpClientLoggerOptions(options)
	// client.httpLogger.BaseLogger.skipCompression = true

	client.Get("https://www.thecocktaildb.com/api/json/v1/1/search.php?s=margarita")

	// client.Get("https://www.thecocktaildb.com/api/json/v1/1/search.php?s=manhattan")
	t.Error()
	// client.Get("")
}
