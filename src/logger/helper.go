package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Helper struct {
	DEMO_URL string
	MOCK_AGENT string
	MOCK_HTML string
	MOCK_HTML2 string
	MOCK_HTML3 string
	MOCK_HTML4 string
	MOCK_HTML5 string
	MOCK_JSON string
	MOCK_JSON_ESCAPED string
	MOCK_NOW int64
	MOCK_QUERY_STRING string
	MOCK_URL string
	MOCK_URLS_DENIED []string
	MOCK_URLS_INVALID []string

}

type Article struct {
	Title string `json:"Title"`
	Desc string `json:"desc"`
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
		DEMO_URL: "https://demo.resurface.io/ping",

		MOCK_AGENT: "helper.java",

		MOCK_HTML: "<html>Hello World!</html>",

		MOCK_HTML2: "<html>Hola Mundo!</html>",

		MOCK_HTML3: "<html>1 World 2 World Red World Blue World!</html>",

		MOCK_HTML4: "<html>1 World\n2 World\nRed World \nBlue World!\n</html>",

		MOCK_HTML5: `<html>\n
		<input type=\"hidden\">SENSITIVE</input>\n
		<input class='foo' type=\"hidden\">\n
		SENSITIVE\n
		</input>\n
		</html>`,

		MOCK_JSON: "{ \"hello\" : \"world\" }",

		MOCK_JSON_ESCAPED: "{ \\'hello\\' : \\'world\\' }",

		MOCK_NOW: 1455908640173,

		MOCK_QUERY_STRING: "foo=bar",

		MOCK_URL: "http://something.com:3000/index.html",

		MOCK_URLS_DENIED: []string{ "https://demo.resurface.io/ping",
			"/noway3is5this1valid2",
			"https://www.noway3is5this1valid2.com/"},

		MOCK_URLS_INVALID: []string{"",
			"noway3is5this1valid2",
			"ftp:\\www.noway3is5this1valid2.com/",
			"urn:ISSN:1535–3613"},

	}

	return &newHelper
}