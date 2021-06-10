package main

import (
	"../service"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type DefaultConfig struct {
	db       string
	user     string
	password string
}

var Config = DefaultConfig{
	`the_best_scraper`,
	`root`,
	`root`,
}

func main() {
	var checkEnter string
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", Config.user, Config.password, Config.db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	checkErr(err)

	s := service.New(db)

	s.InitServices()

	showrightEnter()
	for {
		fmt.Scan(&checkEnter)
		switch {
		case checkEnter == "-" || checkEnter == "exit":
			return
		case checkEnter == "rand":
			s.RandomRow()
		case checkEnter == "max":
			s.MaxTimeRow()
		case checkEnter == "min":
			s.MinTimeRow()
		case checkEnter == "show":
			s.ShowTable()
		default:
			showrightEnter()
		}
	}

	go s.UpdateServicesDB()
}

func showrightEnter() {
	fmt.Println("Pleace enter\n- exit\nrand\nmax\nmin\nshow")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
