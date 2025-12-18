package main

import (
	"crypto/tls"
	_ "unsafe"

	"github.com/agiledragon/gomonkey/v2"
)

//lint:ignore U1000 tlsKeyShare is used in tlsClientHelloMsg
type tlsKeyShare struct {
	group tls.CurveID
	data  []byte
}

//lint:ignore U1000 tlsPSKIdentity is used in tlsClientHelloMsg
type tlsPSKIdentity struct {
	label               []byte
	obfuscatedTicketAge uint32
}

//lint:ignore U1000 tlsClientHelloMsg is a mirror of crypto/tls.keyShare 
type tlsClientHelloMsg struct {
	original                         []byte
	vers                             uint16
	random                           []byte
	sessionId                        []byte
	cipherSuites                     []uint16
	compressionMethods               []uint8
	serverName                       string
	ocspStapling                     bool
	supportedCurves                  []tls.CurveID
	supportedPoints                  []uint8
	ticketSupported                  bool
	sessionTicket                    []uint8
	supportedSignatureAlgorithms     []tls.SignatureScheme
	supportedSignatureAlgorithmsCert []tls.SignatureScheme
	secureRenegotiationSupported     bool
	secureRenegotiation              []byte
	extendedMasterSecret             bool
	alpnProtocols                    []string
	scts                             bool
	supportedVersions                []uint16
	cookie                           []byte
	keyShares                        []tlsKeyShare
	earlyData                        bool
	pskModes                         []uint8
	pskIdentities                    []tlsPSKIdentity
	pskBinders                       [][]byte
	quicTransportParameters          []byte
	encryptedClientHello             []byte
	// extensions are only populated on the server-side of a handshake
	extensions []uint16
}

//go:linkname clientHelloMsgMarshal crypto/tls.(*clientHelloMsg).marshal
func clientHelloMsgMarshal(m *tlsClientHelloMsg) ([]byte, error)

//go:linkname clientHelloMsgMarshalMsg crypto/tls.(*clientHelloMsg).marshalMsg
func clientHelloMsgMarshalMsg(m *tlsClientHelloMsg, echInner bool) ([]byte, error)

func clientHelloMsgMarshalPatched(m *tlsClientHelloMsg) ([]byte, error) {
	tmp := *m
	tmp.serverName = ""
	return clientHelloMsgMarshalMsg(&tmp, false)
}

func patchTLSServerName() *gomonkey.Patches {
	return gomonkey.ApplyFunc(clientHelloMsgMarshal, clientHelloMsgMarshalPatched)
}