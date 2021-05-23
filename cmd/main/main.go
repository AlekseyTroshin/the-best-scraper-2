package main

import (
	"../../api"
	"../../web"
)

func main() {
	api.InitServices()
	go api.UpdateServicesDB()

	web.Web()
}
