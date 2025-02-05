package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var receiverEmails = []string{"rishavkumar2700@gmail.com"}

func sendEmail(name, email, message, service string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderPass := os.Getenv("SENDER_PASS")
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	subject := "New Form Submission"
	body := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s\nService: %s", name, email, message, service)
	msg := []byte("Subject: " + subject + "\n\n" + body)

	// Send email with TLS
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, receiverEmails, msg)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))
	app.Post("/submit", func(c *fiber.Ctx) error {
		form := struct {
			Name    string `form:"name"`
			Email   string `form:"email"`
			Enquiry string `form:"enquiry"`
			Service string `form:"service"`
		}{}

		if err := c.BodyParser(&form); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid form submission")
		}

		if err := sendEmail(form.Name, form.Email, form.Enquiry, form.Service); err != nil {
			log.Println("Error sending email:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to send email")
		}

		return c.SendString("Form submitted successfully!")
	})

	log.Fatal(app.Listen(":8010"))
}
