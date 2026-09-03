package breno

import (
	"context"
	"net/url"
)

// httpDoer é o subconjunto de *api.Client que o Service usa. Interface local
// para permitir fakes nos testes.
type httpDoer interface {
	REST(ctx context.Context, method, path string, query url.Values, body, out any) error
	RPC(ctx context.Context, fn string, body, out any) error
}

// Identity fornece o id do usuário autenticado, disparando login se necessário.
// Satisfeito por *auth.Manager.
type Identity interface {
	EnsureUserID(ctx context.Context) (string, error)
}

// Service implementa Financas sobre o client HTTP.
type Service struct {
	client httpDoer
	id     Identity
}

// NewService monta o Service.
func NewService(client httpDoer, id Identity) *Service {
	return &Service{client: client, id: id}
}

var _ Financas = (*Service)(nil)

func eq(v string) string { return "eq." + v }
