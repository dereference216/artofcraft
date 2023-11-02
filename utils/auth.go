package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (Acc *Account) SignUP() error {
	Acc.Config.PrintDebug("Username", Acc.Email)
	Acc.Config.PrintDebug("Password", Acc.Password)
	var resp *http.Response
	var rawBytes []byte
	var data = strings.NewReader(`{"password":"` + Acc.Password + `","email":"` + Acc.Email + `","consents":[]}`)
	req, err := http.NewRequest("POST", "https://api.picsart.com/user-account/auth/signup", data)
	if err != nil {
		return fmt.Errorf("SignUP err:%s", err)
	}
	req.Header.Set("Host", "api.picsart.com")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"`)
	req.Header.Set("G-Recaptcha-Action", "signup")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Platform", "website")
	req.Header.Set("G-Recaptcha-Token", Acc.RecaptchaToken)
	req.Header.Set("Token", "hI3KTwebV5yr2Gj")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Android"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://picsart.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://picsart.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	for i := 0; i <= Acc.Config.RequestRetry; i++ {
		resp, err = (*Acc.Client).Do(req)
		if err != nil {
			if i == Acc.Config.RequestRetry {
				return fmt.Errorf("TimedOut (SignUP) (MaxRetry)")
			}
			Acc.Config.PrintWarn(fmt.Sprintf("TimeOut (SignUP) (%d)", i), err.Error())
			continue
		}
		break
	}
	rawBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error In ReadAll (SignUP): %s", err)
	}
	defer resp.Body.Close()
	if bytes.Contains(rawBytes, []byte(`access_token`)) {
		var cred struct {
			Token struct {
				AccessToken string `json:"access_token"`
			} `json:"token"`
			Key string `json:"key"`
		}
		err = json.Unmarshal(rawBytes, &cred)
		if err != nil {
			return fmt.Errorf("jsonUnmarshal(SignUP) err: %s", err)
		}
		Acc.AuthToken = cred.Token.AccessToken
		Acc.XApiKey = cred.Key
		Acc.Etag = resp.Header.Get("Etag")
		Acc.Config.PrintDebug("Etag", Acc.Etag)
		Acc.Config.PrintDebug("AuthToken", Acc.AuthToken)
		Acc.Config.PrintDebug("XApi-Key", Acc.XApiKey)
		return nil
	}
	Acc.Config.PrintNetworkError(resp.StatusCode, "SignUP", string(rawBytes))
	return fmt.Errorf("unknown error occured(SignUP)")
}

func (Acc *Account) Link() error {
	var rawBytes []byte
	var resp *http.Response
	req, err := http.NewRequest("GET", "https://api.picsart.com/discord/link", nil)
	if err != nil {
		return fmt.Errorf("Link err:%s", err)
	}
	req.Header.Set("Host", "api.picsart.com")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Authorization", "Bearer "+Acc.AuthToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Platform", "website")
	req.Header.Set("X-Api-Key", Acc.XApiKey)
	req.Header.Set("Sec-Ch-Ua-Platform", `"Android"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://picsart.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://picsart.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("If-None-Match", Acc.Etag)
	for i := 0; i <= Acc.Config.RequestRetry; i++ {
		resp, err = (*Acc.Client).Do(req)
		if err != nil {
			if i == Acc.Config.RequestRetry {
				return fmt.Errorf("TimedOut (Link) (MaxRetry)")
			}
			Acc.Config.PrintWarn(fmt.Sprintf("TimeOut (Link) (%d)", i), err.Error())
			continue
		}
		break
	}
	defer resp.Body.Close()
	rawBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error In ReadAll (Link): %s", err)
	}
	defer resp.Body.Close()
	if bytes.Contains(rawBytes, []byte(`"status":"success",`)) {
		var prom struct {
			Response string `json:"response"`
		}
		err = json.Unmarshal(rawBytes, &prom)
		if err != nil {
			return fmt.Errorf("jsonUnmarshal(Link) err: %s", err)
		}
		Acc.ClaimLink = prom.Response
		Acc.Config.PrintGen("ClaimLink", Acc.ClaimLink)
		return nil
	}
	Acc.Config.PrintNetworkError(resp.StatusCode, "Link", string(rawBytes))
	return fmt.Errorf("unknown error occured(Link)")
}
