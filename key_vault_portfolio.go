package KeyVault

import (
	"fmt"
	core "github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/portfolios/hd"
	"github.com/google/uuid"
	util "github.com/wealdtech/go-eth2-util"
)

// CreateAccount creates a new account in the wallet.
// This will error if an account with the name already exists.
// Will push to the new wallet the lock policy
func (portfolio *KeyVault) CreateWallet(name string) (core.Wallet, error) {
	var retWallet *hd.HDWallet

	// create wallet
	id := len(portfolio.walletIds)
	path := fmt.Sprintf("m/12381/3600/%d",id)
	nodeBytes,err := util.PrivateKeyFromSeedAndPath(portfolio.key.Seed(),path)
	if err != nil {
		return nil,err
	}
	lockableKey := core.NewEncryptableSeed(nodeBytes.Marshal(),portfolio.context.Encryptor)
	retWallet,err = hd.NewHDWallet(name,
		lockableKey,
		path,
		portfolio.context,
	)
	if err != nil {
		return nil,err
	}

	// register new wallet and save portfolio
	reset := func() {
		portfolio.walletsIndexer.Remove(retWallet.ID(),name)
		portfolio.walletIds = portfolio.walletIds[:len(portfolio.walletIds)-1]
	}
	portfolio.walletIds = append(portfolio.walletIds,retWallet.ID()) // we first save indexes so they could be saved as part of the portfolio
	portfolio.walletsIndexer.Add(retWallet.ID(),name)
	err = portfolio.context.Storage.SaveWallet(retWallet)
	if err != nil {
		reset()
		return nil,err
	}
	err = portfolio.context.Storage.SavePortfolio(portfolio)
	if err != nil {
		reset()
		return nil,err
	}

	return retWallet,nil
}

// Accounts provides all accounts in the wallet.
func (portfolio *KeyVault) Wallets() (<-chan core.Wallet,error) {
	ch := make (chan core.Wallet,1024) // TODO - handle more?
	go func() {
		for i := range portfolio.walletIds {
			id := portfolio.walletIds[i]
			wallet,err := portfolio.WalletByID(id)
			if err != nil {
				continue
			}
			ch <- wallet
		}
	}()

	return ch,nil
}

// AccountByID provides a single account from the wallet given its ID.
// This will error if the account is not found.
func (portfolio *KeyVault) WalletByID(id uuid.UUID) (core.Wallet, error) {
	return portfolio.context.Storage.OpenWallet(id)
}

// AccountByName provides a single account from the wallet given its name.
// This will error if the account is not found.
func (portfolio *KeyVault) WalletByName(name string) (core.Wallet, error) {
	id,exists := portfolio.walletsIndexer.ID(name)
	if !exists {
		return nil, fmt.Errorf("no wallet found")
	}

	return portfolio.WalletByID(id)
}