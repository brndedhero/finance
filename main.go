package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/brndedhero/finance/config"
	"github.com/brndedhero/finance/models"
	"github.com/brndedhero/finance/router"

	"github.com/sirupsen/logrus"
)

func main() {
	config.DB = config.ConnectDb()
	config.DB.AutoMigrate(&models.Account{})
	config.Redis = config.ConnectRedis()
	config.Log = config.SetupLogger()
	config.Opensearch = config.ConnectOpensearch()

	http.Handle("/", router.SetupRouter())
	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))
	message := fmt.Sprintf("Listening for requests at http://%s:%d", os.Getenv("HTTP_HOST"), httpPort)
	config.Log.WithFields(logrus.Fields{
		"app":  "finance",
		"func": "main",
	}).Info(message)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}
