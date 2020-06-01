package goose

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestStaticRoutes(t *testing.T) {
	// tmp folder/files setup
	testDir, _ := os.Getwd()

	f, err := ioutil.TempFile(testDir, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())

	want := "I love goose!"
	_, err = f.WriteString(want)
	if err != nil {
		t.Error(err)
	}
	f.Close()

	dir, filename := filepath.Split(f.Name())

	// goose setup
	g := New()
	g.Static("/assets", dir)
	g.StaticFile("/favicon", f.Name())

	w1 := sendRequest(g, http.MethodGet, "/assets/"+filename)
	w2 := sendRequest(g, http.MethodGet, "/favicon")

	if !reflect.DeepEqual(w1, w2) {
		t.Error("Both ResponseRecorder should have the same contents.")
	}

	if w1.Code != http.StatusOK {
		t.Errorf("Got status code %d instead of 200.", w1.Code)
	}

	got := w1.Body.String()
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}

	want = "text/plain; charset=utf-8"
	got = w1.Header().Get("Content-Type")
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}
}
