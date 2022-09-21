# GoGLPI


GLPI is a library that allows you to connect to the GLPI API. 
The first version can only query the computers part with the possibility of accessing the fields of the plugins.

## Examples:

1) Get informations of computer with id

	computer, err := glpi.GetComputer( IdOfComputer)

2) Get all computer into GLPI
	computers, err := glpi.GetAllComputers()

	
3)  Find a computer with criteria, here I am looking for all computers that contains "srv"

	listcriterias := make([]CriteriaSearch, 1, 1)
	listcriterias[0] = CriteriaSearch{
		Link: "AND",
		Field: "Name",
		Searchtype: "contains"
		Name: "srv"
	}
	computersSearch, nbreComputer, err := glpi.SearchComputers(listcriterias, Range{Min: 0, Max: 10})
	
4) Find a computer with fields from a plugin, here the fields had id 12230 ( Use a part 5 for find "ID" field )
	
	type ComputerAddField struct {
		Id		int `json:"2"`
		Name	string `json:"1"`
		Company string `json:"12230"` 
	}
	forcedisplay := make([]interface{}, 3)
	forcedisplay[0] = 1
	forcedisplay[0] = 2
	forcedisplay[2] = 76665
	
	interfaceListComputerInterface, nbreComputer, err := glpi.SearchWithName("Computer", ComputerAddField{}, listcriterias, forcedisplay, Range{Min: 0, Max: 10})
	
5) Find a fields of items, here the fields had id 12230

	mapOptions,err:=glpi.ListSearchOptions(Computer{})
	for id, fieldFeature := range mapOptions {
		fmt.Println(id, ": ", fieldFeature.Name)
	}
	
	
Now it's your turn.
My library does not use all the items but it is relatively easy to complete ;)


## License
© VALMIR, 2022~time.Now

Released under the [MIT License]
