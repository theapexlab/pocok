package main

import (
	"os"
	"path"
	"pocok/src/amp/summary_email_template"
	"pocok/src/consumers/email_sender/create_email"
	"runtime"
)

func main() {
	writeFileRelative(summary_email_template.Get(), "/templates/summary_email.html")
	email_content, _ := create_email.GetHtmlSummary("https://test.com")
	writeFileRelative(email_content, "/emails/summary_email.html")
}

func writeFileRelative(content string, filepath string) {
	_, filename, _, _ := runtime.Caller(0)
	currentPath := path.Dir(filename)
	os.WriteFile(currentPath+filepath, []byte(content), 0644)
}
