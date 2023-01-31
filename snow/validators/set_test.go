// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
)

<<<<<<< HEAD
<<<<<<< HEAD
func TestSetAddZeroWeight(t *testing.T) {
	require := require.New(t)

	s := NewSet()
	err := s.Add(ids.GenerateTestNodeID(), nil, ids.Empty, 0)
	require.ErrorIs(err, errZeroWeight)
}

func TestSetAddDuplicate(t *testing.T) {
	require := require.New(t)
=======
func TestSamplerSample(t *testing.T) {
	vdr0 := ids.GenerateTestNodeID()
	vdr1 := ids.GenerateTestNodeID()
>>>>>>> 1437bfe45 (Remove validators.Set#Set from the interface (#2275))

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
	err := s.Add(nodeID, nil, ids.Empty, 1)
	require.NoError(err)

	err = s.Add(nodeID, nil, ids.Empty, 1)
	require.ErrorIs(err, errDuplicateValidator)
}

=======
func TestSetAddZeroWeight(t *testing.T) {
	require := require.New(t)

	s := NewSet()
	err := s.Add(ids.GenerateTestNodeID(), 0)
	require.ErrorIs(err, errZeroWeight)
}

func TestSetAddDuplicate(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
	err := s.Add(nodeID, 1)
	require.NoError(err)

	err = s.Add(nodeID, 1)
	require.ErrorIs(err, errDuplicateValidator)
}

>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
func TestSetAddOverflow(t *testing.T) {
	require := require.New(t)

	s := NewSet()
<<<<<<< HEAD
	err := s.Add(ids.GenerateTestNodeID(), nil, ids.Empty, 1)
	require.NoError(err)

	err = s.Add(ids.GenerateTestNodeID(), nil, ids.Empty, math.MaxUint64)
=======
	err := s.Add(ids.GenerateTestNodeID(), 1)
	require.NoError(err)

	err = s.Add(ids.GenerateTestNodeID(), math.MaxUint64)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.Error(err)

	weight := s.Weight()
	require.EqualValues(1, weight)
}

func TestSetAddWeightZeroWeight(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
<<<<<<< HEAD
	err := s.Add(nodeID, nil, ids.Empty, 1)
=======
	err := s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	err = s.AddWeight(nodeID, 0)
	require.ErrorIs(err, errZeroWeight)
}

func TestSetAddWeightOverflow(t *testing.T) {
	require := require.New(t)

	s := NewSet()

<<<<<<< HEAD
	err := s.Add(ids.GenerateTestNodeID(), nil, ids.Empty, 1)
	require.NoError(err)

	nodeID := ids.GenerateTestNodeID()
	err = s.Add(nodeID, nil, ids.Empty, 1)
=======
	err := s.Add(ids.GenerateTestNodeID(), 1)
	require.NoError(err)

	nodeID := ids.GenerateTestNodeID()
	err = s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	err = s.AddWeight(nodeID, math.MaxUint64-1)
	require.Error(err)

	weight := s.Weight()
	require.EqualValues(2, weight)
}

func TestSetGetWeight(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
	weight := s.GetWeight(nodeID)
	require.Zero(weight)

<<<<<<< HEAD
	err := s.Add(nodeID, nil, ids.Empty, 1)
=======
	err := s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	weight = s.GetWeight(nodeID)
	require.EqualValues(1, weight)
}

func TestSetSubsetWeight(t *testing.T) {
	require := require.New(t)

	nodeID0 := ids.GenerateTestNodeID()
	nodeID1 := ids.GenerateTestNodeID()
	nodeID2 := ids.GenerateTestNodeID()

	weight0 := uint64(93)
	weight1 := uint64(123)
	weight2 := uint64(810)

<<<<<<< HEAD
	subset := set.Set[ids.NodeID]{}
=======
	subset := ids.NodeIDSet{}
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	subset.Add(nodeID0)
	subset.Add(nodeID1)

	s := NewSet()

<<<<<<< HEAD
	err := s.Add(nodeID0, nil, ids.Empty, weight0)
	require.NoError(err)

	err = s.Add(nodeID1, nil, ids.Empty, weight1)
	require.NoError(err)

	err = s.Add(nodeID2, nil, ids.Empty, weight2)
=======
	err := s.Add(nodeID0, weight0)
	require.NoError(err)

	err = s.Add(nodeID1, weight1)
	require.NoError(err)

	err = s.Add(nodeID2, weight2)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	expectedWeight := weight0 + weight1
	subsetWeight := s.SubsetWeight(subset)
	require.Equal(expectedWeight, subsetWeight)
}

func TestSetRemoveWeightZeroWeight(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
<<<<<<< HEAD
	err := s.Add(nodeID, nil, ids.Empty, 1)
=======
	err := s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	err = s.RemoveWeight(nodeID, 0)
	require.ErrorIs(err, errZeroWeight)
}

func TestSetRemoveWeightMissingValidator(t *testing.T) {
	require := require.New(t)

	s := NewSet()

<<<<<<< HEAD
	err := s.Add(ids.GenerateTestNodeID(), nil, ids.Empty, 1)
=======
	err := s.Add(ids.GenerateTestNodeID(), 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	err = s.RemoveWeight(ids.GenerateTestNodeID(), 1)
	require.ErrorIs(err, errMissingValidator)
}

func TestSetRemoveWeightUnderflow(t *testing.T) {
	require := require.New(t)

	s := NewSet()

<<<<<<< HEAD
	err := s.Add(ids.GenerateTestNodeID(), nil, ids.Empty, 1)
	require.NoError(err)

	nodeID := ids.GenerateTestNodeID()
	err = s.Add(nodeID, nil, ids.Empty, 1)
=======
	err := s.Add(ids.GenerateTestNodeID(), 1)
	require.NoError(err)

	nodeID := ids.GenerateTestNodeID()
	err = s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	err = s.RemoveWeight(nodeID, 2)
	require.Error(err)

	weight := s.Weight()
	require.EqualValues(2, weight)
}

func TestSetGet(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
	_, ok := s.Get(nodeID)
	require.False(ok)

<<<<<<< HEAD
	sk, err := bls.NewSecretKey()
	require.NoError(err)

	pk := bls.PublicFromSecretKey(sk)
	err = s.Add(nodeID, pk, ids.Empty, 1)
=======
	err := s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	vdr0, ok := s.Get(nodeID)
	require.True(ok)
<<<<<<< HEAD
	require.Equal(nodeID, vdr0.NodeID)
	require.Equal(pk, vdr0.PublicKey)
	require.EqualValues(1, vdr0.Weight)
=======
	require.Equal(nodeID, vdr0.ID())
	require.EqualValues(1, vdr0.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))

	err = s.AddWeight(nodeID, 1)
	require.NoError(err)

	vdr1, ok := s.Get(nodeID)
	require.True(ok)
<<<<<<< HEAD
	require.Equal(nodeID, vdr0.NodeID)
	require.Equal(pk, vdr0.PublicKey)
	require.EqualValues(1, vdr0.Weight)
	require.Equal(nodeID, vdr1.NodeID)
	require.Equal(pk, vdr1.PublicKey)
	require.EqualValues(2, vdr1.Weight)
=======
	require.Equal(nodeID, vdr0.ID())
	require.EqualValues(1, vdr0.Weight())
	require.Equal(nodeID, vdr1.ID())
	require.EqualValues(2, vdr1.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
}

func TestSetContains(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	nodeID := ids.GenerateTestNodeID()
	contains := s.Contains(nodeID)
	require.False(contains)

<<<<<<< HEAD
	err := s.Add(nodeID, nil, ids.Empty, 1)
=======
	err := s.Add(nodeID, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	contains = s.Contains(nodeID)
	require.True(contains)

	err = s.RemoveWeight(nodeID, 1)
	require.NoError(err)

	contains = s.Contains(nodeID)
	require.False(contains)
}

func TestSetLen(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	len := s.Len()
	require.Zero(len)

	nodeID0 := ids.GenerateTestNodeID()
<<<<<<< HEAD
	err := s.Add(nodeID0, nil, ids.Empty, 1)
=======
	err := s.Add(nodeID0, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	len = s.Len()
	require.Equal(1, len)

	nodeID1 := ids.GenerateTestNodeID()
<<<<<<< HEAD
	err = s.Add(nodeID1, nil, ids.Empty, 1)
=======
	err = s.Add(nodeID1, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	len = s.Len()
	require.Equal(2, len)

	err = s.RemoveWeight(nodeID1, 1)
	require.NoError(err)

	len = s.Len()
	require.Equal(1, len)

	err = s.RemoveWeight(nodeID0, 1)
	require.NoError(err)

	len = s.Len()
	require.Zero(len)
}

func TestSetList(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	list := s.List()
	require.Empty(list)

<<<<<<< HEAD
	sk, err := bls.NewSecretKey()
	require.NoError(err)

	pk := bls.PublicFromSecretKey(sk)
	nodeID0 := ids.GenerateTestNodeID()
	err = s.Add(nodeID0, pk, ids.Empty, 2)
=======
	nodeID0 := ids.GenerateTestNodeID()
	err := s.Add(nodeID0, 2)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	list = s.List()
	require.Len(list, 1)

	node0 := list[0]
<<<<<<< HEAD
	require.Equal(nodeID0, node0.NodeID)
	require.Equal(pk, node0.PublicKey)
	require.EqualValues(2, node0.Weight)

	nodeID1 := ids.GenerateTestNodeID()
	err = s.Add(nodeID1, nil, ids.Empty, 1)
=======
	require.Equal(nodeID0, node0.ID())
	require.EqualValues(2, node0.Weight())

	nodeID1 := ids.GenerateTestNodeID()
	err = s.Add(nodeID1, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	list = s.List()
	require.Len(list, 2)

	node0 = list[0]
<<<<<<< HEAD
	require.Equal(nodeID0, node0.NodeID)
	require.Equal(pk, node0.PublicKey)
	require.EqualValues(2, node0.Weight)

	node1 := list[1]
	require.Equal(nodeID1, node1.NodeID)
	require.Nil(node1.PublicKey)
	require.EqualValues(1, node1.Weight)

	err = s.RemoveWeight(nodeID0, 1)
	require.NoError(err)
	require.Equal(nodeID0, node0.NodeID)
	require.Equal(pk, node0.PublicKey)
	require.EqualValues(2, node0.Weight)
=======
	require.Equal(nodeID0, node0.ID())
	require.EqualValues(2, node0.Weight())

	node1 := list[1]
	require.Equal(nodeID1, node1.ID())
	require.EqualValues(1, node1.Weight())

	err = s.RemoveWeight(nodeID0, 1)
	require.NoError(err)
	require.Equal(nodeID0, node0.ID())
	require.EqualValues(2, node0.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))

	list = s.List()
	require.Len(list, 2)

	node0 = list[0]
<<<<<<< HEAD
	require.Equal(nodeID0, node0.NodeID)
	require.Equal(pk, node0.PublicKey)
	require.EqualValues(1, node0.Weight)

	node1 = list[1]
	require.Equal(nodeID1, node1.NodeID)
	require.Nil(node1.PublicKey)
	require.EqualValues(1, node1.Weight)
=======
	require.Equal(nodeID0, node0.ID())
	require.EqualValues(1, node0.Weight())

	node1 = list[1]
	require.Equal(nodeID1, node1.ID())
	require.EqualValues(1, node1.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))

	err = s.RemoveWeight(nodeID0, 1)
	require.NoError(err)

	list = s.List()
	require.Len(list, 1)

	node0 = list[0]
<<<<<<< HEAD
	require.Equal(nodeID1, node0.NodeID)
	require.Nil(node0.PublicKey)
	require.EqualValues(1, node0.Weight)
=======
	require.Equal(nodeID1, node0.ID())
	require.EqualValues(1, node0.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))

	err = s.RemoveWeight(nodeID1, 1)
	require.NoError(err)

	list = s.List()
	require.Empty(list)
}

func TestSetWeight(t *testing.T) {
	require := require.New(t)

	vdr0 := ids.NodeID{1}
	weight0 := uint64(93)
	vdr1 := ids.NodeID{2}
	weight1 := uint64(123)

	s := NewSet()
<<<<<<< HEAD
	err := s.Add(vdr0, nil, ids.Empty, weight0)
	require.NoError(err)

	err = s.Add(vdr1, nil, ids.Empty, weight1)
=======
	err := s.Add(vdr0, weight0)
	require.NoError(err)

	err = s.Add(vdr1, weight1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	setWeight := s.Weight()
	expectedWeight := weight0 + weight1
	require.Equal(expectedWeight, setWeight)
}

func TestSetSample(t *testing.T) {
	require := require.New(t)

	s := NewSet()

	sampled, err := s.Sample(0)
	require.NoError(err)
	require.Empty(sampled)

<<<<<<< HEAD
	sk, err := bls.NewSecretKey()
	require.NoError(err)

	nodeID0 := ids.GenerateTestNodeID()
	pk := bls.PublicFromSecretKey(sk)
	err = s.Add(nodeID0, pk, ids.Empty, 1)
=======
	nodeID0 := ids.GenerateTestNodeID()
	err = s.Add(nodeID0, 1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	sampled, err = s.Sample(1)
	require.NoError(err)
<<<<<<< HEAD
	require.Equal([]ids.NodeID{nodeID0}, sampled)
=======
	require.Len(sampled, 1)

	sampled0 := sampled[0]
	require.Equal(nodeID0, sampled0.ID())
	require.EqualValues(1, sampled0.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))

	_, err = s.Sample(2)
	require.Error(err)

	nodeID1 := ids.GenerateTestNodeID()
<<<<<<< HEAD
	err = s.Add(nodeID1, nil, ids.Empty, math.MaxInt64-1)
=======
	err = s.Add(nodeID1, math.MaxInt64-1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	sampled, err = s.Sample(1)
	require.NoError(err)
<<<<<<< HEAD
	require.Equal([]ids.NodeID{nodeID1}, sampled)

	sampled, err = s.Sample(2)
	require.NoError(err)
	require.Equal([]ids.NodeID{nodeID1, nodeID1}, sampled)

	sampled, err = s.Sample(3)
	require.NoError(err)
	require.Equal([]ids.NodeID{nodeID1, nodeID1, nodeID1}, sampled)
=======
	require.Len(sampled, 1)

	sampled0 = sampled[0]
	require.Equal(nodeID1, sampled0.ID())
	require.EqualValues(math.MaxInt64-1, sampled0.Weight())

	sampled, err = s.Sample(2)
	require.NoError(err)
	require.Len(sampled, 2)

	sampled0 = sampled[0]
	require.Equal(nodeID1, sampled0.ID())
	require.EqualValues(math.MaxInt64-1, sampled0.Weight())

	sampled1 := sampled[1]
	require.Equal(nodeID1, sampled1.ID())
	require.EqualValues(math.MaxInt64-1, sampled1.Weight())

	sampled, err = s.Sample(3)
	require.NoError(err)
	require.Len(sampled, 3)

	sampled0 = sampled[0]
	require.Equal(nodeID1, sampled0.ID())
	require.EqualValues(math.MaxInt64-1, sampled0.Weight())

	sampled1 = sampled[1]
	require.Equal(nodeID1, sampled1.ID())
	require.EqualValues(math.MaxInt64-1, sampled1.Weight())

	sampled2 := sampled[2]
	require.Equal(nodeID1, sampled2.ID())
	require.EqualValues(math.MaxInt64-1, sampled2.Weight())
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
}

func TestSetString(t *testing.T) {
	require := require.New(t)

	nodeID0 := ids.EmptyNodeID
	nodeID1 := ids.NodeID{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	s := NewSet()
<<<<<<< HEAD
	err := s.Add(nodeID0, nil, ids.Empty, 1)
	require.NoError(err)

	err = s.Add(nodeID1, nil, ids.Empty, math.MaxInt64-1)
=======
	err := s.Add(nodeID0, 1)
	require.NoError(err)

	err = s.Add(nodeID1, math.MaxInt64-1)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	expected := "Validator Set: (Size = 2, Weight = 9223372036854775807)\n" +
		"    Validator[0]: NodeID-111111111111111111116DBWJs, 1\n" +
		"    Validator[1]: NodeID-QLbz7JHiBTspS962RLKV8GndWFwdYhk6V, 9223372036854775806"
	result := s.String()
	require.Equal(expected, result)
}

var _ SetCallbackListener = (*callbackListener)(nil)

type callbackListener struct {
	t         *testing.T
<<<<<<< HEAD
	onAdd     func(ids.NodeID, *bls.PublicKey, ids.ID, uint64)
=======
	onAdd     func(ids.NodeID, uint64)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	onWeight  func(ids.NodeID, uint64, uint64)
	onRemoved func(ids.NodeID, uint64)
}

func (c *callbackListener) OnValidatorAdded(nodeID ids.NodeID, pk *bls.PublicKey, txID ids.ID, weight uint64) {
	if c.onAdd != nil {
		c.onAdd(nodeID, pk, txID, weight)
	} else {
		c.t.Fail()
	}
}

func (c *callbackListener) OnValidatorRemoved(nodeID ids.NodeID, weight uint64) {
	if c.onRemoved != nil {
		c.onRemoved(nodeID, weight)
	} else {
		c.t.Fail()
	}
}

func (c *callbackListener) OnValidatorWeightChanged(nodeID ids.NodeID, oldWeight, newWeight uint64) {
	if c.onWeight != nil {
		c.onWeight(nodeID, oldWeight, newWeight)
	} else {
		c.t.Fail()
	}
}

func TestSetAddCallback(t *testing.T) {
	require := require.New(t)

	nodeID0 := ids.NodeID{1}
<<<<<<< HEAD
	sk0, err := bls.NewSecretKey()
	require.NoError(err)
	pk0 := bls.PublicFromSecretKey(sk0)
	txID0 := ids.GenerateTestID()
=======
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	weight0 := uint64(1)

	s := NewSet()
	callCount := 0
	s.RegisterCallbackListener(&callbackListener{
		t: t,
<<<<<<< HEAD
		onAdd: func(nodeID ids.NodeID, pk *bls.PublicKey, txID ids.ID, weight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Equal(pk0, pk)
			require.Equal(txID0, txID)
=======
		onAdd: func(nodeID ids.NodeID, weight uint64) {
			require.Equal(nodeID0, nodeID)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
			require.Equal(weight0, weight)
			callCount++
		},
	})
<<<<<<< HEAD
	err = s.Add(nodeID0, pk0, txID0, weight0)
=======
	err := s.Add(nodeID0, weight0)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)
	require.Equal(1, callCount)
}

func TestSetAddWeightCallback(t *testing.T) {
	require := require.New(t)

	nodeID0 := ids.NodeID{1}
<<<<<<< HEAD
	txID0 := ids.GenerateTestID()
=======
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	weight0 := uint64(1)
	weight1 := uint64(93)

	s := NewSet()
<<<<<<< HEAD
	err := s.Add(nodeID0, nil, txID0, weight0)
=======
	err := s.Add(nodeID0, weight0)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	callCount := 0
	s.RegisterCallbackListener(&callbackListener{
		t: t,
<<<<<<< HEAD
		onAdd: func(nodeID ids.NodeID, pk *bls.PublicKey, txID ids.ID, weight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Nil(pk)
			require.Equal(txID0, txID)
=======
		onAdd: func(nodeID ids.NodeID, weight uint64) {
			require.Equal(nodeID0, nodeID)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
			require.Equal(weight0, weight)
			callCount++
		},
		onWeight: func(nodeID ids.NodeID, oldWeight, newWeight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Equal(weight0, oldWeight)
			require.Equal(weight0+weight1, newWeight)
			callCount++
		},
	})
	err = s.AddWeight(nodeID0, weight1)
	require.NoError(err)
	require.Equal(2, callCount)
}

func TestSetRemoveWeightCallback(t *testing.T) {
	require := require.New(t)

	nodeID0 := ids.NodeID{1}
<<<<<<< HEAD
	txID0 := ids.GenerateTestID()
=======
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	weight0 := uint64(93)
	weight1 := uint64(92)

	s := NewSet()
<<<<<<< HEAD
	err := s.Add(nodeID0, nil, txID0, weight0)
=======
	err := s.Add(nodeID0, weight0)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	callCount := 0
	s.RegisterCallbackListener(&callbackListener{
		t: t,
<<<<<<< HEAD
		onAdd: func(nodeID ids.NodeID, pk *bls.PublicKey, txID ids.ID, weight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Nil(pk)
			require.Equal(txID0, txID)
=======
		onAdd: func(nodeID ids.NodeID, weight uint64) {
			require.Equal(nodeID0, nodeID)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
			require.Equal(weight0, weight)
			callCount++
		},
		onWeight: func(nodeID ids.NodeID, oldWeight, newWeight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Equal(weight0, oldWeight)
			require.Equal(weight0-weight1, newWeight)
			callCount++
		},
	})
	err = s.RemoveWeight(nodeID0, weight1)
	require.NoError(err)
	require.Equal(2, callCount)
}

func TestSetValidatorRemovedCallback(t *testing.T) {
	require := require.New(t)

	nodeID0 := ids.NodeID{1}
<<<<<<< HEAD
	txID0 := ids.GenerateTestID()
	weight0 := uint64(93)

	s := NewSet()
	err := s.Add(nodeID0, nil, txID0, weight0)
=======
	weight0 := uint64(93)

	s := NewSet()
	err := s.Add(nodeID0, weight0)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
	require.NoError(err)

	callCount := 0
	s.RegisterCallbackListener(&callbackListener{
		t: t,
<<<<<<< HEAD
		onAdd: func(nodeID ids.NodeID, pk *bls.PublicKey, txID ids.ID, weight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Nil(pk)
			require.Equal(txID0, txID)
=======
		onAdd: func(nodeID ids.NodeID, weight uint64) {
			require.Equal(nodeID0, nodeID)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
			require.Equal(weight0, weight)
			callCount++
		},
		onRemoved: func(nodeID ids.NodeID, weight uint64) {
			require.Equal(nodeID0, nodeID)
			require.Equal(weight0, weight)
			callCount++
		},
	})
<<<<<<< HEAD
<<<<<<< HEAD
	err = s.RemoveWeight(nodeID0, weight0)
	require.NoError(err)
	require.Equal(2, callCount)
=======
	err = s.RemoveWeight(vdr0, weight0)
	require.NoError(t, err)
	require.Equal(t, 2, callcount)
>>>>>>> 1437bfe45 (Remove validators.Set#Set from the interface (#2275))
=======
	err = s.RemoveWeight(nodeID0, weight0)
	require.NoError(err)
	require.Equal(2, callCount)
>>>>>>> 749a0d8e9 (Add validators.Set#Add function and report errors (#2276))
}
