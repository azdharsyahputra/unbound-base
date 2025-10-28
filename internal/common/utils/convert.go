package utils

import "fmt"

func ToUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}
func ToUintPtr(s string) *uint {
	u := ToUint(s)
	return &u
}