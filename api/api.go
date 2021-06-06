package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type service struct {
	db *gorm.DB
}

type Service interface {
	RandomRow() Result
}

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

func (s *service) RandomRow() Result {
	var result Result
	s.db.Raw("SELECT * FROM services ORDER BY RAND() LIMIT 1;").Scan(&result)
	return addCountUser(s.db, result)
}

func main() {
	var s Service

	dsn := "root:root@tcp(127.0.0.1:3306)/the_best_scraper?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	checkErr(err)

	s = &service{db}

	fmt.Println(s.RandomRow())
}

func addCountUser(db *gorm.DB, result Result) Result {
	r := jsonToMap(result)
	db.Model(&Services{}).Where("domain = ?", r.Domain).Update("request_count", r.Request_count+1)
	return r
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
