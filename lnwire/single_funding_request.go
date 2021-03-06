package lnwire

import (
	"io"

	"github.com/roasbeef/btcd/btcec"
	"github.com/roasbeef/btcutil"
)

// SingleFundingRequest is the message Alice sends to Bob if we should like
// to create a channel with Bob where she's the sole provider of funds to the
// channel. Single funder channels simplify the initial funding workflow, are
// supported by nodes backed by SPV Bitcoin clients, and have a simpler
// security models than dual funded channels.
//
// NOTE: In order to avoid a slow loris like resource exhaustion attack, the
// responder of a single funding channel workflow *should not* watch the
// blockchain for the funding transaction. Instead, it is the initiator's job
// to provide the responder with an SPV proof of funding transaction inclusion
// after a sufficient number of confirmations.
type SingleFundingRequest struct {
	// PendingChannelID serves to uniquely identify the future channel
	// created by the initiated single funder workflow.
	PendingChannelID [32]byte

	// ChannelType represents the type of channel this request would like
	// to open. At this point, the only supported channels are type 0
	// channels, which are channels with regular commitment transactions
	// utilizing HTLCs for payments.
	ChannelType uint8

	// CoinType represents which blockchain the channel will be opened
	// using. By default, this field should be set to 0, indicating usage
	// of the Bitcoin blockchain.
	CoinType uint64

	// FeePerKb is the required number of satoshis per KB that the
	// requester will pay at all timers, for both the funding transaction
	// and commitment transaction. This value can later be updated once the
	// channel is open.
	FeePerKb btcutil.Amount

	// FundingAmount is the number of satoshis the initiator would like
	// to commit to the channel.
	FundingAmount btcutil.Amount

	// PushSatoshis is the number of satoshis that will be pushed to the
	// responder as a part of the first commitment state.
	PushSatoshis btcutil.Amount

	// CsvDelay is the number of blocks to use for the relative time lock
	// in the pay-to-self output of both commitment transactions.
	CsvDelay uint32

	// CommitmentKey is key the initiator of the funding workflow wishes to
	// use within their versino of the commitment transaction for any
	// delayed (CSV) or immediate outputs to them.
	CommitmentKey *btcec.PublicKey

	// ChannelDerivationPoint is an secp256k1 point which will be used to
	// derive the public key the initiator will use for the half of the
	// 2-of-2 multi-sig. Using the channel derivation point (CDP), and the
	// initiators identity public key (A), the channel public key is
	// computed as: C = A + CDP. In order to be valid all CDP's MUST have
	// an odd y-coordinate.
	ChannelDerivationPoint *btcec.PublicKey

	// DeliveryPkScript defines the public key script that the initiator
	// would like to use to receive their balance in the case of a
	// cooperative close. Only the following script templates are
	// supported: P2PKH, P2WKH, P2SH, and P2WSH.
	DeliveryPkScript PkScript

	// DustLimit is the threshold below which no HTLC output should be
	// generated for our commitment transaction; ie. HTLCs below
	// this amount are not enforceable onchain from our point view.
	DustLimit btcutil.Amount

	// ConfirmationDepth is the number of confirmations that the initiator
	// of a funding workflow is requesting be required before the channel
	// is considered fully open.
	ConfirmationDepth uint32
}

// NewSingleFundingRequest creates, and returns a new empty SingleFundingRequest.
func NewSingleFundingRequest(chanID [32]byte, chanType uint8, coinType uint64,
	fee btcutil.Amount, amt btcutil.Amount, delay uint32, ck,
	cdp *btcec.PublicKey, deliveryScript PkScript,
	dustLimit btcutil.Amount, pushSat btcutil.Amount,
	confDepth uint32) *SingleFundingRequest {

	return &SingleFundingRequest{
		PendingChannelID:       chanID,
		ChannelType:            chanType,
		CoinType:               coinType,
		FeePerKb:               fee,
		FundingAmount:          amt,
		CsvDelay:               delay,
		CommitmentKey:          ck,
		ChannelDerivationPoint: cdp,
		DeliveryPkScript:       deliveryScript,
		DustLimit:              dustLimit,
		PushSatoshis:           pushSat,
		ConfirmationDepth:      confDepth,
	}
}

// Decode deserializes the serialized SingleFundingRequest stored in the passed
// io.Reader into the target SingleFundingRequest using the deserialization
// rules defined by the passed protocol version.
//
// This is part of the lnwire.Message interface.
func (c *SingleFundingRequest) Decode(r io.Reader, pver uint32) error {
	return readElements(r,
		c.PendingChannelID[:],
		&c.ChannelType,
		&c.CoinType,
		&c.FeePerKb,
		&c.FundingAmount,
		&c.PushSatoshis,
		&c.CsvDelay,
		&c.CommitmentKey,
		&c.ChannelDerivationPoint,
		&c.DeliveryPkScript,
		&c.DustLimit,
		&c.ConfirmationDepth)
}

// Encode serializes the target SingleFundingRequest into the passed io.Writer
// implementation. Serialization will observe the rules defined by the passed
// protocol version.
//
// This is part of the lnwire.Message interface.
func (c *SingleFundingRequest) Encode(w io.Writer, pver uint32) error {
	return writeElements(w,
		c.PendingChannelID[:],
		c.ChannelType,
		c.CoinType,
		c.FeePerKb,
		c.FundingAmount,
		c.PushSatoshis,
		c.CsvDelay,
		c.CommitmentKey,
		c.ChannelDerivationPoint,
		c.DeliveryPkScript,
		c.DustLimit,
		c.ConfirmationDepth)
}

// MsgType returns the uint32 code which uniquely identifies this message as a
// SingleFundingRequest on the wire.
//
// This is part of the lnwire.Message interface.
func (c *SingleFundingRequest) MsgType() MessageType {
	return MsgSingleFundingRequest
}

// MaxPayloadLength returns the maximum allowed payload length for a
// SingleFundingRequest. This is calculated by summing the max length of all
// the fields within a SingleFundingRequest. To enforce a maximum
// DeliveryPkScript size, the size of a P2PKH public key script is used.
//
// This is part of the lnwire.Message interface.
func (c *SingleFundingRequest) MaxPayloadLength(uint32) uint32 {
	var length uint32

	// PendingChannelID - 32 bytes
	length += 32

	// ChannelType - 1 byte
	length++

	// CoinType - 8 bytes
	length += 8

	// FeePerKb - 8 bytes
	length += 8

	// FundingAmount - 8 bytes
	length += 8

	// PushSatoshis - 8 bytes
	length += 8

	// CsvDelay - 4 bytes
	length += 4

	// CommitmentKey - 33 bytes
	length += 33

	// ChannelDerivationPoint - 33 bytes
	length += 33

	// DeliveryPkScript - 34 bytes + 1 byte varInt
	length += 35

	// DustLimit - 8 bytes
	length += 8

	// ConfirmationDepth - 4 bytes
	length += 4

	return length
}
