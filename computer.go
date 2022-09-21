package goglpi

import (
	"net/http"
	//"time"
	"fmt"
)

// GLPI struct Computer
type Computers []Computer
type Computer struct {
	Id                   int `json:"2"`
	Entities_id          int
	Entities             string `json:"80"`
	Name                 string `json:"1"`
	Serial               string
	Otherserial          string
	Contact              string
	Contact_num          string //
	Users_id_tech        int
	Groups_id_tech       int
	Comment              string
	Date_mod             string // Beware full date had not a good format
	Autoupdatesystems_id int
	Locations_id         int
	Networks_id          int
	Computermodels_id    int
	Computertypes_id     int
	Is_template          int
	Template_name        string //
	Manufacturers_id     int
	Is_deleted           int
	Is_dynamic           int
	Users_id             int
	Groups_id            int
	States_id            int
	Ticket_tco           string
	Uuid                 string
	Date_creation        string // Beware full date had not a good format
	Links                Links
}

// GetComputer: Get info of computer
func (glpi GLPI) GetComputer(id int) (*Computer, error) {

	request, err := http.NewRequest("GET", glpi.url+fmt.Sprintf("/Computer/%v?expand_dropdowns=%v", id, false), nil)
	if err != nil {
		return nil, err
	}

	var computer Computer
	err = glpi.executeResquest(request, &computer)

	return &computer, err
}

// GetAllComputers: Get all computers Beacare full Entities membre is not complete ( use search )
func (glpi GLPI) GetAllComputers() (*Computers, error) {

	request, err := http.NewRequest("GET", glpi.url+fmt.Sprintf("/Computer?expand_dropdowns=%v", false), nil)
	if err != nil {
		return nil, err
	}

	var computers Computers
	err = glpi.executeResquest(request, &computers)

	return &computers, err
}

// Search computer by criterias, return Computer with IdComputer,Name and Entity
func (glpi GLPI) SearchComputers(criterias []CriteriaSearch, rangeLimit Range) ([]Computer, int, error) {

	interfaceListComputerInterface, nbreItem, err := glpi.Search(Computer{}, criterias, []interface{}{1, 2, 80}, rangeLimit)
	if err != nil {
		return nil, 0, err
	}

	return *interfaceListComputerInterface.(*[]Computer), nbreItem, nil
}
