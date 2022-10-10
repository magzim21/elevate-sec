package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// albums slice to seed record album data.

var test = map[string]Priorities{"1": {Low: Severity{Count: 1, Incidents: []Incident{Incident{Priority: "low", Employee_id: 1, Timestamp: 1234567890}}}, Medium: Severity{Count: 1, Incidents: []Incident{Incident{Priority: "medium", Employee_id: 1, Timestamp: 1234567890}}}, High: Severity{Count: 1, Incidents: []Incident{Incident{Priority: "high", Employee_id: 1, Timestamp: 1234567890}}}, Critical: Severity{Count: 1, Incidents: []Incident{Incident{Priority: "critical", Employee_id: 1, Timestamp: 1234567890}}}}}

// var test = []employee{
// 	{Employee: ["1": [Count: 1, Incidents: [Priority: "low", Employee_id: 1, Timestamp: 1234567890]]]},
// 	{Employee: ["2": [Count: 1, Incidents: [Priority: "low", Employee_id: 2, Timestamp: 1234567890]]]},
// 	{Employee: ["3": [Count: 1, Incidents: [Priority: "low", Employee_id: 3, Timestamp: 1234567890]]]}
// }

func main() {
	// router := gin.Default()
	// router.GET("/incidents", getIncidents)
	// router.GET("/albums/:id", getAlbumByID)
	// router.POST("/albums", postAlbums)

	// router.Run("localhost:8080")
	// fmt.Println(parseIncidents())
	mergeIncidents()

}

// getAlbums responds with the list of all albums as JSON.
func getIncidents(c *gin.Context) {

	mergeIncidents()
	c.IndentedJSON(http.StatusOK, mergeIncidents())
	// c.IndentedJSON(http.StatusOK, test)
}

type Response struct {
	Results []Incident
}
type Incident struct {
	Priority    string  `json:"priority"`
	Employee_id int64   `json:"employee_id"` // or reference to employee struct
	Timestamp   float64 `json:"timestamp"`   // or float32 ?
}

// response represents data about a record response. ? Struct tags
type Priorities struct {
	// todo / declare default values
	Low      Severity `json:"low"`
	Medium   Severity `json:"medium"`
	High     Severity `json:"high"`
	Critical Severity `json:"critical"`
}

// func (p Priorities) fillDefaults()  {

// //   return &Priorities{Low: Severity{}, Medium: Severity{}, High: Severity{}, Critical: Severity{}}
// }

type Severity struct {
	Count     int32      `json:"count"`
	Incidents []Incident `json:"incidents"`
}
type Merged map[string]Priorities

func addIncidnetToEmployee(incident Incident, employee_id string) {
	

}

func mergeIncidents() map[string]Priorities {
	// parse incidents from url1
	// parse incidents from url2
	// parse incidents from url3
	// add to test
	// provide these in config file
	// url1:="https://incident-api.use1stag.elevatesecurity.io/identities/"
	url2 := "https://incident-api.use1stag.elevatesecurity.io/incidents/misuse/"
	// url3 := "https://incident-api.use1stag.elevatesecurity.io/incidents/unauthorized/"

	// declaring a variable for storing incidents of types "misuse" and "unauthorized"
	merged := Merged{}
	// merged := make(map[string]Priorities)
	// var merged []string
	var response1 Response
	// var response2 Response
	json.Unmarshal(querryIncidents(url2), &response1)
	// json.Unmarshal(querryIncidents(url3), &response2 )
	
	// iterating over incidents in the original response
	for _, incident := range response1.Results {
		employee_id_string := strconv.FormatInt(incident.Employee_id, 10)

		// initialize employee_id if not present
		if _, ok := merged[employee_id_string]; !ok {
			merged[employee_id_string] = Priorities{}
		}
		switch incident.Priority {
		case "low":
			merged[employee_id_string].Low.Count++
		case "medium":
			// merged[employee_id_string].Medium.Count++
		case "high":
			// merged[employee_id_string].High.Count++
		case "critical":
			// merged[employee_id_string].Critical.Count++
			
		}

		fmt.Printf("%+v\n", merged[employee_id_string])

		
		// merged = append(merged, employee_id_string)

	}
	// fmt.Printf("Code = %v, Employees = %v", response.Code, response.Employees)
	// fmt.Println(reflect.TypeOf(response1.Results[0].Employee_id))
	// fmt.Println(len(response1.Results))
	// fmt.Printf("%+v\n", response1)
	// fmt.Printf("%+v\n", merged)
	// fmt.Println(reflect.TypeOf(response2))

	//

	return map[string]Priorities{}

}

func querryIncidents(url string) []byte {
	client := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}
	req, err := http.NewRequest("GET", url, nil)
	// todo  / there is a shorter way to catch errors
	if err != nil {
		log.Fatal(err)
	}
	// todo / configure Auth in more elegant way
	req.Header.Set("Authorization", "Basic ")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body

}
