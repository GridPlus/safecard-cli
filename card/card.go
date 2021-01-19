package card

import (
	"errors"
	"fmt"

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

func OpenSecureConnection() (*safecard.CommandSet, error) {
	cs, err := Connect()
	if err != nil {
		fmt.Println("error connecting to card")
		fmt.Println(err)
		return cs, err
	}
	err = cs.Select()
	if err != nil {
		fmt.Println("error selecting applet. err: ", err)
		return cs, err
	}
	err = cs.Pair()
	if err != nil {
		fmt.Println("error pairing with card. err: ", err)
		return cs, err
	}
	err = cs.OpenSecureChannel()
	if err != nil {
		fmt.Println("error opening secure channel. err: ", err)
		return cs, err
	}
	return cs, nil
}
