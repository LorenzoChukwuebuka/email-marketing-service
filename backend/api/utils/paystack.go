package utils

import (
	"fmt"
	"io"
	"net/http"
)

type BankAccountDTO struct {
	AccountNumber string
	BankCode      string
}

func resolveAccount(account BankAccountDTO) (string, error) {
	url := fmt.Sprintf("https://api.paystack.co/bank/resolve?account_number=%s&bank_code=%s", account.AccountNumber, account.BankCode)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "your_paystack_secret_key")) // Replace with your actual Paystack secret key

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("account number invalid")
	}

	return string(body), nil
}

func Hello() {
	account := BankAccountDTO{
		AccountNumber: "1234567890",
		BankCode:      "011",
	}

	response, err := resolveAccount(account)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", response)
	}
}
