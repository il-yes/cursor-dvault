package vault_infrastructure_crypto

import (
	"crypto/rand"
	"crypto/sha512"

	"filippo.io/edwards25519"
	"github.com/stellar/go/strkey"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

type AsymmetricService struct{}

func (a *AsymmetricService) EncryptForRecipient(pub string, data []byte) ([]byte, error) {
	edPub, err := strkey.Decode(strkey.VersionByteAccountID, pub)
	if err != nil {
		return nil, err
	}

	curvePub := Ed25519PubToCurve(edPub)

	return box.SealAnonymous(nil, data, curvePub, rand.Reader)
}

// -------------------------------------
// ED25519 → CURVE25519
// -------------------------------------
// PUBLIC (used for encryption)
func Ed25519PubToCurve(pub []byte) *[32]byte {
	var out [32]byte
	p, err := new(edwards25519.Point).SetBytes(pub)
	Must(err)
	copy(out[:], p.BytesMontgomery())
	return &out
}

func CurvePrivFromStellarSeed(seed []byte) *[32]byte {
	// Ed25519 seed → SHA-512 → first 32 bytes → clamp
	h := sha512.Sum512(seed)

	var priv [32]byte
	copy(priv[:], h[:32])

	// Clamp (required by Curve25519)
	priv[0] &= 248
	priv[31] &= 127
	priv[31] |= 64

	return &priv
}

// Derive Curve25519 public key from private
func CurvePubFromPriv(priv *[32]byte) *[32]byte {
	var pub [32]byte
	curve25519.ScalarBaseMult(&pub, priv)
	return &pub
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}


func (c *AsymmetricService) GenerateSymmetricKey() []byte {
	symKey := make([]byte, 32)
	_, err := rand.Read(symKey)
	Must(err)

	return symKey
}

