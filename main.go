package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

// CodeURL is used to get code.
const CodeURL = "https://api.instagram.com/oauth/authorize/?client_id=%s&redirect_uri=%s&response_type=code"

// AccessTokenURL is URL to get access token.
const AccessTokenURL = "https://api.instagram.com/oauth/access_token"

type (
	// Result is result of request.
	Result struct {
		AccessToken string `json:"access_token"`
		User        User   `json:"user"`
	}
	// User is object in Result.
	User struct {
		ID             string `json:"id"`
		UserName       string `json:"username"`
		FullName       string `json:"full_name"`
		ProfilePicture string `json:"profile_picture"`
	}
	// Error is error of the request.
	Error struct {
		Code         string `json:"code"`
		ErrorMessage string `json:"error_message"`
		Errortype    string `json:"error_type"`
	}
)

func main() {
	cliID, err := getInput("Please input your Client ID ")
	cliSec, err := getInput("Please input your Client secret ")
	redirectURL, err := getInput("Please input your redirect URL ")
	if err != nil {
		panic(err)
	}
	err = openCodeURL(cliID, redirectURL)
	if err != nil {
		panic(err)
	}
	code, err := getInput("Please input your code ([Your reqirect url]/?code=CODE) ")
	if err != nil {
		panic(err)
	}

	accessToken, err := getAccessToken(cliID, cliSec, redirectURL, code)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Access token: %s\n", accessToken)

}

func getInput(msg string) (string, error) {
	fmt.Print(msg, ">")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", errors.New("Failed to get input.")
}

func openCodeURL(cliID, redirectURL string) (err error) {
	url := fmt.Sprintf(CodeURL, cliID, redirectURL)
	cmd := exec.Command("open", url)
	return cmd.Run()
}

func getAccessToken(cliID, cliSec, redirectURL, code string) (accessToken string, err error) {

	v := url.Values{}
	v.Add("client_id", cliID)
	v.Add("client_secret", cliSec)
	v.Add("grant_type", "authorization_code")
	v.Add("redirect_uri", redirectURL)
	v.Add("code", code)

	resp, err := http.PostForm(AccessTokenURL, v)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusOK {
		var result Result
		err = json.Unmarshal(b, &result)
		if err != nil {
			return "", err
		}
		return result.AccessToken, nil
	} else if resp.StatusCode == http.StatusBadRequest {
		var e Error
		err = json.Unmarshal(b, &e)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("Error occured(Type: %s): %s\n", e.Errortype, e.ErrorMessage)
	}
	return "", fmt.Errorf("Undefined error occured.")

}
