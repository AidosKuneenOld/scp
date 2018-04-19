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
type CryptoKeyType int

const (
	KeyTypeED25519 CryptoKeyType = iota
	KeyTypePreAuthTx
	KeyTypeHashX
)

//enum
type PublicKeyType int

const (
	PublicKeyTypeED25519 PublicKeyType = iota
)

//enum
type SignerKeyType int

const (
	SignerKeyTypeED25519   = KeyTypeED25519
	SignerKeyTypePreAuthTx = KeyTypePreAuthTx
	SignerKeyTypeHashX     = KeyTypeHashX
)

type PublicKey struct {
	PublicKeyTypeED25519 struct {
		ed25519 uint256
	}
}

type SignerKey struct {
	SignerKeyTypeED25519 struct {
		ed25519 uint256
	}
	SignerKeyTypePreAuthTx struct {
		/* Hash of Transaction structure */
		preAuthTx uint256
	}
	SignerKeyTypeHashX struct {
		/* Hash of Transaction structure */
		hashX uint256
	}
}

type NodeID = PublicKey

type Hash [32]uint8
type uint256 [32]uint8

type Signature [64]uint8
type SignatureHint [4]uint8

//enum
type Rounding int

const (
	RoundDown Rounding = iota
	RoundUp
)
