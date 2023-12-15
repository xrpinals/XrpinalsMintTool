package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	walletUrl = "http://192.168.1.73:8093"
)

func TestGetAssetInfo(t *testing.T) {
	resp, err := GetAssetInfo(walletUrl, "BTC")
	if err != nil {
		t.Fatal(err.Error())
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(string(respBody))
}

func TestGetAddressBalance(t *testing.T) {
	resp, err := GetAddressBalance(walletUrl, "mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs", "1.3.1")
	if err != nil {
		t.Fatal(err.Error())
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(string(respBody))
}

func TestGetInfo(t *testing.T) {
	resp, err := GetInfo(walletUrl)
	if err != nil {
		t.Fatal(err.Error())
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(string(respBody))
}

func TestGetRefBlockInfo(t *testing.T) {
	refBlockNum, refBlockPrefix, err := GetRefBlockInfo(walletUrl)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(refBlockNum)
	fmt.Println(refBlockPrefix)
}
