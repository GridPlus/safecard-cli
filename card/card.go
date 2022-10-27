package card

import (
	"errors"
	"fmt"

	safecard "github.com/GridPlus/keycard-go"
	"github.com/GridPlus/keycard-go/gridplus"
	"github.com/GridPlus/keycard-go/io"
	"github.com/ebfe/scard"
	"github.com/manifoldco/promptui"
)

func Connect(readerIdx int) (*safecard.CommandSet, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}
	readers, err := ctx.ListReaders()
	if err != nil {
		return nil, err
	}

	if len(readers) <= readerIdx {
		return nil, errors.New(
			fmt.Sprintf(
				"Cannot select card reader (have %d readers, want index=%d)", 
				len(readers), 
				readerIdx,
			),
		)
	}
	card, err := ctx.Connect(readers[readerIdx], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return nil, err
	}
	return safecard.NewCommandSet(io.NewNormalChannel(card)), nil
}

func OpenSecureConnection(readerIdx int) (*safecard.CommandSet, error) {
	cs, err := Connect(readerIdx)
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

func ExportSeed(readerIdx int) ([]byte, error) {
	cs, err := OpenSecureConnection(readerIdx)
	if err != nil {
		fmt.Println("unable to open secure connection with card: ", err)
		return nil, err
	}
	//Prompt user for pin
	prompt := promptui.Prompt{
		Label: "Pin",
		Mask:  '*',
	}
	fmt.Println("Please enter 6 digit pin:")
	result, err := prompt.Run()
	if err != nil {
		fmt.Println("prompt failed: err: ", err)
		return nil, err
	}

	err = cs.VerifyPIN(result)
	if err != nil {
		fmt.Println("error verifying pin. err: ", err)
		return nil, err
	}
	seed, err := cs.ExportSeed()
	if err == gridplus.ErrSeedInvalidLength {
		fmt.Println("card does not appear to have valid exportable seed.")
		return nil, err
	}
	if err != nil {
		fmt.Println("unable to export seed. err: ", err)
		return nil, err
	}
	return seed, nil
}
