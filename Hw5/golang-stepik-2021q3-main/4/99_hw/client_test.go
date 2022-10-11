package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Row struct {
	Text      string `xml:",chardata"`
	ID        int    `xml:"id"`
	Guid      string `xml:"guid"`
	IsActive  string `xml:"isActive"`
	Balance   string `xml:"balance"`
	Picture   string `xml:"picture"`
	Age       int    `xml:"age"`
	EyeColor  string `xml:"eyeColor"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	Company   string `xml:"company"`
	Email     string `xml:"email"`
	Phone     string `xml:"phone"`
	Address   string `xml:"address"`
	About     string `xml:"about"`
}
type Dataset struct {
	XMLName xml.Name `xml:"root"`
	Text    string   `xml:",chardata"`
	Rows    []Row    `xml:"row"`
}
type TestCase struct {
	Request SearchRequest
	Error   string
	IsError bool
}

var dataset Dataset

func SearchServer(resp http.ResponseWriter, req *http.Request) {
	if accessToken := req.Header.Get("AccessToken"); accessToken == "" {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := req.ParseForm(); err != nil {
		message := fmt.Sprintf("Error parsing form: %s", err)
		http.Error(resp, message, http.StatusInternalServerError)
		return
	}
	query := req.Form.Get("query")
	var response []User
	limit, err := strconv.Atoi(req.Form.Get("limit"))
	if err != nil {
		message := fmt.Sprintf("Error parsing limit: %s", err)
		http.Error(resp, message, http.StatusInternalServerError)
		return
	}
	offset, err := strconv.Atoi(req.Form.Get("offset"))
	if err != nil {
		message := fmt.Sprintf("Error parsing offset: %s", err)
		http.Error(resp, message, http.StatusInternalServerError)
		return
	}

	file, err := os.Open("dataset.xml")
	if err != nil {
		http.Error(resp, "Fatal Error, cant open data file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(resp, "Fatal Error, cant read data file", http.StatusInternalServerError)
		return
	}
	err = xml.Unmarshal(data, &dataset)
	if err != nil {
		http.Error(resp, "Fatal Error, cant Unmarshal data file", http.StatusInternalServerError)
		return
	}

	for _, row := range dataset.Rows {
		name := row.FirstName + " " + row.LastName
		if strings.Contains(name, query) || strings.Contains(row.About, query) {
			response = append(response, User{
				Id:     row.ID,
				Name:   name,
				Age:    row.Age,
				About:  row.About,
				Gender: row.Gender,
			})
		}
	}
	orderField := req.Form.Get("order_field")
	orderBy, err := strconv.Atoi(req.Form.Get("order_by"))
	if err != nil || -1 > orderBy || orderBy > 1 {
		message, _ := json.Marshal(SearchErrorResponse{Error: "InvalidSortOrder"})
		http.Error(resp, string(message), http.StatusBadRequest)
		return
	}
	switch orderField {
	case "Id":
		if orderBy == OrderByAsc {
			sort.SliceStable(response, func(i, j int) bool { return response[i].Id < response[j].Id })
		} else if orderBy == OrderByDesc {
			sort.SliceStable(response, func(i, j int) bool { return response[i].Id > response[j].Id })
		}
	case "Age":
		if orderBy == OrderByAsc {
			sort.SliceStable(response, func(i, j int) bool { return response[i].Age < response[j].Age })
		} else if orderBy == OrderByDesc {
			sort.SliceStable(response, func(i, j int) bool { return response[i].Age > response[j].Age })
		}
	case "Name", "":
		if orderBy == OrderByAsc {
			sort.SliceStable(response, func(i, j int) bool { return response[i].Name < response[j].Name })
		} else if orderBy == OrderByDesc {
			sort.SliceStable(response, func(i, j int) bool { return response[i].Name > response[j].Name })
		}
	default:
		message, _ := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
		http.Error(resp, string(message), http.StatusBadRequest)
		return
	}
	if offset > len(response) {
		offset = len(response)
	}
	if limit == 0 {
		limit = len(response)
	}
	if offset+limit > len(response) {
		limit = len(response) - offset
	}
	response = response[offset : limit+offset]
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal to json:", response)
	}
	resp.Write(responseJSON)
}

func TestFindUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	testCases := []TestCase{
		TestCase{
			Request: SearchRequest{
				Limit:      100,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
			IsError: false,
			Error:   "Nope",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
			IsError: false,
			Error:   "Nope",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      100,
				Offset:     0,
				Query:      "",
				OrderField: "Gender",
				OrderBy:    OrderByAsc,
			},
			IsError: true,
			Error:   fmt.Sprintf("OrderFeld %s invalid", "Gender"),
		},
		TestCase{
			Request: SearchRequest{
				Limit:      100,
				Offset:     0,
				Query:      "",
				OrderField: "Age",
				OrderBy:    OrderByAsc,
			},
			IsError: false,
			Error:   "Nope",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      100,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			IsError: false,
			Error:   "Nope",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsc,
			},
			IsError: true,
			Error:   "limit must be > 0",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      0,
				Offset:     -1,
				Query:      "",
				OrderField: "",
				OrderBy:    OrderByAsc,
			},
			IsError: true,
			Error:   "offset must be > 0",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      0,
				Offset:     0,
				Query:      "",
				OrderField: "",
				OrderBy:    -2,
			},
			IsError: true,
			Error:   "unknown bad request error",
		},
	}
	for _, testCase := range testCases {
		client := &SearchClient{
			AccessToken: "secret",
			URL:         ts.URL,
		}
		result, err := client.FindUsers(testCase.Request)

		if err == nil && testCase.IsError {
			t.Errorf("Unexpected success: %#v", result)
		}
		if err != nil && !strings.Contains(err.Error(), testCase.Error) {
			t.Errorf("Unexpected error: %#v", err)
		}
	}
}

func TestAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	testCase := TestCase{
		Request: SearchRequest{
			Limit:      0,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    0,
		},
		IsError: true,
		Error:   "Bad AccessToken",
	}

	client := &SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}
	result, err := client.FindUsers(testCase.Request)

	if err == nil && testCase.IsError {
		t.Errorf("Unexpected success: %#v", result)
	}
	if err != nil && !strings.Contains(err.Error(), testCase.Error) {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestUnknownError(t *testing.T) {
	testCase := TestCase{
		Request: SearchRequest{
			Limit:      100,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    0,
		},
		IsError: true,
		Error:   "unknown error",
	}

	client := &SearchClient{
		AccessToken: "",
		URL:         "1",
	}
	result, err := client.FindUsers(testCase.Request)

	if err == nil && testCase.IsError {
		t.Errorf("Unexpected success: %#v", result)
	}
	if err != nil && !strings.Contains(err.Error(), testCase.Error) {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestBadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Fatal Error Test", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "TestBadResultJSON:{{{test....,}--}.,!!test:1")
	}))
	defer ts.Close()

	testCase := TestCase{
		Request: SearchRequest{
			Limit:      0,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    0,
		},
		IsError: true,
		Error:   "cant unpack error json",
	}

	client := &SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}
	result, err := client.FindUsers(testCase.Request)

	if err == nil && testCase.IsError {
		t.Errorf("Unexpected success: %#v", result)
	}
	if err != nil && !strings.Contains(err.Error(), testCase.Error) {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestBadResultJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "TestBadResultJSON:{{{test....,}--}.,!!test:1")
	}))
	defer ts.Close()

	testCase := TestCase{
		Request: SearchRequest{
			Limit:      0,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    0,
		},
		IsError: true,
		Error:   "cant unpack result json",
	}

	client := &SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}
	result, err := client.FindUsers(testCase.Request)

	if err == nil && testCase.IsError {
		t.Errorf("Unexpected success: %#v", result)
	}
	if err != nil && !strings.Contains(err.Error(), testCase.Error) {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func TestFatalError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Fatal Error Test", http.StatusInternalServerError)
	}))
	defer ts.Close()

	testCase := TestCase{
		Request: SearchRequest{
			Query:      "",
			OrderField: "",
			OrderBy:    0,
		},
		IsError: true,
		Error:   "SearchServer fatal error",
	}

	client := &SearchClient{
		AccessToken: "AccessToken",
		URL:         ts.URL,
	}
	result, err := client.FindUsers(testCase.Request)

	if err == nil && testCase.IsError {
		t.Errorf("Unexpected success: %#v", result)
	}
	if err != nil && !strings.Contains(err.Error(), testCase.Error) {
		t.Errorf("Unexpected error: %#v", err)
	}
}

func SearchServerTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1500 * time.Millisecond)
	var response []User
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	if err != nil {
		log.Println("Failed to marshal to json:", response)
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerTimeout))
	defer ts.Close()

	testCase := TestCase{
		Request: SearchRequest{
			Limit:      0,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    0,
		},
		IsError: true,
		Error:   "timeout for",
	}

	client := &SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}
	result, err := client.FindUsers(testCase.Request)

	if err == nil && testCase.IsError {
		t.Errorf("Unexpected success: %#v", result)
	}
	if err != nil && !strings.Contains(err.Error(), testCase.Error) {
		t.Errorf("Unexpected error: %#v", err)
	}
}
