package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: "Bad AccessToken",
	}

	result, err := client.FindUsers(SearchRequest{Query: "on", OrderField: "Age", OrderBy: 1})
	if result != nil && err.Error() != "Bad AccessToken" {
		t.Errorf("Token auth not working")
	}
}

func TestSearcherRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         "",
	}

	result, err := client.FindUsers(SearchRequest{})
	if result != nil || err == nil {
		t.Errorf("Expected error, but got result %v", result)
	}
}

func TestUnknownError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerUnknownError))
	searchClient := &SearchClient{
		URL: "bad_link",
	}

	_, err := searchClient.FindUsers(SearchRequest{})

	if err == nil {
		t.Error("TestUnknownError :(")
	}

	ts.Close()
}

func TestLimitAndOffset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: "Good access token",
	}

	result, err := client.FindUsers(SearchRequest{Limit: -1})
	if result != nil && err.Error() != "limit must be > 0" {
		t.Errorf("Wrong limit parameter, must be > 0")
	}

	result, err = client.FindUsers(SearchRequest{Offset: -1})
	if result != nil && err.Error() != "offset must be > 0" {
		t.Errorf("Wrong offset parameter, must be > 0")
	}

	result, err = client.FindUsers(SearchRequest{Limit: 35})
	if result != nil && len(result.Users) != 25 {
		t.Errorf("Wrong limit parameter, must be 25")
	}
}

func TestSort(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	client := SearchClient{
		URL:         ts.URL,
		AccessToken: "Good access token",
	}

	_, err := client.FindUsers(SearchRequest{OrderField: "Id", OrderBy: -1})
	if err.Error() == "Bad orderBy value" {
		t.Errorf("Wrong order parameter")
	}
	_, err1 := client.FindUsers(SearchRequest{OrderField: "Id", OrderBy: 0})
	if err1.Error() == "Bad orderBy value" {
		t.Errorf("Wrong order parameter")
	}
	_, err2 := client.FindUsers(SearchRequest{OrderField: "Id", OrderBy: 1})
	if err2.Error() == "Bad orderBy value" {
		t.Errorf("Wrong order parameter")
	}
}

func TestNextPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: "GoodToken",
	}

	result, err := client.FindUsers(SearchRequest{Query: "o"})
	if result == nil && err != nil {
		t.Errorf("Next Page error")
	}

	result, err = client.FindUsers(SearchRequest{Limit: 1})
	if result == nil && err != nil {
		t.Errorf("Next Page error")
	}
	result, err = client.FindUsers(SearchRequest{Limit: 30})
	if result == nil && err != nil {
		t.Errorf("Next Page error")
	}
	result, err = client.FindUsers(SearchRequest{Limit: 25, Offset: 1})
	if result == nil && err != nil {
		t.Errorf("Next Page error")
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	client := SearchClient{
		URL:         ts.URL,
		AccessToken: "Timeout",
	}

	_, err := client.FindUsers(SearchRequest{Limit: 1, Offset: 0, Query: "", OrderField: "", OrderBy: 0})
	if err != nil && err.Error() == "timeout for" {
		t.Errorf("Expected timeout error, got %s", err)
	}
}

// код писать тут
func SearchServer(w http.ResponseWriter, r *http.Request) {
	var query = r.URL.Query().Get("query")
	var orderField = r.URL.Query().Get("orderField")
	var orderBy, _ = strconv.Atoi(r.URL.Query().Get("orderBy"))
	var limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	var offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	var accessToken = r.URL.Query().Get("AccessToken")

	switch accessToken {
	case "Timeout":
		time.Sleep(4 * time.Second)
		return
	}

	var rawData Root

	xmlFile, openFileError := os.Open("dataset.xml")
	if openFileError != nil {
		fmt.Println(openFileError)
		return
	}

	xmlFileBytes, readFileError := io.ReadAll(xmlFile)
	unmarshalErr := xml.Unmarshal(xmlFileBytes, &rawData)

	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
		return
	}

	if readFileError != nil {
		fmt.Println(readFileError)
		return
	}

	result := rawData.Profiles

	if query != "" {
		result = findSubstringInFirstNameLastNameAndAbout(rawData.Profiles, query)
	}

	sortUsers(result, orderField, orderBy)

	result = limitUsers(result, limit)
	result = offsetUsers(result, offset)
}
