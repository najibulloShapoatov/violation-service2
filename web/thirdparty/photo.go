package thirdparty

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"service/pkg/log"
)

func GetPhotoLinks(bId string) []string {

	key := fmt.Sprintf("%x", sha256.Sum256([]byte(bId+"727d56ee6941345481741ae32ca80806e11f3782e827190058a303c765cdcac8")))

	uri := "https://img.bekhatar.tj/getImages/" + bId + "?key=" + key

	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Error("sending SMS error", err)
	}
	q := url.Values{}
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		log.Warn("client.Do(req) ->\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil
	}

	type bodyT struct{ links []string }

	var body map[string][]string
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		log.Warn("json.NewDecoder(response.Body).Decode(body) => ", err)
		return nil
	}

	res, ok := body["links"]
	if !ok {
		return []string{fmt.Sprintf("%+v", body)}
	}
	return res
}
