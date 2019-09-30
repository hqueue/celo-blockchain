// Copyright 2017 The Celo Authors
// This file is part of the celo library.
//
// The celo library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The celo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the celo library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	// "bytes"
	// "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	// mrand "math/rand"
	// "strings"
	"time"

	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/consensus/istanbul"
	// contract_errors "github.com/ethereum/go-ethereum/contract_comm/errors"
	// "github.com/ethereum/go-ethereum/contract_comm/validators"
	// "github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/log"
	// "github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
)

// ==============================================
//
// define the istanbul announce message

type valEnodeShareMessage struct {
	// Address               common.Address
    TestPayload           string
	// View                  *istanbul.View
	// Signature             []byte
}

func (sm *valEnodeShareMessage) String() string {
    return fmt.Sprintf("%+v", sm)
	// return fmt.Sprintf("{Address: %s, View: %v, IncompleteEnodeURL: %v, Signature: %v}", am.Address.String(), am.View, am.IncompleteEnodeURL, hex.EncodeToString(am.Signature))
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes am into the Ethereum RLP format.
func (sm *valEnodeShareMessage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{sm.TestPayload})
}

// DecodeRLP implements rlp.Decoder, and load the am fields from a RLP stream.
func (sm *valEnodeShareMessage) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
        TestPayload string
		// Address               common.Address
		// IncompleteEnodeURL    string
		// EncryptedEndpointData [][][]byte
		// View                  *istanbul.View
		// Signature             []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	sm.TestPayload = msg.TestPayload
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for the istanbul announce sender and handler
func (sm *valEnodeShareMessage) FromPayload(b []byte) error {
	// Decode message
	err := rlp.DecodeBytes(b, &sm)
	return err
}

func (sm *valEnodeShareMessage) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(sm)
}
//
// func (am *valEnodeShareMessage) Sign(signingFn func(data []byte) ([]byte, error)) error {
// 	// Construct and encode a message with no signature
// 	var payloadNoSig []byte
// 	payloadNoSig, err := rlp.EncodeToBytes(&announceMessage{
// 		Address:               am.Address,
// 		IncompleteEnodeURL:    am.IncompleteEnodeURL,
// 		EncryptedEndpointData: am.EncryptedEndpointData,
// 		View:                  am.View,
// 		Signature:             []byte{},
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	am.Signature, err = signingFn(payloadNoSig)
// 	return err
// }
//
// func (am *valEnodeShareMessage) VerifySig() error {
// 	// Construct and encode a message with no signature
// 	var payloadNoSig []byte
// 	payloadNoSig, err := rlp.EncodeToBytes(&announceMessage{
// 		Address:               am.Address,
// 		IncompleteEnodeURL:    am.IncompleteEnodeURL,
// 		EncryptedEndpointData: am.EncryptedEndpointData,
// 		View:                  am.View,
// 		Signature:             []byte{},
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	sigAddr, err := istanbul.GetSignatureAddress(payloadNoSig, am.Signature)
// 	if err != nil {
// 		return err
// 	}
//
// 	if sigAddr != am.Address {
// 		log.Error("Address in the message is different than the address that signed it",
// 			"sigAddr", sigAddr.Hex(),
// 			"msg.Address", am.Address.Hex())
// 		return errInvalidSignature
// 	}
//
// 	return nil
// }

// This function is meant to be run as a goroutine.  It will periodically gossip announce messages
// to the rest of the registered validators to communicate it's enodeURL to them.
func (sb *Backend) sendValEnodeShareMsgs() {
	sb.announceWg.Add(1)
	defer sb.announceWg.Done()

	ticker := time.NewTicker(time.Minute / 5.0)

	for {
		select {
		case <-ticker.C:
			// output the valEnodeTable for debugging purposes
			log.Trace("ValidatorEnodeTable dump", "ValidatorEnodeTable", sb.valEnodeTable.String())
			go sb.sendValEnodeShareMsg()

		case <-sb.announceQuit:
			ticker.Stop()
			return
		}
	}
}

func (sb *Backend) generateValEnodeShareMsg() ([]byte, error) {
	// block := sb.currentBlock()
	// selfEnode := sb.Enode()
    //
	// if selfEnode == nil {
	// 	sb.logger.Error("Enode is nil in sendIstAnnounce")
	// 	return nil, nil
	// }
    //
	// enodeUrl := selfEnode.String()
	// view := sb.core.CurrentView()
	// incompleteEnodeUrl := enodeUrl[:strings.Index(enodeUrl, "@")]
	// endpointData := enodeUrl[strings.Index(enodeUrl, "@"):]
    //
	// regAndActiveVals, err := validators.RetrieveRegisteredValidators(nil, nil)
	// // The validator contract may not be deployed yet.
	// // Even if it is deployed, it may not have any registered validators yet.
	// if err == contract_errors.ErrSmartContractNotDeployed || len(regAndActiveVals) == 0 {
	// 	sb.logger.Trace("Can't retrieve the registered validators.  Only allowing the initial validator set to send announce messages", "err", err, "regAndActiveVals", regAndActiveVals)
	// 	regAndActiveVals = make(map[common.Address]bool)
	// } else if err != nil {
	// 	sb.logger.Error("Error in retrieving the registered validators", "err", err)
	// 	return nil, err
	// }
    //
	// // Add active validators regardless
	// valSet := sb.getValidators(block.Number().Uint64(), block.Hash())
	// for _, val := range valSet.List() {
	// 	regAndActiveVals[val.Address()] = true
	// }
    //
	// encryptedEndpoints := make([][][]byte, 0)
	// for addr := range regAndActiveVals {
	// 	if validatorEnodeEntry, ok := sb.valEnodeTable.valEnodeTable[addr]; ok {
	// 		validatorEnode, err := enode.ParseV4(validatorEnodeEntry.enodeURL)
	// 		pubKey := ecies.ImportECDSAPublic(validatorEnode.Pubkey())
	// 		encryptedEndpoint, err := ecies.Encrypt(rand.Reader, pubKey, []byte(endpointData), nil, nil)
	// 		if err != nil {
	// 			log.Warn("Unable to unmarshal public key", "err", err)
	// 		} else {
	// 			encryptedEndpoints = append(encryptedEndpoints, [][]byte{addr.Bytes(), encryptedEndpoint})
	// 		}
	// 	}
	// }

	// msg := &announceMessage{
	// 	Address:               sb.Address(),
	// 	IncompleteEnodeURL:    incompleteEnodeUrl,
	// 	EncryptedEndpointData: encryptedEndpoints,
	// 	View:                  view,
	// }
    //
    //
	// // Sign the announce message
	// if err := msg.Sign(sb.Sign); err != nil {
	// 	sb.logger.Error("Error in signing an Istanbul Announce Message", "AnnounceMsg", msg.String(), "err", err)
	// 	return nil, err
	// }
    //
    msg := &valEnodeShareMessage{
        TestPayload: "Hey Dude!!!!",
    }

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		sb.logger.Error("Error in converting Validator Enode Share Message to payload", "ValEnodeShareMessage", msg.String(), "err", err)
		return nil, err
	}

	sb.logger.Trace("Broadcasting an announce message", "AnnounceMsg", msg)

	return payload, nil
}

func (sb *Backend) sendValEnodeShareMsg() error {
	payload, err := sb.generateValEnodeShareMsg()
	if err != nil {
		return err
	}

	if payload == nil {
		return nil
	}
	//
	// var hash common.Hash
	// if !ignoreCache {
	// 	hash = istanbul.RLPHash(payload)
	// 	sb.knownMessages.Add(hash, true)
	// }
	//
	// sentries := sb.broadcaster.Get
	//
	// var targets map[common.Address]bool = nil
	//
	// if valSet != nil {
	// 	targets = make(map[common.Address]bool)
	// 	for _, val := range valSet.List() {
	// 		if val.Address() != sb.Address() {
	// 			targets[val.Address()] = true
	// 		}
	// 	}
	// }
	//
	// if sb.broadcaster != nil && ((valSet == nil) || (len(targets) > 0)) {
	// 	ps := sb.broadcaster.FindPeers(targets)
	//
	// 	for addr, p := range ps {
	// 		if !ignoreCache {
	// 			ms, ok := sb.recentMessages.Get(addr)
	// 			var m *lru.ARCCache
	// 			if ok {
	// 				m, _ = ms.(*lru.ARCCache)
	// 				if _, k := m.Get(hash); k {
	// 					// This peer had this event, skip it
	// 					continue
	// 				}
	// 			} else {
	// 				m, _ = lru.NewARC(inmemoryMessages)
	// 			}
	//
	// 			m.Add(hash, true)
	// 			sb.recentMessages.Add(addr, m)
	// 		}
	//
	// 		go p.Send(msgCode, payload)
	// 	}
	// }
	// return nil
	//
	//
	sentryPeers := sb.broadcaster.GetSentryPeers()
	if len(sentryPeers) > 0 {
		sb.logger.Warn("Sentrypeers!!!!!", "sentryPeers", sentryPeers, "len", len(sentryPeers))
		// sentryPeers[0].Send(istanbulValEnodeShareMsg, payload)
		sentryPeers[0].Send(istanbulValEnodeShareMsg, payload)
	} else {
		sb.logger.Warn("No sentry enodes")
	}

	// sb.Gossip(nil, payload, istanbulAnnounceMsg, true)

	return nil
}

func (sb *Backend) handleValEnodeShareMsg(payload []byte) error {
	sb.logger.Warn("Handling an IstanbulAnnounce message")

	msg := new(valEnodeShareMessage)
	// Decode message
	err := msg.FromPayload(payload)
	if err != nil {
		sb.logger.Error("Error in decoding received Istanbul Announce message", "err", err, "payload", hex.EncodeToString(payload))
		return err
	}

	payloadContent, err := msg.Payload()

    sb.logger.Warn("woo! Payload", "payloadContent", payloadContent, "msg.TestPayload", msg.TestPayload, "err", err)

	// // Verify message signature
	// if err := msg.VerifySig(); err != nil {
	// 	sb.logger.Error("Error in verifying the signature of an Istanbul Announce message", "err", err, "AnnounceMsg", msg.String())
	// 	return err
	// // }
    //
	// // If the message is originally from this node, then ignore it
	// if msg.Address == sb.Address() {
	// 	sb.logger.Trace("Received an IstanbulAnnounce message originating from this node. Ignoring it.")
	// 	return nil
	// }
    //
	// // If the message is not within the registered validator set, then ignore it
	// regAndActiveVals, err := validators.RetrieveRegisteredValidators(nil, nil)
    //
	// // The validator contract may not be deployed yet.
	// // Even if it is deployed, it may not have any registered validators yet.
	// if err == contract_errors.ErrSmartContractNotDeployed || len(regAndActiveVals) == 0 {
	// 	sb.logger.Trace("Can't retrieve the registered validators.  Only allowing the initial validator set to send announce messages", "err", err, "regAndActiveVals", regAndActiveVals)
	// 	regAndActiveVals = make(map[common.Address]bool)
	// } else if err != nil {
	// 	sb.logger.Error("Error in retrieving the registered validators", "err", err)
	// 	return err
	// }
    //
	// // Add active validators regardless
	// block := sb.currentBlock()
	// valSet := sb.getValidators(block.Number().Uint64(), block.Hash())
	// for _, val := range valSet.List() {
	// 	regAndActiveVals[val.Address()] = true
	// }
    //
	// if !regAndActiveVals[msg.Address] {
	// 	sb.logger.Warn("Received an IstanbulAnnounce message from a non registered validator. Ignoring it.", "AnnounceMsg", msg.String(), "validators", regAndActiveVals, "err", err)
	// 	return errUnauthorizedAnnounceMessage
	// }
    //
	// // Decrypt the EnodeURL
	// nodeKey := ecies.ImportECDSA(sb.GetNodeKey())
    //
	// encryptedEndpoint := []byte("")
	// for _, entry := range msg.EncryptedEndpointData {
	// 	if bytes.Equal(entry[0], sb.Address().Bytes()) {
	// 		encryptedEndpoint = entry[1]
	// 	}
	// }
	// endpointBytes, err := nodeKey.Decrypt(encryptedEndpoint, nil, nil)
	// if err != nil && len(encryptedEndpoint) > 0 {
	// 	sb.logger.Warn("Error in decrypting endpoint", "err", err, "encryptedEndpoint", encryptedEndpoint)
	// }
	// enodeUrl := msg.IncompleteEnodeURL + string(endpointBytes)
    //
	// // Save in the valEnodeTable if mining
	// if sb.coreStarted {
	// 	block := sb.currentBlock()
	// 	valSet := sb.getValidators(block.Number().Uint64(), block.Hash())
    //
	// 	newValEnode := &validatorEnode{enodeURL: enodeUrl, view: msg.View}
	// 	if err := sb.valEnodeTable.upsert(msg.Address, newValEnode, valSet, sb.Address()); err != nil {
	// 		sb.logger.Warn("Error in upserting a valenode entry", "AnnounceMsg", msg, "error", err)
	// 		return err
	// 	}
	// }
    //
	// // If we gossiped this address/enodeURL within the last 60 seconds, then don't regossip
	// sb.lastAnnounceGossipedMu.RLock()
	// if lastGossipTs, ok := sb.lastAnnounceGossiped[msg.Address]; ok {
	// 	if lastGossipTs.enodeURL == enodeUrl && time.Since(lastGossipTs.timestamp) < time.Minute {
	// 		sb.logger.Trace("Already regossiped the msg within the last minute, so not regossiping.", "AnnounceMsg", msg)
	// 		sb.lastAnnounceGossipedMu.RUnlock()
	// 		return nil
	// 	}
	// }
	// sb.lastAnnounceGossipedMu.RUnlock()
    //
	// sb.logger.Trace("Regossiping the istanbul announce message", "AnnounceMsg", msg)
	// sb.Gossip(nil, payload, istanbulAnnounceMsg, true)
    //
	// sb.lastAnnounceGossipedMu.Lock()
	// defer sb.lastAnnounceGossipedMu.Unlock()
	// sb.lastAnnounceGossiped[msg.Address] = &AnnounceGossipTimestamp{enodeURL: enodeUrl, timestamp: time.Now()}
    //
	// // prune non registered validator entries in the valEnodeTable, reverseValEnodeTable, and lastAnnounceGossiped tables about 5% of the times that an announce msg is handled
	// if (mrand.Int() % 100) <= 5 {
	// 	for remoteAddress := range sb.lastAnnounceGossiped {
	// 		if !regAndActiveVals[remoteAddress] {
	// 			log.Trace("Deleting entry from the lastAnnounceGossiped table", "address", remoteAddress, "gossip timestamp", sb.lastAnnounceGossiped[remoteAddress])
	// 			delete(sb.lastAnnounceGossiped, remoteAddress)
	// 		}
	// 	}
    //
	// 	sb.valEnodeTable.pruneEntries(regAndActiveVals)
	// }

	return nil
}
