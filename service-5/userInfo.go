package main

import (
	wapi "github.com/tera-insights/go-win64api"
	"github.com/tera-insights/go-win64api/shared"
)

func getLoggedUser() ([]shared.SessionDetails, error) {

	userDesc, err := wapi.ListLoggedInUsers()
	if err != nil {
		return nil, err
	}

	return userDesc, nil

}
