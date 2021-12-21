package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/opensearch-project/opensearch-go"
	"github.com/sirupsen/logrus"
)

var Opensearch *opensearch.Client

func ConnectOpensearch() *opensearch.Client {
	opensearchPort, _ := strconv.Atoi(os.Getenv("OPENSEARCH_PORT"))
	opensearchUri := fmt.Sprintf("https://%s:%d", os.Getenv("OPENSEARCH_HOST"), opensearchPort)

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{opensearchUri},
		Username:  os.Getenv("OPENSEARCH_USER"),
		Password:  os.Getenv("OPENSEARCH_PASSWORD"),
	})
	if err != nil {
		Log.WithFields(logrus.Fields{
			"app":  "opensearch",
			"func": "connectOpensearch",
		}).Fatal(err)
	}
	return client
}
