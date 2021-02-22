package logger

import "sync"

var once sync.Once

type helper struct {
	demoURL         string
	mockAgent       string
	mockHTML        string
	mockHTML2       string
	mockHTML3       string
	mockHTML4       string
	mockHTML5       string
	mockJSON        string
	mockJSONescaped string
	mockNow         int64
	mockQueryString string
	mockURL         string
	mockURLSdenied  []string
	mockURLSinvalid []string
}

var testHelper *helper

func GetTestHelper() *helper {
	once.Do(func() {
		testHelper = &helper{
			demoURL: "https://demo.resurface.io/ping",

			mockAgent: "helper.java",

			mockHTML: "<html>Hello World!</html>",

			mockHTML2: "<html>Hola Mundo!</html>",

			mockHTML3: "<html>1 World 2 World Red World Blue World!</html>",

			mockHTML4: "<html>1 World\n2 World\nRed World \nBlue World!\n</html>",

			mockHTML5: `<html>\n
			<input type=\"hidden\">SENSITIVE</input>\n
			<input class='foo' type=\"hidden\">\n
			SENSITIVE\n
			</input>\n
			</html>`,

			mockJSON: "{ \"hello\" : \"world\" }",

			mockJSONescaped: "{ \\'hello\\' : \\'world\\' }",

			mockNow: 1455908640173,

			mockQueryString: "foo=bar",

			mockURL: "http://something.com:3000/index.html",

			mockURLSdenied: []string{"https://demo.resurface.io/ping",
				"/noway3is5this1valid2",
				"https://www.noway3is5this1valid2.com/"},

			mockURLSinvalid: []string{"",
				"noway3is5this1valid2",
				"ftp:\\www.noway3is5this1valid2.com/",
				"urn:ISSN:1535–3613"},
		}
	})

	return testHelper
}
