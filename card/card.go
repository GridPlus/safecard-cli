package card

import (
	"errors"

	safecard "github.com/GridPlus/keycard-go"
	"github.com/GridPlus/keycard-go/io"
	"github.com/ebfe/scard"
)

func Connect() (*safecard.CommandSet, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}
	readers, err := ctx.ListReaders()
	if err != nil {
		return nil, err
	}

	if len(readers) > 0 {
		card, err := ctx.Connect(readers[0], scard.ShareShared, scard.ProtocolAny)
		if err != nil {
			return nil, err
		}
		return safecard.NewCommandSet(io.NewNormalChannel(card)), nil
	}
	return nil, errors.New("no card reader found")
}
