package lnd

import (
	"crypto/tls"
	"io/ioutil"
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
	macaroonBytes, err := ioutil.ReadFile(c.Macaroon)
	if err != nil {
		return nil, err
	}

	mac := &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(macaroonBytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	macConstraints := []macaroons.Constraint{
		macaroons.TimeoutConstraint(c.MacaroonTimeOut),
		macaroons.IPLockConstraint(c.MacaroonIP),
	}

	constrainedMac, err := macaroons.AddConstraints(mac, macConstraints...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cred credentials.TransportCredentials
	if c.Cert != "" {
		cred, err = credentials.NewClientTLSFromFile(c.Cert, "")
		if err != nil {
			return nil, err
		}
	} else {
		cred = credentials.NewTLS(&tls.Config{})
	}

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}

	macaroon, err := macaroons.NewMacaroonCredential(constrainedMac)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithPerRPCCredentials(macaroon),
		grpc.WithContextDialer(lncfg.ClientAddressDialer(u.Port())),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.MaxMsgRecvSize)),
	}

	conn, err := grpc.Dial(u.Hostname(), opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return conn, nil
}
