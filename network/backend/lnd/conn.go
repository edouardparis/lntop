package lnd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	macaroon "gopkg.in/macaroon.v2"

	"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/macaroons"

	"github.com/edouardparis/lntop/config"
)

func newClientConn(c *config.Network) (*grpc.ClientConn, error) {
	macaroonBytes, err := hex.DecodeString(c.Macaroon)
	if err != nil {
		return nil, err
	}

	mac := &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(macaroonBytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	macConstraints := []macaroons.Constraint{
		// We add a time-based constraint to prevent replay of the
		// macaroon. It's good for 60 seconds by default to make up for
		// any discrepancy between client and server clocks, but leaking
		// the macaroon before it becomes invalid makes it possible for
		// an attacker to reuse the macaroon. In addition, the validity
		// time of the macaroon is extended by the time the server clock
		// is behind the client clock, or shortened by the time the
		// server clock is ahead of the client clock (or invalid
		// altogether if, in the latter case, this time is more than 60
		// seconds).
		macaroons.TimeoutConstraint(c.MacaroonTimeOut),

		// Lock macaroon down to a specific IP address.
		macaroons.IPLockConstraint(c.MacaroonIP),

		// ... Add more constraints if needed.
	}

	// Apply constraints to the macaroon.
	constrainedMac, err := macaroons.AddConstraints(mac, macConstraints...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cred, err := newCredentialsFromCert(c.Cert)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(constrainedMac)),
		// We need to use a custom dialer so we can also connect to unix sockets
		// and not just TCP addresses.
		grpc.WithDialer(lncfg.ClientAddressDialer(u.Port())),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.MaxMsgRecvSize)),
	}

	conn, err := grpc.Dial(u.Hostname(), opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return conn, nil
}

func newCredentialsFromCert(cert string) (credentials.TransportCredentials, error) {
	b, err := hex.DecodeString(cert)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{ServerName: "", RootCAs: cp}), nil
}
