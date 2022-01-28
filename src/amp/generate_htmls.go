package main

import (
	"fmt"
	"os"
	"path"
	"pocok/src/amp/summary_email_template"
	"pocok/src/consumers/email_sender/create_email"
	"runtime"
)

func main() {
	writeFileRelative(summary_email_template.Get(), "/templates/summary_email.html")
	testUrl := "https://test.com"
	testLogoUrl := "https://github.com/theapexlab/pocok/raw/master/assets/pocok-logo.png"
	email_content, _ := create_email.GetHtmlSummary(testUrl, testLogoUrl)
	writeFileRelative(email_content, "/emails/summary_email.html")

	fmt.Println("⚡️ Succesfully generated HTML files.")
}

func writeFileRelative(content string, filepath string) {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	os.WriteFile(currentPath+filepath, []byte(content), 0644)
}
