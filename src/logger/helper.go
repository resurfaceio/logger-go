package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Helper struct {
<<<<<<< HEAD:src/logger/helper.go
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
=======
	DEMO_URL          string
	MOCK_AGENT        string
	MOCK_HTML         string
	MOCK_HTML2        string
	MOCK_HTML3        string
	MOCK_HTML4        string
	MOCK_HTML5        string
	MOCK_JSON         string
	MOCK_JSON_ESCAPED string
	MOCK_NOW          int64
	MOCK_QUERY_STRING string
	MOCK_URL          string
	MOCK_URLS_DENIED  []string
	MOCK_URLS_INVALID []string
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:test/helper.go
}

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

func allArticles(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Test Title", Desc: "Test Description", Content: "<html>Hello World!</html>"},
	}

	fmt.Println("Endpoint Hit: All Articles Endpoint")
	json.NewEncoder(w).Encode(articles)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Enpoint Hit")
}

func (h *Helper) MockCustomApp() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", allArticles)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func NewTestHelper() *Helper {
	newHelper := Helper{
<<<<<<< HEAD:src/logger/helper.go
		demoURL: "https://demo.resurface.io/ping",
=======
		DEMO_URL: "https://demo.resurface.io/ping",
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:test/helper.go

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

<<<<<<< HEAD:src/logger/helper.go
		mockURLSdenied: []string{"https://demo.resurface.io/ping",
=======
		MOCK_URLS_DENIED: []string{"https://demo.resurface.io/ping",
>>>>>>> 7c6f55721f57030eaff4cd21234ba3cece9a0c17:test/helper.go
			"/noway3is5this1valid2",
			"https://www.noway3is5this1valid2.com/"},

		mockURLSinvalid: []string{"",
			"noway3is5this1valid2",
			"ftp:\\www.noway3is5this1valid2.com/",
			"urn:ISSN:1535–3613"},
	}

	return &newHelper
}
