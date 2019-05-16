package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	cmdOpts := getCmdOpts()

	fmt.Println("Type the *BIP39* mnemonic seed. End with an empty line.")
	bipSeed := readSeedInput()

	pass := ""
	if cmdOpts.ReadPass {
		pass = readPass()
	}

	bts, err := bip39.NewSeedWithErrorChecking(bipSeed, pass)
	defer zeroBytes(bts)
	orPanic(err)

	fmt.Printf("Resulting *DECRED* hex seed:\n%s\n\n", hex.EncodeToString(bts))

	if !cmdOpts.ShowAccountKey && !cmdOpts.ShowAddresses {
		os.Exit(0)
	}

	chainParams := &chaincfg.MainNetParams
	if cmdOpts.TestNet {
		chainParams = &chaincfg.TestNet3Params
	}

	coinType := chainParams.SLIP0044CoinType
	account := uint32(0)
	branch := uint32(0)
	maxKeys := 5
	purpose := uint32(44)

	rootKey, err := hdkeychain.NewMaster(bts, chainParams)
	orPanic(err)

	purposeKey, err := rootKey.Child(purpose + hdkeychain.HardenedKeyStart)
	orPanic(err)

	coinTypeKey, err := purposeKey.Child(coinType + hdkeychain.HardenedKeyStart)
	orPanic(err)

	accountKey, err := coinTypeKey.Child(account + hdkeychain.HardenedKeyStart)
	orPanic(err)
	accountMasterPub, err := accountKey.Neuter()
	orPanic(err)

	if cmdOpts.ShowAccountKey {
		fmt.Printf("Account %d MasterPubKey:\n%s\n\n", account, accountMasterPub.String())
	}

	if cmdOpts.ShowAddresses {
		branchKey, err := accountKey.Child(branch)
		orPanic(err)

		for i := 0; i < maxKeys; i++ {
			addressKey, err := branchKey.Child(uint32(i))
			orPanic(err)

			addr, err := addressKey.Address(chainParams)
			orPanic(err)
			fmt.Printf("Address %d: %s\n", i, addr.EncodeAddress())
		}
	}
}
