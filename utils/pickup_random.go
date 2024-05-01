package utils

import (
	"math/rand"
	"strconv"

	"github.com/ddosify/go-faker/faker"
)

func PickupRandom(list []string) string {
	return list[rand.Intn(len(list))]
}

func FakeCheckout() map[string]string {
	fackerGeneration := faker.NewFaker()

	fake_checkout := map[string]string{
		"email":                        fackerGeneration.RandomEmail(),
		"street_address":               fackerGeneration.RandomAddressStreetAddress(),
		"zip_code":                     strconv.Itoa(rand.Intn(10001-1000) + 1000),
		"city":                         fackerGeneration.RandomAddressCity(),
		"state":                        fackerGeneration.RandomCountryCode(),
		"country":                      fackerGeneration.RandomCountryCode(),
		"credit_card_number":           fackerGeneration.RandomBankAccountIban(),
		"credit_card_expiration_month": strconv.Itoa(rand.Intn(12) + 1),
		"credit_card_expiration_year":  strconv.Itoa(rand.Intn(70) + 2023),
		"credit_card_cvv":              strconv.Itoa(rand.Intn(999) + 100),
	}
	return fake_checkout
}
