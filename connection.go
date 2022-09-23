package goglpi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

// GLPI struct
type GLPI struct {
	Session_token string

	url        string
	app_token  string
	user_token string
}

// Init: Init glpi connexion
func Init(url, apptoken, usertoken string) (*GLPI, error) {
	var glpi GLPI
	glpi.url = url
	glpi.app_token = apptoken
	glpi.user_token = usertoken

	request, err := http.NewRequest("GET", url+"/initSession", nil)
	if err != nil {
		return nil, err
	}

	err = glpi.executeResquest(request, &glpi)
	if err != nil {
		return nil, err
	}

	if glpi.Session_token == "" {
		return nil, fmt.Errorf("GLPI don't return session_token")
	}

	return &glpi, nil
}

// executeResquest: internal function for executing an HTTP request to GLPI
func (glpi *GLPI) executeResquest(request *http.Request, objectReturn interface{}) error {
	client := &http.Client{}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("App-Token", glpi.app_token)
	if glpi.Session_token == "" {
		request.Header.Add("Authorization", "user_token "+glpi.user_token)
	} else {
		request.Header.Add("Session-Token", glpi.Session_token)
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if string(body[0:8]) == "[\"ERROR_" {
		return fmt.Errorf("GLPI response error:%v", string(body))
	}

	fmt.Println("======================")
	fmt.Println(request)
	fmt.Println(string(body))
	fmt.Println("======================")

	err = json.Unmarshal(body, objectReturn)
	if err != nil {
		return fmt.Errorf("%v:%v", err, string(body[:100]))
	}

	return nil
}

// CriteriaSearch:
type CriteriaSearch struct {
	Link       string
	Field      interface{}
	Searchtype string
	Value      interface{}
}

// ResultSearch: result information from search
type resultSearch struct {
	Totalcount int
	Count      int
	Sort       string
	Order      string
	Data       interface{}
}

// Range of search
type Range struct {
	Min int
	Max int
}

// Search: Search information for item ( Computer, ...)
func (glpi GLPI) Search(item Item, criterias []CriteriaSearch, forcedisplays []interface{}, rangeLimit Range) (interface{}, int, error) {
	// Get type of Item
	itemType := reflect.TypeOf(item)

	return glpi.SearchWithName(itemType.Name(), item, criterias, forcedisplays, rangeLimit)
}

// Search: Search information for customer item ( Computer with Entity, ...)
func (glpi GLPI) SearchWithName(itemName string, item Item, criterias []CriteriaSearch, forcedisplays []interface{}, rangeLimit Range) (interface{}, int, error) {
	// Get type of Item
	itemType := reflect.TypeOf(item)

	//Create slice and get pointer
	sliceItem := reflect.MakeSlice(reflect.SliceOf(itemType), 0, 0)
	ptrSliceItem := reflect.New(sliceItem.Type())
	ptrSliceItem.Elem().Set(sliceItem)

	// Create parameter
	var values url.Values = url.Values{}
	for i, criteria := range criterias {
		values.Add(fmt.Sprintf("criteria[%v][field]", i), fmt.Sprintf("%v", criteria.Field))
		values.Add(fmt.Sprintf("criteria[%v][searchtype]", i), fmt.Sprintf("%v", criteria.Searchtype))
		values.Add(fmt.Sprintf("criteria[%v][value]", i), fmt.Sprintf("%v", criteria.Value))
	}
	for i, forcedisplay := range forcedisplays {
		values.Add(fmt.Sprintf("forcedisplay[%v]", i), fmt.Sprintf("%v", forcedisplay))
	}
	values.Add("range", rangeLimit.String())

	// Create url with parameter
	urlGLPIraw, _ := url.Parse(glpi.url)
	urlGLPIraw.Path += fmt.Sprintf("/search/%v/", itemName)
	urlGLPIraw.RawQuery = values.Encode()

	// Create Request
	request, err := http.NewRequest("GET", urlGLPIraw.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	// Send request into Agent
	var result resultSearch
	result.Data = ptrSliceItem.Interface()
	err = glpi.executeResquest(request, &result)

	return result.Data, result.Totalcount, err

}

type SearchOption struct {
	Name     string
	Table    string
	Field    string
	Datatype string
	Uid      string
}

// ListSearchOptions: Get all fields for requesting tiem
func (glpi GLPI) ListSearchOptions(item Item) (map[int]SearchOption, error) {
	// Get type of Item
	itemType := reflect.TypeOf(item)

	// Create url with parameter
	urlGLPIraw, _ := url.Parse(glpi.url)
	urlGLPIraw.Path += fmt.Sprintf("/listSearchOptions/%v", itemType.Name())

	// Create Request
	request, err := http.NewRequest("GET", urlGLPIraw.String(), nil)
	if err != nil {
		return nil, err
	}

	// Send request into Agent
	result := make(map[string]interface{})
	err = glpi.executeResquest(request, &result)

	resultReturn := make(map[int]SearchOption)
	for idString, interfaceMap := range result {

		id, _ := strconv.Atoi(idString)
		if id > 0 {
			var searchOption SearchOption

			cfg := &mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   &searchOption,
				TagName:  "json",
			}

			decoder, _ := mapstructure.NewDecoder(cfg)
			decoder.Decode(interfaceMap)

			resultReturn[id] = searchOption
		}
	}

	return resultReturn, err
}

type Item interface{}
type Items []Item

// Link: Struct link information GLPI
type Links []Link
type Link struct {
	Rel  string
	Href string
}

func (rangeLimit Range) String() string {
	return fmt.Sprintf("%v-%v", rangeLimit.Min, rangeLimit.Max)
}
