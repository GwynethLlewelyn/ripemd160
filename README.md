![RIPE MD160 Logo]
# RIPE MD160 implementation in Go

RIPEMD-160 is designed by Hans Dobbertin, [Antoon Bosselaers], and Bart Preneel with specifications available at [The RIPEMD-160 page].

This implementation is by by [Henrik Hautakoski], deployed on [antelope-go].

You can compare this implementation with its [pseudocode].

`go test` will successfully run the [standard test vectors] through the RIPEMD-160 algorithm.

A very simple [command-line interface] can be found under the `cmd/` folder, loosely inspired by [IBM]'s own `rmd160` tool (and similar to other commands such as `sha256sum`).

## Bibliography

As provided by [Antoon Bosselaers] on [The RIPEMD-160 page].

1. H. Dobbertin, A. Bosselaers, B. Preneel, "[RIPEMD-160, a strengthened version of RIPEMD]," Fast Software Encryption, LNCS 1039, D. Gollmann, Ed., Springer-Verlag, 1996, pp. 71-82.
2. H. Dobbertin, "Digitale Fingerabdrücke; Sichere Hashfunktionen für digitale Signaturen," Datenschutz und Datensicherheit, Vol. 21, No. 2, 1997, pp. 82-87.
3. [ISO/IEC 10118-3:2018], "Information technology - Security techniques - Hash-functions - Part 3: Dedicated hash-functions," International Organization for Standardization, Geneva, Switzerland, 2018.
4. A. Menezes, P. van Oorschot, S. Vanstone, *[Handbook of Applied Cryptography]*, CRC press, 1996, Section 9.4.2, pp. 349-351 (pdf).
5. A. Bosselaers, H. Dobbertin, B. Preneel, "The RIPEMD-160 cryptographic hash function," Dr. Dobb's Journal, Vol. 22, No. 1, January 1997, pp. 24-28.
6. B. Preneel, A. Bosselaers, H. Dobbertin, "[The cryptographic hash function RIPEMD-160]," CryptoBytes, Vol. 3, No. 2, 1997, pp. 9-14.

[![go-test](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/ci.yml/badge.svg)](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/ci.yml) [![CodeQL Advanced](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/codeql.yml/badge.svg)](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/codeql.yml)

[RIPE MD160 Logo]: ./assets/rmd160-logo-small.png
[Antoon Bosselaers]: https://homes.esat.kuleuven.be/~bosselae/
[The RIPEMD-160 page]: https://homes.esat.kuleuven.be/~bosselae/ripemd160.html
[Henrik Hautakoski]: https://github.com/pnx
[antelope-go]: https://github.com/antelope-go/ripemd160
[pseudocode]: pseudocode.md
[standard test vectors]: test-vectors.md
[command-line interface]: ./cmd/README.md
[IBM]: https://www.ibm.com/docs/en/zos/3.1.0?topic=descriptions-rmd160-calculate-check-ripemd-160-cryptographic-hashes
[RIPEMD-160, a strengthened version of RIPEMD]: https://homes.esat.kuleuven.be/~bosselae/ripemd160/pdf/AB-9601/AB-9601.pdf
[ISO/IEC 10118-3:2018]: https://www.iso.org/standard/67116.html
[Handbook of Applied Cryptography]: https://cacr.uwaterloo.ca/hac/
[The cryptographic hash function RIPEMD-160]: https://www.networkdls.com/Articles/crypto3n2.pdf
