// Package recaptcha allows you to interact with the Google reCAPTCHA API to verify user responses to the challenge.
package recaptcha

/*
 * The MIT License (MIT)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"
)

// defaultVerificationURL is The default URL that's used to verify the user's response to the challenge.
// @see https://developers.google.com/recaptcha/docs/verify#api-request
const defaultVerificationURL = "https://www.google.com/recaptcha/api/siteverify"

// RecaptchaErrorMap is the list of error codes mapped to a human-readable error code.
// @see https://developers.google.com/recaptcha/docs/verify#error-code-reference
var RecaptchaErrorMap = map[string]error{
	"missing-input-secret":   errors.New("The secret parameter is missing"),
	"invalid-input-secret":   errors.New("The secret parameter is invalid or malformed"),
	"missing-input-response": errors.New("The response parameter is missing"),
	"invalid-input-response": errors.New("The response parameter is invalid or malformed"),
}

// Response is the JSON structure that is returned by the verification API after a challenge response is verified.
// @see https://developers.google.com/recaptcha/docs/verify#api-response
type Response struct {
	Success bool `json:"success"`

	// Timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Challenge string `json:"challenge_ts"`

	// The hostname of the site where the reCAPTCHA was solved
	Hostname string `json:"hostname"`

	// Optional list of error codes returned by the service
	ErrorCodes []string `json:"error-codes"`
}

// The Recaptcha main structure. Its only purpose is to verify the user's response to a challenge with Google.
// You should initialize the structure with the Private Key that was supplied to you in the documentation.
type Recaptcha struct {
	PrivateKey string
	URL        string
}

// Verify the users's response to the reCAPTCHA challenge with the API server.
//
// The parameter response is obtained after the user successfully solves the challenge presented by the JS widget. The
// remoteip parameter is optional; just send it empty if you don't want to use it.
//
// This function will return a boolean that will have the final result returned by the API as well as an optional list
// of errors. They might be useful for logging purposed but you don't have to show them to the user.
func (r Recaptcha) Verify(response string, remoteip string) (Response, []error) {
	verificationURL := defaultVerificationURL
	if len(r.URL) != 0 {
		verificationURL = r.URL
	}

	params := url.Values{}

	if len(r.PrivateKey) > 0 {
		params.Set("secret", r.PrivateKey)
	}

	if len(response) > 0 {
		params.Set("response", response)
	}

	if net.ParseIP(remoteip) != nil {
		params.Set("remoteip", remoteip)
	}

	jsonResponse := Response{Success: false}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, httpError := httpClient.PostForm(verificationURL, params)

	if httpError != nil {
		return jsonResponse, []error{httpError}
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		return jsonResponse, []error{errors.New(httpResponse.Status)}
	}

	json.NewDecoder(httpResponse.Body).Decode(&jsonResponse)

	apiErrors := make([]error, len(jsonResponse.ErrorCodes))
	for i, singleError := range jsonResponse.ErrorCodes {
		if apiErrors[i] = errors.New(singleError); RecaptchaErrorMap[singleError] != nil {
			apiErrors[i] = RecaptchaErrorMap[singleError]
		}
	}

	return jsonResponse, apiErrors
}
