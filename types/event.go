package types

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tmclient "github.com/tendermint/tendermint/rpc/client"

	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	TypeKey      EventKey = "tm.event"
	ActionKey    EventKey = "action"
	SenderKey    EventKey = "sender"
	RecipientKey EventKey = "recipient"

	TxValue EventValue = "Tx"
)

type WSClient interface {
	SubscribeNewBlock(builder *EventQueryBuilder, handler EventNewBlockHandler) (Subscription, Error)
	SubscribeTx(builder *EventQueryBuilder, handler EventTxHandler) (Subscription, Error)
	SubscribeNewBlockHeader(handler EventNewBlockHeaderHandler) (Subscription, Error)
	SubscribeValidatorSetUpdates(handler EventValidatorSetUpdatesHandler) (Subscription, Error)
	Unsubscribe(subscription Subscription) Error
}

type TmClient interface {
	tmclient.ABCIClient
	tmclient.SignClient
	WSClient
}

type EventType string
type EventKey string
type EventValue string

type Subscription struct {
	Ctx   context.Context `json:"-"`
	Query string          `json:"query"`
	ID    string          `json:"id"`
}

type EventHandler func(data EventData)

// EventData for SubscribeAny
type EventData interface{}

// EventDataTx for SubscribeTx
type EventDataTx struct {
	Hash   string   `json:"hash"`
	Height int64    `json:"height"`
	Index  uint32   `json:"index"`
	Tx     StdTx    `json:"tx"`
	Result TxResult `json:"result"`
}

func (tx EventDataTx) MarshalJson() []byte {
	bz, _ := json.Marshal(tx)
	return bz
}

type TxResult struct {
	Code      uint32 `json:"code"`
	Log       string `json:"log"`
	GasWanted int64  `json:"gas_wanted"`
	GasUsed   int64  `json:"gas_used"`
	Tags      Tags   `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Tags []Tag

func ParseTags(pairs []cmn.KVPair) (tags []Tag) {
	if pairs == nil || len(pairs) == 0 {
		return tags
	}
	for _, pair := range pairs {
		key := string(pair.Key)
		value := string(pair.Value)
		tags = append(tags, Tag{
			Key:   key,
			Value: value,
		})
	}
	return
}
func (t Tags) GetValues(key string) (values []string) {
	for _, tag := range t {
		if tag.Key == key {
			values = append(values, tag.Value)
		}
	}
	return
}

func (t Tags) GetValue(key string) string {
	for _, tag := range t {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
}

func (t Tags) String() string {
	var tags = make([]string, len(t))
	for i, tag := range t {
		tags[i] = fmt.Sprintf("%s=%s", tag.Key, tag.Value)
	}
	return strings.Join(tags, ",")
}

type EventTxHandler func(EventDataTx)

//EventDataNewBlock for SubscribeNewBlock
type EventDataNewBlock struct {
	Block            Block            `json:"block"`
	ResultBeginBlock ResultBeginBlock `json:"result_begin_block"`
	ResultEndBlock   ResultEndBlock   `json:"result_end_block"`
}

type ValidatorUpdate struct {
	PubKey PubKey `json:"pub_key"`
	Power  int64  `json:"power"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EventNewBlockHandler func(EventDataNewBlock)

//EventDataNewBlockHeader for SubscribeNewBlockHeader
type EventDataNewBlockHeader struct {
	Header tmtypes.Header `json:"header"`

	ResultBeginBlock ResultBeginBlock `json:"result_begin_block"`
	ResultEndBlock   ResultEndBlock   `json:"result_end_block"`
}

type EventNewBlockHeaderHandler func(EventDataNewBlockHeader)

//EventDataValidatorSetUpdates for SubscribeValidatorSetUpdates
type Validator struct {
	Bech32Address    string `json:"bech32_address"`
	Bech32PubKey     string `json:"bech32_pubkey"`
	Address          string `json:"address"`
	PubKey           PubKey `json:"pub_key"`
	VotingPower      int64  `json:"voting_power"`
	ProposerPriority int64  `json:"proposer_priority"`
}
type EventDataValidatorSetUpdates struct {
	ValidatorUpdates []Validator `json:"validator_updates"`
}

type EventValidatorSetUpdatesHandler func(EventDataValidatorSetUpdates)

//EventQueryBuilder for build query string
type condition struct {
	key   EventKey
	value EventValue
	op    string
}

func Cond(key EventKey) *condition {
	return &condition{
		key: key,
	}
}

func (c *condition) LTE(v EventValue) *condition {
	return c.fill(v, "<=")
}

func (c *condition) GTE(v EventValue) *condition {
	return c.fill(v, ">=")
}

func (c *condition) LE(v EventValue) *condition {
	return c.fill(v, "<")
}

func (c *condition) GE(v EventValue) *condition {
	return c.fill(v, ">")
}

func (c *condition) EQ(v EventValue) *condition {
	return c.fill(v, "=")
}

func (c *condition) Contains(v EventValue) *condition {
	return c.fill(v, "CONTAINS")
}

func (c *condition) fill(v EventValue, op string) *condition {
	c.value = v
	c.op = op
	return c
}

func (c *condition) String() string {
	if len(c.key) == 0 || len(c.value) == 0 || len(c.op) == 0 {
		return ""
	}
	return fmt.Sprintf("%s %s '%s'", c.key, c.op, c.value)
}

//EventQueryBuilder is responsible for constructing listening conditions
type EventQueryBuilder struct {
	conditions []string
}

func NewEventQueryBuilder() *EventQueryBuilder {
	return &EventQueryBuilder{
		conditions: []string{},
	}
}

//AddCondition is responsible for adding listening conditions
func (eqb *EventQueryBuilder) AddCondition(c *condition) *EventQueryBuilder {
	if c == nil {
		return nil
	}
	eqb.conditions = append(eqb.conditions, c.String())
	return eqb
}

//Build is responsible for constructing the listening condition into a listening instruction identified by tendermint
func (eqb *EventQueryBuilder) Build() string {
	var buf bytes.Buffer
	for _, condition := range eqb.conditions {
		if buf.Len() > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(condition)
	}
	return buf.String()
}
