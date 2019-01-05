package plurgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
	//"time"
)

type PlurkCredentials struct {
	ConsumerToken  string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

var baseURL = "https://www.plurk.com"

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://www.plurk.com/OAuth/request_token",
	ResourceOwnerAuthorizationURI: "https://www.plurk.com/OAuth/authorize",
	TokenRequestURI:               "https://www.plurk.com/OAuth/access_token",
}

var plurkOAuth PlurkCredentials
var signinOAuthClient oauth.Client

func ReadCredentials(credPath string) (*PlurkCredentials, error) {
	return readCredentials(credPath)
}

func readCredentials(credPath string) (*PlurkCredentials, error) {
	b, err := ioutil.ReadFile(credPath)
	if err != nil {
		return nil, err
	}
	var cred PlurkCredentials
	err = json.Unmarshal(b, &cred)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func doAuth(requestToken *oauth.Credentials) (*oauth.Credentials, error) {
	_url := oauthClient.AuthorizationURL(requestToken, nil)
	fmt.Println("Open the following URL and authorize it:", _url)

	var pinCode string
	fmt.Print("Input the PIN code: ")
	fmt.Scan(&pinCode)
	accessToken, _, err := oauthClient.RequestToken(http.DefaultClient, requestToken, pinCode)
	if err != nil {
		log.Fatal("failed to request token:", err)
	}
	return accessToken, nil
}

func getAccessToken(impl_ func(*PlurkCredentials) (*oauth.Credentials, bool, error),
	cred *PlurkCredentials) (*oauth.Credentials, bool, error) {
	return impl_(cred)
}

func GetAccessToken(cred *PlurkCredentials) (*oauth.Credentials, bool, error) {
	return getAccessToken(getAccessToken_, cred)
}

func getAccessToken_(cred *PlurkCredentials) (*oauth.Credentials, bool, error) {
	oauthClient.Credentials.Token = cred.ConsumerToken
	oauthClient.Credentials.Secret = cred.ConsumerSecret

	authorized := false
	var token *oauth.Credentials
	if cred.AccessToken != "" && cred.AccessSecret != "" {
		token = &oauth.Credentials{cred.AccessToken, cred.AccessSecret}
	} else {
		requestToken, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, "", nil)
		if err != nil {
			log.Printf("failed to request temporary credentials: %v", err)
			return nil, false, err
		}
		token, err = doAuth(requestToken)
		if err != nil {
			log.Printf("failed to request temporary credentials: %v", err)
			return nil, false, err
		}

		cred.AccessToken = token.Token
		cred.AccessSecret = token.Secret
		authorized = true
	}
	return token, authorized, nil
}

func callAPI(impl_ func(*oauth.Credentials, string, map[string]string) ([]byte, error),
	token *oauth.Credentials, _url string, opt map[string]string) ([]byte, error) {
	return impl_(token, _url, opt)
}

func CallAPI(token *oauth.Credentials, _url string, opt map[string]string) ([]byte, error) {
	return callAPI(callAPI_, token, _url, opt)
}

func callAPI_(token *oauth.Credentials, _url string, opt map[string]string) ([]byte, error) {
	var apiURL = baseURL + _url
	param := make(url.Values)
	for k, v := range opt {
		param.Set(k, v)
	}
	oauthClient.SignParam(token, "POST", apiURL, param)
	res, err := http.PostForm(apiURL, url.Values(param))
	if err != nil {
		log.Println("failed to call API:", err, apiURL, param)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("failed to get response:", err)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Println("failed to call API err=200:", err, apiURL, param)
		return nil, fmt.Errorf("%s", string(body))
	}
	return body, nil
}
