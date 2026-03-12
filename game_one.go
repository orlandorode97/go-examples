package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DestinationAddress struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
	Street1    string `json:"street_1"`
}

type Shipping struct {
	DestinationAddress DestinationAddress `json:"destination_address"`
	Fee                int                `json:"fee"`
	Method             string             `json:"method"`
	Name               string             `json:"name"`
}

type PriceAddition struct {
	Amount        int    `json:"amount"`
	IsFundraising bool   `json:"is_fundraising"`
	IsHidden      bool   `json:"is_hidden"`
	Name          string `json:"name"`
}

type ProductDesign struct {
	DecorationLocation string `json:"decoration_location"`
	DecorationType     string `json:"decoration_type"`
	DesignID           string `json:"design_id"`
}

type Variant struct {
	DesignImage  interface{} `json:"design_image"`
	DesignImages []string    `json:"design_images"`
	SKU          string      `json:"sku"`
}

type Product struct {
	BasePrice      int             `json:"base_price"`
	Category       string          `json:"category"`
	IsMandatory    bool            `json:"is_mandatory"`
	Model          string          `json:"model"`
	PriceAdditions []PriceAddition `json:"price_additions"`
	ProductDesign  []ProductDesign `json:"product_design"`
	TaxCategory    string          `json:"tax_category"`
	Variants       []Variant       `json:"variants"`
}

type Store struct {
	CoachEmail       string    `json:"coach_email"`
	CustomStoreType  string    `json:"custom_store_type"`
	CustomerID       string    `json:"customer_id"`
	EdgewareOrderID  string    `json:"edgeware_order_id"`
	ExpiresAt        string    `json:"expires_at"`
	Group            string    `json:"group"`
	Industry         string    `json:"industry"`
	Name             string    `json:"name"`
	OpensAt          string    `json:"opens_at"`
	Organization     string    `json:"organization"`
	Products         []Product `json:"products"`
	SalespersonEmail string    `json:"salesperson_email"`
	Shipping         Shipping  `json:"shipping"`
	StoreType        string    `json:"store_type"`
	Subdomain        string    `json:"subdomain"`
	Template         string    `json:"template"`
	WorkOrderNumber  string    `json:"work_order_number"`
}

type RequestBody struct {
	Store Store `json:"store"`
}

// decorationTypes is the list of values to cycle through for decoration_type.
// Add or modify values as needed.
var decorationTypes = []string{
	"SCR",
	"EMB",
	"DTF",
	"HTV",
	"SUB",
}

func main() {
	url := "https://app.lions-qe.omg.pub/stores_api/"

	for i := 201; i <= 230; i++ {
		decorationType := decorationTypes[(i-1)%len(decorationTypes)]

		body := RequestBody{
			Store: Store{
				CoachEmail:      "erica.dregits@betheluniversity.edu",
				CustomStoreType: "Custom+",
				CustomerID:      "1113191",
				EdgewareOrderID: "17803891",
				ExpiresAt:       "2024-08-19T06:59:00.000Z",
				Group:           "Volleyball",
				Industry:        "Colleges & College Sports",
				Name:            fmt.Sprintf("Store%d", i),
				OpensAt:         "2024-08-04T08:00:00.000Z",
				Organization:    "Bethel University",
				Products: []Product{
					{
						BasePrice:   2600,
						Category:    "Apparel",
						IsMandatory: false,
						Model:       "312448",
						PriceAdditions: []PriceAddition{
							{
								Amount:        300,
								IsFundraising: true,
								IsHidden:      true,
								Name:          "Fundraising",
							},
						},
						ProductDesign: []ProductDesign{
							{
								DecorationLocation: "Full Front",
								DecorationType:     decorationType,
								DesignID:           "1018899",
							},
						},
						TaxCategory: "clothing",
						Variants: []Variant{
							{
								DesignImage:  nil,
								DesignImages: []string{},
								SKU:          "312448FR8A1300",
							},
						},
					},
				},
				SalespersonEmail: "dealer+127@ordermygear.com",
				Shipping: Shipping{
					DestinationAddress: DestinationAddress{
						City:       "Swanton",
						Country:    "US",
						FirstName:  "Erica",
						LastName:   "Dregits",
						PostalCode: "43558",
						State:      "OH",
						Street1:    "11500 Tailwinds Dr",
					},
					Fee:    1200,
					Method: "Flat Fee",
					Name:   "Ship to Home",
				},
				StoreType:       "Spirit Wear",
				Subdomain:       fmt.Sprintf("domain%d", i),
				Template:        "te-custom",
				WorkOrderNumber: "none",
			},
		}

		jsonData, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("[iteration %d] failed to marshal JSON: %v\n", i, err)
			continue
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("[iteration %d] failed to create request: %v\n", i, err)
			continue
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "insomnia/12.2.0")
		req.Header.Add("Authorization", "Bearer 25f33d36-1afe-4621-bb61-7e1ee8abc5dd")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("[iteration %d] request failed: %v\n", i, err)
			continue
		}

		responseBody, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Printf("[iteration %d] failed to read response: %v\n", i, err)
			continue
		}

		fmt.Printf("[iteration %d] name=Store%d decoration_type=%s subdomain=subdomain%d status=%s\n",
			i, i, decorationType, i, res.Status)
		fmt.Println(string(responseBody))
	}
}
