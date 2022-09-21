package goglpi

import (
	"encoding/json"
	"strings"
	"testing"

	"fmt"
	"io/ioutil"
	"log"
)

/*
 *  For testing create config-test.json with this format ( 3:Localisation ):
 *
 {
	"url":"https://$URL_GLPI/apirest.php",
	"apptoken":"CREATE APP TOKEN INTO GLPI",
	"usertoken":"CREATE USER TOKEN INTO GLPI",

	"searchComputerField": "3",
	"searchComputerType": "contains",
	"searchComputerValue": "scaleway"
}
	and start go test
 *
*/
func TestConnectionGLPI(t *testing.T) {

	// get information for testing on your glpi
	content, err := ioutil.ReadFile("./config-test.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Load into map
	conf := make(map[string]string)
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	glpi, err := Init(conf["url"], conf["apptoken"], conf["usertoken"])
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Test 1: get computer ->")
	computer, err := glpi.GetComputer(21876)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		if computer.Id == 0 {
			fmt.Printf("Failed: ID is null\n")
		} else {
			fmt.Printf("OK\n")
		}
	}

	fmt.Printf("Test 2: get lists computers ->")
	computers, err := glpi.GetAllComputers()
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		if len(*computers) == 0 {
			fmt.Printf("There is no computer recovering.\n")
		} else {
			fmt.Printf("OK\n")
		}
	}

	fmt.Printf("Test 3: get search computers ->")
	listcriterias := make([]CriteriaSearch, 1, 1)
	listcriterias[0] = CriteriaSearch{Link: "AND", Field: conf["searchComputerField"], Searchtype: conf["searchComputerType"], Value: conf["searchComputerValue"]}
	computersSearch, _, err := glpi.SearchComputers(listcriterias, Range{Min: 0, Max: 10})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		if len(computersSearch) == 0 {
			fmt.Printf("There is no computer recovering.\n")
		} else {
			fmt.Printf("OK\n")
		}
	}

	fmt.Printf("Test 4: get search computers customer ->")
	type ComputerAddField struct {
		Computer
		Costumer string `json:"16665"` // If you have plugins adapt ID from your field
	}
	forcedisplay := make([]interface{}, 3)
	forcedisplay[0] = 1
	forcedisplay[1] = 2
	forcedisplay[2] = 16665 // If you have plugins adapt ID from your field
	interfaceListComputerInterface, nbreComputer, err := glpi.SearchWithName("Computer", ComputerAddField{}, listcriterias, forcedisplay, Range{Min: 0, Max: 10})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		if len(*computers) == 0 {
			fmt.Printf("There is no computer recovering.\n")
		} else {
			fmt.Println("OK\nThere are :", nbreComputer, "computer(s)")
			for i, computer := range *interfaceListComputerInterface.(*[]ComputerAddField) {
				fmt.Println(i, ": ", computer.Name, "16665 =>", strings.Replace(computer.Costumer, "\r\n", "\\r\\n", -1))
			}
		}
	}

	fmt.Printf("Test 5: get Option search ->")
	mapOptions,err:=glpi.ListSearchOptions(Computer{})
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		if len(mapOptions) == 0 {
			fmt.Printf("There is no computer recovering.\n")
		} else {
			fmt.Println("OK\nList Options fields:")
			for id, fieldFeature := range mapOptions {
				fmt.Println(id, ": ", fieldFeature.Name)
			}
		}
	}	
}
