package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

type Helper struct {
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
	err := json.NewEncoder(w).Encode(articles)
	if err != nil {
		fmt.Println("Helper json encoding failed")
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Enpoint Hit")
}

func (h *Helper) MockCustomApp() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", allArticles)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

//This is a rough outline.
//I am thinking we need to cover the requests rather than the responses.
//MockGetRequest covers a get request to compare against loggin.
func (h *Helper) MockGetRequest(url, html) {

}

func (h *Helper) MockDoRequest() {

}

func (h *Helper) MockHeadRequest() {

}

func (h *Helper) MockPostRequest() {

}

func (h *Helper) MockPostFormRequest() {

}

// https://github.com/gorilla/mux
// This could be server side though only examples I can find for client
// is if we initilize a struct http client
func handleMockRequest(){
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/articles", allArticles).
		Host(url).
		Methods("GET").
		Schemes("http")

}

func NewTestHelper() *Helper {
	newHelper := Helper{
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

	return &newHelper
}
