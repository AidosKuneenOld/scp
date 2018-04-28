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
	"time"

	"github.com/scp/types"
)

/**
 * The Slot object is in charge of maintaining the state of the SCP protocol
 * for a given slot index.
 */

type Slot struct {
	mSlotIndex uint64
	mSCP       SCP
	//mBallotProtocol     BallotProtocol
	//mNominationProtocol NominationProtocol
	mStatementsHistory []HistoricalStatement
	mFullyValidated    bool
}

// keeps track of all statements seen so far for this slot.
// it is used for debugging purpose
type HistoricalStatement struct {
	mWhen      time.Time
	mStatement types.SCPStatement
	mValidated bool
}

func (ns *Slot) NSlot(slotIndex uint64, scp SCP) {
	ns.mSlotIndex = slotIndex
	ns.mSCP = scp
	//ns.mBallotProtocol(*ns)
	//ns.mNominationProtocol(*ns)
	//ns.mFullyValidated = scp.getLocalNode.isValidator
}
