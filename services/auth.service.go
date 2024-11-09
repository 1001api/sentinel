package services

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthService interface {
	RegisterFirstUser(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type AuthServiceImpl struct {
	UtilService  UtilService
	UserService  UserService
	StateStore   *session.Store
	SessionStore *session.Store
}

func (s *AuthServiceImpl) RegisterFirstUser(c *fiber.Ctx) error {
	fullname := c.FormValue("fullname")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	payload := dto.CreateUserInput{
		Fullname:        fullname,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}

	if err := s.UtilService.ValidateInput(&payload); err != "" {
		return c.SendString(err)
	}

	adminExist, err := s.UserService.CheckAdminExist()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// if admin exist?
	// prevent from creating a new user from this endpoint.
	if adminExist {
		return c.SendString("New user registration is disabled. Contact your admin.")
	}

	sess, sessErr := s.SessionStore.Get(c)
	if sessErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// find current user by provided email,
	// if the user not found in the database, then register that user.
	isExist, err := s.UserService.FindByEmail(email)
	if isExist == nil || err != nil {
		// hash password
		hashed := s.UtilService.GenerateHash(password)
		if hashed == "" {
			log.Printf("Failed to hash password: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to register user")
		}

		// register user and save their data into database
		result, err := s.UserService.CreateUser(&gen.CreateUserParams{
			Fullname:       fullname,
			Email:          email,
			PasswordHashed: hashed,
			ProfileUrl: pgtype.Text{
				String: fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", fullname),
				Valid:  true,
			},
			RootUser: true,
		})
		if err != nil {
			log.Printf("Failed to register user: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to register user")
		}

		// Store the user's id in the session
		sess.Set("ID", result.ID.String())

		// Save into memory session and.
		// saving also set a session cookie containing session_id
		if err := sess.Save(); err != nil {
			log.Printf("Failed to save user session: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save user session")
		}

		// return immediately
		c.Set("HX-Redirect", "/events")
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendString("User with the same email already exist")
}

func (s *AuthServiceImpl) Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	sess, sessErr := s.SessionStore.Get(c)
	if sessErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	isExist, err := s.UserService.FindByEmailWithHash(email)
	if isExist != nil && err == nil {
		// compare hashed password
		validHashed := s.UtilService.CompareHash(password, isExist.PasswordHashed)
		if !validHashed {
			return c.SendString("Wrong email or password")
		}

		// Store the user's id in the session
		sess.Set("ID", isExist.ID.String())

		// Save into memory session and.
		// saving also set a session cookie containing session_id
		if err := sess.Save(); err != nil {
			log.Printf("Failed to save user session: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save user session")
		}

		// return immediately
		c.Set("HX-Redirect", "/events")
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendString("Wrong email or password")
}

func (s *AuthServiceImpl) Logout(c *fiber.Ctx) error {
	sess, err := s.SessionStore.Get(c)
	if err != nil {
		log.Println(err.Error())
	}

	// destroy current user session
	sess.Destroy()

	return c.Status(fiber.StatusOK).Redirect("/")
}
