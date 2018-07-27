package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"

	flags "github.com/btcsuite/go-flags"
	bip39 "github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// collapseSpace takes a string and replaces any repeated areas of whitespace
// with a single space character.
func collapseSpace(in string) string {
	whiteSpace := false
	out := ""
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !whiteSpace {
				out = out + " "
			}
			whiteSpace = true
		} else {
			out = out + string(c)
			whiteSpace = false
		}
	}
	return out
}

func readSeedInput() string {
	var seedStr string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		seedStr += " " + line
	}

	seedStrTrimmed := strings.TrimSpace(seedStr)
	seedStrTrimmed = collapseSpace(seedStrTrimmed)

	_, err := bip39.MnemonicToByteArray(seedStrTrimmed)
	orPanic(err)

	return seedStrTrimmed
}

func readPass() string {
	for {
		fmt.Print("Type the *BIP39* password: ")
		var pass []byte
		var err error
		pass, err = terminal.ReadPassword(int(os.Stdin.Fd()))
		orPanic(err)
		fmt.Print("\n")

		pass = bytes.TrimSpace(pass)
		if len(pass) == 0 {
			return ""
		}

		fmt.Print("Confirm the password: ")
		confirm, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		orPanic(err)
		fmt.Print("\n")

		confirm = bytes.TrimSpace(confirm)
		if !bytes.Equal(pass, confirm) {
			fmt.Println("The typed passwords do not match")
			continue
		}

		return string(pass)
	}
}

func zeroBytes(bts []byte) {
	for i := range bts {
		bts[i] = 0
	}
}

type opts struct {
	ReadPass       bool `long:"readpass" description:"Whether to read a bip39 password for seed derivation"`
	ShowAccountKey bool `long:"showaccountkey" description:"Whether to show the default account master pub key"`
	ShowAddresses  bool `long:"showaddresses" description:"Whether to show the first 5 external addresses of the default account key"`
	TestNet        bool `long:"testnet" description:"Whether to derive public key and addresses for testnet"`
}

func getCmdOpts() *opts {
	cmdOpts := &opts{}
	parser := flags.NewParser(cmdOpts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		orPanic(err)
	}

	return cmdOpts
}
