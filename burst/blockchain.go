package burst

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

type record struct {
	Type string
	Data string
	TTL  uint32
}

type response struct {
	Account   string
	AccountRS string
	AliasName string
	AliasURI  string
	Timestamp uint64
	Alias     string
}

var httpClient = &http.Client{Timeout: 2 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetRecords(aliasName string) ([]record, error) {
	resp := response{}

	nodeURL := viper.GetString("node.url")
	err := getJson(nodeURL+"/burst?requestType=getAlias&aliasName="+aliasName, &resp)
	if err != nil {
		log.Printf("Node error: %s", err)
		return nil, err
	}

	records := make([]record, 0)
	err = json.Unmarshal([]byte(resp.AliasURI), &records)
	if err != nil {
		log.Printf("Malformed records json: %s. AliasURI: >%s<", err, resp.AliasURI)
		return nil, err
	}

	return records, nil
}
