// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"errors"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// IstanbulDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from Istanbul consensus engine
	IstanbulDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	IstanbulExtraVanity        = 32                       // Fixed number of extra-data bytes reserved for validator vanity
	IstanbulExtraCommittedSeal = blscrypto.SIGNATUREBYTES // Fixed number of extra-data bytes reserved for validator seal
	IstanbulExtraSeal          = 65                       // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidIstanbulHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidIstanbulHeaderExtra = errors.New("invalid istanbul header extra-data")
)

type IstanbulSeal struct {
	// Bitmap is a bitmap having an active bit for each validator that signed this block
	Bitmap *big.Int
	// Signature is an aggregated BLS signature resulting from signatures by each validator that signed this block
	Signature []byte
	// Round is the round in which the signature was created.
	Round *big.Int
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *IstanbulSeal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.Bitmap,
		ist.Signature,
		ist.Round,
	})
}

// DecodeRLP implements rlp.Decoder, and load the istanbul fields from a RLP stream.
func (ist *IstanbulSeal) DecodeRLP(s *rlp.Stream) error {
	var istanbulSeal struct {
		Bitmap    *big.Int
		Signature []byte
		Round     *big.Int
	}
	if err := s.Decode(&istanbulSeal); err != nil {
		return err
	}
	ist.Bitmap, ist.Signature, ist.Round = istanbulSeal.Bitmap, istanbulSeal.Signature, istanbulSeal.Round
	return nil
}

type IstanbulExtra struct {
	// AddedValidators are the validators that have been added in the block
	AddedValidators []common.Address
	// AddedValidatorsPublicKeys are the BLS public keys for the validators added in the block
	AddedValidatorsPublicKeys [][]byte
	// RemovedValidators is a bitmap having an active bit for each removed validator in the block
	RemovedValidators *big.Int
	// Seal is an ECDSA signature by the proposer
	Seal []byte
	// CommittedSeal contains the aggregated BLS signature created via IBFT consensus.
	CommittedSeal IstanbulSeal
	// EpochData is a SNARK-friendly encoding of the validator set diff (WIP)
	EpochData []byte
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *IstanbulExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.AddedValidators,
		ist.AddedValidatorsPublicKeys,
		ist.RemovedValidators,
		ist.Seal,
		ist.CommittedSeal,
		ist.EpochData,
	})
}

// DecodeRLP implements rlp.Decoder, and load the istanbul fields from a RLP stream.
func (ist *IstanbulExtra) DecodeRLP(s *rlp.Stream) error {
	var istanbulExtra struct {
		AddedValidators           []common.Address
		AddedValidatorsPublicKeys [][]byte
		RemovedValidators         *big.Int
		Seal                      []byte
		CommittedSeal             IstanbulSeal
		EpochData                 []byte
	}
	if err := s.Decode(&istanbulExtra); err != nil {
		return err
	}
	ist.AddedValidators, ist.AddedValidatorsPublicKeys, ist.RemovedValidators, ist.Seal, ist.CommittedSeal, ist.EpochData = istanbulExtra.AddedValidators, istanbulExtra.AddedValidatorsPublicKeys, istanbulExtra.RemovedValidators, istanbulExtra.Seal, istanbulExtra.CommittedSeal, istanbulExtra.EpochData
	return nil
}

// ExtractIstanbulExtra extracts all values of the IstanbulExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractIstanbulExtra(h *Header) (*IstanbulExtra, error) {
	if len(h.Extra) < IstanbulExtraVanity {
		return nil, ErrInvalidIstanbulHeaderExtra
	}

	var istanbulExtra *IstanbulExtra
	err := rlp.DecodeBytes(h.Extra[IstanbulExtraVanity:], &istanbulExtra)
	if err != nil {
		return nil, err
	}
	return istanbulExtra, nil
}

// IstanbulFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Istanbul hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func IstanbulFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	istanbulExtra, err := ExtractIstanbulExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		istanbulExtra.Seal = []byte{}
	}
	istanbulExtra.CommittedSeal = IstanbulSeal{}

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:IstanbulExtraVanity], payload...)

	return newHeader
}
