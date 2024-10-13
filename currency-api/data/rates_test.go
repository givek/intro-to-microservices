package data

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewRates(t *testing.T) {

	testLogger := log.New(os.Stdout, "TestNewRates", log.LstdFlags)

	exchangeRates, err := NewExchangeRates(testLogger)

	if err != nil {

		t.Fatal(err)

	}

	fmt.Printf("%#v", exchangeRates.rates)

}
