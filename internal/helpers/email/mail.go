package mail_helper

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

func readHTMLFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func processHTML(htmlContent string, body string, title string) string {
	htmlContent = strings.Replace(htmlContent, "{{body}}", body, -1)
	htmlContent = strings.Replace(htmlContent, "{{title}}", title, -1)
	return htmlContent
}

func SendMail(to []string, message []byte, title []byte) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "./templates/template.html")

	htmlContent, err := readHTMLFile(path)
	if err != nil {
		log.Fatal("Error reading HTML file:", err)
		return
	}

	htmlContent = processHTML(htmlContent, string(message), string(title))

	userName, password := os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD")
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_USERNAME"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", htmlContent)

	d := gomail.NewDialer("smtp.gmail.com", 587, userName, password)

	d.DialAndSend(m)
}
