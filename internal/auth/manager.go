package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// expirySkew é a folga com que consideramos um token "perto de expirar".
const expirySkew = 60 * time.Second

// Manager entrega um access token válido, reautenticando pela rota de senha
// quando o token em cache expira. É seguro para uso concorrente.
type Manager struct {
	hc      *http.Client
	baseURL string
	apiKey  string
	store   Store

	mu    sync.Mutex
	creds Credentials
	token Token
}

// Options configura o Manager.
type Options struct {
	// HTTPClient usado só para os endpoints /auth/v1/*. Se nil, usa um client
	// com timeout padrão. Nunca passe aqui o client com auto-refresh.
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
	Store      Store
	// Credentials têm precedência sobre as que estiverem no Store.
	Credentials Credentials
}

// NewManager carrega o estado persistido (token + credenciais) e devolve um
// Manager pronto para uso.
func NewManager(opts Options) (*Manager, error) {
	if opts.BaseURL == "" || opts.APIKey == "" {
		return nil, fmt.Errorf("auth: BaseURL e APIKey são obrigatórios")
	}
	if opts.Store == nil {
		return nil, fmt.Errorf("auth: Store é obrigatório")
	}
	hc := opts.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 20 * time.Second}
	}

	m := &Manager{
		hc:      hc,
		baseURL: opts.BaseURL,
		apiKey:  opts.APIKey,
		store:   opts.Store,
		creds:   opts.Credentials,
	}

	f, err := opts.Store.Load()
	if err != nil {
		return nil, err
	}
	if !m.creds.complete() {
		m.creds = Credentials{Email: f.Email, Password: f.Password}
	}
	if f.Token != nil {
		m.token = *f.Token
	}
	return m, nil
}

// Token devolve um access token válido, reautenticando se necessário.
func (m *Manager) Token(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.token.Valid(expirySkew) {
		return m.token.AccessToken, nil
	}
	return m.login(ctx)
}

// ForceLogin ignora o cache e reautentica. Usado quando a API responde 401.
func (m *Manager) ForceLogin(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.login(ctx)
}

// Login reautentica com credenciais explícitas e as persiste. Usado pelo
// comando `breno auth login`.
func (m *Manager) Login(ctx context.Context, creds Credentials) (Token, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.creds = creds
	if _, err := m.login(ctx); err != nil {
		return Token{}, err
	}
	return m.token, nil
}

// Logout apaga o estado persistido e o cache em memória.
func (m *Manager) Logout() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.token = Token{}
	m.creds = Credentials{}
	return m.store.Clear()
}

// Status devolve o token em cache e se ele ainda é válido.
func (m *Manager) Status() (Token, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.token, m.token.Valid(expirySkew)
}

// UserID devolve o id do usuário do último login (vazio se nunca autenticou).
func (m *Manager) UserID() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.token.UserID
}

// EnsureUserID garante um token válido (autenticando se preciso) e devolve o id
// do usuário. Implementa breno.Identity.
func (m *Manager) EnsureUserID(ctx context.Context) (string, error) {
	if _, err := m.Token(ctx); err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.token.UserID, nil
}

// login executa o fluxo de senha e persiste o resultado. Requer m.mu travado.
func (m *Manager) login(ctx context.Context) (string, error) {
	if !m.creds.complete() {
		return "", fmt.Errorf("sem credenciais: defina API_USERNAME/API_PASSWORD ou rode `breno auth login`")
	}

	auth, err := passwordLogin(ctx, m.hc, m.baseURL, m.apiKey, m.creds)
	if err != nil {
		return "", err
	}
	userID, err := fetchUserID(ctx, m.hc, m.baseURL, m.apiKey, auth.AccessToken)
	if err != nil {
		return "", err
	}

	m.token = Token{
		AccessToken: auth.AccessToken,
		ExpiresAt:   auth.expiry(),
		UserID:      userID,
	}

	if err := m.persist(); err != nil {
		return "", err
	}
	return m.token.AccessToken, nil
}

// persist grava credenciais + token, preservando os demais campos do arquivo.
func (m *Manager) persist() error {
	f, err := m.store.Load()
	if err != nil {
		return err
	}
	f.Email = m.creds.Email
	f.Password = m.creds.Password
	tok := m.token
	f.Token = &tok
	return m.store.Save(f)
}
