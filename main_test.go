package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	payload1 := "something"
	res1 := putBlob(t, ts.URL, "data", payload1)

	if int(res1.Size) != len(payload1) {
		t.Errorf("expected a size of %d, got %d", len(payload1), res1.Size)
	}

	payload2 := "other"
	res2 := putBlob(t, ts.URL, "data", payload2)

	if int(res2.Size) != len(payload2) {
		t.Errorf("expected a size of %d, got %d", len(payload2), res2.Size)
	}

	if res1.OID == res2.OID {
		t.Errorf("expected to have unique oid")
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	content1 := "something"
	res1 := putBlob(t, ts.URL, "my-repo", content1)

	body, status := getBlob(t, ts.URL, "my-repo", res1.OID)
	if status != http.StatusOK {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusOK, status)
	}
	if body != content1 {
		t.Errorf("expected content of %s, got %s", content1, body)
	}
}

func TestGetNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	_, status := getBlob(t, ts.URL, "my-repo", "missing")
	if status != http.StatusNotFound {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusNotFound, status)
	}
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	content := "something"
	res := putBlob(t, ts.URL, "my-repo", content)

	dupContent := "something"
	dupRes := putBlob(t, ts.URL, "other-repo", dupContent)

	status := deleteBlob(t, ts.URL, "my-repo", res.OID)
	if status != http.StatusOK {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusOK, status)
	}

	_, status = getBlob(t, ts.URL, "my-repo", res.OID)
	if status != http.StatusNotFound {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusNotFound, status)
	}

	dupBlob, status := getBlob(t, ts.URL, "other-repo", dupRes.OID)
	if status != http.StatusOK {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusOK, status)
	}
	if dupBlob != dupContent {
		t.Errorf("expected %s got %s", dupContent, dupBlob)
	}
}

func TestDeleteNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	status := deleteBlob(t, ts.URL, "my-repo", "missing")
	if status != http.StatusNotFound {
		t.Errorf("expected HTTP status of %d, got %d", http.StatusNotFound, status)
	}
}

type _response struct {
	OID  string `json:"oid"`
	Size int64  `json:"size"`
}

func getBlob(t *testing.T, serverURL, repo, oid string) (string, int) {
	path := fmt.Sprintf("%s/data/%s/%s", serverURL, repo, oid)
	res, err := http.Get(path)
	if err != nil {
		t.Fatalf("error making GET request: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading GET response: %s", err)
	}

	return string(body), res.StatusCode
}

func putBlob(t *testing.T, serverURL, repo, payload string) *_response {
	path := fmt.Sprintf("%s/data/%s", serverURL, repo)
	req, err := http.NewRequest("PUT", path, strings.NewReader(payload))
	if err != nil {
		t.Fatalf("creating req %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making PUT request: %s", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected response code 201, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %s", contentType)
	}

	var data _response
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	return &data
}

func deleteBlob(t *testing.T, serverURL, repo, oid string) int {
	path := fmt.Sprintf("%s/data/%s/%s", serverURL, repo, oid)
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		t.Fatalf("creating req %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making DELETE request: %s", err)
	}

	return res.StatusCode
}
