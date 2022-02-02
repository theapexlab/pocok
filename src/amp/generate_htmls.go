package main

import (
	"os"
	"path"
	"pocok/src/consumers/email_sender/create_email"
	"pocok/src/utils"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	testUrl := os.Getenv("API_URL")
	testLogoUrl := "https://github.com/theapexlab/pocok/raw/master/assets/pocok-logo.png"
	if testUrl == "" {
		testUrl = "https://test.com"
	}

	email_content, _ := create_email.GetHtmlSummary(testUrl, testLogoUrl)
	writeFileRelative(email_content, "/emails/summary_email.html")

	utils.Log("⚡️ Succesfully generated HTML files.")
}

func writeFileRelative(content string, filepath string) {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	os.WriteFile(currentPath+filepath, []byte(content), 0644)
}
