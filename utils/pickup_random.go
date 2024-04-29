package utils

import "math/rand"

func PickupRandom(list []string) string {
	return list[rand.Intn(len(list))]
}
