package data

import "testing"

func TestCheckValidation(t *testing.T) {

	p := &Product{
		Name:  "Water",
		Price: 1,
		SKU:   "asd-ad-asd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
