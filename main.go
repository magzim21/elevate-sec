package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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
	Employee_id int64   `json:"employee_id"`
	Timestamp   float64 `json:"timestamp"`
}

type Priorities struct {
	// todo / declare default values
	Low      Severity `json:"low"`
	Medium   Severity `json:"medium"`
	High     Severity `json:"high"`
	Critical Severity `json:"critical"`
}

func (p *Priorities) sortSeverities() {
	p.Low.sortIncidents()
	p.Medium.sortIncidents()
	p.High.sortIncidents()
	p.Critical.sortIncidents()
}

type Severity struct {
	Count     int32      `json:"count"`
	Incidents []Incident `json:"incidents"`
}

func (s *Severity) addIncident(incident Incident) {
	s.Count++
	s.Incidents = append(s.Incidents, incident)
}
func (s *Severity) sortIncidents() {
	sort.SliceStable(s.Incidents, func(i, j int) bool {
		return s.Incidents[i].Timestamp < s.Incidents[j].Timestamp
	})
}

func main() {
	initConfig()

	// setting up gin (web framework)
	router := gin.Default()
	router.GET("/incidents", getIncidents)
	// router.SetTrustedProxies([]string{})
	router.Run(fmt.Sprintf(":%v", configuration.Port))

}

// getAlbums responds with the list of all albums as JSON.
func getIncidents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mergeIncidents())
}

func mergeIncidents() map[int64]Priorities {
	url1 := fmt.Sprint(configuration.ConnectionString, "/incidents/misuse/")
	url2 := fmt.Sprint(configuration.ConnectionString, "/incidents/unauthorized/")

	// declaring a variable for storing incidents of types "misuse" and "unauthorized"
	prepare := map[int64]*Priorities{}
	var response1 Response
	var response2 Response
	json.Unmarshal(querryIncidents(url1), &response1)
	json.Unmarshal(querryIncidents(url2), &response2)

	allIncidents := append(response1.Results, response2.Results...)

	for _, incident := range allIncidents {

		// initialize employee_id if not present
		if _, ok := prepare[incident.Employee_id]; !ok {
			prepare[incident.Employee_id] = &Priorities{}
			// prepare[incident.Employee_id].setDefaultValues()
		}
		switch incident.Priority {
		case "low":
			prepare[incident.Employee_id].Low.addIncident(incident)
		case "medium":
			prepare[incident.Employee_id].Medium.addIncident(incident)
		case "high":
			prepare[incident.Employee_id].High.addIncident(incident)
		case "critical":
			prepare[incident.Employee_id].Critical.addIncident(incident)
		}
		prepare[incident.Employee_id].sortSeverities()

	}

	// building the same structure but without pointers
	merged := map[int64]Priorities{}
	for employee_id, priorities := range prepare {
		merged[employee_id] = *priorities
	}
	return merged
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
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Printf("Can not decode config file\n")
		os.Exit(1)
	}
	secretsPath := flag.String("secrets-path", "./config/secrets", "a path where auth_username and auth_password secrets are stored as files.")
	AuthUsernameFile, err := os.ReadFile(fmt.Sprint(*secretsPath, "/auth_username"))
	if err != nil {
		fmt.Printf("Can not open file with auth_username secret \n")
		os.Exit(1)
	}
	AuthPasswordFile, err := os.ReadFile(fmt.Sprint(*secretsPath, "/auth_password"))
	if err != nil {
		fmt.Printf("Can not open file with auth_password secret \n")
		os.Exit(1)
	}

	configuration.AuthUsername = string(AuthUsernameFile)
	configuration.AuthPassword = string(AuthPasswordFile)

}
