package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DoPayment(email string) (string, error) {

	url := "https://api.chapa.co/v1/transaction/initialize"
	method := "POST"
	plString := fmt.Sprintf(`{
		"amount":"69",
		 "currency": "ETB",
		 "email": "%v",
		 "first_name": "Bilen",
		 "last_name": "Gizachew",
		 "phone_number": "0912345678",
		 "tx_ref": "chewatatest-8182",
		 "callback_url": "https://webhook.site/077164d6-29cb-40df-ba29-8a00e59a7e60",
		 "return_url": "https://www.github.com/mahider-t",
		 "customization[title]": "Payment for my favorite merchant",
		 "customization[description]": "I love online payments"
		 }`, email)

	// fmt.Println("My plString is :", plString)

	payload := strings.NewReader(plString)

	fmt.Println("My payload is :", payload)

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
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshalling json", err)
		return "", err
	}
	checkoutURL, ok := jsonBody["data"].(map[string]interface{})["checkout_url"].(string)
	if ok {
		fmt.Println("Checkout URL:", checkoutURL)
	} else {
		fmt.Println("Error: Could not access checkout URL")
		return "", err
	}
	fmt.Print(jsonBody)
	return checkoutURL, nil
}

func VerifyPayment() {

	url := "https://api.chapa.co/v1/transaction/verify/chewatatest-8182"
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer CHASECK_TEST-7rRq6kCN5opIeUPSEauZjDyVHESdPJoJ")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
