package KeyVault

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"math/big"
	"testing"
)

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

func TestNewKeyVault(t *testing.T) {
	// setup vault
	options := &PortfolioOptions{}
	options.EnableSimpleSigner(true)
	options.SetStorage(inmemStorage())
	options.SetEncryptor(keystorev4.New())
	options.SetPassword("password")
	v,err := NewKeyVault(options)
	require.NoError(t,err)

	testVault(t,v)
}

func TestImportKeyVault(t *testing.T) {
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	options := &PortfolioOptions{}
	options.EnableSimpleSigner(true)
	options.SetStorage(inmemStorage())
	options.SetSeed(seed)
	options.SetEncryptor(keystorev4.New())
	options.SetPassword("password")
	v,err := ImportKeyVault(options)
	require.NoError(t,err)

	// test common tests
	testVault(t,v)

	// test specific derivation
	w,err := v.WalletByName("wallet1")
	require.NoError(t,err)
	require.NotNil(t,w)
	val,err := w.AccountByName("val1")
	require.NoError(t,err)
	require.NotNil(t,val)
	with,err := w.GetWithdrawalAccount()
	require.NoError(t,err)
	require.NotNil(t,with)

	expectedValKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218").Bytes())
	require.NoError(t,err)
	expectedWithdrawalKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("51023953445614749789943419502694339066585011438324100967164633618358653841358").Bytes())
	require.NoError(t,err)

	assert.Equal(t,expectedValKey.PublicKey().Marshal(),val.PublicKey().Marshal())
	assert.Equal(t,"m/12381/3600/0/0/0",val.Path())
	assert.Equal(t,expectedWithdrawalKey.PublicKey().Marshal(),with.PublicKey().Marshal())
	assert.Equal(t,"m/12381/3600/0/0",with.Path())
}

func TestOpenKeyVault(t *testing.T) {
	// set up existing seed
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	storage := inmemStorage()
	storage.SetEncryptor(keystorev4.New(), []byte("password"))
	require.NoError(t,storage.SecurelySavePortfolioSeed(seed))

	// setup vault
	options := &PortfolioOptions{}
	options.EnableSimpleSigner(true)
	options.SetStorage(storage)
	options.SetEncryptor(keystorev4.New())
	options.SetPassword("password")
	v,err := OpenKeyVault(options)
	require.NoError(t,err)

	// test common tests
	testVault(t,v)

	// test specific derivation
	w,err := v.WalletByName("wallet1")
	require.NoError(t,err)
	require.NotNil(t,w)
	val,err := w.AccountByName("val1")
	require.NoError(t,err)
	require.NotNil(t,val)
	with,err := w.GetWithdrawalAccount()
	require.NoError(t,err)
	require.NotNil(t,with)

	expectedValKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218").Bytes())
	require.NoError(t,err)
	expectedWithdrawalKey,err := e2types.BLSPrivateKeyFromBytes(_bigInt("51023953445614749789943419502694339066585011438324100967164633618358653841358").Bytes())
	require.NoError(t,err)

	assert.Equal(t,expectedValKey.PublicKey().Marshal(),val.PublicKey().Marshal())
	assert.Equal(t,"m/12381/3600/0/0/0",val.Path())
	assert.Equal(t,expectedWithdrawalKey.PublicKey().Marshal(),with.PublicKey().Marshal())
	assert.Equal(t,"m/12381/3600/0/0",with.Path())
}

func testVault(t *testing.T, v *KeyVault) {
	// create wallet
	w,err := v.CreateWallet("wallet1")
	require.NoError(t,err)
	// fetch wallet
	w1,err := v.WalletByName("wallet1")
	require.NoError(t,err)
	w2,err := v.WalletByID(w.ID())
	require.NoError(t,err)
	require.NotNil(t,w1)
	require.NotNil(t,w2)
	require.Equal(t,w.ID().String(),w1.ID().String())
	require.Equal(t,w.ID().String(),w2.ID().String())
	require.Equal(t,w.Name(),w1.Name())
	require.Equal(t,w.Name(),w2.Name())

	// create and fetch validator account
	val,err := w.CreateValidatorAccount("val1")
	require.NoError(t,err)
	val1,err := w.AccountByName("val1")
	require.NoError(t,err)
	val2,err := w.AccountByID(val.ID())
	require.NoError(t,err)
	require.NotNil(t,val1)
	require.NotNil(t,val2)
	require.Equal(t,val.ID().String(),val1.ID().String())
	require.Equal(t,val.ID().String(),val2.ID().String())
	require.Equal(t,val.Name(),val1.Name())
	require.Equal(t,val.Name(),val2.Name())

	// create and fetch withdrawal account
	with,err := w.GetWithdrawalAccount()
	require.NoError(t,err)
	require.NotNil(t,with)
}