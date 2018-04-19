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
	"math"
	"math/big"

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

func (nl *LocalNode) QuorumSetHash() types.Hash {
	return nl.mQSetHash
}

// returns the quorum set {{X}}
func SingletonQSet(nodeID types.NodeID) *types.SCPQuorumSet {
	return buildSingletonQSet(nodeID)
}

// called recursively
func forAllNodesInternal(qset types.SCPQuorumSet, proc func(nodeID types.NodeID)) {

	for _, n := range qset.Validators {
		proc(n)
	}
	for _, q := range qset.InnerSets {
		forAllNodesInternal(q, proc)
	}
}

// runs proc over all nodes contained in qset
func ForAllNodes(qset types.SCPQuorumSet, proc func(nodeID types.NodeID)) {
	var done map[types.NodeID]struct{}
	forAllNodesInternal(qset, func(n types.NodeID) {
		if _, exist := done[n]; !exist {
			// n was not present
			done[n] = struct{}{}
			proc(n)
		}
	})
}

// returns the weight of the node within the qset
// normalized between 0-UINT64_MAX
// *if a validator is repeated multiple times its weight is only the
// weight of the first occurrence
func GetNodeWeight(nodeID types.NodeID, qset types.SCPQuorumSet) uint64 {
	n := uint64(qset.Threshold)
	d := uint64(len(qset.InnerSets) + len(qset.Validators))

	for _, qsetNode := range qset.Validators {
		if qsetNode == nodeID {
			return bigDivideReturn(math.MaxUint64, n, d, types.RoundDown)
		}
	}
	for _, q := range qset.InnerSets {
		leafW := GetNodeWeight(nodeID, q)
		if leafW != 0 {
			return bigDivideReturn(leafW, n, d, types.RoundDown)
		}
	}

	return 0
}

func isQuorumSliceInternal(qset types.SCPQuorumSet, nodeSet []types.NodeID) bool {
	thresholdLeft := qset.Threshold
	for _, validator := range qset.Validators {
		index := SliceIndex(len(qset.Validators), func(i int) bool { return validator == nodeSet[i] })
		if index != -1 {
			//found
			thresholdLeft--
			if thresholdLeft <= 0 {
				return true
			}
		}
	}
	for _, inner := range qset.InnerSets {
		if isQuorumSliceInternal(inner, nodeSet) {
			thresholdLeft--
			if thresholdLeft <= 0 {
				return true
			}
		}
	}
	return false
}

// Tests this node against nodeSet for the specified qSethash.
func IsQuorumSlice(qSet types.SCPQuorumSet, nodeSet []types.NodeID) bool {

	log.Printf("SCP: LocalNode IsQuorumSlice len(nodeSet): %v", len(nodeSet))

	return isQuorumSliceInternal(qSet, nodeSet)
}

//Return slice index for an element
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

//util
func bigDivideReturn(A uint64, B uint64, C uint64, rounding types.Rounding) uint64 {

	a := new(big.Int).SetUint64(A)
	b := new(big.Int).SetUint64(B)
	c := new(big.Int).SetUint64(C)
	x := big.NewInt(0)
	if rounding == types.RoundDown {
		x.Mul(a, b)
		x.Div(x, c)
	} else {
		x.Mul(a, b)
		x.Add(x, new(big.Int).Sub(c, big.NewInt(1))).Div(x, c)
	}

	return x.Uint64()
}

//util
func bigDivide(result uint64, A uint64, B uint64, C uint64, rounding types.Rounding) bool {

	a := new(big.Int).SetUint64(A)
	b := new(big.Int).SetUint64(B)
	c := new(big.Int).SetUint64(C)
	x := big.NewInt(0)
	if rounding == types.RoundDown {
		x.Mul(a, b)
		x.Div(x, c)
	} else {
		x.Mul(a, b)
		x.Add(x, new(big.Int).Sub(c, big.NewInt(1))).Div(x, c)
	}

	result = x.Uint64()

	if x.Cmp(new(big.Int).SetUint64(math.MaxUint64)) <= 0 {
		// if x fits in uint64
		return true
	}
	return false
}
