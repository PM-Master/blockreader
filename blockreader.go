package blockreader

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
)

func GetBlockData(block *common.Block) (BlockData, error) {
	blockData := block.Data.Data

	// First Get the Envelope from the BlockData
	envelope, err := GetEnvelopeFromBlock(blockData[0])
	if err != nil {
		return BlockData{}, err
	}

	// Retrieve the Payload from the Envelope
	payload := &common.Payload{}
	err = proto.Unmarshal(envelope.Payload, payload)
	if err != nil {
		return BlockData{}, err
	}

	payloadJson, err := GetPayloadJson(payload)
	if err != nil {
		return BlockData{}, err
	}

	// Read the Transaction from the Payload Data
	transaction := &peer.Transaction{}
	err = proto.Unmarshal(payload.Data, transaction)
	if err != nil {
		return BlockData{}, err
	}

	// Payload field is marshalled object of ChaincodeActionPayload
	chaincodeActionPayload := &peer.ChaincodeActionPayload{}
	err = proto.Unmarshal(transaction.Actions[0].Payload, chaincodeActionPayload)
	if err != nil {
		return BlockData{}, err
	}

	transactionJson, err := GetTransactionJson(chaincodeActionPayload)
	if err != nil {
		return BlockData{}, err
	}

	blockDataJson := BlockData{
		Envelope: Envelope{
			Header: Header{
				Payload: payloadJson,
			},
			Data: Data{
				Transaction: transactionJson,
			},
		},
	}

	return blockDataJson, nil
}

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {

	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	return env, nil
}

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
