package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Services struct {
	ID            uint `gorm:"primaryKey"`
	Domain        string
	Request_time  int
	Request_count int
	Created_at    time.Time
	Updated_at    time.Time
}

type Result struct {
	Domain        string
	Request_time  int
	Request_count int
}

type service struct {
	db *gorm.DB
}

func New(db *gorm.DB) Service {
	return &service{db}
}

type Service interface {
	InitServices()
	RandomRow()
	MinTimeRow()
	MaxTimeRow()
	ShowTable()
	UpdateServicesDB()
}

func showResult(jsonService Result) {
	fmt.Printf("| %20s | %10d | %10d |\n", jsonService.Domain, jsonService.Request_time, jsonService.Request_count)
}

func (s *service) MinTimeRow() {
	var result Result
	s.db.Raw("SELECT * FROM services WHERE request_time <> 0 ORDER BY request_time LIMIT 1;").Scan(&result)

	jsonResult := addCountUser(s.db, result)
	showResult(jsonResult)
}

func (s *service) MaxTimeRow() {
	var result Result
	s.db.Raw("SELECT * FROM services ORDER BY request_time DESC LIMIT 1;").Scan(&result)

	jsonResult := addCountUser(s.db, result)
	showResult(jsonResult)
}

func (s *service) ShowTable() {
	var servicesArr []Services
	s.db.Find(&servicesArr)

	for _, service := range servicesArr {
		fmt.Printf("| %20s | %10d | %10d | %40s | %40s |\n", service.Domain, service.Request_time, service.Request_count, service.Created_at, service.Updated_at)
	}
}

func (s *service) RandomRow() {
	var result Result
	s.db.Raw("SELECT * FROM services ORDER BY RAND() LIMIT 1;").Scan(&result)
	jsonResult := addCountUser(s.db, result)
	showResult(jsonResult)
}

func (s *service) InitServices() {
	if boolTable := s.db.Migrator().HasTable(&Services{}); boolTable {
		return
	}
	createTableServices(s.db)
	services := make(chan Services)
	urls := getStrings("../api/sites.txt")
	begin := makeTimestamp()
	for _, url := range urls {
		go initService(url, services, begin)
	}

	for i := 0; i < len(urls); i++ {
		service := <-services
		addRow(s.db, service.Domain, service.Request_time)
	}
}

func (s *service) UpdateServicesDB() {
	for {
		updateServices(s.db)
		time.Sleep(60 * time.Second)
	}
}

func updateServices(db *gorm.DB) {
	services := make(chan Services)
	urls := getStrings("../api/sites.txt")
	begin := makeTimestamp()
	for _, url := range urls {
		go initService(url, services, begin)
	}

	for i := 0; i < len(urls); i++ {
		service := <-services
		updateRow(db, service.Domain, service.Request_time)
	}
}

func updateRow(db *gorm.DB, domain string, requestTime int) {
	services := Services{
		Domain:       domain,
		Request_time: requestTime,
		Updated_at:   time.Now(),
	}

	db.Model(&Services{}).Where("domain = ?", domain).Updates(services)
}

func initService(url string, channel chan Services, loadingBegin int) {
	var end int
	response, err := http.Get("https://www." + url)
	if err == nil {
		defer response.Body.Close()
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		end = makeTimestamp() - loadingBegin
	}
	channel <- Services{Domain: url, Request_time: end}
}

func jsonToMap(r Result) Result {
	jsonData, _ := json.Marshal(r)
	service := Result{}
	err := json.Unmarshal(jsonData, &service)
	if err != nil {
		log.Println(err)
	}

	return service
}

func jsonToMapAllTable(r []Services) []Services {
	jsonData, _ := json.Marshal(r)
	services := []Services{}
	err := json.Unmarshal(jsonData, &services)
	if err != nil {
		log.Println(err)
	}

	return services
}

func addCountUser(db *gorm.DB, result Result) Result {
	r := jsonToMap(result)
	db.Model(&Services{}).Where("domain = ?", r.Domain).Update("request_count", r.Request_count+1)
	return r
}

func addRow(db *gorm.DB, domain string, requestTime int) {
	services := Services{
		Domain:        domain,
		Request_time:  requestTime,
		Request_count: 0,
		Created_at:    time.Now(),
		Updated_at:    time.Now(),
	}

	db.Create(&services)
}

func createTableServices(db *gorm.DB) {
	db.Migrator().CreateTable(&Services{})
}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	checkErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	checkErr(scanner.Err())
	return lines
}

func makeTimestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
