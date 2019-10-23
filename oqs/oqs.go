// Package oqs provides a GO wrapper for the C liboqs quantum-resistant library
package oqs // import "github.com/open-quantum-safe/liboqs-go/oqs"

/*
#cgo pkg-config: liboqs
#include <oqs/oqs.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

/**************** Misc functions ****************/

// MemCleanse sets to zero the content of a byte slice by invoking the liboqs
// OQS_MEM_cleanse() function. Use it to clean "hot" memory areas, such as
// secret keys etc.
func MemCleanse(v []byte) {
	C.OQS_MEM_cleanse(unsafe.Pointer(&v[0]),
		C.size_t(len(v)))
}

/**************** END Misc functions ****************/

/**************** KEMs ****************/

// List of enabled KEMs, populated by init().
var enabledKEMs []string

// List of supported KEMs, populated by init().
var supportedKEMs []string

// MaxNumberKEMs returns the maximum number of supported KEMs.
func MaxNumberKEMs() int {
	return int(C.OQS_KEM_alg_count())
}

// IsKEMEnabled returns true if a KEM is enabled, and false otherwise.
func IsKEMEnabled(algName string) bool {
	result := C.OQS_KEM_alg_is_enabled(C.CString(algName))
	return result != 0
}

// IsKEMSupported returns true if a KEM is supported, and false otherwise.
func IsKEMSupported(algName string) bool {
	for i := range supportedKEMs {
		if supportedKEMs[i] == algName {
			return true
		}
	}
	return false
}

// GetKEMName returns the KEM name from its corresponding numerical id.
func GetKEMName(algID int) string {
	if algID >= MaxNumberKEMs() {
		panic(errors.New("algorithm ID out of range"))
	}
	return C.GoString(C.OQS_KEM_alg_identifier(C.size_t(algID)))
}

// GetSupportedKEMs returns the list of supported KEMs.
func GetSupportedKEMs() []string {
	return supportedKEMs
}

// GetEnabledKEMs returns the list of enabled KEMs.
func GetEnabledKEMs() []string {
	return enabledKEMs
}

// Initializes the lists of enabledKEMs and supportedKEMs.
func init() {
	for i := 0; i < MaxNumberKEMs(); i++ {
		KEMName := GetKEMName(i)
		supportedKEMs = append(supportedKEMs, KEMName)
		if IsKEMEnabled(KEMName) {
			enabledKEMs = append(enabledKEMs, KEMName)
		}
	}
}

/**************** END KEMs ****************/

/**************** KeyEncapsulation ****************/

// KeyEncapsulationDetails defines the KEM algorithm details.
type KeyEncapsulationDetails struct {
	ClaimedNISTLevel   int
	IsINDCCA           bool
	LengthCiphertext   int
	LengthPublicKey    int
	LengthSecretKey    int
	LengthSharedSecret int
	Name               string
	Version            string
}

// KeyEncapsulation defines the KEM main data structure.
type KeyEncapsulation struct {
	kem        *C.OQS_KEM
	algName    string
	secretKey  []byte
	algDetails KeyEncapsulationDetails
}

// Init initializes the KEM data structure with an algorithm name and a secret
// key. If the secret key is null, then the user must invoke the
//// KeyEncapsulation.GenerateKeyPair method to generate the pair of
// secret key/public key.
func (kem *KeyEncapsulation) Init(algName string, secretKey []byte) {
	if !IsKEMEnabled(algName) {
		// perhaps it's supported
		if IsKEMSupported(algName) {
			panic(errors.New(`"` + algName + `" is not enabled by OQS`))
		} else {
			panic(errors.New(`"` + algName + `" is not supported by OQS`))
		}
	}
	kem.kem = C.OQS_KEM_new(C.CString(algName))
	kem.algName = algName
	kem.secretKey = secretKey
	kem.algDetails.Name = C.GoString(kem.kem.method_name)
	kem.algDetails.Version = C.GoString(kem.kem.alg_version)
	kem.algDetails.ClaimedNISTLevel = int(kem.kem.claimed_nist_level)
	kem.algDetails.IsINDCCA = bool(kem.kem.ind_cca)
	kem.algDetails.LengthPublicKey = int(kem.kem.length_public_key)
	kem.algDetails.LengthSecretKey = int(kem.kem.length_secret_key)
	kem.algDetails.LengthCiphertext = int(kem.kem.length_ciphertext)
	kem.algDetails.LengthSharedSecret = int(kem.kem.length_shared_secret)
}

// GetDetails returns the KEM algorithm details.
func (kem *KeyEncapsulation) GetDetails() KeyEncapsulationDetails {
	return kem.algDetails
}

// GenerateKeypair generates a pair of secret key/public key and returns the
// public key. The secret key is stored inside the kem receiver. The secret key
// is not directly accessible, unless one exports it with
// KeyEncapsulation.ExportSecretKey method.
func (kem *KeyEncapsulation) GenerateKeypair() []byte {
	publicKey := make([]byte, kem.algDetails.LengthPublicKey)
	kem.secretKey = make([]byte, kem.algDetails.LengthSecretKey)

	rv := C.OQS_KEM_keypair(kem.kem, (*C.uint8_t)(&publicKey[0]),
		(*C.uint8_t)(&kem.secretKey[0]))
	if rv != C.OQS_SUCCESS {
		panic(errors.New("can not generate keypair"))
	}

	return publicKey
}

// ExportSecretKey exports the corresponding secret key from the kem receiver.
func (kem *KeyEncapsulation) ExportSecretKey() []byte {
	return kem.secretKey
}

// EncapSecret encapsulates a secret using a public key and returns the
// corresponding ciphertext and shared secret.
func (kem *KeyEncapsulation) EncapSecret(publicKey []byte) (ciphertext,
	sharedSecret []byte) {
	if len(publicKey) != kem.algDetails.LengthPublicKey {
		panic(errors.New("incorrect public key length"))
	}

	ciphertext = make([]byte, kem.algDetails.LengthCiphertext)
	sharedSecret = make([]byte, kem.algDetails.LengthSharedSecret)

	rv := C.OQS_KEM_encaps(kem.kem, (*C.uint8_t)(&ciphertext[0]),
		(*C.uint8_t)(&sharedSecret[0]), (*C.uint8_t)(&publicKey[0]))

	if rv != C.OQS_SUCCESS {
		panic(errors.New("can not encapsulate secret"))
	}

	return ciphertext, sharedSecret
}

// DecapSecret decapsulates a ciphertexts and returns the corresponding shared
// secret.
func (kem *KeyEncapsulation) DecapSecret(ciphertext []byte) []byte {
	if len(ciphertext) != kem.algDetails.LengthCiphertext {
		panic(errors.New("incorrect ciphertext length"))
	}

	if len(kem.secretKey) != kem.algDetails.LengthSecretKey {
		panic(errors.New("incorrect secret key length, make sure you " +
			"specify one in Init() or run GenerateKeypair()"))

	}

	sharedSecret := make([]byte, kem.algDetails.LengthSharedSecret)
	rv := C.OQS_KEM_decaps(kem.kem, (*C.uint8_t)(&sharedSecret[0]),
		(*C.uchar)(&ciphertext[0]), (*C.uint8_t)(&kem.secretKey[0]))

	if rv != C.OQS_SUCCESS {
		panic(errors.New("can not decapsulate secret"))
	}

	return sharedSecret
}

// Clean zeroes-in the stored secret key and resets the kem receiver. One can
// reuse the KEM by re-initializing it with the KeyEncapsulation.Init method.
func (kem *KeyEncapsulation) Clean() {
	if len(kem.secretKey) > 0 {
		MemCleanse(kem.secretKey)
	}
	C.OQS_KEM_free(kem.kem)
	*kem = KeyEncapsulation{}
}

/**************** END KeyEncapsulation ****************/

/**************** SIGs ****************/

// List of enabled signatures, populated by init().
var enabledSIGs []string

// List of supported signatures, populated by init().
var supportedSIGs []string

// MaxNumberSIGs returns the maximum number of supported signatures.
func MaxNumberSIGs() int {
	return int(C.OQS_SIG_alg_count())
}

// IsSIGEnabled returns true if a signature is enabled, and false otherwise.
func IsSIGEnabled(algName string) bool {
	result := C.OQS_SIG_alg_is_enabled(C.CString(algName))
	return result != 0
}

// IsSIGSupported returns true if a signature is supported, and false otherwise.
func IsSIGSupported(algName string) bool {
	for i := range supportedSIGs {
		if supportedSIGs[i] == algName {
			return true
		}
	}
	return false
}

// GetSIGName returns the signature name from its corresponding numerical id.
func GetSIGName(algID int) string {
	if algID >= MaxNumberSIGs() {
		panic(errors.New("algorithm ID out of range"))
	}
	return C.GoString(C.OQS_SIG_alg_identifier(C.size_t(algID)))
}

// GetSupportedSIGs returns the list of supported signatures.
func GetSupportedSIGs() []string {
	return supportedSIGs
}

// GetEnabledSIGs returns the list of enabled signatures.
func GetEnabledSIGs() []string {
	return enabledSIGs
}

// Initializes the lists of enabledSIGs and supportedSIGs.
func init() {
	for i := 0; i < MaxNumberSIGs(); i++ {
		SIGName := GetSIGName(i)
		supportedSIGs = append(supportedSIGs, SIGName)
		if IsSIGEnabled(SIGName) {
			enabledSIGs = append(enabledSIGs, SIGName)
		}
	}
}

/**************** END SIGs ****************/

/**************** Signature ****************/

// SignatureDetails defines the signature algorithm details.
type SignatureDetails struct {
	ClaimedNISTLevel   int
	IsEUFCMA           bool
	LengthPublicKey    int
	LengthSecretKey    int
	MaxLengthSignature int
	Name               string
	Version            string
}

// Signature defines the signature main data structure.
type Signature struct {
	sig        *C.OQS_SIG
	algName    string
	secretKey  []byte
	algDetails SignatureDetails
}

// Init initializes the signature data structure with an algorithm name and a
// secret key. If the secret key is null, then the user must invoke the
// Signature.GenerateKeyPair method to generate the pair of secret key/public
// key.
func (sig *Signature) Init(algName string, secretKey []byte) {
	if !IsSIGEnabled(algName) {
		// perhaps it's supported
		if IsSIGSupported(algName) {
			panic(errors.New(`"` + algName + `" is not enabled by OQS`))
		} else {
			panic(errors.New(`"` + algName + `" is not supported by OQS`))
		}
	}
	sig.sig = C.OQS_SIG_new(C.CString(algName))
	sig.algName = algName
	sig.secretKey = secretKey
	sig.algDetails.Name = C.GoString(sig.sig.method_name)
	sig.algDetails.Version = C.GoString(sig.sig.alg_version)
	sig.algDetails.ClaimedNISTLevel = int(sig.sig.claimed_nist_level)
	sig.algDetails.IsEUFCMA = bool(sig.sig.euf_cma)
	sig.algDetails.LengthPublicKey = int(sig.sig.length_public_key)
	sig.algDetails.LengthSecretKey = int(sig.sig.length_secret_key)
	sig.algDetails.MaxLengthSignature = int(sig.sig.length_signature)
}

// GetDetails returns the signature algorithm details.
func (sig *Signature) GetDetails() SignatureDetails {
	return sig.algDetails
}

// GenerateKeypair generates a pair of secret key/public key and returns the
// public key. The secret key is stored inside the sig receiver. The secret key
// is not directly accessible, unless one exports it with
// Signature.ExportSecretKey method.
func (sig *Signature) GenerateKeypair() []byte {
	publicKey := make([]byte, sig.algDetails.LengthPublicKey)
	sig.secretKey = make([]byte, sig.algDetails.LengthSecretKey)

	rv := C.OQS_SIG_keypair(sig.sig, (*C.uint8_t)(&publicKey[0]),
		(*C.uint8_t)(&sig.secretKey[0]))
	if rv != C.OQS_SUCCESS {
		panic(errors.New("can not generate keypair"))
	}

	return publicKey
}

// ExportSecretKey exports the corresponding secret key from the sig receiver.
func (sig *Signature) ExportSecretKey() []byte {
	return sig.secretKey
}

// Sign signs a message and returns the corresponding signature.
func (sig *Signature) Sign(message []byte) []byte {
	if len(sig.secretKey) != sig.algDetails.LengthSecretKey {
		panic(errors.New("incorrect secret key length, make sure you " +
			"specify one in Init() or run GenerateKeypair()"))
	}

	signature := make([]byte, sig.algDetails.MaxLengthSignature)
	var lenSig int
	rv := C.OQS_SIG_sign(sig.sig, (*C.uint8_t)(&signature[0]),
		(*C.size_t)(unsafe.Pointer(&lenSig)), (*C.uint8_t)(&message[0]),
		C.size_t(len(message)), (*C.uint8_t)(&sig.secretKey[0]))

	if rv != C.OQS_SUCCESS {
		panic(errors.New("can not sign message"))
	}

	return signature[:lenSig]
}

// Verify verifies the validity of a signed message, returning true if the
// signature is valid, and false otherwise.
func (sig *Signature) Verify(message []byte, signature []byte,
	publicKey []byte) bool {
	if len(publicKey) != sig.algDetails.LengthPublicKey {
		panic(errors.New("incorrect public key length"))
	}

	if len(signature) > sig.algDetails.MaxLengthSignature {
		panic(errors.New("incorrect signature size"))
	}

	rv := C.OQS_SIG_verify(sig.sig, (*C.uint8_t)(&message[0]),
		C.size_t(len(message)), (*C.uint8_t)(&signature[0]),
		C.size_t(len(signature)), (*C.uint8_t)(&publicKey[0]))

	if rv != C.OQS_SUCCESS {
		return false
	}

	return true
}

// Clean zeroes-in the stored secret key and resets the sig receiver. One can
// reuse the signature by re-initializing it with the Signature.Init method.
func (sig *Signature) Clean() {
	if len(sig.secretKey) > 0 {
		MemCleanse(sig.secretKey)
	}
	C.OQS_SIG_free(sig.sig)
	*sig = Signature{}
}

/**************** END Signature ****************/
