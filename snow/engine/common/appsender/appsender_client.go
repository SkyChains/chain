// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package appsender

import (
	"context"

	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/snow/engine/common"

	appsenderpb "github.com/luxdefi/luxd/proto/pb/appsender"
)

var _ common.AppSender = (*Client)(nil)

type Client struct {
	client appsenderpb.AppSenderClient
}

// NewClient returns a client that is connected to a remote AppSender.
func NewClient(client appsenderpb.AppSenderClient) *Client {
	return &Client{client: client}
}

func (c *Client) SendCrossChainAppRequest(ctx context.Context, chainID ids.ID, requestID uint32, appRequestBytes []byte) error {
	_, err := c.client.SendCrossChainAppRequest(
		ctx,
		&appsenderpb.SendCrossChainAppRequestMsg{
			ChainId:   chainID[:],
			RequestId: requestID,
			Request:   appRequestBytes,
		},
	)
	return err
}

func (c *Client) SendCrossChainAppResponse(ctx context.Context, chainID ids.ID, requestID uint32, appResponseBytes []byte) error {
	_, err := c.client.SendCrossChainAppResponse(
		ctx,
		&appsenderpb.SendCrossChainAppResponseMsg{
			ChainId:   chainID[:],
			RequestId: requestID,
			Response:  appResponseBytes,
		},
	)
	return err
}

func (c *Client) SendAppRequest(ctx context.Context, nodeIDs ids.NodeIDSet, requestID uint32, request []byte) error {
	nodeIDsBytes := make([][]byte, nodeIDs.Len())
	i := 0
	for nodeID := range nodeIDs {
		nodeID := nodeID // Prevent overwrite in next iteration
		nodeIDsBytes[i] = nodeID[:]
		i++
	}
	_, err := c.client.SendAppRequest(
		ctx,
		&appsenderpb.SendAppRequestMsg{
			NodeIds:   nodeIDsBytes,
			RequestId: requestID,
			Request:   request,
		},
	)
	return err
}

func (c *Client) SendAppResponse(ctx context.Context, nodeID ids.NodeID, requestID uint32, response []byte) error {
	_, err := c.client.SendAppResponse(
		ctx,
		&appsenderpb.SendAppResponseMsg{
			NodeId:    nodeID[:],
			RequestId: requestID,
			Response:  response,
		},
	)
	return err
}

func (c *Client) SendAppGossip(ctx context.Context, msg []byte) error {
	_, err := c.client.SendAppGossip(
		ctx,
		&appsenderpb.SendAppGossipMsg{
			Msg: msg,
		},
	)
	return err
}

func (c *Client) SendAppGossipSpecific(ctx context.Context, nodeIDs ids.NodeIDSet, msg []byte) error {
	nodeIDsBytes := make([][]byte, nodeIDs.Len())
	i := 0
	for nodeID := range nodeIDs {
		nodeID := nodeID // Prevent overwrite in next iteration
		nodeIDsBytes[i] = nodeID[:]
		i++
	}
	_, err := c.client.SendAppGossipSpecific(
		ctx,
		&appsenderpb.SendAppGossipSpecificMsg{
			NodeIds: nodeIDsBytes,
			Msg:     msg,
		},
	)
	return err
}
