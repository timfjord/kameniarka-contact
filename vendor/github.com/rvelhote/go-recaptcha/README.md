[![](https://godoc.org/github.com/rvelhote/go-recaptcha?status.svg)](https://godoc.org/github.com/rvelhote/go-recaptcha) [![Build Status](https://travis-ci.org/rvelhote/go-recaptcha.svg?branch=master)](https://travis-ci.org/rvelhote/go-recaptcha) [![Code Climate](https://codeclimate.com/github/rvelhote/go-recaptcha/badges/gpa.svg)](https://codeclimate.com/github/rvelhote/go-recaptcha) [![Issue Count](https://codeclimate.com/github/rvelhote/go-recaptcha/badges/issue_count.svg)](https://codeclimate.com/github/rvelhote/go-recaptcha)

# reCAPTCHA Verification Package
This is a Go package that allows you to verify user response of reCAPTCHA challenged against the verification API. There are already a [few](https://github.com/HiFX/go-recaptcha) [packages](https://github.com/haisum/recaptcha) [available](https://github.com/dpapathanasiou/go-recaptcha) that you can use as alternatives. This package was specifically implemented as a learning experience as well as to be used in [another one of my projects](https://github.com/rvelhote/dnspropagation).

The package is only meant to help you with verifying the user's response with the API. It will not handle the form submission for you. That is up to you to deal with.

If you never used reCAPTCHA and want to know how to set it up in your project I recommend that you visit the [official documentation](https://developers.google.com/recaptcha/intro).

## Installation
Install this package as you would with any Go package by using `go get`.

```
go get github.com/rvelhote/go-recaptcha
```

## Usage
As mentioned in the beginning, this package only facilitates that interaction with the API; it will not handle the form submission itself.

Here is an example of how to use this package in a web application. Please note that the private key being used is a test key as defined by the [documentation FAQ](https://developers.google.com/recaptcha/docs/faq).

```
func verify(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	
	challenge := req.PostFormValue("g-recaptcha-response")
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)

	instance := recaptcha.Recaptcha{ PrivateKey: "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe" }
	response, err := instance.Verify(challenge, ip)

	log.Println(response.Success)
	log.Println(response.Challenge)
	log.Println(response.Hostname)
	log.Println(response.ErrorCodes)
	log.Println(err)
}
```

The `Verify` function that is part of the package will return a `boolean` with the end-result and a list of errors, if any, that might have occurred during the processing.

## Remote IP
In the example before I simply used the `RemoteAddr` from the request object. I believe it's not within the scope of this package to handle that so be advised that this is not the most complete/best way to get an IP Address if your application is behind a proxy or a load balancer. For a more complete verification you should also check `X-Forwarded-For` and/or `X-Real-IP`. There is also a `Forwarded` HTTP header specified in [RFC-7239](https://tools.ietf.org/html/rfc7239) however I'm not sure if it's widely used.

## Contributing
Contributions, suggestions and requests are welcome via Issue Tracker and via Pull Requests. I will do my best to reply and discuss.

Thank you!