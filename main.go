package domain

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	domains, err := readDomains()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	avaliabileDomains := checkDomains(domains)
}


