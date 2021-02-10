package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RegistrationReq struct {
    Credential map[string]interface{} `json:"credential"`
    FirstName  string `json:"firstName"`
    LastName   string `json:"lastName"`
    Email      string `json:"email"`
}

type RegistrationRes struct {
	User *User `json:"user"`
}

func handleRegister(res http.ResponseWriter, req *http.Request) {
	defer func() {
        if err := recover(); err != nil {
            fmt.Printf("Caught panic: %v\n", err)
			requestFailed(res, 500, "Request failed")
			return
        }
    }()

	fmt.Println("Handling /register")
	res.Header().Set("Content-Type", "application/json")

	var body RegistrationReq
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		requestFailed(res, 400, "Invalid Request")
		return
	}
	fmt.Printf("body = %v\n", body)

	email := sanitizeEmail(body.Email)
	if len(email) <= 0 {
		requestFailed(res, 400, "Invalid Request")
		return
	}

	_, err = ctx.FindByEmail(email)
	if err == nil {
		requestFailed(res, 400, "Email address already in use")
		return
	}

	createCredRes, credentialBody, err := ctx.api.CredentialCreate(body.Credential)
	if createCredRes.StatusCode != 201 {
		requestFailed(res, 500, "Unable to create credential")
		return
	}

	user := &User{
		NuID: credentialBody.NuID,
		Email: email,
		FirstName: body.FirstName,
		LastName: body.LastName,
	}
	dbRes := ctx.db.Create(user)
	if dbRes.RowsAffected != 1 {
		requestFailed(res, 500, "Unable to create user")
		return
	}
	res.WriteHeader(201)
	result := &RegistrationRes{User: user}
	json.NewEncoder(res).Encode(result)
}
