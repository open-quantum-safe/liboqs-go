// liboqs GO wrapper
package oqs

/*
#cgo CFLAGS: -I/Users/vlad/liboqs/include
#cgo LDFLAGS: -L/usr/local/lib -loqs
#include <stdio.h>
#include <stdlib.h>
#include <oqs/oqs.h>
*/
import "C"

import (
    "bytes"
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
    return C.GoString(C.OQS_KEM_alg_identifier(C.ulong(algID)))
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
    return C.GoString(C.OQS_SIG_alg_identifier(C.ulong(algID)))
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
type Signature struct {
    algName    string
    secretKey  bytes.Buffer
    algDetails struct {
        ClaimedNISTLevel   int
        IsEUFCMA           bool
        LengthPublicKey    int
        LengthSecretKey    int
        MaxLengthSignature int
        Name               string
        Version            string
    }
}

/**************** END Signature ****************/

/**************** KeyEncapsulation ****************/
type KeyEncapsulation struct {
    algName    string
    secretKey  bytes.Buffer
    algDetails struct {
        ClaimedNISTLevel int
        IsINDCCA         bool
        LengthCiphertext int
        LengthPublicKey  int
        LengthSecretKey  int
        MaxSharedSecret  int
        Name             string
        Version          string
    }
}

/**************** END KeyEncapsulation ****************/
