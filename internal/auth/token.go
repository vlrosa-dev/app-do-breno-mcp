package auth

import "time"

// Token é o resultado de um login bem-sucedido que vale a pena reaproveitar
// entre execuções. O refresh_token do provedor é intencionalmente ignorado:
// a re-autenticação usa sempre a rota de senha.
type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	UserID      string    `json:"user_id"`
}

// Valid informa se o token ainda pode ser usado, com uma folga (skew) para
// evitar usar um token que expira no meio de uma requisição.
func (t Token) Valid(skew time.Duration) bool {
	if t.AccessToken == "" {
		return false
	}
	if t.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().Add(skew).Before(t.ExpiresAt)
}

// Credentials são o e-mail e a senha usados no grant_type=password.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"senha"`
}

func (c Credentials) complete() bool {
	return c.Email != "" && c.Password != ""
}
