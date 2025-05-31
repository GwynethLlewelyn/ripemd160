# RIPEMD160 Test Vectors

The table below shows the standard test as provided by [Antoon Bosselaers] on [The RIPEMD-160 page]. The messages are given in ASCII format, while the corresponding hash results are in hexadecimal format.

| Message               | Hash result using RIPEMD-160               |
| --------------------- | ------------------------------------------ |
| "" (empty string)     | `9c1185a5c5e9fc54612808977ee8f548b2258d31` |
| "a"                   | `0bdc9d2d256b3ee9daae347be6f4dc835a467ffe` |
| "abc"                 | `8eb208f7e05d987a9b044a8e98c6b087f15a0bfc` |
| "message digest"      | `5d0689ef49d2fae572b881b123a85ffa21595f36` |
| "a...z"[^1]           | `f71c27109c692c1b56bbdceb5b9d2865b3708dbc` |
| "abcdbcde...nopq"[^2] | `12a053384a9c0c88e405a06c27dcf49ada62eb2b` |
| "A...Za...z0...9"[^3] | `b0e20b6e3116640286ed3a87a5713079b21f5189` |
| 8 times "1234567890"  | `9b752e45573d4b39f4dbd3323cab82bf63326bfb` |
| 1 million times "a"   | `52783243c1697bdbe16d37f97f68f08325dc1528` |

[^1]: "abcdefghijklmnopqrstuvwxyz"

[^2]: "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"

[^3]: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

`go test` will correctly run all the above and compute them successfully.

[Antoon Bosselaers]: https://homes.esat.kuleuven.be/~bosselae/
[The RIPEMD-160 page]: https://homes.esat.kuleuven.be/~bosselae/ripemd160.html
