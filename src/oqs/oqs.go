// Package oqs provides a GO wrapper for the C liboqs quantum-resistant library
package oqs

/*
#cgo pkg-config: liboqs
#include <oqs/oqs.h>
*/
import "C"

import (
    "unsafe"
)

/**************** KEMs ****************/

var enabledKEMs []string
var supportedKEMs []string

func MaxNumberKEMs() int {
    return int(C.OQS_KEM_alg_count())
}

func IsKEMEnabled(algName string) bool {
    result := C.OQS_KEM_alg_is_enabled(C.CString(algName))
    return result != 0
}

func IsKEMSupported(algName string) bool {
    for i := range supportedKEMs {
        if supportedKEMs[i] == algName {
            return true
        }
    }
    return false
}

func GetKEMName(algID int) string {
    if algID >= MaxNumberKEMs() {
        panic("Algorithm ID out of range")
    }
    return C.GoString(C.OQS_KEM_alg_identifier(C.size_t(algID)))
}

func GetSupportedKEMs() []string {
    return supportedKEMs
}

func GetEnabledKEMs() []string {
    return enabledKEMs
}

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

type keyEncapsulationDetails struct {
    ClaimedNISTLevel   int
    IsINDCCA           bool
    LengthCiphertext   int
    LengthPublicKey    int
    LengthSecretKey    int
    LengthSharedSecret int
    Name               string
    Version            string
}

type KeyEncapsulation struct {
    kem        *C.OQS_KEM
    algName    string
    secretKey  []byte
    algDetails keyEncapsulationDetails
}

func (kem *KeyEncapsulation) Init(algName string, secretKey []byte) {
    if !IsKEMEnabled(algName) {
        // perhaps it's supported
        if IsKEMSupported(algName) {
            panic(`"` + algName + `" is not enabled by OQS`)
        } else {
            panic(`"` + algName + `" is not supported by OQS`)
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

func (kem *KeyEncapsulation) GetDetails() keyEncapsulationDetails {
    return kem.algDetails
}

func (kem *KeyEncapsulation) GenerateKeypair() []byte {
    publicKey := make([]byte, kem.algDetails.LengthPublicKey)
    kem.secretKey = make([]byte, kem.algDetails.LengthSecretKey)

    rv := C.OQS_KEM_keypair(kem.kem, (*C.uint8_t)(&publicKey[0]),
        (*C.uint8_t)(&kem.secretKey[0]))
    if rv != C.OQS_SUCCESS {
        panic("Can not generate keypair")
    }

    return publicKey
}

func (kem *KeyEncapsulation) ExportSecretKey() []byte {
    return kem.secretKey
}

func (kem *KeyEncapsulation) EncapSecret(publicKey []byte) (ciphertext, sharedSecret []byte) {
    if len(publicKey) != kem.algDetails.LengthPublicKey {
        panic("Incorrect public key length")
    }

    ciphertext = make([]byte, kem.algDetails.LengthCiphertext)
    sharedSecret = make([]byte, kem.algDetails.LengthSharedSecret)

    rv := C.OQS_KEM_encaps(kem.kem, (*C.uint8_t)(&ciphertext[0]),
        (*C.uint8_t)(&sharedSecret[0]), (*C.uint8_t)(&publicKey[0]))

    if rv != C.OQS_SUCCESS {
        panic("Can not encapsulate secret")
    }

    return ciphertext, sharedSecret
}

func (kem *KeyEncapsulation) DecapSecret(ciphertext []byte) []byte {
    if len(ciphertext) != kem.algDetails.LengthCiphertext {
        panic("Incorrect ciphertext length")
    }

    if len(kem.secretKey) != kem.algDetails.LengthSecretKey {
        panic("Incorrect secret key length, make sure you specify one in " +
            "Init() or run GenerateKeypair()")

    }

    sharedSecret := make([]byte, kem.algDetails.LengthSharedSecret)
    rv := C.OQS_KEM_decaps(kem.kem, (*C.uint8_t)(&sharedSecret[0]),
        (*C.uchar)(&ciphertext[0]), (*C.uint8_t)(&kem.secretKey[0]))

    if rv != C.OQS_SUCCESS {
        panic("Can not decapsulate secret")
    }

    return sharedSecret
}

func (kem *KeyEncapsulation) Close() {
    if len(kem.secretKey) > 0 {
        C.OQS_MEM_cleanse(unsafe.Pointer(&kem.secretKey[0]),
            C.size_t(len(kem.secretKey)))
    }
    C.OQS_KEM_free(kem.kem)
}

/**************** END KeyEncapsulation ****************/

/**************** SIGs ****************/

var enabledSIGs []string
var supportedSIGs []string

func MaxNumberSIGs() int {
    return int(C.OQS_SIG_alg_count())
}

func IsSIGEnabled(algName string) bool {
    result := C.OQS_SIG_alg_is_enabled(C.CString(algName))
    return result != 0
}

func IsSIGSupported(algName string) bool {
    for i := range supportedSIGs {
        if supportedSIGs[i] == algName {
            return true
        }
    }
    return false
}

func GetSIGName(algID int) string {
    if algID >= MaxNumberSIGs() {
        panic("Algorithm ID out of range")
    }
    return C.GoString(C.OQS_SIG_alg_identifier(C.size_t(algID)))
}

func GetSupportedSIGs() []string {
    return supportedSIGs
}

func GetEnabledSIGs() []string {
    return enabledSIGs
}

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

type signatureDetails struct {
    ClaimedNISTLevel   int
    IsEUFCMA           bool
    LengthPublicKey    int
    LengthSecretKey    int
    MaxLengthSignature int
    Name               string
    Version            string
}

type Signature struct {
    sig        *C.OQS_SIG
    algName    string
    secretKey  []byte
    algDetails signatureDetails
}

func (sig *Signature) Init(algName string, secretKey []byte) {
    if !IsSIGEnabled(algName) {
        // perhaps it's supported
        if IsSIGSupported(algName) {
            panic(`"` + algName + `" is not enabled by OQS`)
        } else {
            panic(`"` + algName + `" is not supported by OQS`)
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

func (sig *Signature) GetDetails() signatureDetails {
    return sig.algDetails
}

func (sig *Signature) GenerateKeypair() []byte {
    publicKey := make([]byte, sig.algDetails.LengthPublicKey)
    sig.secretKey = make([]byte, sig.algDetails.LengthSecretKey)

    rv := C.OQS_SIG_keypair(sig.sig, (*C.uint8_t)(&publicKey[0]),
        (*C.uint8_t)(&sig.secretKey[0]))
    if rv != C.OQS_SUCCESS {
        panic("Can not generate keypair")
    }

    return publicKey
}

func (sig *Signature) ExportSecretKey() []byte {
    return sig.secretKey
}

func (sig *Signature) Sign(message []byte) []byte {
    if len(sig.secretKey) != sig.algDetails.LengthSecretKey {
        panic("Incorrect secret key length, make sure you specify one in " +
            "Init() or run GenerateKeypair()")
    }

    signature := make([]byte, sig.algDetails.MaxLengthSignature)
    var lenSig int
    rv := C.OQS_SIG_sign(sig.sig, (*C.uint8_t)(&signature[0]),
        (*C.size_t)(unsafe.Pointer(&lenSig)), (*C.uint8_t)(&message[0]),
        C.size_t(len(message)), (*C.uint8_t)(&sig.secretKey[0]))

    if rv != C.OQS_SUCCESS {
        panic("Can not sign message")
    }

    return signature[:lenSig]
}

func (sig *Signature) Verify(message []byte, signature []byte,
    publicKey []byte) bool {
    if len(publicKey) != sig.algDetails.LengthPublicKey {
        panic("Incorrect public key length")
    }

    if len(signature) > sig.algDetails.MaxLengthSignature {
        panic("Incorrect signature size")
    }

    rv := C.OQS_SIG_verify(sig.sig, (*C.uint8_t)(&message[0]),
        C.size_t(len(message)), (*C.uint8_t)(&signature[0]),
        C.size_t(len(signature)), (*C.uint8_t)(&publicKey[0]))

    if rv != C.OQS_SUCCESS {
        return false
    }

    return true
}

func (sig *Signature) Close() {
    if len(sig.secretKey) > 0 {
        C.OQS_MEM_cleanse(unsafe.Pointer(&sig.secretKey[0]),
            C.size_t(len(sig.secretKey)))
    }
    C.OQS_SIG_free(sig.sig)
}

/**************** END Signature ****************/
