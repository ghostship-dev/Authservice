package datatypes

type OTPRequest struct {
	Action string `json:"action"`
	OTP    string `json:"otp"`
}

func (r *OTPRequest) Validate() map[string]string {
	var errors map[string]string = make(map[string]string)
	if r.Action != "enable" && r.Action != "disable" && r.Action != "verify" {
		errors["action"] = "action not allowed! available actions: enable, disable, verify"
	}
	if r.OTP == "" && r.Action != "enable" {
		errors["otp"] = "otp is required for this action"
	}
	return errors
}
