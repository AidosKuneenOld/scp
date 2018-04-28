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
	"log"

	"github.com/scp/types"
)

type SCP struct {
	mDriver    SCPDriver
	mLocalNode *LocalNode
}

func (ns *SCP) nSCP(driver SCPDriver, nodeID types.NodeID, isValidator bool,
	qSetLocal types.SCPQuorumSet) {
	ns.mDriver = driver
	ns.mLocalNode.NLocalNode(nodeID, isValidator, qSetLocal, ns)
}

// this is the main entry point of the SCP library
// it processes the envelope, updates the internal state and
// invokes the appropriate methods
func (ns *SCP) receiveEnvelope(envelope types.SCPEnvelope) EnvelopeState {
	// If the envelope is not correctly signed, we ignore it.
	if !ns.mDriver.verifyEnvelope(envelope) {
		log.Println("DEBUG SCP: receiveEnvelope invalid")
		return Invalid
	}
	//slotIndex := envelope.Statement.SlotIndex
	//TODO create function for return
	return Valid
}

//enum
type EnvelopeState int32

const (
	Invalid EnvelopeState = iota
	Valid
)

var envelopeStateMap = map[int32]string{
	0: "Invalid",
	1: "Valid",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EnvelopeState
func (e EnvelopeState) ValidEnum(v int32) bool {
	_, ok := envelopeStateMap[v]
	return ok
}

// String returns the name of `e`
func (e EnvelopeState) String() string {
	name, _ := envelopeStateMap[int32(e)]
	return name
}

//enum
type TriBool int32

const (
	TBTrue TriBool = iota
	TBFalse
	TBMaybe
)

var triBoolMap = map[int32]string{
	0: "TBTrue",
	1: "TBFalse",
	2: "TBMaybe",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ScpStatementType
func (e TriBool) ValidEnum(v int32) bool {
	_, ok := triBoolMap[v]
	return ok
}

// String returns the name of `e`
func (e TriBool) String() string {
	name, _ := triBoolMap[int32(e)]
	return name
}
