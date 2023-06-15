package logger

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat/combin"
)

var logger *HttpLogger

// helpers
func setup() *HttpLogger {
	if logger == nil {
		var err error
		logger, err = NewHttpLogger(Options{
			Queue: make([]string, 0),
			Rules: "include debug",
		})
		if err != nil {
			log.Panicln(err)
		}
	}
	return logger
}

func retrieveLastMessage(t *testing.T, logger *HttpLogger) map[string]string {
	msg := make([][]string, 1)
	queue := logger.Queue()

	err := json.Unmarshal([]byte(queue[len(queue)-1]), &msg)
	if err != nil {
		log.Panicln(err)
	}

	msgMap := make(map[string]string)
	for _, m := range msg {
		msgMap[m[0]] = m[1]
	}

	return msgMap
}

func deflateIt(payload []byte) bytes.Buffer {
	var deflateBuffer bytes.Buffer
	deflateWriter := zlib.NewWriter(&deflateBuffer)
	_, err := deflateWriter.Write(payload)
	if err != nil {
		log.Panicln(err)
	}
	deflateWriter.Flush()
	deflateWriter.Close()

	return deflateBuffer
}

func gzipIt(payload []byte) bytes.Buffer {
	var gzipBuffer bytes.Buffer

	gzipWriter := gzip.NewWriter(&gzipBuffer)
	_, err := gzipWriter.Write(payload)
	if err != nil {
		log.Panicln(err)
	}
	gzipWriter.Flush()
	gzipWriter.Close()

	return gzipBuffer
}

func brIt(payload []byte) bytes.Buffer {
	var brBuffer bytes.Buffer

	brWriter := brotli.NewWriter(&brBuffer)
	_, err := brWriter.Write(payload)
	if err != nil {
		log.Panicln(err)
	}
	brWriter.Flush()
	brWriter.Close()

	return brBuffer
}

// tests
func TestIdentityPlainText(t *testing.T) {
	logger := setup()

	expectedBody := string(bodies["plain"])

	req := MockGetNoBodyRequest()
	res := MockGetPlainTextResponse(&req)
	res.Header.Add("Content-Encoding", "identity")

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

	res = MockGetPlainTextResponse(&req)
	res.Header.Add("Content-Encoding", "")

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap = retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestIdentityPlainTextNoHeader(t *testing.T) {
	logger := setup()

	expectedBody := string(bodies["plain"])

	req := MockGetNoBodyRequest()
	res := MockGetPlainTextResponse(&req)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestIdentityJSON(t *testing.T) {
	logger := setup()

	expectedBody := string(bodies["json"])

	req := MockGetNoBodyRequest()
	res := MockGetJSONResponse(&req)
	res.Header.Add("Content-Encoding", "identity")

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

	res = MockGetJSONResponse(&req)
	res.Header.Add("Content-Encoding", "")

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap = retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestIdentityJSONNoHeader(t *testing.T) {
	logger := setup()

	expectedBody := string(bodies["json"])

	req := MockGetNoBodyRequest()
	res := MockGetJSONResponse(&req)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestDeflate(t *testing.T) {
	logger := setup()

	decoded, err := zlib.NewReader(bytes.NewReader(bodies["deflate"]))
	if err != nil {
		log.Panicln(err)
	}
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()
	res := MockGetDeflateResponse(&req)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestDeflateNoHeader(t *testing.T) {
	logger := setup()

	decoded, err := zlib.NewReader(bytes.NewReader(bodies["deflate"]))
	if err != nil {
		log.Panicln(err)
	}
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()
	res := MockGetDeflateResponse(&req)

	res.Header.Del("Content-Encoding")
	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestGzip(t *testing.T) {
	logger := setup()

	decoded, err := gzip.NewReader(bytes.NewReader(bodies["gzip"]))
	if err != nil {
		log.Panicln(err)
	}
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()
	res := MockGetGzipResponse(&req)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestGzipNoHeader(t *testing.T) {
	logger := setup()

	decoded, err := gzip.NewReader(bytes.NewReader(bodies["gzip"]))
	if err != nil {
		log.Panicln(err)
	}
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()
	res := MockGetGzipResponse(&req)
	res.Header.Del("Content-Encoding")
	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestBrotli(t *testing.T) {
	logger := setup()

	decoded := brotli.NewReader(bytes.NewReader(bodies["br"]))
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()
	res := MockGetBrotliResponse(&req)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestBrotliNoHeader(t *testing.T) {
	logger := setup()

	decoded := brotli.NewReader(bytes.NewReader(bodies["br"]))
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()
	res := MockGetBrotliResponse(&req)
	res.Header.Del("Content-Encoding")

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	if brotliFallbackEnabled {
		assert.Equal(t, expectedBody, msgMap["response_body"])
	} else {
		// The way algo is implemented, br can never be decoded if Content-Encoding header != "br"
		assert.NotEqual(t, expectedBody, msgMap["response_body"])
	}

}

func TestMisleadingMagic(t *testing.T) {
	logger := setup()

	// payload is actually br but it looks like gzip (case where magic bytes are wrong)
	payload := bodies["misleading"]
	decoded := brotli.NewReader(bytes.NewReader(payload))
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()

	headers := map[string][]string{
		"Content-Length":   {strconv.Itoa(len(payload))},
		"Content-Type":     {"application/json"},
		"Content-Encoding": {"gzip"},
	}
	res := mockGetCustomResponse(&req, headers, payload)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

	// double misleading (CEH say it's deflate, magic bytes look like gzip, but it's actually br)
	headers["Content-Encoding"][0] = "deflate"
	res = mockGetCustomResponse(&req, headers, payload)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap = retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

	// only magic bytes are misleading
	headers["Content-Encoding"][0] = "br"
	res = mockGetCustomResponse(&req, headers, payload)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap = retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

}

func TestMisleadingMagicIdentity(t *testing.T) {
	logger := setup()

	// payload is actually not encoded but it looks like gzip (case where magic bytes are wrong)
	payload := []byte{31, 139}
	payload = append(payload, bodies["plain"]...)

	// json.Unmarshall replaces invalid UTF-8 characters with Unicode replacement character U+FFFD (65533 in decimal)
	expectedBody := string([]rune{31, 65533}) + string(payload)[2:]
	req := MockGetNoBodyRequest()

	headers := map[string][]string{
		"Content-Length":   {strconv.Itoa(len(payload))},
		"Content-Type":     {"application/json"},
		"Content-Encoding": {"gzip"},
	}
	res := mockGetCustomResponse(&req, headers, payload)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

	// double misleading (CEH say it's deflate, magic bytes look like gzip, but it's actually not encoded)
	headers["Content-Encoding"][0] = "deflate"
	res = mockGetCustomResponse(&req, headers, payload)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap = retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

	// only magic bytes are misleading
	headers["Content-Encoding"][0] = "identity"
	res = mockGetCustomResponse(&req, headers, payload)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap = retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])
}

func TestUnreadable(t *testing.T) {
	logger := setup()

	payload := []string{string(bodies["jpg"])}
	b, e := json.Marshal(payload)
	if e != nil {
		log.Println(e)
	}

	var expectedPayload []string
	e = json.Unmarshal(b, &expectedPayload)
	if e != nil {
		log.Println(e)
	}
	expectedBody := expectedPayload[0]

	req := MockGetNoBodyRequest()
	res := MockGetJPEGResponse(&req)

	SendHttpMessage(logger, &res, &req, 0, 0, nil)
	msgMap := retrieveLastMessage(t, logger)

	assert.Equal(t, expectedBody, msgMap["response_body"])

}

func TestMultipleEncodings(t *testing.T) {
	logger := setup()

	expectedBody := string(bodies["json"])

	encodings := []string{"identity", "deflate", "gzip", "br"}

	// initialize permutations slice
	n := len(encodings)
	var sum int // sum(nPr) with n = len(encodings) for 1 <= r <= n
	for r := 1; r <= n; r++ {
		sum += combin.NumPermutations(n, r)
	}
	permutations := make([][]string, sum)

	// compute permutations
	var offset int
	for r := 1; r <= n; r++ {
		nPr := combin.Permutations(n, r)
		for i, row := range nPr {
			permutations[offset+i] = make([]string, r)
			for j, index := range row {
				permutations[offset+i][j] = encodings[index]
			}
		}
		offset += len(nPr)
	}

	// append each encoding to all items to have at least one repeated encoding
	repeated := make([][]string, n*(len(permutations)-n))
	i := -1
	for _, encoding := range permutations[:n] {
		ie := i + 1
		for ip, permutation := range permutations[n:] {
			i = ie + ip
			repeated[i] = append(permutation, encoding[0])
		}
	}

	for _, perm := range append(repeated, permutations...) {
		payload := bodies["json"]

		for _, enc := range perm {
			var buff bytes.Buffer
			switch enc {
			case "deflate":
				buff = deflateIt(payload)
			case "gzip":
				buff = gzipIt(payload)
			case "br":
				buff = brIt(payload)
			case "", "identity":
				continue
			}
			payload = buff.Bytes()
		}

		headers := map[string][]string{
			"Content-Length":   {strconv.Itoa(len(payload))},
			"Content-Type":     {"application/json"},
			"Content-Encoding": {strings.Join(perm, ",")},
		}

		req := MockGetNoBodyRequest()
		res := mockGetCustomResponse(&req, headers, payload)

		SendHttpMessage(logger, &res, &req, 0, 0, nil)
		msgMap := retrieveLastMessage(t, logger)

		assert.Equal(t, expectedBody, msgMap["response_body"])
	}

}

func TestMisconfiguredIdentity(t *testing.T) {
	logger := setup()

	expectedBody := string(bodies["json"])

	req := MockGetNoBodyRequest()

	for _, enc := range []string{"deflate", "gzip", "br"} {
		res := MockGetJSONResponse(&req)
		res.Header["Content-Encoding"] = []string{enc}

		SendHttpMessage(logger, &res, &req, 0, 0, nil)

		msgMap := retrieveLastMessage(t, logger)
		assert.Equal(t, expectedBody, msgMap["response_body"], "Content-Encoding header value: "+enc)
	}

}

func TestMisconfiguredDeflate(t *testing.T) {
	logger := setup()

	decoded, err := zlib.NewReader(bytes.NewReader(bodies["deflate"]))
	if err != nil {
		log.Panicln(err)
	}
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()

	for _, enc := range []string{"identity", "gzip", "br"} {
		res := MockGetDeflateResponse(&req)
		res.Header["Content-Encoding"] = []string{enc}

		SendHttpMessage(logger, &res, &req, 0, 0, nil)

		msgMap := retrieveLastMessage(t, logger)
		assert.Equal(t, expectedBody, msgMap["response_body"], "Content-Encoding header value: "+enc)
	}

}

func TestMisconfiguredGzip(t *testing.T) {
	logger := setup()

	decoded, err := gzip.NewReader(bytes.NewReader(bodies["gzip"]))
	if err != nil {
		log.Panicln(err)
	}
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()

	for _, enc := range []string{"identity", "deflate", "br"} {
		res := MockGetGzipResponse(&req)
		res.Header["Content-Encoding"] = []string{enc}

		SendHttpMessage(logger, &res, &req, 0, 0, nil)

		msgMap := retrieveLastMessage(t, logger)
		assert.Equal(t, expectedBody, msgMap["response_body"], "Content-Encoding header value: "+enc)
	}

}

func TestMisconfiguredBrotli(t *testing.T) {
	logger := setup()

	decoded := brotli.NewReader(bytes.NewReader(bodies["br"]))
	decodedBytes, err := ioutil.ReadAll(decoded)
	if err != nil {
		log.Panicln(err)
	}
	expectedBody := string(decodedBytes)

	req := MockGetNoBodyRequest()

	for _, enc := range []string{"identity", "gzip", "deflate"} {
		res := MockGetBrotliResponse(&req)
		res.Header["Content-Encoding"] = []string{enc}

		SendHttpMessage(logger, &res, &req, 0, 0, nil)

		msgMap := retrieveLastMessage(t, logger)
		if brotliFallbackEnabled {
			assert.Equal(t, expectedBody, msgMap["response_body"], "Content-Encoding header value: "+enc)
		} else {
			// The way algo is implemented, br can never be decoded if Content-Encoding header != "br"
			assert.NotEqual(t, expectedBody, msgMap["response_body"], "Content-Encoding header value: "+enc)
		}
	}
}
