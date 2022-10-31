// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"context"

	"github.com/luxdefi/luxd/ids"
)

// Send a query composed partially of push queries and partially of pull queries.
// The validators in [vdrs] will be queried.
// This function sends at most [numPushTo] push queries. The rest are pull queries.
// If [numPushTo] > len(vdrs), len(vdrs) push queries are sent.
// [containerID] and [container] are the ID and body of the container being queried.
// [sender] is used to actually send the queries.
func SendMixedQuery(
	ctx context.Context,
	sender Sender,
	vdrs []ids.NodeID,
	numPushTo int,
	reqID uint32,
	containerID ids.ID,
	container []byte,
) {
	if numPushTo > len(vdrs) {
		numPushTo = len(vdrs)
	}
	if numPushTo > 0 {
		sendPushQueryTo := ids.NewNodeIDSet(numPushTo)
		sendPushQueryTo.Add(vdrs[:numPushTo]...)
		sender.SendPushQuery(ctx, sendPushQueryTo, reqID, container)
	}
	if numPullTo := len(vdrs) - numPushTo; numPullTo > 0 {
		sendPullQueryTo := ids.NewNodeIDSet(numPullTo)
		sendPullQueryTo.Add(vdrs[numPushTo:]...)
		sender.SendPullQuery(ctx, sendPullQueryTo, reqID, containerID)
	}
}
