package nodecommon

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"
)

func Decode(encoded common.EncodedRollup) (*Rollup, error) {
	r := Rollup{}
	err := rlp.DecodeBytes(encoded, &r)

	return &r, err
}

func EncodeRollup(r *Rollup) common.EncodedRollup {
	encoded, err := r.encode()
	if err != nil {
		panic(err)
	}

	return encoded
}

func DecodeRollup(rollup common.EncodedRollup) *Rollup {
	r, err := Decode(rollup)
	if err != nil {
		panic(err)
	}

	return r
}

func (r Rollup) encode() (common.EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}
