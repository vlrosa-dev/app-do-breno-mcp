package breno

import (
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

var tiposValidos = map[TipoMovimentacao]bool{
	Diario: true, Saida: true, Cartao: true, Entrada: true, Economia: true,
}

var frequenciasValidas = map[Frequencia]bool{
	Mensalmente: true, Semanalmente: true, Diariamente: true, Parcelado: true,
}

// normalizeTipo aceita o valor em qualquer caixa ("diario", "DIARIO", ...) e
// devolve a forma canônica.
func normalizeTipo(s string) (TipoMovimentacao, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DIARIO":
		return Diario, nil
	case "SAIDA":
		return Saida, nil
	case "CARTAO":
		return Cartao, nil
	case "ENTRADA":
		return Entrada, nil
	case "ECONOMIA":
		return Economia, nil
	default:
		return "", fmt.Errorf("tipo inválido %q: use Diario, Saida, Cartao, Entrada ou Economia", s)
	}
}

// normalizeFrequencia aceita o valor em qualquer caixa e devolve a forma canônica.
func normalizeFrequencia(s string) (Frequencia, error) {
	f := Frequencia(strings.ToLower(strings.TrimSpace(s)))
	if !frequenciasValidas[f] {
		return "", fmt.Errorf("frequência inválida %q: use mensalmente, semanalmente, diariamente ou parcelado", s)
	}
	return f, nil
}

// parseData valida o formato YYYY-MM-DD. Vazio => data de hoje.
func parseData(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now().Format(dateLayout), nil
	}
	if _, err := time.Parse(dateLayout, s); err != nil {
		return "", fmt.Errorf("data %q inválida: use o formato YYYY-MM-DD", s)
	}
	return s, nil
}

// mustDate valida o formato YYYY-MM-DD e exige valor presente.
func mustDate(field, s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("%s é obrigatória (formato YYYY-MM-DD)", field)
	}
	return parseData(s)
}

// normalize valida e preenche defaults da entrada de criação, devolvendo os
// valores canônicos já resolvidos.
func (in *NovaMovimentacao) normalize() error {
	tipo, err := normalizeTipo(string(in.Tipo))
	if err != nil {
		return err
	}
	in.Tipo = tipo

	if in.Valor <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	if strings.TrimSpace(in.Descricao) == "" {
		return fmt.Errorf("descrição é obrigatória")
	}

	data, err := parseData(in.Data)
	if err != nil {
		return err
	}
	in.Data = data

	if in.Recorrencia != nil {
		freq, err := normalizeFrequencia(string(in.Recorrencia.Frequencia))
		if err != nil {
			return err
		}
		in.Recorrencia.Frequencia = freq
		if in.Recorrencia.Infinita {
			in.Recorrencia.Repeticoes = 24
		}
		if in.Recorrencia.Repeticoes <= 0 {
			return fmt.Errorf("número de repetições deve ser maior que zero")
		}
	}
	return nil
}

func (in *AtualizacaoMovimentacao) normalize() error {
	if strings.TrimSpace(in.ID) == "" {
		return fmt.Errorf("id da movimentação é obrigatório")
	}
	tipo, err := normalizeTipo(string(in.Tipo))
	if err != nil {
		return err
	}
	in.Tipo = tipo
	if in.Valor <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	if strings.TrimSpace(in.Descricao) == "" {
		return fmt.Errorf("descrição é obrigatória")
	}
	data, err := mustDate("data", in.Data)
	if err != nil {
		return err
	}
	in.Data = data
	return nil
}

func (in *AtualizacaoSerie) normalize() error {
	if strings.TrimSpace(in.RecorrenciaID) == "" {
		return fmt.Errorf("id da recorrência é obrigatório")
	}
	tipo, err := normalizeTipo(string(in.Tipo))
	if err != nil {
		return err
	}
	in.Tipo = tipo
	if in.Valor <= 0 {
		return fmt.Errorf("valor deve ser maior que zero")
	}
	if strings.TrimSpace(in.Descricao) == "" {
		return fmt.Errorf("descrição é obrigatória")
	}
	de, err := mustDate("data de origem (de)", in.De)
	if err != nil {
		return err
	}
	in.De = de
	nova, err := mustDate("nova data", in.NovaData)
	if err != nil {
		return err
	}
	in.NovaData = nova
	return nil
}

// resolve valida o filtro e devolve as datas efetivas. Quando De == Ate, o
// limite superior vira o dia seguinte (intervalo semiaberto).
func (f *FiltroMovimentacoes) resolve() (tipo TipoMovimentacao, de, ate, sel string, err error) {
	tipo, err = normalizeTipo(string(f.Tipo))
	if err != nil {
		return
	}
	de, err = mustDate("data inicial (de)", f.De)
	if err != nil {
		return
	}
	ate, err = mustDate("data final (ate)", f.Ate)
	if err != nil {
		return
	}
	if de == ate {
		d, _ := time.Parse(dateLayout, ate)
		ate = d.AddDate(0, 0, 1).Format(dateLayout)
	}
	sel = strings.TrimSpace(f.Select)
	if sel == "" {
		sel = "*"
	}
	return
}
