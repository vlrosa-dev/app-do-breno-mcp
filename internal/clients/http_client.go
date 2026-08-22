package clients

import (
	"net"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewApiClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,              // total de conexões idle no pool
		MaxIdleConnsPerHost: 10,               // idle por host — default da stdlib é só 2!
		MaxConnsPerHost:     0,                // 0 = sem limite de conexões simultâneas por host
		IdleConnTimeout:     90 * time.Second, // quanto tempo uma conexão idle fica no pool
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // timeout de estabelecer a conexão TCP
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second, // tempo até o primeiro byte do cabeçalho de resposta
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // timeout total da requisição
	}

	return client
}
