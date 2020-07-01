package recaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	fmt.Println("vim-go")
}

type recaptchaResult struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

const apiURL = "https://www.google.com/recaptcha/api/siteverify"

var ERR_EMPTY_ID = fmt.Errorf("Passed response ID is empty")

type Client struct {
	secret string
}

func New(secret string) Client {
	return Client{secret: secret}
}

func (c Client) Verify(responseID string) (bool, error) {
	if responseID == "" {
		return false, ERR_EMPTY_ID
	}
	values := url.Values{
		"secret":   {c.secret},
		"response": {responseID},
	}

	resp, err := http.DefaultClient.PostForm(apiURL, values)

	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}

	result := recaptchaResult{}

	if err = json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	if len(result.ErrorCodes) > 0 {
		return false, fmt.Errorf(strings.Join(result.ErrorCodes, " "))
	}

	return result.Success, nil
}
