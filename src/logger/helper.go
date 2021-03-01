package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	mockFormData    string
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
//https://appdividend.com/2019/12/02/golang-http-example-get-post-http-requests-in-golang/
func MockGetRequest() http.Request {
	helper := NewTestHelper()
	resp, err := http.Get(helper.demoURL)
	request := resp.Request
	return *request
}

func MockDoRequest() http.Request {
	helper := NewTestHelper()
	// No do method.
	//resp, err := http.
	request := resp.Request
	return *request
}

func MockHeadRequest() http.Request {
	helper := NewTestHelper()
	resp, err := http.Head(helper.demoURL)
	request := resp.Request
	return *request
}

func MockPostRequest() http.Request {
	helper := NewTestHelper()
	resp, err := http.Post(helper.demoURL, "html", bytes.NewBuffer([]byte(helper.mockJSON)))
	request := resp.Request
	return *request
}

func (h *Helper) MockPostFormRequest() {
	helper := NewTestHelper()
	resp, err := http.PostForm(helper.demoURL, helper.mockFormData)
	request := resp.Request
	return *request
}

// // https://github.com/gorilla/mux
// // This could be server side though only examples I can find for client
// // is if we initilize a struct http client
// func handleMockRequest() {
// 	r := mux.NewRouter().StrictSlash(true)
// 	r.HandleFunc("/", homePage)
// 	r.HandleFunc("/articles", allArticles).
// 		Host(url).
// 		Methods("GET").
// 		Schemes("http")

// }

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

		mockFormData: "\"username\": { \" ResurfaceIO \" ",
	}

	return &newHelper
}
