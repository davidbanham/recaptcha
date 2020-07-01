package main

import (
	"net/http"
	"os"

	"github.com/davidbanham/recaptcha"
)

var recaptchaClient recaptcha.Client

func main() {
	recaptchaClient = recaptcha.New(os.Getenv("RECAPTCHA_SECRET"))

	http.HandleFunc("/", ServeForm)
	http.HandleFunc("/verified", CheckResponse)
	http.ListenAndServe(":8080", nil)
}

func ServeForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>
  <head>
    <title>reCAPTCHA demo: Simple page</title>
    <script src="https://www.google.com/recaptcha/api.js" async defer></script>
  </head>
  <body>
    <form action="?" method="POST">
      <div class="g-recaptcha" data-sitekey="` + os.Getenv("RECAPTCHA_SITE_KEY") + `"></div>
      <br/>
      <input type="submit" value="Submit">
    </form>
  </body>
</html>`))
}

func CheckResponse(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	verified, err := recaptchaClient.Verify(r.FormValue("g-recaptcha-response"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if verified {
		w.Write([]byte("yeah!"))
	} else {
		w.Write([]byte("nah"))
	}
}
