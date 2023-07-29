package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mailerservice/internal/config"
	"mailerservice/internal/utils"
	"net/http"
)

type Log struct {
	Collection     string      `json:"collection"`
	Source         string      `json:"source"`
	Method         string      `json:"method"`
	Request        interface{} `json:"request"`
	RequestHeader  interface{} `json:"request_header"`
	Response       interface{} `json:"response"`
	ResponseHeader interface{} `json:"response_header"`
	Duration       string      `json:"duration"`
	Status         int         `json:"status"`
}

// SendLog sends log to mongo
func SendLog(successLog Log) {

	url := config.EnvConfigs.LoggerMongoUrl
	successLog.Collection = "mailerservice"

	data, err := json.Marshal(successLog)
	utils.LogErr("error marshalling log", err)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	utils.LogErr("error creating request", err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	utils.LogErr("error sending request", err)

	// Close request body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.LogErr("Error while closing request body: ", err)
			return
		}
	}(resp.Body)

	log.Default().Println("Log sent successfully!")

}