package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	// "errors"
	"github.com/gin-gonic/gin"
)

//  todo / add unauthorized , misuse, disclosure, modification, destruction
type Configuration struct {
	Port             int    `json:"port"`
	ConnectionString string `json:"connection_string"`
	AuthUsername     string `json:"auth_username"`
	AuthPassword     string `json:"auth_password"`
	Timeout          int    `json:"timeout"`
}

var configuration Configuration

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

type Severity struct {
	Count     int32      `json:"count"`
	Incidents []Incident `json:"incidents"`
}

type Merged map[string]*Priorities

func main() {
	initConfig()

	// setting up gin (web framework)
	router := gin.Default()
	router.GET("/incidents", getIncidents)
	// router.SetTrustedProxies([]string{})
	router.Run(fmt.Sprintf("localhost:%v", configuration.Port))

}

// getAlbums responds with the list of all albums as JSON.
func getIncidents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mergeIncidents())
}

func mergeIncidents() map[string]Priorities {
	url1 := fmt.Sprint(configuration.ConnectionString, "/incidents/misuse/")
	url2 := fmt.Sprint(configuration.ConnectionString, "/incidents/unauthorized/")

	// declaring a variable for storing incidents of types "misuse" and "unauthorized"
	merged := Merged{}
	var response1 Response
	var response2 Response
	json.Unmarshal(querryIncidents(url1), &response1)
	json.Unmarshal(querryIncidents(url2), &response2)

	// iterating over incidents in original responses
	for _, incident := range append(response1.Results, response2.Results...) {
		employee_id_string := strconv.FormatInt(incident.Employee_id, 10)

		// initialize employee_id if not present
		if _, ok := merged[employee_id_string]; !ok {
			merged[employee_id_string] = &Priorities{}
		}
		switch incident.Priority {
		case "low":
			merged[employee_id_string].Low.Count++
			merged[employee_id_string].Low.Incidents = append(merged[employee_id_string].Low.Incidents, incident)
		case "medium":
			merged[employee_id_string].Medium.Count++
			merged[employee_id_string].Medium.Incidents = append(merged[employee_id_string].Medium.Incidents, incident)
		case "high":
			merged[employee_id_string].High.Count++
			merged[employee_id_string].High.Incidents = append(merged[employee_id_string].High.Incidents, incident)
		case "critical":
			merged[employee_id_string].Critical.Count++
			merged[employee_id_string].Critical.Incidents = append(merged[employee_id_string].Critical.Incidents, incident)

		}

	}
	// todo / rethink data model
	normal_merged := map[string]Priorities{}
	for employee_id, priorities := range merged {
		normal_merged[employee_id] = *priorities
	}
	return normal_merged
}

func querryIncidents(url string) []byte {
	client := http.Client{Timeout: time.Second * time.Duration(configuration.Timeout)}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+basicAuth(configuration.AuthUsername, configuration.AuthPassword))
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Printf("Can not read remote data. Response code: %v . Error: %v  ", res.StatusCode, err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Can not read response body. Error: %v  ", err)
		os.Exit(1)
	}
	return body

}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func initConfig() {
	// getting configuration values from the config file
	configPath := flag.String("config", "./config/config-development.json", "a path to configuration file")
	file, err := os.Open(*configPath)
	if err != nil {
		fmt.Printf("Can not open config file\n")
		return
	}
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Printf("Can not decode config file\n")
		return
	}
	secretsPath := flag.String("secrets-path", "./config/secrets", "a path where auth_username and auth_password secrets are stored as files.")
	AuthUsernameFile, _ := os.ReadFile(fmt.Sprint(*secretsPath, "/auth_username"))
	AuthPasswordFile, _ := os.ReadFile(fmt.Sprint(*secretsPath, "/auth_password"))

	configuration.AuthUsername = string(AuthUsernameFile)
	configuration.AuthPassword = string(AuthPasswordFile)

}
