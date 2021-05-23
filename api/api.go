package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
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

func AddRow(domain string, requestTime int) {
	database := callDB()

	services := Services{
		Domain:        domain,
		Request_time:  requestTime,
		Request_count: 0,
		Created_at:    time.Now(),
		Updated_at:    time.Now(),
	}

	database.Create(&services)
}

func CreateTableServices() {
	database := callDB()
	database.Migrator().CreateTable(&Services{})
}

func ChoiseRow(domain string) Result {
	database := callDB()
	var result Result
	database.Raw("SELECT Domain, Request_time, Request_count FROM services WHERE Domain = ?", domain).Scan(&result)
	return addCountUser(database, result)
}

func RandomRow() Result {
	database := callDB()
	var result Result
	database.Raw("SELECT * FROM services ORDER BY RAND() LIMIT 1;").Scan(&result)
	return addCountUser(database, result)
}

func MinTimeRow() Result {
	database := callDB()
	var result Result
	database.Raw("SELECT * FROM services WHERE request_time = (SELECT MIN(NULLIF(request_time, 0)) FROM services)").Scan(&result)
	return addCountUser(database, result)
}

func MaxTimeRow() Result {
	database := callDB()
	var result Result
	database.Raw("SELECT * FROM services WHERE request_time = (SELECT MAX(request_time)  FROM services)").Scan(&result)
	return addCountUser(database, result)
}

func addCountUser(db *gorm.DB, result Result) Result {
	r := jsonToMap(result)
	db.Model(&Services{}).Where("domain = ?", r.Domain).Update("request_count", r.Request_count+1)
	return r
}

func ShowTable() {
	database := callDB()
	var servicesArr []Services
	database.Find(&servicesArr)

	for _, service := range servicesArr {
		fmt.Println(service)
	}
}

func GetDomains() []Services {
	database := callDB()
	var servicesArr []Services
	database.Find(&servicesArr)

	return jsonToMapAllTable(servicesArr)
}

func UpdateServicesDB() {
	for {
		UpdateServices()
		time.Sleep(60 * time.Second)
	}
}

func UpdateServices() {
	services := make(chan Services)
	urls := getStrings("../../api/sites.txt")
	begin := makeTimestamp()
	for _, url := range urls {
		go initService(url, services, begin)
	}

	for i := 0; i < len(urls); i++ {
		service := <-services
		UpdateRow(service.Domain, service.Request_time)
	}
}

func UpdateRow(domain string, requestTime int) {
	database := callDB()

	services := Services{
		Domain:        domain,
		Request_time:  requestTime,
		Updated_at:    time.Now(),
	}

	database.Model(&Services{}).Where("domain = ?", domain).Updates(services)
}

func InitServices() {
	database := callDB()
	if boolTable := database.Migrator().HasTable(&Services{}); boolTable {
		return
	}
	CreateTableServices()
	services := make(chan Services)
	urls := getStrings("../../api/sites.txt")
	begin := makeTimestamp()
	for _, url := range urls {
		go initService(url, services, begin)
	}

	for i := 0; i < len(urls); i++ {
		service := <-services
		AddRow(service.Domain, service.Request_time)
	}
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
	channel<- Services{Domain: url, Request_time: end}
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

func callDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/the_best_scraper?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	checkErr(err)
	return database
}

func makeTimestamp() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
