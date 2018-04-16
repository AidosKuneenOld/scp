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
	"github.com/scp/types"
)

type QuorumSetSanityChecker struct {
	mExtraChecks bool
	mKnownNodes  map[types.PublicKey]struct{}
	mIsSane      bool
	mCount       int
}

func (nq *QuorumSetSanityChecker) IsSane() bool {
	return nq.mIsSane
}

func (nq *QuorumSetSanityChecker) NQuorumSetSanityChecker(qSet types.SCPQuorumSet, extraChecks bool) {
	nq.mIsSane = nq.CheckSanity(qSet, 0) && nq.mCount >= 1 && nq.mCount <= 1000
	nq.mExtraChecks = extraChecks
}

func (nq *QuorumSetSanityChecker) CheckSanity(qSet types.SCPQuorumSet, depth int) bool {
	if depth > 2 {
		return false
	}
	if qSet.Threshold < 1 {
		return false
	}

	v := qSet.Validators
	i := qSet.InnerSets

	totEntries1 := len(v) + len(i)
	totEntries := uint32(totEntries1)
	vBlockingSize := totEntries - qSet.Threshold + 1
	nq.mCount += len(v)

	if qSet.Threshold > totEntries {
		return false
	}
	// threshold is within the proper range
	if nq.mExtraChecks && qSet.Threshold < vBlockingSize {
		return false
	}

	for _, n := range v {
		{
			if _, exist := nq.mKnownNodes[n]; exist {
				// n was already present
				return false
			}
			// insert
			nq.mKnownNodes[n] = struct{}{}
		}

		for _, iSet := range i {
			if !nq.CheckSanity(iSet, depth+1) {
				return false
			}
		}
	}
	return true
}

func IsQuorumSetSane(qSet types.SCPQuorumSet, extraChecks bool) bool {
	checker := QuorumSetSanityChecker{}
	checker.NQuorumSetSanityChecker(qSet, extraChecks)
	return checker.IsSane()
}

// helper function that:
//  * simplifies singleton inner set into outerset
//      { t: n, v: { ... }, { t: 1, X }, ... }
//        into
//      { t: n, v: { ..., X }, .... }
//  * simplifies singleton innersets
//      { t:1, { innerSet } } into innerSet
func NormalizeQSet(qSet *types.SCPQuorumSet) {
	v := qSet.Validators
	iS := qSet.InnerSets
	for i := len(iS) - 1; i >= 0; i-- {
		NormalizeQSet(&iS[i])
		// merge singleton inner sets into validator list
		if iS[i].Threshold == 1 && len(iS[i].Validators) == 1 &&
			len(iS[i].InnerSets) == 0 {
			v = append(v, iS[i].Validators[0])
			iS = append(iS[:i], iS[i+1:]...)
		}
	}

	// simplify quorum set if needed
	if qSet.Threshold == 1 && len(v) == 0 && len(iS) == 1 {

		//t := qSet.InnerSets[len(qSet.InnerSets)-1]
		// or t := qSet.InnerSets[0] ?
		qSet = &iS[0]
	}
}
