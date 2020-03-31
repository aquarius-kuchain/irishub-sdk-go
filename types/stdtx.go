package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"

	json2 "github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	maxMemoCharacters = 100
	txSigLimit        = 7

	Sync   BroadcastMode = "sync"
	Async  BroadcastMode = "async"
	Commit BroadcastMode = "commit"
)

type BroadcastMode string

// Transactions messages must fulfill the Msg
type Msg interface {
	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// QueryAndRefreshAccount the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []AccAddress
}

type Msgs []Msg

func (m Msgs) Len() int {
	return len(m)
}

func (m Msgs) Sub(begin, end int) SplitAble {
	return m[begin:end]
}

// Transactions objects must fulfill the Tx
type Tx interface {

	// Gets the all the transaction's messages.
	GetMsgs() []Msg
	// ValidateBasic does a simple and lightweight validation check that doesn't
	// require access to any other information.
	ValidateBasic() error
}

// StdFee includes the amount of coins paid in fees and the maximum
// Gas to be used by the transaction. The ratio yields an effective "gasprice",
// which must be above some miminum to be accepted into the mempool.
type StdFee struct {
	Amount Coins  `json:"amount"`
	Gas    uint64 `json:"gas"`
}

func NewStdFee(gas uint64, amount ...Coin) StdFee {
	return StdFee{
		Amount: amount,
		Gas:    gas,
	}
}

// Fee bytes for signing later
func (fee StdFee) Bytes() []byte {
	if len(fee.Amount) == 0 {
		fee.Amount = Coins{}
	}
	bz, err := NewAminoCodec().MarshalJSON(fee)
	if err != nil {
		panic(err)
	}
	return bz
}

type Signature struct {
	crypto.PubKey `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
}

// Standard Signature
type StdSignature struct {
	crypto.PubKey `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
	AccountNumber uint64           `json:"account_number"`
	Sequence      uint64           `json:"sequence"`
}

// StdSignMsg is a convenience structure for passing along
// a Msg with the other requirements for a StdSignDoc before
// it is signed. For use in the CLI.
type StdSignMsg struct {
	ChainID       string `json:"chain_id"`
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
	Fee           StdFee `json:"fee"`
	Msgs          []Msg  `json:"msgs"`
	Memo          string `json:"memo"`
}

// get message bytes
func (msg StdSignMsg) Bytes(cdc Codec) []byte {
	var msgsBytes []json.RawMessage
	for _, msg := range msg.Msgs {
		msgsBytes = append(msgsBytes, json.RawMessage(msg.GetSignBytes()))
	}
	bz, err := cdc.MarshalJSON(StdSignDoc{
		AccountNumber: msg.AccountNumber,
		ChainID:       msg.ChainID,
		Fee:           json.RawMessage(msg.Fee.Bytes()),
		Memo:          msg.Memo,
		Msgs:          msgsBytes,
		Sequence:      msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return json2.MustSort(bz)
}

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	AccountNumber uint64            `json:"account_number"`
	ChainID       string            `json:"chain_id"`
	Fee           json.RawMessage   `json:"fee"`
	Memo          string            `json:"memo"`
	Msgs          []json.RawMessage `json:"msgs"`
	Sequence      uint64            `json:"sequence"`
}

// StdTx is a standard way to wrap a Msg with Fee and Signatures.
// NOTE: the first signature is the fee payer (Signatures must not be nil).
type StdTx struct {
	Msgs       []Msg          `json:"msg"`
	Fee        StdFee         `json:"fee"`
	Signatures []StdSignature `json:"signatures"`
	Memo       string         `json:"memo"`
}

func NewStdTx(msgs []Msg, fee StdFee, sigs []StdSignature, memo string) StdTx {
	return StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: sigs,
		Memo:       memo,
	}
}

//nolint
// GetMsgs returns the all the transaction's messages.
func (tx StdTx) GetMsgs() []Msg { return tx.Msgs }
func (tx StdTx) GetSignBytes() []string {
	var bz []string
	for _, msg := range tx.Msgs {
		bz = append(bz, string(msg.GetSignBytes()))
	}
	return bz
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() error {
	stdSigs := tx.GetSignatures()
	if tx.Fee.Amount.IsAnyNegative() {
		return errors.New(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.Amount))
	}
	if len(stdSigs) == 0 {
		return errors.New("no signers")
	}
	if len(stdSigs) != len(tx.GetSigners()) {
		return errors.New("wrong number of signers")
	}
	if len(tx.GetMemo()) > maxMemoCharacters {
		return errors.New(
			fmt.Sprintf(
				"maximum number of characters is %d but received %d characters",
				maxMemoCharacters, len(tx.GetMemo()),
			),
		)
	}
	sigCount := 0
	for i := 0; i < len(stdSigs); i++ {
		sigCount += countSubKeys(stdSigs[i].PubKey)
		if sigCount > txSigLimit {
			return errors.New(
				fmt.Sprintf("ssdk.ErrTooManySignaturesignatures: %d, limit: %d", sigCount, txSigLimit),
			)
		}
	}
	return nil
}
func countSubKeys(pub crypto.PubKey) int {
	v, ok := pub.(multisig.PubKeyMultisigThreshold)
	if !ok {
		return 1
	}
	numKeys := 0
	for _, subkey := range v.PubKeys {
		numKeys += countSubKeys(subkey)
	}
	return numKeys
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigners() []AccAddress {
	seen := map[string]bool{}
	var signers []AccAddress
	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}
	return signers
}

//nolint
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetSignatures returns the signature of signers who signed the Msg.
// CONTRACT: Length returned is same as length of
// pubkeys returned from MsgKeySigners, and the order
// matches.
// CONTRACT: If the signature is missing (ie the Msg is
// invalid), then the corresponding signature is
// .Empty().
func (tx StdTx) GetSignatures() []StdSignature { return tx.Signatures }

type BaseTx struct {
	From     string        `json:"from"`
	Password string        `json:"password"`
	Gas      uint64        `json:"gas"`
	Fee      DecCoins      `json:"fee"`
	Memo     string        `json:"memo"`
	Mode     BroadcastMode `json:"broadcast_mode"`
	Simulate bool          `json:"simulate"`
}

// ResultTx encapsulates the return result of the transaction. When the transaction fails,
// it is an empty object. The specific error information can be obtained through the Error interface.
type ResultTx struct {
	GasWanted int64  `json:"gas_wanted"`
	GasUsed   int64  `json:"gas_used"`
	Tags      Tags   `json:"tags"`
	Hash      string `json:"hash"`
	Height    int64  `json:"height"`
}

// ResultQueryTx is used to prepare info to display
type ResultQueryTx struct {
	Hash      string   `json:"hash"`
	Height    int64    `json:"height"`
	Tx        Tx       `json:"tx"`
	Result    TxResult `json:"result"`
	Timestamp string   `json:"timestamp"`
}

// ResultSearchTxs defines a structure for querying txs pageable
type ResultSearchTxs struct {
	Total int             `json:"total"` // Count of all txs
	Txs   []ResultQueryTx `json:"txs"`   // List of txs in current page
}
