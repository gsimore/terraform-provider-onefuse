// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const ApiVersion = "/api/v3/"
const ApiNamespace = "onefuse"
const NamingResourceType = "customNames"
const WorkspaceResourceType = "workspaces"
const MicrosoftAdPolicyResourceType = "microsoftActiveDirectoryPolicies"
const ModuleEndpointResourceType = "endpoints"

type OneFuseAPIClient struct {
	config *Config
}

type Workspace struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type CustomName struct {
	Id        int
	Version   int
	Name      string
	DnsSuffix string
}

type WorkspacesListResponse struct {
	Embedded struct {
		Workspaces []Workspace `json:"workspaces"`
	} `json:"_embedded"`
}

type MicrosoftEndpoint struct {
	Links struct {
		Workspace `json:"workspace"`
	} `json:"_links"`
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"string"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	SSL              bool   `json:"ssl"`
	MicrosoftVersion int    `json:"microsoftVersion"`
}

type MicrosoftAdPolicy struct {
	Links struct {
		// TODO: just model what we get; use the workspace URL.
		Workspace         `json:"workspace"`
		MicrosoftEndpoint struct {
			// TODO Update
			Name string
			URL  string
		}
	} `json:"_links"`
	Name                   string `json:"name"`
	ID                     int    `json:"id"`
	Description            string `json:"description"`
	MicrosoftEndpoint      string `json:"microsoftEndpoint"`
	ComputerNameLetterCase string `json:"computerNameLetterCase"`
	OU                     string `json:"ou"`
}

func (c *Config) NewOneFuseApiClient() *OneFuseAPIClient {
	return &OneFuseAPIClient{
		config: c,
	}
}

func (apiClient *OneFuseAPIClient) GenerateCustomName(dnsSuffix string, namingPolicyID string, workspaceID string,
	templateProperties map[string]interface{}) (result *CustomName, err error) {

	config := apiClient.config
	url := collectionURL(config, NamingResourceType)
	log.Println("reserving custom name from " + url + "  dnsSuffix=" + dnsSuffix)

	if templateProperties == nil {
		templateProperties = make(map[string]interface{})
	}
	if workspaceID == "" {
		workspaceID, err = findDefaultWorkspaceID(config)
		if err != nil {
			return
		}
	}

	postBody := map[string]interface{}{
		"namingPolicy":       fmt.Sprintf("%s%s/namingPolicies/%s/", ApiVersion, ApiNamespace, namingPolicyID),
		"templateProperties": templateProperties,
		"workspace":          fmt.Sprintf("%s%s/workspaces/%s/", ApiVersion, ApiNamespace, workspaceID),
	}
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(postBody)
	requestBody := string(jsonBytes)
	if err != nil {
		err = errors.New("unable to marshal request body to JSON")
		return
	}
	payload := strings.NewReader(requestBody)

	log.Println("CONFIG:")
	log.Println(config)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return
	}
	log.Println("HTTP PAYLOAD to " + url + ":")
	log.Println(postBody)

	setHeaders(req, config)

	client := getHttpClient(config)
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}

	checkForErrors(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println("HTTP POST RESULTS:")
	log.Println(string(body))
	json.Unmarshal(body, &result)
	res.Body.Close()

	if result == nil {
		err = errors.New("invalid response " + strconv.Itoa(res.StatusCode) + " while generating a custom name: " + string(body))
		return
	}

	log.Println("custom name reserved: " +
		"custom_name_id=" + strconv.Itoa(result.Id) +
		" name=" + result.Name +
		" dnsSuffix=" + result.DnsSuffix)
	fmt.Printf("Complete!")
	return
}

func (apiClient *OneFuseAPIClient) GetCustomName(id int) (result CustomName, err error) {
	config := apiClient.config
	url := itemURL(config, NamingResourceType, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	setHeaders(req, config)

	log.Println("REQUEST:")
	log.Println(req)
	client := getHttpClient(config)
	res, _ := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println("HTTP GET RESULTS:")
	log.Println(string(body))

	json.Unmarshal(body, &result)
	res.Body.Close()
	return
}

func (apiClient *OneFuseAPIClient) DeleteCustomName(id int) error {
	config := apiClient.config
	url := itemURL(config, NamingResourceType, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	setHeaders(req, config)
	client := getHttpClient(config)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	return checkForErrors(res)
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftEndpoint(newEndpoint MicrosoftEndpoint) (MicrosoftEndpoint, error) {
	endpoint := MicrosoftEndpoint{}
	err := errors.New("Not implemented yet")
	return endpoint, err
}

func (apiClient *OneFuseAPIClient) GetMicrosoftEndpoint(id int) (MicrosoftEndpoint, error) {
	endpoint := MicrosoftEndpoint{}
	err := errors.New("Not implemented yet")
	return endpoint, err
}

func (apiClient *OneFuseAPIClient) GetMicrosoftEndpointByName(name string) (MicrosoftEndpoint, error) {
	endpoint := MicrosoftEndpoint{}
	config := apiClient.config
	url := collectionURL(config, ModuleEndpointResourceType)
	url += fmt.Sprintf("?filter=name:%s;type:microsoft", name)

	fmt.Print(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return endpoint, err
	}

	setHeaders(req, config)

	log.Println("REQUEST:")
	log.Println(req)
	client := getHttpClient(config)
	res, _ := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return endpoint, err
	}
	log.Println("HTTP GET RESULTS:")
	log.Println(string(body))

	json.Unmarshal(body, &endpoint)
	log.Println("Endpoint:")
	log.Println(endpoint)
	res.Body.Close()
	return endpoint, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftEndpoint(id int, updatedEndpoint MicrosoftEndpoint) (MicrosoftEndpoint, error) {
	endpoint := MicrosoftEndpoint{}
	err := errors.New("Not implemented yet")
	return endpoint, err
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftEndpoint(id int) error {
	return errors.New("Not implemented yet")
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftAdPolicy(newPolicy MicrosoftAdPolicy) (MicrosoftAdPolicy, error) {
	policy := MicrosoftAdPolicy{}
	config := apiClient.config

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/microsoftADPolicies/
	url := collectionURL(config, MicrosoftAdPolicyResourceType)

	var jsonBytes []byte
	jsonBytes, err := json.Marshal(newPolicy)
	requestBody := string(jsonBytes)
	if err != nil {
		err = errors.New("unable to marshal request body to JSON")
		return policy, err
	}
	payload := strings.NewReader(requestBody)

	// Create the DELETE request
	req, _ := http.NewRequest("POST", url, payload)

	setHeaders(req, config)

	client := getHttpClient(config)

	// Make the delete request
	res, err := client.Do(req)

	// Return err if it went poorly
	if err != nil {
		return policy, err
	}

	err = checkForErrors(res)
	if err != nil {
		return policy, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return policy, err
	}

	err = res.Body.Close()
	if err != nil {
		return policy, err
	}

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return policy, err
	}
	log.Println(policy)

	return policy, nil
}

func (apiClient *OneFuseAPIClient) GetMicrosoftAdPolicy(id int) (MicrosoftAdPolicy, error) {
	policy := MicrosoftAdPolicy{}
	config := apiClient.config
	url := itemURL(config, MicrosoftAdPolicyResourceType, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return policy, err
	}

	setHeaders(req, config)

	log.Println("REQUEST:")
	log.Println(req)
	client := getHttpClient(config)
	res, _ := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return policy, err
	}
	log.Println("HTTP GET RESULTS:")
	log.Println(string(body))

	json.Unmarshal(body, &policy)
	res.Body.Close()
	return policy, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftAdPolicy(id int, updatedPolicy MicrosoftAdPolicy) (MicrosoftAdPolicy, error) {
	policy := MicrosoftAdPolicy{}
	err := errors.New("Not implemented yet")

	return policy, err
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftAdPolicy(id int) error {
	config := apiClient.config

	// Construct a URL we are going to DELETE to
	// /api/v3/onefuse/microsoftADPolicy/<id>/
	url := itemURL(config, MicrosoftAdPolicyResourceType, id)

	// Create the DELETE request
	req, _ := http.NewRequest("DELETE", url, nil)

	setHeaders(req, config)

	client := getHttpClient(config)

	// Make the delete request
	res, err := client.Do(req)

	// Return err if it went poorly
	if err != nil {
		return err
	}

	return checkForErrors(res)
}

func findDefaultWorkspaceID(config *Config) (workspaceID string, err error) {
	filter := "filter=name.exact:Default"
	url := fmt.Sprintf("%s?%s", collectionURL(config, WorkspaceResourceType), filter)
	req, clientErr := http.NewRequest("GET", url, nil)
	if clientErr != nil {
		err = clientErr
		return
	}

	setHeaders(req, config)

	client := getHttpClient(config)
	res, clientErr := client.Do(req)
	if clientErr != nil {
		err = clientErr
		return
	}

	checkForErrors(res)

	body, clientErr := ioutil.ReadAll(res.Body)
	if clientErr != nil {
		err = clientErr
		return
	}

	var data WorkspacesListResponse
	json.Unmarshal(body, &data)
	res.Body.Close()

	workspaces := data.Embedded.Workspaces
	if len(workspaces) == 0 {
		panic("Unable to find default workspace.")
	}
	workspaceID = workspaces[0].ID
	return
}

func getHttpClient(config *Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.verifySSL},
	}
	return &http.Client{Transport: tr}
}

func checkForErrors(res *http.Response) error {
	// TODO: How to handle 40X errors?
	if res.StatusCode >= 500 {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}

func setHeaders(req *http.Request, config *Config) {
	setStandardHeaders(req)
	req.Header.Add("Host", config.address+":"+config.port)
	req.Header.Add("SOURCE", "Terraform")
	req.SetBasicAuth(config.user, config.password)
}

func collectionURL(config *Config, resourceType string) string {
	address := config.address
	port := config.port
	return config.scheme + "://" + address + ":" + port + ApiVersion + ApiNamespace + "/" + resourceType + "/"
}

func itemURL(config *Config, resourceType string, id int) string {
	address := config.address
	port := config.port
	idString := strconv.Itoa(id)
	return config.scheme + "://" + address + ":" + port + ApiVersion + ApiNamespace + "/" + resourceType + "/" + idString + "/"
}
