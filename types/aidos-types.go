// Copyright (c) 2018 Aidos Developer

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Adapted from C++ code by 2014 Stellar Development Foundation and contributors

package types

//enum
type CryptoKeyType int32

const (
	KeyTypeED25519 CryptoKeyType = iota
	KeyTypePreAuthTx
	KeyTypeHashX
)

var cryptoKeyTypeMap = map[int32]string{
	0: "KeyTypeEd25519",
	1: "KeyTypePreAuthTx",
	2: "KeyTypeHashX",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CryptoKeyType
func (e CryptoKeyType) ValidEnum(v int32) bool {
	_, ok := cryptoKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e CryptoKeyType) String() string {
	name, _ := cryptoKeyTypeMap[int32(e)]
	return name
}

//enum
type PublicKeyType int32

const (
	PublicKeyTypeED25519 PublicKeyType = iota
)

var publicKeyTypeMap = map[int32]string{
	0: "PublicKeyTypeEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PublicKeyType
func (e PublicKeyType) ValidEnum(v int32) bool {
	_, ok := publicKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e PublicKeyType) String() string {
	name, _ := publicKeyTypeMap[int32(e)]
	return name
}

//enum
type SignerKeyType int32

const (
	SignerKeyTypeED25519   = KeyTypeED25519
	SignerKeyTypePreAuthTx = KeyTypePreAuthTx
	SignerKeyTypeHashX     = KeyTypeHashX
)

var signerKeyTypeMap = map[int32]string{
	0: "SignerKeyTypeEd25519",
	1: "SignerKeyTypePreAuthTx",
	2: "SignerKeyTypeHashX",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SignerKeyType
func (e SignerKeyType) ValidEnum(v int32) bool {
	_, ok := signerKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SignerKeyType) String() string {
	name, _ := signerKeyTypeMap[int32(e)]
	return name
}

type PublicKey struct {
	Type    PublicKeyType
	Ed25519 *Uint256
}

type SignerKey struct {
	Type      SignerKeyType
	Ed25519   *Uint256
	PreAuthTx *Uint256
	HashX     *Uint256
}

type NodeID = PublicKey

type Hash [32]uint8
type Uint256 [32]uint8

type Signature [64]uint8
type SignatureHint [4]uint8

//enum
type Rounding int

const (
	RoundDown Rounding = iota
	RoundUp
)
