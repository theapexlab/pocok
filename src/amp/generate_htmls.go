package main

import (
	"fmt"
	"os"
	"path"
	"pocok/src/consumers/email_sender/create_email"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	url := os.Getenv("API_URL")
	if url == "" {
		url = "https://test.com"
	}

	email_content, _ := create_email.GetHtmlSummary(url)
	writeFileRelative(email_content, "/emails/summary_email.html")

	fmt.Println("⚡️ Succesfully generated HTML files.")
}

func writeFileRelative(content string, filepath string) {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	os.WriteFile(currentPath+filepath, []byte(content), 0644)
}
