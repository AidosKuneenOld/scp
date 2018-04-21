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
	"sort"

	"github.com/scp/types"
)

type LocalNode struct {
	mNodeID      types.NodeID
	mIsValidator bool
	mQSet        types.SCPQuorumSet
	mQSetHash    types.Hash

	// alternative qSet used during externalize {{mNodeID}}
	gSingleQSetHash types.Hash          // hash of the singleton qSet
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
func forAllNodesInternal(qSet types.SCPQuorumSet, proc func(nodeID types.NodeID)) {

	for _, n := range qSet.Validators {
		proc(n)
	}
	for _, q := range qSet.InnerSets {
		forAllNodesInternal(q, proc)
	}
}

// runs proc over all nodes contained in qSet
func ForAllNodes(qSet types.SCPQuorumSet, proc func(nodeID types.NodeID)) {
	done := make(map[types.NodeID]struct{})
	forAllNodesInternal(qSet, func(n types.NodeID) {
		if _, exist := done[n]; !exist {
			// n was not present
			done[n] = struct{}{}
			proc(n)
		}
	})
}

// returns the weight of the node within the qSet
// normalized between 0-UINT64_MAX
// *if a validator is repeated multiple times its weight is only the
// weight of the first occurrence
func GetNodeWeight(nodeID types.NodeID, qSet types.SCPQuorumSet) uint64 {
	n := uint64(qSet.Threshold)
	d := uint64(len(qSet.InnerSets) + len(qSet.Validators))

	for _, qSetNode := range qSet.Validators {
		if qSetNode == nodeID {
			return bigDivideReturn(math.MaxUint64, n, d, types.RoundDown)
		}
	}
	for _, q := range qSet.InnerSets {
		leafW := GetNodeWeight(nodeID, q)
		if leafW != 0 {
			return bigDivideReturn(leafW, n, d, types.RoundDown)
		}
	}

	return 0
}

func isQuorumSliceInternal(qSet types.SCPQuorumSet, nodeSet []types.NodeID) bool {
	thresholdLeft := qSet.Threshold
	for _, validator := range qSet.Validators {
		index := SliceIndex(len(qSet.Validators), func(i int) bool { return validator == nodeSet[i] })
		if index != -1 {
			//found
			thresholdLeft--
			if thresholdLeft <= 0 {
				return true
			}
		}
	}
	for _, inner := range qSet.InnerSets {
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

// called recursively
func isVBlockingInternal(qSet types.SCPQuorumSet, nodeSet []types.NodeID) bool {
	// There is no v-blocking set for {\empty}
	if qSet.Threshold == 0 {
		return false
	}

	leftTillBlock := (1 + len(qSet.Validators) + len(qSet.InnerSets)) - int(qSet.Threshold)
	for _, validator := range qSet.Validators {
		index := SliceIndex(len(qSet.Validators), func(i int) bool { return validator == nodeSet[i] })
		if index != -1 {
			//found
			leftTillBlock--
			if leftTillBlock <= 0 {
				return true
			}
		}
	}
	for _, inner := range qSet.InnerSets {
		if isVBlockingInternal(inner, nodeSet) {
			leftTillBlock--
			if leftTillBlock <= 0 {
				return true
			}
		}
	}
	return false
}

// Tests this node against a map of nodeID -> T for the specified qSetHash.
func IsVBlocking(qSet types.SCPQuorumSet, nodeSet []types.NodeID) bool {

	log.Printf("SCP: LocalNode IsVBlocking len(nodeSet): %v", len(nodeSet))

	return isVBlockingInternal(qSet, nodeSet)
}

// `isVBlocking` tests if the filtered nodes V are a v-blocking set for
// this node.
func IsVBlockingF(qSet types.SCPQuorumSet, map1 map[types.NodeID]types.SCPEnvelope,
	filter func(types.SCPStatement) bool) bool {

	var pNodes []types.NodeID

	for k, v := range map1 {
		if filter(v.Statement) {
			pNodes = append(pNodes, k)
		}
	}

	return IsVBlocking(qSet, pNodes)
}

// `isQuorum` tests if the filtered nodes V form a quorum
// (meaning for each v \in V there is q \in Q(v)
// included in V and we have quorum on V for qSetHash). `qfun` extracts the
// SCPQuorumSetPtr from the SCPStatement for its associated node in map
// (required for transitivity)
func IsQuorum(qSet types.SCPQuorumSet, map1 map[types.NodeID]types.SCPEnvelope,
	qFun func(types.SCPStatement) *types.SCPQuorumSet,
	filter func(types.SCPStatement) bool) bool {

	var pNodes []types.NodeID

	for k, v := range map1 {
		if filter(v.Statement) {
			pNodes = append(pNodes, k)
		}
	}

	count := 0
	//exec at least once (do..while)
	for while := true; while; while = (count != len(pNodes)) {
		count = len(pNodes)
		var fNodes []types.NodeID

		quroumFilter := func(nodeID types.NodeID) bool {

			qSetPtr := qFun(map1[nodeID].Statement)
			if qSetPtr != nil {
				return IsQuorumSlice(*qSetPtr, pNodes)
			}
			return false
		}

		for _, p := range pNodes {
			if quroumFilter(p) {
				fNodes = append(fNodes, p)
			}
		}
		pNodes = fNodes
	}

	return IsQuorumSlice(qSet, pNodes)
}

// computes the distance to the set of v-blocking sets given
// a set of nodes that agree (but can fail)
// excluded, if set will be skipped altogether
func FindClosestVBlocking(qSet types.SCPQuorumSet, map1 map[types.NodeID]types.SCPEnvelope,
	filter func(types.SCPStatement) bool, excluded *types.NodeID) []types.NodeID {

	s := make(map[types.NodeID]struct{})

	for k, v := range map1 {
		if filter(v.Statement) {
			s[k] = struct{}{}
		}
	}
	return findClosestVBlockingF(qSet, s, excluded)
}

func findClosestVBlockingF(qSet types.SCPQuorumSet, nodes map[types.NodeID]struct{},
	excluded *types.NodeID) []types.NodeID {

	leftTillBlock := (1 + len(qSet.Validators) + len(qSet.InnerSets)) - int(qSet.Threshold)

	var res []types.NodeID

	// first, compute how many top level items need to be blocked
	for _, validator := range qSet.Validators {
		if excluded == nil || !(validator == *excluded) {

			if _, exist := nodes[validator]; !exist {
				// n was not present
				leftTillBlock--
				if leftTillBlock == 0 {
					// already blocked
					return []types.NodeID{}
				}
			} else {
				// save this for later
				res = append(res, validator)
			}
		}
	}

	//SliceStable(slice interface{}, less func(i, j int) bool)
	var resInternals [][]types.NodeID

	for _, inner := range qSet.InnerSets {
		v := findClosestVBlockingF(inner, nodes, excluded)
		if len(v) == 0 {
			leftTillBlock--
			if leftTillBlock == 0 {
				// already blocked
				return []types.NodeID{}
			}
		} else {
			resInternals = append(resInternals, v)
		}
	}
	//sort by length (stable) after all are inserted
	sort.SliceStable(resInternals, func(i, j int) bool { return len(resInternals[i]) < len(resInternals[j]) })

	//use the top level validators to get closer
	if len(res) > leftTillBlock {
		res = res[:leftTillBlock-1]
	}
	leftTillBlock -= len(res)

	// use subsets to get closer, using the smallest ones first
	for it := 0; leftTillBlock != 0 && it < len(resInternals); it++ {
		res = append(res, resInternals[it]...)
		leftTillBlock--
	}

	return res
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
