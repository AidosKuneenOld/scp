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
	"log"

	"github.com/scp/types"
)

type LocalNode struct {
	mNodeID      types.NodeID
	mIsValidator bool
	mQSet        types.SCPQuorumSet
	mQSetHash    types.Hash

	// alternative qset used during externalize {{mNodeID}}
	gSingleQSetHash types.Hash          // hash of the singleton qset
	mSingleQSet     *types.SCPQuorumSet // {{mNodeID}}

	mSCP *SCP
}

func (nl *LocalNode) NLocalNode(nodeID types.NodeID, isValidator bool,
	qSet types.SCPQuorumSet, scp *SCP) {
	nl.mNodeID = nodeID
	nl.mIsValidator = isValidator
	nl.mQSet = qSet
	nl.mSCP = scp

	NormalizeQSet(&nl.mQSet)
	nl.mQSetHash = sha256.Sum256(types.Pack(nl.mQSet))

	log.Printf("SCP: LocalNode@%.6s qSet: %.3X", nl.mNodeID, nl.mQSetHash)

	nl.mSingleQSet = buildSingletonQSet(nl.mNodeID)
	nl.gSingleQSetHash = sha256.Sum256(types.Pack(nl.mSingleQSet))
}

// returns a quorum set {{ nodeID }}
func buildSingletonQSet(nodeID types.NodeID) *types.SCPQuorumSet {
	return &types.SCPQuorumSet{
		Threshold:  1,
		Validators: []types.NodeID{nodeID}}
}

func (nl *LocalNode) UpdateQuorumSet(qSet types.SCPQuorumSet) {
	nl.mQSetHash = sha256.Sum256(types.Pack(qSet))
	nl.mQSet = qSet
}

func (nl *LocalNode) QuorumSet() types.SCPQuorumSet {
	return nl.mQSet
}
