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
const MicrosoftADPolicyResourceType = "microsoftADPolicies"
const ModuleEndpointResourceType = "endpoints"

type OneFuseAPIClient struct {
	config *Config
}

type CustomName struct {
	Id        int
	Version   int
	Name      string
	DnsSuffix string
}

type LinkRef struct {
	Href  string `json:"href,omitempty"`
	Title string `json:"href,omitempty"`
}

type Workspace struct {
	Links *struct {
		Self LinkRef `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type WorkspacesListResponse struct {
	Embedded struct {
		Workspaces []Workspace `json:"workspaces"`
	} `json:"_embedded"`
}

type MicrosoftEndpoint struct {
	Links *struct {
		Self       LinkRef `json:"self,omitempty"`
		Workspace  LinkRef `json:"workspace,omitempty"`
		Credential LinkRef `json:"credential,omitempty"`
	} `json:"_links,omitempty"`
	ID               int    `json:"id,omitempty"`
	Type             string `json:"type,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"string,omitempty"`
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
	SSL              bool   `json:"ssl,omitempty"`
	MicrosoftVersion int    `json:"microsoftVersion,omitempty"`
}

type MicrosoftADPolicy struct {
	Links *struct {
		Self              LinkRef `json:"self,omitempty"`
		Workspace         LinkRef `json:"workspace,omitempty"`
		MicrosoftEndpoint LinkRef `json:"microsoftEndpoint,omitempty"`
	} `json:"_links,omitempty"`
	Name                   string `json:"name,omitempty"`
	ID                     int    `json:"id,omitempty"`
	Description            string `json:"description,omitempty"`
	MicrosoftEndpointID    int    `json:"microsoftEndpointId,omitempty"`
	MicrosoftEndpoint      string `json:"microsoftEndpoint,omitempty"`
	ComputerNameLetterCase string `json:"computerNameLetterCase,omitempty"`
	WorkspaceURL           string `json:"workspace,omitempty"`
	OU                     string `json:"ou,omitempty"`
}

func (c *Config) NewOneFuseApiClient() *OneFuseAPIClient {
	return &OneFuseAPIClient{
		config: c,
	}
}

func (apiClient *OneFuseAPIClient) GenerateCustomName(
	dnsSuffix string,
	namingPolicyID string,
	workspaceID string,
	templateProperties map[string]interface{},
) (
	result *CustomName,
	err error,
) {

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

func (apiClient *OneFuseAPIClient) CreateMicrosoftADPolicy(newPolicy *MicrosoftADPolicy) (MicrosoftADPolicy, error) {
	policy := MicrosoftADPolicy{}
	config := apiClient.config

	if newPolicy.Name == "" || newPolicy.WorkspaceURL == "" || newPolicy.MicrosoftEndpointID == 0 {
		return policy, errors.New("Microsoft AD Policy Updates Require a Name and, MicrosoftEndpointID, WorkspaceURL")
	}

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/microsoftADPolicies/
	url := collectionURL(config, MicrosoftADPolicyResourceType)

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

func (apiClient *OneFuseAPIClient) GetMicrosoftADPolicy(id int) (MicrosoftADPolicy, error) {
	policy := MicrosoftADPolicy{}
	config := apiClient.config
	url := itemURL(config, MicrosoftADPolicyResourceType, id)
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

func (apiClient *OneFuseAPIClient) UpdateMicrosoftADPolicy(id int, updatedPolicy *MicrosoftADPolicy) (MicrosoftADPolicy, error) {
	policy := MicrosoftADPolicy{}
	config := apiClient.config
	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	if updatedPolicy.Name == "" || updatedPolicy.WorkspaceURL == "" {
		return policy, errors.New("Microsoft AD Policy Updates Require a Name and WorkspaceURL")
	}

	jsonBytes, err := json.Marshal(updatedPolicy)
	if err != nil {
		return policy, err
	}

	requestBody := string(jsonBytes)
	if err != nil {
		err = errors.New("unable to marshal request body to JSON")
		return policy, err
	}

	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return policy, err
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return policy, err
	}

	checkForErrors(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return policy, err
	}
	res.Body.Close()

	err = json.Unmarshal(body, &policy)
	if err != nil {
		return policy, err
	}

	return policy, err
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftADPolicy(id int) error {
	config := apiClient.config

	// Construct a URL we are going to DELETE to
	// /api/v3/onefuse/microsoftADPolicy/<id>/
	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	// Make the DELETE request
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
	workspaceID = strconv.Itoa(workspaces[0].ID)
	return
}

func getHttpClient(config *Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.verifySSL},
	}
	return &http.Client{Transport: tr}
}

func checkForErrors(res *http.Response) error {
	if res.StatusCode >= 500 {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	} else if res.StatusCode >= 400 {
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
