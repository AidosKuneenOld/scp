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

// type Value struct {
// 	value []string
// }

type SCPBallot struct {
	counter uint32      // n
	value   interface{} // x
}

type SCPStatementType int32

const (
	SCPStPrepare SCPStatementType = iota
	SCPStConfirm
	SCPStExternalize
	SCPStNominate
)

var scpStatementTypeMap = map[int32]string{
	0: "ScpStPrepare",
	1: "ScpStConfirm",
	2: "ScpStExternalize",
	3: "StNominate",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ScpStatementType
func (e SCPStatementType) ValidEnum(v int32) bool {
	_, ok := scpStatementTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SCPStatementType) String() string {
	name, _ := scpStatementTypeMap[int32(e)]
	return name
}

type SCPNomination struct {
	quorumSetHash Hash          // D
	votes         []interface{} // X
	accepted      []interface{} // Y
}

type SCPStatement struct {
	NodeID    NodeID // v
	SlotIndex uint64 // i

	SCPStPrepare struct {
		quorumSetHash Hash       // D
		ballot        SCPBallot  // b
		prepared      *SCPBallot // p
		preparedPrime *SCPBallot // p'
		nC            uint32     // c.n
		nH            uint32     // h.n
	}
	SCPStConfirm struct {
		ballot        SCPBallot // b
		nPrepared     uint32    // p.n
		nCommit       uint32    // c.n
		nH            uint32    // h.n
		quorumSetHash Hash      // D
	}
	SCPStExternalize struct {
		commit              SCPBallot // c
		nH                  uint32    // h.n
		commitQuorumSetHash Hash      // D used before EXTERNALIZE
	}
	SCPStNominate struct {
		nominate SCPNomination
	}
}

type SCPEnvelope struct {
	Statement SCPStatement
	Signature Signature
}

// supports things like: A,B,C,(D,E,F),(G,H,(I,J,K,L))
// only allows 2 levels of nesting
type SCPQuorumSet struct {
	Threshold  uint32      `json:"t"`
	Validators []PublicKey `json:"v"`
	InnerSets  []SCPQuorumSet
}

type Json struct {
	T uint32        `json:"t"`
	V []interface{} `json:"v"` //PublicKey or SCPQuorumSet
}

const HexAbbrev uint8 = 3
const StrAbbrev uint8 = 6
