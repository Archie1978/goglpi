package goglpi

import (
	"net/http"
	//"time"
	"fmt"
)



// GLPI struct Certificate
type Certificates []Certificate
type Certificate struct {
	Id                   int `json:"ID"`
	Name		string
	Otherserial	string
	Entities_id	int
	Is_recursive	int
	Comment	string
	Is_deleted	int
	Is_demplate	int
	Template_name string
	Certificatetypes_Id	int
	DnsName	string
	Dns_suffix	string
	Users_id_tech int
	Groups_id_tech int
	Locations_id	int
	Manufacturers_id	int
	Contact	string
	Contact_num string
	Users_id	int
	Is_autosign	int
	Date_expiration string	// not format standard json
	States_id	int
	Command		string
	Certificate_request	string
	Certificate_item	string
	Date_creation	string
	Date_mod	string
}


// GetCertificate: Get info of Certificate
func (glpi GLPI) GetCertificate(id int) (*Certificate, error) {

	request, err := http.NewRequest("GET", glpi.url+fmt.Sprintf("/Certificate/%v", id), nil)
	if err != nil {
		return nil, err
	}

	var certificate Certificate
	err = glpi.executeResquest(request, &certificate)

	return &certificate, err
}


// GetCertificate: Get all Certificate
func (glpi GLPI) GetAllCertificates() (*Certificates, error) {

	request, err := http.NewRequest("GET", glpi.url+fmt.Sprintf("/Certificate?expand_dropdowns=%v", false), nil)
	if err != nil {
		return nil, err
	}

	var certificates Certificates
	err = glpi.executeResquest(request, &certificates)

	return &certificates, err
}
