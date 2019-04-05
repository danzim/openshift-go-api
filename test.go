package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"

	"github.com/tidwall/gjson"
)

const (
	//token             = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InB2aW5zdGFsbGVyLXRva2VuLW5wenZtIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6InB2aW5zdGFsbGVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiYjRmNjA4ZjMtNTQ3NS0xMWU5LTk1MWItN2EzODJhMjJkYmU5Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmRlZmF1bHQ6cHZpbnN0YWxsZXIifQ.SCGxkd-KUuAuCC4qOqzhpD86LqmPqTra43I0Ym2uat3u6wQtawztc2q6d_Qb866oIzPfXMLQPUJzeMZONHq713vDIvJaBd23YKE0ONJijL0jIRSWYFs7KM-0VH7o3amZww0yexC-lJw6ESfZzIGZEW5pn5rTZGUGZxyOxBrYm6TLaImy5ZH3A2QA4jRvw66IfVK0JX-Y9247fXMQVhYFD10AmphqPrewYSdM-W37YSKhqH9tVfZrjBUHzDj3AYcMrp7_ZzIuActe6bMUX9_c7mEelMpPIJyhFt_fU3eqC64YL6NksE8dXBfdEipLx3TeCKgMxj7G3gYtEf9wMCzYFQ"
	token             = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJ0ZXN0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InByb2otcm9ib3QtdG9rZW4tamhoZHgiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoicHJvai1yb2JvdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImI3YTJmZjIzLTU2YmUtMTFlOS1iZTZjLTdhMzgyYTIyZGJlOSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDp0ZXN0OnByb2otcm9ib3QifQ.rN9l-PJWpaM36lwebaKKIa_CuNPJvfEVQkJf3DeyKlot6ImVIkALON-cjea2BoKAwojVcLmODlGq42kLiwA23fptXSOI1v_ZUj-Z2lj7IK52rtrBxJZXBkzRdLHixDoqompfTxXNGrpqW9MMv8jWS-ZFvBXtFnDGs0RJMWJi0RXnZdP7lnuPekD7DQv2gggDR6g4h3A7asrvkS_wjrazCjrwd2ZZlTHz1QpRnTbS8YbWomzKYP1RpyOuwakFwKVXzjbZCzADEbjFMFZGUecB0Up_DGB8I9DFFEFAiZdK4KQkEeZPGWjXouZQBibUzQvezvCXqpoEbvvm1Q2GOUJ3EA"
	applicationHeader = "application/json"
	baseURL           = "https://192.168.64.4:8443/oapi/v1/"
	projectTemplate   = `{"kind":"ProjectRequest","apiVersion":"v1","metadata":{"name":"{{.Name}}","annotations":{"openshift.io/description":"Das ist ein Test-Project","openshift.io/display-name":"{{.DisplayName}}"}}}`
	projectTemplate2  = `{"kind":"ProjectRequest","apiVersion":"v1","metadata":{"name":"bla","annotations":{"openshift.io/description":"Das ist ein Test-Project","openshift.io/display-name":"blubb"}}}`
)

func main() {
	fmt.Println("Starting the application...")
	http.HandleFunc("/postproject/create/", func(w http.ResponseWriter, r *http.Request) {
		paramName := r.URL.Query().Get("name")
		if paramName != "" {
			paramDisplayName := r.URL.Query().Get("displayname")
			if paramDisplayName != "" {
				fmt.Printf("Create project...\nname: %s\nDisplay Name: %s\n", paramName, paramDisplayName)
				urlPost := fmt.Sprintf("%sprojectrequests/", baseURL)
				fmt.Printf(urlPost)
				match, err := regexp.MatchString(`ci-[0-9]{8}`, paramName)
				if err != nil {
					log.Fatal("Regex is false...", err)
				}
				if match {
					oscpPost(urlPost, paramName, paramDisplayName)
				}

			}
		}
	})

	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		paramProject := r.URL.Query().Get("project")
		if paramProject != "" {
			fmt.Println(paramProject)
			urlGet := fmt.Sprintf("%sprojects/%s/", baseURL, paramProject)
			fmt.Printf(urlGet)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func project(url string) {
	fmt.Println("Show selected project...")
	urlProject := fmt.Sprintf("https://192.168.64.4:8443/apis/project.openshift.io/v1/projects/%s/", url)
	dataJSON := oscpGet(urlProject)
	name := gjson.GetBytes(dataJSON, "metadata.name")
	fmt.Printf("Name des Projekts: %s\n", name)
}

func oscpPost(url string, name string, displayName string) {
	type Info struct {
		Name        string
		DisplayName string
	}

	var buffer bytes.Buffer
	var data io.Reader
	info := Info{name, displayName}
	//t := template.New("Test")
	//text := `Hallo {{.Name}} du wirst angezeigt als {{.DisplayName}}`
	temp := template.New("ProjectRequest")
	temp.Parse(projectTemplate)
	temp.Execute(&buffer, info)

	fmt.Println(buffer.String())

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	bearer := fmt.Sprintf("bearer %s", token)
	//data := []io.Reader{buffer}
	data = &buffer
	//data := buffer.Bytes()

	client := &http.Client{Transport: transport}
	//req, err := http.NewRequest("POST", url, &buffer)
	req, err := http.NewRequest("POST", url, data)

	if err != nil {
		log.Fatal("Error reading request:\n", err)
	}

	//set headers
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", applicationHeader)
	req.Header.Add("Content-Type", applicationHeader)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("The HTTP POST request failed with error:\n", err)
	}
	responseData, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}

func oscpGet(url string) []byte {
	//skip certificate check
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	bearer := fmt.Sprintf("bearer %s", token)

	//set http client and request
	client := &http.Client{Transport: transport}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error reading request:\n", err)
	}

	//set headers
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", applicationHeader)
	req.Header.Add("Content-Type", applicationHeader)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("The HTTP request failed with error:\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	return data
}
