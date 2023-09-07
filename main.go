package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Domain struct {
	Available  bool    `json:"available"`
	Currency   string  `json:"currency"`
	Definitive bool    `json:"definitive"`
	Domain     string  `json:"domain"`
	Period     int     `json:"period"`
	Price      float64 `json:"price"`
}

func main() {
	domains, err := readDomains()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	savedAvailaableDomains, err := readSavedDomains()

	availaableDomains := []string{}

	var savedIndex int = 0
	for index, domain := range domains {
		if savedAvailaableDomains[savedIndex] == domain {
			availaableDomains = append(availaableDomains, domain)
			if savedIndex < len(savedAvailaableDomains)-1 {
				savedIndex++
			}
			continue
		}
		fmt.Printf("Checking domain %d/%d, domain name: %s\n", index+1, len(domains), domain)
		availaable, err := checkDomain(domain)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if availaable {
			availaableDomains = append(availaableDomains, domain)
		}
	}

	err = writeDomains(availaableDomains)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeDomains(domains []string) error {
	err := os.WriteFile("availibleDomains.txt", []byte(strings.Join(domains, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func checkDomain(domainName string) (bool, error) {

	url := fmt.Sprintf("https://api.ote-godaddy.com/v1/domains/available?domain=%s&checkType=FAST&forTransfer=false", domainName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "sso-key :")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	var domain Domain
	err = json.Unmarshal(body, &domain)
	if err != nil {
		return false, err
	}

	return domain.Available, nil
}

func readDomains() ([]string, error) {
	file, err := os.ReadFile("domains.txt")
	if err != nil {
		return nil, err
	}

	domains := strings.Split(string(file), "\n")

	return domains, nil
}

func readSavedDomains() ([]string, error) {
	file, err := os.ReadFile("availibleDomains.txt")
	if err != nil {
		return nil, err
	}

	domains := strings.Split(string(file), "\n")

	return domains, nil
}
