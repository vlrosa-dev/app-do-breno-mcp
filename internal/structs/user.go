package structs

import "time"

type UserInfoResponse struct {
	ID               string    `json:"id"`
	Aud              string    `json:"aud"`
	Role             string    `json:"role"`
	Email            string    `json:"email"`
	EmailConfirmedAt time.Time `json:"email_confirmed_at"`
	Phone            string    `json:"phone"`
	ConfirmedAt      time.Time `json:"confirmed_at"`
	LastSignInAt     time.Time `json:"last_sign_in_at"`
}
