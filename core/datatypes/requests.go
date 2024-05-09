package datatypes

type LoginRequestData struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
	Scope       string `json:"scope"`
}

func (r *LoginRequestData) Validate() map[string]string {
	var errors map[string]string = make(map[string]string)
	if r.Email == "" {
		errors["email"] = "email is required"
	}
	if r.Password == "" {
		errors["password"] = "password is required"
	}
	if r.ClientId == "" {
		errors["client_id"] = "client_id is required"
	}
	if r.RedirectUri == "" {
		errors["redirect_uri"] = "redirect_uri is required"
	}
	if r.Scope == "" {
		errors["scope"] = "scope is required"
	}
	return errors
}

type RegisterRequestData struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
}

func (r *RegisterRequestData) Validate() map[string]string {
	var errors map[string]string = make(map[string]string)
	if r.Email == "" {
		errors["email"] = "email is required"
	}
	if r.Username == "" {
		errors["username"] = "username is required"
	}
	if r.Password == "" {
		errors["password"] = "password is required"
	}
	if r.ClientId == "" {
		errors["client_id"] = "client_id is required"
	}
	return errors
}

type IntrospectRequestData struct {
	Token string `json:"token"`
}

func (r *IntrospectRequestData) Validate() map[string]string {
	var errors map[string]string = make(map[string]string)
	if r.Token == "" {
		errors["token"] = "token is required"
	}
	return errors
}
