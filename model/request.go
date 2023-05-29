package model

type DomainRequest struct {
	Name string `json:"domain,omitempty"`
	Ipv4 string `json:"ipv4,omitempty"`
}

type PasswordRequest struct {
	OldPassword  string `json:"oldPassword"`
	NewPassword  string `json:"newPassword"`
	NewRPassword string `json:"newRPassword"`
}
