package jwt

import "crypto/rsa"

type RS256Config interface {
	PrivateRSAKey() *rsa.PrivateKey
	PublicRSAKey() *rsa.PublicKey
}
