package keys_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
)

type KeysTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

func (kts *KeysTestSuite) SetupTest() {
	tc := test.GetMock()
	kts.MockClient = tc
}

func (kts *KeysTestSuite) TestKeys() {
	name, password := kts.RandStringOfLength(20), kts.RandStringOfLength(8)

	address, mnemonic, err := kts.Keys().Add(name, password)
	require.NoError(kts.T(), err)
	require.NotEmpty(kts.T(), address)
	require.NotEmpty(kts.T(), mnemonic)

	address1, err := kts.Keys().Show(name)
	require.NoError(kts.T(), err)
	require.Equal(kts.T(), address, address1)

	newPwd := kts.RandStringOfLength(8)
	keystore, err := kts.Keys().Export(name, password, newPwd)
	require.NoError(kts.T(), err)

	err = kts.Keys().Delete(name)
	require.NoError(kts.T(), err)

	address2, err := kts.Keys().Import(name, newPwd, keystore)
	require.NoError(kts.T(), err)
	require.Equal(kts.T(), address, address2)

	err = kts.Keys().Delete(name)
	require.NoError(kts.T(), err)

	address3, err := kts.Keys().Recover(name, newPwd, mnemonic)
	require.NoError(kts.T(), err)
	require.Equal(kts.T(), address, address3)
}
