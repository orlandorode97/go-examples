package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OrderResponse struct {
	Orders []struct {
		ID string `json:"id"`
	}
}

type OrderRequest struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID    string `json:"id"`
	Links Link   `json:"links"`
}

type Link struct {
	Sale string `json:"sale"`
}

type IndexerRequest struct {
	DocType string   `json:"docType"`
	DocIDs  []string `json:"docIDs"`
}

func main() {

	for i := 0; i < 1000; i++ {
		url := "https://api.lions-qe.omg.pub/tdo/orders/"

		orderRequest := &OrderRequest{
			Orders: []Order{{
				Links: Link{
					Sale: "4WWTW",
				},
			}},
		}

		payload, _ := json.Marshal(orderRequest)
		req, _ := http.NewRequest("POST", url, strings.NewReader(string(payload)))

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		resp := &OrderResponse{}
		err := json.Unmarshal(body, resp)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)

		url = fmt.Sprintf("https://api.lions-qe.omg.pub/tdo/orders/%v/actions/confirm", resp.Orders[0].ID)

		payloadReader := strings.NewReader("{\n\t\"gift_cards\": []\n}")

		req, _ = http.NewRequest("POST", url, payloadReader)

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		res, _ = http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ = io.ReadAll(res.Body)
		fmt.Println(string(body))

		url = "https://api.lions-qe.omg.pub/indexer/index"

		indexerRequest := &IndexerRequest{
			DocType: "tdo_order",
			DocIDs:  []string{resp.Orders[0].ID},
		}
		payload, _ = json.Marshal(indexerRequest)

		req, _ = http.NewRequest("POST", url, strings.NewReader(string(payload)))

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "insomnia/10.2.0")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
		res, _ = http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ = io.ReadAll(res.Body)
		fmt.Println(string(body))

		time.Sleep(2 * time.Millisecond)

	}

}
