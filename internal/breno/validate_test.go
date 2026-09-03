package breno

import (
	"testing"
	"time"
)

func TestNormalizeTipo(t *testing.T) {
	cases := map[string]struct {
		want TipoMovimentacao
		ok   bool
	}{
		"diario":   {Diario, true},
		"DIARIO":   {Diario, true},
		" Saida ":  {Saida, true},
		"cartao":   {Cartao, true},
		"ENTRADA":  {Entrada, true},
		"economia": {Economia, true},
		"":         {"", false},
		"foo":      {"", false},
	}
	for in, exp := range cases {
		got, err := normalizeTipo(in)
		if exp.ok && (err != nil || got != exp.want) {
			t.Errorf("normalizeTipo(%q) = %q, %v; quer %q", in, got, err, exp.want)
		}
		if !exp.ok && err == nil {
			t.Errorf("normalizeTipo(%q) deveria falhar", in)
		}
	}
}

func TestNormalizeFrequencia(t *testing.T) {
	for _, ok := range []string{"mensalmente", "SEMANALMENTE", "Diariamente", "parcelado"} {
		if _, err := normalizeFrequencia(ok); err != nil {
			t.Errorf("normalizeFrequencia(%q) inesperado: %v", ok, err)
		}
	}
	if _, err := normalizeFrequencia("anual"); err == nil {
		t.Error("normalizeFrequencia(anual) deveria falhar")
	}
}

func TestParseData(t *testing.T) {
	hoje := time.Now().Format(dateLayout)
	if got, _ := parseData(""); got != hoje {
		t.Errorf("parseData(\"\") = %q, quer %q", got, hoje)
	}
	if _, err := parseData("2026-13-40"); err == nil {
		t.Error("data inválida deveria falhar")
	}
	if got, err := parseData("2026-09-03"); err != nil || got != "2026-09-03" {
		t.Errorf("parseData válido = %q, %v", got, err)
	}
}

func TestNovaMovimentacaoNormalize(t *testing.T) {
	in := NovaMovimentacao{Tipo: "diario", Valor: 1.5, Descricao: "Padaria"}
	if err := in.normalize(); err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if in.Tipo != Diario || in.Data == "" {
		t.Fatalf("normalize não canonizou: %+v", in)
	}

	bad := []NovaMovimentacao{
		{Tipo: "diario", Valor: 0, Descricao: "x"},
		{Tipo: "diario", Valor: 1, Descricao: " "},
		{Tipo: "xpto", Valor: 1, Descricao: "x"},
		{Tipo: "diario", Valor: 1, Descricao: "x", Data: "03/09/2026"},
	}
	for i, b := range bad {
		if err := (&b).normalize(); err == nil {
			t.Errorf("caso inválido %d passou", i)
		}
	}
}

func TestNovaMovimentacaoNormalizeSerie(t *testing.T) {
	inf := NovaMovimentacao{
		Tipo: "saida", Valor: 100, Descricao: "Aluguel",
		Recorrencia: &Recorrencia{Frequencia: "mensalmente", Infinita: true},
	}
	if err := inf.normalize(); err != nil {
		t.Fatalf("normalize infinita: %v", err)
	}
	if inf.Recorrencia.Repeticoes != 24 {
		t.Fatalf("infinita deveria virar 24, veio %d", inf.Recorrencia.Repeticoes)
	}

	zero := NovaMovimentacao{
		Tipo: "saida", Valor: 100, Descricao: "Aluguel",
		Recorrencia: &Recorrencia{Frequencia: "mensalmente", Repeticoes: 0},
	}
	if err := zero.normalize(); err == nil {
		t.Fatal("repeticoes 0 deveria falhar")
	}
}

func TestFiltroResolveBumpsEqualDates(t *testing.T) {
	f := FiltroMovimentacoes{Tipo: "diario", De: "2026-09-03", Ate: "2026-09-03"}
	_, de, ate, sel, err := f.resolve()
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if de != "2026-09-03" || ate != "2026-09-04" {
		t.Fatalf("de/ate = %s/%s, quer 2026-09-03/2026-09-04", de, ate)
	}
	if sel != "*" {
		t.Fatalf("select default = %q", sel)
	}
}
