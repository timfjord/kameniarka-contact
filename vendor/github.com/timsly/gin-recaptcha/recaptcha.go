/*
 */
package recaptcha

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type Config struct {
	Secret string
}

/*
*/
func Middleware(config Config) gin.HandlerFunc {
	// Create the Middleware function
	return func(c *gin.Context) {
		log.Println("rc middleware")
		challenge := c.Request.FormValue("g-recaptcha-response")
		ip := c.ClientIP()

		if challenge == "" || ValidateCaptcha(config.Secret, challenge, ip) == nil {
			c.Next()
		} else {
			c.AbortWithStatus(400)
		}
	}
}

// Validate recaptcha
func ValidateCaptcha(secret string, response string, remoteip string) error {
	var recaptchaResponse RecaptchaResponse
	resp, err := http.PostForm(recaptchaServerName,
		url.Values{"secret": {secret}, "response": {response}, "remoteip": {remoteip}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &recaptchaResponse)
	if err != nil {
		return err
	}
	if !recaptchaResponse.Success {
		return errors.New(strings.Join(recaptchaResponse.ErrorCodes[:], ","))
	}
	return nil
}
