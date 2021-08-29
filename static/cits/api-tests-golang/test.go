package main

import (
    "bytes"
    "flag"
    "log"
    "fmt"
    "os"
    "path/filepath"
    "net/http"
    "regexp"
    "io/ioutil"
    "strings"
    "strconv"
)

func main() {
    testFiles := flag.String("t", "./data/", "path to test files")
    flag.Parse()

    client := &http.Client{}
    n_err := 0
    files, _ := WalkMatch(*testFiles, "*.http")
    for _, file := range files {
        tc := parse_http_file(file)
        if ok := tc.Test(client); !ok {
            n_err += 1
        }
    }

    if n_err > 0 {
        fmt.Printf("\n\n%d Errors caught\n", n_err)
        os.Exit(1)
    } else {
        os.Exit(0)
    }
}

type HttpTestCase struct {
    Name string
    Request *http.Request
    Expected struct {
        StatusCode int
        Headers map[string]string
        Body string
    }
}

func (tc *HttpTestCase) Pass() bool {
    fmt.Printf("✓ API:%s\n", tc.Name)
    return true
}

func (tc *HttpTestCase) Fail(reason string, expected interface{}, actual interface{}) bool {
    fmt.Printf("✘ %s: Failed because '%s'\n\tExpected: '%v'\n\tActual: '%v'\n", tc.Name, reason, expected, actual)
    return false
}

func (tc *HttpTestCase) Test(client *http.Client) bool {
    if tc.Request == nil {
        return tc.Fail("Could not parse test case", "", "")
    }

    resp, _ := client.Do(tc.Request)

    if resp.StatusCode != tc.Expected.StatusCode {
        return tc.Fail("Incorrect status code", tc.Expected.StatusCode, resp.StatusCode)
    }

    for header, expected := range tc.Expected.Headers {
        if actual := strings.ToLower(resp.Header.Get(header)); actual != expected {
            return tc.Fail("Header mismatch for " + header, expected, actual)
        }
    }

    if tc.Expected.Body != "" {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return tc.Fail("Failed to parse body", tc.Expected.Body, "")
        }
        actualBody := strings.TrimSpace(string(bodyBytes))
        if actualBody != tc.Expected.Body {
            return tc.Fail("Body does not match", tc.Expected.Body, actualBody)
        }
    }

    return tc.Pass()
}

func parseHeaders(headers string) map[string]string {
    ret := make(map[string]string)
    for _, header_line := range strings.Split(headers, "\n") {
        header_part := strings.Split(header_line, ":")
        if len(header_part) != 2 {
            continue
        }
        ret[strings.ToLower(strings.TrimSpace(header_part[0]))] = strings.ToLower(strings.TrimSpace(header_part[1]))
    }
    return ret
}

type MatchGroup struct {
    matches []string
    re *regexp.Regexp
}

func (m *MatchGroup) Get(key string) string {
    return strings.TrimSpace(m.matches[m.re.SubexpIndex(key)])
}

func (m *MatchGroup) Match(target string) bool {
    if m.re.MatchString(target) {
        m.matches = m.re.FindStringSubmatch(target)
        return true
    } else {
        return false
    }
}

func parse_http_file(fname string) *HttpTestCase {

    file, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Fatal(err)
    }

    // I just concatenated the patterns in the table above here
    // and gave the important capture groups some names with the ?P<Name> syntax
    http_file_format := regexp.MustCompile(`^(?P<Method>GET|POST|PUT|DELETE)\s+(?P<URL>[^\s]*)\s?(HTTP/1.1)?\n(?P<ReqHeader>(.*:.*\n)*)(?P<ReqBody>(.|\n)*)HTTP/1.1\s?(?P<StatusCode>[0-9]*)\s?.*\n(?P<RespHeader>(.*:.*\n)*)(?P<RespBody>(.|\n)*)`)

    // Matching and Getting
    tc := &HttpTestCase{Name: fname}
    mg := &MatchGroup{re: http_file_format}

    if mg.Match(string(file)) {
        reqBody := []byte(mg.Get("ReqBody"))
        tc.Request, _ = http.NewRequest(mg.Get("Method"), mg.Get("URL"), bytes.NewBuffer(reqBody))

        for k,v := range parseHeaders(mg.Get("ReqHeader")) {
            tc.Request.Header.Set(k,v)
        }

        tc.Expected.StatusCode, _ = strconv.Atoi(mg.Get("StatusCode"))
        tc.Expected.Headers = parseHeaders(mg.Get("RespHeader"))
        tc.Expected.Body = mg.Get("RespBody")
    }
    return tc
}

func WalkMatch(root, pattern string) ([]string, error) {
    var matches []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
            return err
        } else if matched {
            matches = append(matches, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return matches, nil
}