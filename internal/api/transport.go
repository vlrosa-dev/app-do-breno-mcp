package api

import (
	"io"
	"net/http"
)

// authTransport injeta `apikey` + `Authorization: Bearer <token>` em toda
// requisição e, ao receber 401, força um novo login e repete a requisição uma
// única vez. Mantém o Client.do() livre de lógica de autenticação.
type authTransport struct {
	base   http.RoundTripper
	auth   TokenProvider
	apiKey string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.auth.Token(req.Context())
	if err != nil {
		return nil, err
	}

	res, err := t.base.RoundTrip(t.decorated(req, token))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusUnauthorized {
		return res, nil
	}

	// 401 — tenta reautenticar. Só dá para repetir se o corpo for rebobinável.
	retry, rerr := rewind(req)
	if rerr != nil {
		return res, nil
	}
	token, err = t.auth.ForceLogin(req.Context())
	if err != nil {
		return res, nil
	}
	drain(res)
	return t.base.RoundTrip(t.decorated(retry, token))
}

// decorated devolve uma cópia do request com os headers de autenticação.
func (t *authTransport) decorated(req *http.Request, token string) *http.Request {
	r := req.Clone(req.Context())
	r.Header.Set("apikey", t.apiKey)
	r.Header.Set("Authorization", "Bearer "+token)
	return r
}

// rewind devolve um clone do request com o corpo reposicionado no início.
func rewind(req *http.Request) (*http.Request, error) {
	r := req.Clone(req.Context())
	if req.Body == nil || req.Body == http.NoBody {
		return r, nil
	}
	if req.GetBody == nil {
		return nil, io.ErrUnexpectedEOF
	}
	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	r.Body = body
	return r, nil
}

func drain(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 4<<10))
	_ = res.Body.Close()
}
