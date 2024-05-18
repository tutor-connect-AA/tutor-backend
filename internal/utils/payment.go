package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DoPayment(email, tx_ref, return_url string, amount int) (string, error) {

	url := "https://api.chapa.co/v1/transaction/initialize"
	method := "POST"
	plString := fmt.Sprintf(`{
		"amount":"%v",
		 "currency": "ETB",
		 "email": "%v",
		 "tx_ref": "%v",
		 "callback_url": "https://webhook.site/077164d6-29cb-40df-ba29-8a00e59a7e60",
		 "return_url": "%v",
		 "customization[title]": "Payment for my favourite merchant",
  "customization[description]": "I love online payments"
		 }`, amount, email, tx_ref, return_url)

	// fmt.Println("My plString is :", plString)

	payload := strings.NewReader(plString)

	fmt.Println("\nMy payload is :\n", payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer CHASECK_TEST-7rRq6kCN5opIeUPSEauZjDyVHESdPJoJ")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("body at DoPayment", string(body))
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshalling json", err)
		return "", err
	}
	if data, ok := jsonBody["data"].(map[string]interface{}); ok {
		if checkoutURL, ok := data["checkout_url"].(string); ok {
			fmt.Printf("Checkout URL at DoPayment is %v", checkoutURL)
			return checkoutURL, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}

}

func VerifyPayment(tx_ref string) (string, error) {

	// url := `https://api.chapa.co/v1/transaction/verify/${tx_ref}`
	url := fmt.Sprintf(`https://api.chapa.co/v1/transaction/verify/%v`, tx_ref)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer CHASECK_TEST-7rRq6kCN5opIeUPSEauZjDyVHESdPJoJ")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(string(body))
	return string(body), nil
}
