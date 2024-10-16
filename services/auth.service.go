package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/hubkudev/sentinel/configs"
	"github.com/hubkudev/sentinel/dto"
)

type AuthService interface {
	GoogleLogin(ctx *fiber.Ctx) error
	GoogleCallback(ctx *fiber.Ctx) error
	ConvertToken(accessToken string) (*dto.GooglePayload, error)
	Logout(ctx *fiber.Ctx) error
}

type AuthServiceImpl struct {
	UtilService  UtilService
	UserService  UserService
	StateStore   *session.Store
	SessionStore *session.Store
}

func (s *AuthServiceImpl) GoogleLogin(c *fiber.Ctx) error {
	// create google config first
	conf := configs.GoogleConfig()

	// generate random 32-long for state identification
	generated := s.UtilService.GenerateRandomID(32)

	sess, _ := s.StateStore.Get(c)
	sess.Set("session_state", generated)
	sess.Save()

	// create url for auth process.
	// we can pass state as someway to identify
	// and validate the login process.
	URL := conf.AuthCodeURL(generated)

	// redirect to the google authentication URL
	return c.Redirect(URL)

}

func (s *AuthServiceImpl) GoogleCallback(c *fiber.Ctx) error {
	// get session store for current context
	// NOTE: stateSess must be first, then sess, the second, I am not sure why is that,
	// there must be something wrong here.
	stateSess, stateErr := s.StateStore.Get(c)
	sess, sessErr := s.SessionStore.Get(c)
	if sessErr != nil || stateErr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// state from the session storage
	savedState := stateSess.Get("session_state")

	conf := configs.GoogleConfig()

	// get state from the callback
	state := c.Query("state")
	code := c.Query("code")

	// compare the state that is coming from the callback
	// with the one that is stored inside the session storage.
	if state != savedState {
		// Handle the invalid state error
		return c.Status(fiber.StatusBadRequest).SendString("Invalid state")
	}

	// exchange code that retrieved from google via
	// URL query parameter into token, this token
	// can be used later to query information of current user
	// from respective provider.
	token, err := conf.Exchange(c.Context(), code)
	if err != nil {
		log.Printf("Unable to exchange token: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := s.ConvertToken(token.AccessToken)
	if err != nil {
		log.Printf("Unable to convert token: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// find current user by provided email,
	// if the user found in the database, then we can just logged in,
	// if not, then register that user.
	isExist, err := s.UserService.FindByEmail(profile.Email)
	if isExist == nil || err != nil {
		// register user and save their data into database
		result, err := s.UserService.CreateUser(profile)
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
		return c.Status(fiber.StatusOK).Redirect("/misc/auth-redirect")
	}

	// Store the existed user's id in the session
	sess.Set("ID", isExist.ID.String())

	if err := sess.Save(); err != nil {
		log.Printf("Failed to save user session: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save user session")
	}

	return c.Status(fiber.StatusOK).Redirect("/misc/auth-redirect")
}

func (s *AuthServiceImpl) ConvertToken(accessToken string) (*dto.GooglePayload, error) {
	resp, httpErr := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken))
	if httpErr != nil {
		return nil, httpErr
	}
	defer resp.Body.Close()

	respBody, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return nil, bodyErr
	}

	// Unmarshal raw response body to a map
	var body map[string]interface{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, err
	}

	// if json body containing error,
	// then the token is indeed invalid. return invalid token err
	if body["error"] != nil {
		return nil, errors.New("Invalid token")
	}

	// Bind JSON into struct
	var data dto.GooglePayload
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
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
