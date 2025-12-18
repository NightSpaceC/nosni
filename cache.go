package main

import (
	"crypto/tls"
	"sync"
)

// CertStorage is a simple certificate cache that keeps
// everything in memory.
type certStorage struct {
	certs map[string]*tls.Certificate
	mtx   sync.RWMutex
}

func (cs *certStorage) Fetch(hostname string, gen func() (*tls.Certificate, error)) (*tls.Certificate, error) {
	cs.mtx.RLock()
	cert, ok := cs.certs[hostname]
	cs.mtx.RUnlock()
	if ok {
		return cert, nil
	}

	cert, err := gen()
	if err != nil {
		return nil, err
	}

	cs.mtx.Lock()
	cs.certs[hostname] = cert
	cs.mtx.Unlock()

	return cert, nil
}

func newCertStorage() *certStorage {
	return &certStorage{
		certs: make(map[string]*tls.Certificate),
	}
}