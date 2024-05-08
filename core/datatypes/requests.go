package datatypes

import "errors"

type LoginRequestData struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
	Scope       string `json:"scope"`
}

func (r *LoginRequestData) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.ClientId == "" {
		return errors.New("client_id is required")
	}
	if r.RedirectUri == "" {
		return errors.New("redirect_uri is required")
	}
	if r.Scope == "" {
		return errors.New("scope is required")
	}
	return nil
}

type RegisterRequestData struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
}

func (r *RegisterRequestData) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.ClientId == "" {
		return errors.New("client_id is required")
	}
	return nil
}
