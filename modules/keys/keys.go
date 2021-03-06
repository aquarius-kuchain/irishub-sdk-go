package keys

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const ModuleName = "keys"

type keysClient struct {
	sdk.KeyManager
}

func Create(keyManager sdk.KeyManager) rpc.Keys {
	return keysClient{
		KeyManager: keyManager,
	}
}

func (k keysClient) Add(name, password string) (string, string, sdk.Error) {
	address, mnemonic, err := k.Insert(name, password)
	return address, mnemonic, sdk.Wrap(err)
}

func (k keysClient) Recover(name, password, mnemonic string) (string, sdk.Error) {
	address, err := k.KeyManager.Recover(name, password, mnemonic)
	return address, sdk.Wrap(err)
}

func (k keysClient) Import(name, password, keystore string) (string, sdk.Error) {
	address, err := k.KeyManager.Import(name, password, keystore)
	return address, sdk.Wrap(err)
}

func (k keysClient) Export(name, srcPwd, dstPwd string) (string, sdk.Error) {
	keystore, err := k.KeyManager.Export(name, srcPwd, dstPwd)
	return keystore, sdk.Wrap(err)
}

func (k keysClient) Delete(name string) sdk.Error {
	err := k.KeyManager.Delete(name)
	return sdk.Wrap(err)
}

func (k keysClient) Show(name string) (string, sdk.Error) {
	address, err := k.KeyManager.Query(name)
	if err != nil {
		return "", sdk.Wrap(err)
	}
	return address.String(), nil
}

func (k keysClient) RegisterCodec(_ sdk.Codec) {
	//do nothing
}

func (k keysClient) Name() string {
	return ModuleName
}
