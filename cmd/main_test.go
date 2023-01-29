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

const cType = "application/json; charset=utf-8"

func TestCreate(t *testing.T) {
	OpenFileJson()
	s := httptest.NewServer(http.HandlerFunc(Create))
	defer s.Close()
	var str User = User{Name: "Ilya", Age: 23}
	k, _ := json.Marshal(str)
	post, err := http.Post(s.URL+"/create", cType, strings.NewReader(string(k)))
	if err != nil {
		t.Fatal(err)
	}
	var f *User
	json.Unmarshal(k, &f)
	z := StorageMap[Id]

	if z.Name != f.Name {
		t.Errorf("Запрос не соответствует")
	}
	if z.Age != f.Age {
		t.Errorf("Запрос не соответствует")
	}
	if post.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", post.StatusCode)
	}
}

func TestMakeFriends(t *testing.T) {
	OpenFileJson()
	s := httptest.NewServer(handler())
	defer s.Close()
	var str Frends = Frends{SourceId: 2, TargetId: 3}
	k, _ := json.Marshal(str)
	post, err := http.Post(s.URL+"/make_friends", cType, strings.NewReader(string(k)))
	if err != nil {
		t.Fatal(err)
	}
	ok := StorageMap[2].Frend
	ko := StorageMap[3].Frend

	if len(ok) == 0 {
		t.Fatalf("Don`t Friends 2")
	}
	if len(ko) == 0 {
		t.Fatalf("Don`t Friends 3")
	}
	if post.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", post.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	OpenFileJson()
	s := httptest.NewServer(http.HandlerFunc(DeleteUser))
	defer s.Close()
	var sz Delete = Delete{Id: 4}
	r := strings.NewReader(`{"id":4}`)
	resp, err := http.NewRequest("DELETE", s.URL+"/user", r)
	if err != nil {
		t.Fatal(err)
	}
	textBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer resp.Body.Close()
	var text Delete
	errr := json.Unmarshal(textBytes, &text)
	if errr != nil {
		t.Log(err)
		t.Fail()
	}
	if sz.Id != text.Id {
		t.Log(err)
		t.Fail()
	}

}
func TestUpdateAge(t *testing.T) {
	OpenFileJson()
	s := httptest.NewServer(http.HandlerFunc(UpdateAge))
	defer s.Close()
	var sz NewAge = NewAge{NovelAge: 42}
	r := strings.NewReader(`{"new":42}`)
	req, err := http.NewRequest("PUT", s.URL+"/2", r)
	if err != nil {
		t.Fatal(err)
	}
	textBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(string(textBytes))
	defer req.Body.Close()
	var text NewAge
	errr := json.Unmarshal(textBytes, &text)
	if errr != nil {
		t.Log(err)
		t.Fail()
	}
	if sz.NovelAge != text.NovelAge {
		t.Log(err)
		t.Fail()
	}
	CloseFileJson()
	OpenFileJson()
	if StorageMap[2].Age != text.NovelAge {
		t.Log(err)
		t.Fail()
	}
}
