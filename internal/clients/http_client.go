package clients

import (
	"net"
	"net/http"
	"time"
)

// NewTransport devolve um *http.Transport com pool de conexões e timeouts
// afinados para a API do App do Breno.
func NewTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        100,              // total de conexões idle no pool
		MaxIdleConnsPerHost: 10,               // idle por host — default da stdlib é só 2!
		MaxConnsPerHost:     0,                // 0 = sem limite de conexões simultâneas por host
		IdleConnTimeout:     90 * time.Second, // quanto tempo uma conexão idle fica no pool
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // timeout de estabelecer a conexão TCP
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second, // tempo até o primeiro byte do cabeçalho de resposta
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewApiClient devolve um *http.Client usando o transport afinado.
func NewApiClient() *http.Client {
	return &http.Client{
		Transport: NewTransport(),
		Timeout:   30 * time.Second, // timeout total da requisição
	}
}
