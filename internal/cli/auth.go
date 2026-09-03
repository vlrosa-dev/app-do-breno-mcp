package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"app-do-breno-agent/internal/auth"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// newAuthGroupCmd cria o grupo `auth` (login/status/logout).
func newAuthGroupCmd(d *deps) *cobra.Command {
	group := &cobra.Command{Use: "auth", Short: "Autenticação (login por e-mail e senha)"}

	var email, senha string
	login := &cobra.Command{
		Use:   "login",
		Short: "Autentica e grava credenciais + token em ~/.breno/config.json",
		RunE: func(cmd *cobra.Command, _ []string) error {
			m, err := d.authManager()
			if err != nil {
				return err
			}
			if email == "" {
				email = firstNonEmpty(os.Getenv("API_USERNAME"), promptLine(cmd, "E-mail: "))
			}
			if senha == "" {
				senha = firstNonEmpty(os.Getenv("API_PASSWORD"), promptSecret(cmd, "Senha: "))
			}
			if email == "" || senha == "" {
				return fmt.Errorf("e-mail e senha são obrigatórios")
			}
			tok, err := m.Login(cmdContext(cmd), auth.Credentials{Email: email, Password: senha})
			if err != nil {
				return err
			}
			return okMessage(cmd, d.g, fmt.Sprintf("autenticado como %s (token expira em %s)",
				email, tok.ExpiresAt.Format("2006-01-02 15:04")))
		},
	}
	login.Flags().StringVar(&email, "email", "", "e-mail da conta")
	login.Flags().StringVar(&senha, "senha", "", "senha da conta (evite; prefira o prompt)")

	status := &cobra.Command{
		Use:   "status",
		Short: "Mostra o estado da autenticação",
		RunE: func(cmd *cobra.Command, _ []string) error {
			m, err := d.authManager()
			if err != nil {
				return err
			}
			tok, ok := m.Status()
			if tok.AccessToken == "" {
				return okMessage(cmd, d.g, "não autenticado — rode `breno auth login`")
			}
			estado := "válido"
			if !ok {
				estado = "expirado (será renovado no próximo uso)"
			}
			return okMessage(cmd, d.g, fmt.Sprintf("usuário %s — token %s (expira em %s)",
				tok.UserID, estado, tok.ExpiresAt.Format("2006-01-02 15:04")))
		},
	}

	logout := &cobra.Command{
		Use:   "logout",
		Short: "Apaga credenciais e token salvos",
		RunE: func(cmd *cobra.Command, _ []string) error {
			m, err := d.authManager()
			if err != nil {
				return err
			}
			if err := m.Logout(); err != nil {
				return err
			}
			return okMessage(cmd, d.g, "sessão encerrada")
		},
	}

	group.AddCommand(login, status, logout)
	return group
}

func promptLine(cmd *cobra.Command, label string) string {
	fmt.Fprint(cmd.ErrOrStderr(), label)
	r := bufio.NewReader(cmd.InOrStdin())
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func promptSecret(cmd *cobra.Command, label string) string {
	fmt.Fprint(cmd.ErrOrStderr(), label)
	if f, ok := cmd.InOrStdin().(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		b, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(cmd.ErrOrStderr())
		if err == nil {
			return strings.TrimSpace(string(b))
		}
	}
	// Fallback (não-TTY, ex.: testes): lê uma linha comum.
	return promptLine(cmd, "")
}
