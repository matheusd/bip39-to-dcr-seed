# BIP39 to Decred Seed

This is a small tool to convert
[BIP39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) mnemonic
seeds (such as those generated by Trezor devices) into hex seeds importable in
standard [Decred](https://github.com/decred/dcrd) software.

It supports the use of password-derived BIP39 seeds.

Due to how Decred seeds currently work, the resulting hex seed has double the
size of standard decred seeds (33 word mnemonic seeds), therefore it's not
currently possible to generate them (the mnemonics based on the PGP wordlist)
from the BIP39 source seed.

## Running

Go 1.12+ with modules support required.

```
git clone https://github.com/matheusd/bip39-to-dcr-seed
cd bip39-to-dcr-seed
go run . --help
```
