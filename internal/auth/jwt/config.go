package jwt

import "crypto/rsa"

type RSA256Config interface {
	PrivateRSAKey() *rsa.PrivateKey
	PublicRSAKey() *rsa.PublicKey
}
