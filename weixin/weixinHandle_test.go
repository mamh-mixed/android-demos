package weixin

import (
	"fmt"
	"testing"
)

func TestGetAuthAccessToken(t *testing.T) {
	code := "041225b0ee79ba44394562d16282ac64"
	authAccessToken, err := GetAuthAccessToken(code)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("authAccessToken:::", *authAccessToken)
}
func TestGetAuthUserInfo(t *testing.T) {
	authAccessToken := "OezXcEiiBSKSxW0eoylIeHxPF3CuDDLvXrcgueDRsIXSbH71vV2XCjkG1QFj6jxxfYqOxLGYbrCvcEzC9jRq4OO23QXlIxTb9jb7WYI-uuIh0v-j6kqoR_gpWO3iQk2LZ_IMgQH44jJ5BikUzAUB8Q"
	openId := "oKwbgt2UZy7AMhWV0fcYPR5_QWjQ"
	authUserInfo, err := GetAuthUserInfo(authAccessToken, openId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("authUserInfo:::", *authUserInfo)
}

func TestRefreshAuthAccessToken(t *testing.T) {
	appid := appID
	refreshToken := "OezXcEiiBSKSxW0eoylIeHxPF3CuDDLvXrcgueDRsIXSbH71vV2XCjkG1QFj6jxxbE-A2oJSyi_gpzn23R-wweh9lUvu19RQhh3BdInwU-7DpAkoZNhlXn2oJO19DdrR03MCwVuUAFFKZTNBsPAtzg"
	refreshAATokenResp, err := RefreshAuthAccessToken(appid, refreshToken)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("refreshAATokenResp:::", *refreshAATokenResp)
}
