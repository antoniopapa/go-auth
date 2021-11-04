package controllers

import (
	"math/rand"
	"net/smtp"

	"../database"
	"../models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Forgot(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	token := RandStringRunes(12)

	passwordReset := models.PasswordReset{
		Email: data["email"],
		Token: token,
	}

	database.DB.Create(&passwordReset)

	from := "admin@example.com"

	to := []string{
		data["email"],
	}

	url := "http://localhost:3000/reset/" + token

	message := []byte("Click <a href=\"" + url + "\">here</a> to reset your password!")

	err := smtp.SendMail("0.0.0.0:1025", nil, from, to, message)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Reset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	var passwordReset = models.PasswordReset{}

	if err := database.DB.Where("token = ?", data["token"]).Last(&passwordReset); err.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid token!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	database.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
