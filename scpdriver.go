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

package scp

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/scp/types"
)

// Users of the SCP library should inherit from SCPDriver and implement the
// virtual methods which are called by the SCP implementation to
// abstract the transport layer used from the implementation of the SCP
// protocol.

//virtual
type SCPDriver struct{}

// `getValueString` is used for debugging
// default implementation is the hash of the value
func getValueString(value interface{}) string {
	valueHash := sha256.Sum256(types.Pack(value))
	//trim to HexAbbrev length
	return hex.EncodeToString(valueHash[:])[:types.HexAbbrev]
}

// `toShortString` converts to the common name of a key if found
func toShortString(pk types.PublicKey) string {
	//TODO strKey(key)
	//return string(pk)[:types.StringAbbrev]
	return "key"
}

// values used to switch hash function between priority and neighborhood checks
const hashN uint8 = 1
const hashP uint8 = 2
const hashK uint8 = 3

func hashHelper(slotIndex uint64, prev interface{}) uint64 {
	//res := binary.LittleEndian.Uint64(t[:8])
	return 0
}

const maxTimeoutSeconds uint32 = (30 * 60)

// `computeTimeout` computes a timeout given a round number
// it should be sufficiently large such that nodes in a
// quorum can exchange 4 messages
func computeTimeout(roundNumber uint32) int64 {
	// straight linear timeout
	// starting at 1 second and capping at MAX_TIMEOUT_SECONDS
	var timeoutInSeconds uint32
	if roundNumber > maxTimeoutSeconds {
		timeoutInSeconds = maxTimeoutSeconds
	} else {
		timeoutInSeconds = roundNumber
	}
	//sec to millisec
	return int64(timeoutInSeconds * 1000)
}

func (nD *SCPDriver) verifyEnvelope(envelope types.SCPEnvelope) bool {
	return true
}
