package utils

import (
	"log"
)

func Log(i ...interface{}) {
	log.Println(i...)
}

func Logf(s string, i ...interface{}) {
	log.Printf(s, i...)
}

func LogFatal(i ...interface{}) {
	log.Fatal(i...)
}

func LogFatalf(s string, i ...interface{}) {
	log.Fatalf(s, i...)
}

func LogError(desc string, e error) {
	log.Printf("❌ %s: %s", desc, e.Error())
}
