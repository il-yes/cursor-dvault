package stellar_recovery_domain

import "errors"

type StellarSecretKey string

func NewStellarSecretKey(k string) (StellarSecretKey, error) {
	if len(k) < 20 {
		return "", errors.New("invalid stellar secret key")
	}
	return StellarSecretKey(k), nil
}

func (s StellarSecretKey) String() string { return string(s) }
