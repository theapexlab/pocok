package utils

import (
	"fmt"
)

func LogError(desc string, e error) {
	fmt.Printf("❌ %s: %s", desc, e.Error())
}
