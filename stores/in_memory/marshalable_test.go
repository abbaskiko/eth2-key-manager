package in_memory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/wallet_hd"
)

func TestMarshaling(t *testing.T) {
	store := NewInMemStore(core.MainNetwork)

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage: store})
	err := store.SaveWallet(wallet)
	require.NoError(t, err)

	// account
	acc, err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), nil)
	require.NoError(t, err)
	err = store.SaveAccount(acc)
	require.NoError(t, err)

	// attestation
	att := &core.BeaconAttestation{
		Slot:            1,
		CommitteeIndex:  1,
		BeaconBlockRoot: []byte("A"),
		Source: &core.Checkpoint{
			Epoch: 1,
			Root:  []byte("A"),
		},
		Target: &core.Checkpoint{
			Epoch: 2,
			Root:  []byte("A"),
		},
	}
	store.SaveAttestation(acc.ValidatorPublicKey(), att)

	// proposal
	prop := &core.BeaconBlockHeader{
		Slot:          1,
		ProposerIndex: 1,
		ParentRoot:    []byte("A"),
		StateRoot:     []byte("A"),
		BodyRoot:      []byte("A"),
	}
	store.SaveProposal(acc.ValidatorPublicKey(), prop)

	// marshal
	byts, err := json.Marshal(store)
	require.NoError(t, err)

	// un-marshal
	var store2 InMemStore
	require.NoError(t, json.Unmarshal(byts, &store2))

	// verify
	t.Run("verify wallet", func(t *testing.T) {
		wallet2, err := store2.OpenWallet()
		require.NoError(t, err)
		require.Equal(t, wallet.ID().String(), wallet2.ID().String())
	})
	t.Run("verify acc", func(t *testing.T) {
		wallet2, err := store2.OpenWallet()
		require.NoError(t, err)
		acc2, err := wallet2.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")
		require.NoError(t, err)
		require.Equal(t, acc.ID().String(), acc2.ID().String())
	})
	t.Run("verify attestation", func(t *testing.T) {
		att2, err := store.RetrieveAttestation(acc.ValidatorPublicKey(), 2)
		require.NoError(t, err)
		require.Equal(t, att.BeaconBlockRoot, att2.BeaconBlockRoot)
	})
	t.Run("verify proposal", func(t *testing.T) {
		prop2, err := store.RetrieveProposal(acc.ValidatorPublicKey(), 1)
		require.NoError(t, err)
		require.Equal(t, prop.StateRoot, prop2.StateRoot)
	})
}
