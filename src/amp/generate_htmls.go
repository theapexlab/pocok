package main

import (
	"fmt"
	"os"
	"path"
	"pocok/src/consumers/email_sender/create_email"
	"pocok/src/utils"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvError := godotenv.Load(".env.local")
	if loadEnvError != nil {
		utils.LogError("Error loading env", loadEnvError)
	}
	testUrl := os.Getenv("API_URL")
	if testUrl == "" {
		testUrl = "https://test.com"
	}
	testLogoUrl := "https://github.com/theapexlab/pocok/raw/master/assets/pocok-logo.png"

	email_content, _ := create_email.GetHtmlSummary(testUrl, testLogoUrl)
	writeFileRelative(email_content, "/emails/summary_email.html")

	fmt.Println("⚡️ Succesfully generated HTML files.")
}

func writeFileRelative(content string, filepath string) {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	writeError := os.WriteFile(currentPath+filepath, []byte(content), 0644)
	if writeError != nil {
		utils.LogError("Error writing file", writeError)
	}
}
