package crypto

import "testing"
import "bytes"

func Test_Open(t *testing.T) {
	var aead AEAD
	var err error

	// Test Vectors taken from RFC7539 Annexe 5 : http://tools.ietf.org/html/rfc7539

	/*
	   The key:
	   000  1c 92 40 a5 eb 55 d3 8a f3 33 88 86 04 f6 b5 f0  ..@..U...3......
	   016  47 39 17 c1 40 2b 80 09 9d ca 5c bc 20 70 75 c0  G9..@+....\. pu.

	   Ciphertext:
	   000  64 a0 86 15 75 86 1a f4 60 f0 62 c7 9b e6 43 bd  d...u...`.b...C.
	   016  5e 80 5c fd 34 5c f3 89 f1 08 67 0a c7 6c 8c b2  ^.\.4\....g..l..
	   032  4c 6c fc 18 75 5d 43 ee a0 9e e9 4e 38 2d 26 b0  Ll..u]C....N8-&.
	   048  bd b7 b7 3c 32 1b 01 00 d4 f0 3b 7f 35 58 94 cf  ...<2.....;.5X..
	   064  33 2f 83 0e 71 0b 97 ce 98 c8 a8 4a bd 0b 94 81  3/..q......J....
	   080  14 ad 17 6e 00 8d 33 bd 60 f9 82 b1 ff 37 c8 55  ...n..3.`....7.U
	   096  97 97 a0 6e f4 f0 ef 61 c1 86 32 4e 2b 35 06 38  ...n...a..2N+5.8
	   112  36 06 90 7b 6a 7c 02 b0 f9 f6 15 7b 53 c8 67 e4  6..{j|.....{S.g.
	   128  b9 16 6c 76 7b 80 4d 46 a5 9b 52 16 cd e7 a4 e9  ..lv{.MF..R.....
	   144  90 40 c5 a4 04 33 22 5e e2 82 a1 b0 a0 6c 52 3e  .@...3"^.....lR>
	   160  af 45 34 d7 f8 3f a1 15 5b 00 47 71 8c bc 54 6a  .E4..?..[.Gq..Tj
	   176  0d 07 2b 04 b3 56 4e ea 1b 42 22 73 f5 48 27 1a  ..+..VN..B"s.H'.
	   192  0b b2 31 60 53 fa 76 99 19 55 eb d6 31 59 43 4e  ..1`S.v..U..1YCN
	   208  ce bb 4e 46 6d ae 5a 10 73 a6 72 76 27 09 7a 10  ..NFm.Z.s.rv'.z.
	   224  49 e6 17 d9 1d 36 10 94 fa 68 f0 ff 77 98 71 30  I....6...h..w.q0
	   240  30 5b ea ba 2e da 04 df 99 7b 71 4d 6c 6f 2c 29  0[.......{qMlo,)
	   256  a6 ad 5c b4 02 2b 02 70 9b                       ..\..+.p.

	   The nonce:
	   000  00 00 00 00 01 02 03 04 05 06 07 08              ............

	   The AAD:
	   000  f3 33 88 86 00 00 00 00 00 00 4e 91              .3........N.

	   Received Tag:
	   000  ee ad 9d 67 89 0c bb 22 39 23 36 fe a1 85 1f 38  ...g..."9#6....8

	   Poly1305 one-time key:
	   000  bd f0 4a a9 5c e4 de 89 95 b1 4b b6 a1 8f ec af  ..J.\.....K.....
	   016  26 47 8f 50 c0 54 f5 63 db c0 a2 1e 26 15 72 aa  &G.P.T.c....&.r.

	    Next, we construct the AEAD buffer

	   Poly1305 Input:
	   000  f3 33 88 86 00 00 00 00 00 00 4e 91 00 00 00 00  .3........N.....
	   016  64 a0 86 15 75 86 1a f4 60 f0 62 c7 9b e6 43 bd  d...u...`.b...C.
	   032  5e 80 5c fd 34 5c f3 89 f1 08 67 0a c7 6c 8c b2  ^.\.4\....g..l..
	   048  4c 6c fc 18 75 5d 43 ee a0 9e e9 4e 38 2d 26 b0  Ll..u]C....N8-&.
	   064  bd b7 b7 3c 32 1b 01 00 d4 f0 3b 7f 35 58 94 cf  ...<2.....;.5X..
	   080  33 2f 83 0e 71 0b 97 ce 98 c8 a8 4a bd 0b 94 81  3/..q......J....
	   096  14 ad 17 6e 00 8d 33 bd 60 f9 82 b1 ff 37 c8 55  ...n..3.`....7.U
	   112  97 97 a0 6e f4 f0 ef 61 c1 86 32 4e 2b 35 06 38  ...n...a..2N+5.8
	   128  36 06 90 7b 6a 7c 02 b0 f9 f6 15 7b 53 c8 67 e4  6..{j|.....{S.g.
	   144  b9 16 6c 76 7b 80 4d 46 a5 9b 52 16 cd e7 a4 e9  ..lv{.MF..R.....
	   160  90 40 c5 a4 04 33 22 5e e2 82 a1 b0 a0 6c 52 3e  .@...3"^.....lR>
	   176  af 45 34 d7 f8 3f a1 15 5b 00 47 71 8c bc 54 6a  .E4..?..[.Gq..Tj
	   192  0d 07 2b 04 b3 56 4e ea 1b 42 22 73 f5 48 27 1a  ..+..VN..B"s.H'.
	   208  0b b2 31 60 53 fa 76 99 19 55 eb d6 31 59 43 4e  ..1`S.v..U..1YCN
	   224  ce bb 4e 46 6d ae 5a 10 73 a6 72 76 27 09 7a 10  ..NFm.Z.s.rv'.z.
	   240  49 e6 17 d9 1d 36 10 94 fa 68 f0 ff 77 98 71 30  I....6...h..w.q0
	   256  30 5b ea ba 2e da 04 df 99 7b 71 4d 6c 6f 2c 29  0[.......{qMlo,)
	   272  a6 ad 5c b4 02 2b 02 70 9b 00 00 00 00 00 00 00  ..\..+.p........
	   288  0c 00 00 00 00 00 00 00 09 01 00 00 00 00 00 00  ................

	   Calculated Tag:
	   000  ee ad 9d 67 89 0c bb 22 39 23 36 fe a1 85 1f 38  ...g..."9#6....8

	    Finally, we decrypt the ciphertext

	   Plaintext::
	   000  49 6e 74 65 72 6e 65 74 2d 44 72 61 66 74 73 20  Internet-Drafts
	   016  61 72 65 20 64 72 61 66 74 20 64 6f 63 75 6d 65  are draft docume
	   032  6e 74 73 20 76 61 6c 69 64 20 66 6f 72 20 61 20  nts valid for a
	   048  6d 61 78 69 6d 75 6d 20 6f 66 20 73 69 78 20 6d  maximum of six m
	   064  6f 6e 74 68 73 20 61 6e 64 20 6d 61 79 20 62 65  onths and may be
	   080  20 75 70 64 61 74 65 64 2c 20 72 65 70 6c 61 63   updated, replac
	   096  65 64 2c 20 6f 72 20 6f 62 73 6f 6c 65 74 65 64  ed, or obsoleted
	   112  20 62 79 20 6f 74 68 65 72 20 64 6f 63 75 6d 65   by other docume
	   128  6e 74 73 20 61 74 20 61 6e 79 20 74 69 6d 65 2e  nts at any time.
	   144  20 49 74 20 69 73 20 69 6e 61 70 70 72 6f 70 72   It is inappropr
	   160  69 61 74 65 20 74 6f 20 75 73 65 20 49 6e 74 65  iate to use Inte
	   176  72 6e 65 74 2d 44 72 61 66 74 73 20 61 73 20 72  rnet-Drafts as r
	   192  65 66 65 72 65 6e 63 65 20 6d 61 74 65 72 69 61  eference materia
	   208  6c 20 6f 72 20 74 6f 20 63 69 74 65 20 74 68 65  l or to cite the
	   224  6d 20 6f 74 68 65 72 20 74 68 61 6e 20 61 73 20  m other than as
	   240  2f e2 80 9c 77 6f 72 6b 20 69 6e 20 70 72 6f 67  /...work in prog
	   256  72 65 73 73 2e 2f e2 80 9d                       ress./...
	*/
	key := []byte{
		0x1c, 0x92, 0x40, 0xa5, 0xeb, 0x55, 0xd3, 0x8a, 0xf3, 0x33, 0x88, 0x86, 0x04, 0xf6, 0xb5, 0xf0,
		0x47, 0x39, 0x17, 0xc1, 0x40, 0x2b, 0x80, 0x09, 0x9d, 0xca, 0x5c, 0xbc, 0x20, 0x70, 0x75, 0xc0}
	noncePrefix := []byte{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8}
	aad := []byte{0xf3, 0x33, 0x88, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4e, 0x91}
	cipherText := []byte{
		0x64, 0xa0, 0x86, 0x15, 0x75, 0x86, 0x1a, 0xf4, 0x60, 0xf0, 0x62, 0xc7, 0x9b, 0xe6, 0x43, 0xbd,
		0x5e, 0x80, 0x5c, 0xfd, 0x34, 0x5c, 0xf3, 0x89, 0xf1, 0x08, 0x67, 0x0a, 0xc7, 0x6c, 0x8c, 0xb2,
		0x4c, 0x6c, 0xfc, 0x18, 0x75, 0x5d, 0x43, 0xee, 0xa0, 0x9e, 0xe9, 0x4e, 0x38, 0x2d, 0x26, 0xb0,
		0xbd, 0xb7, 0xb7, 0x3c, 0x32, 0x1b, 0x01, 0x00, 0xd4, 0xf0, 0x3b, 0x7f, 0x35, 0x58, 0x94, 0xcf,
		0x33, 0x2f, 0x83, 0x0e, 0x71, 0x0b, 0x97, 0xce, 0x98, 0xc8, 0xa8, 0x4a, 0xbd, 0x0b, 0x94, 0x81,
		0x14, 0xad, 0x17, 0x6e, 0x00, 0x8d, 0x33, 0xbd, 0x60, 0xf9, 0x82, 0xb1, 0xff, 0x37, 0xc8, 0x55,
		0x97, 0x97, 0xa0, 0x6e, 0xf4, 0xf0, 0xef, 0x61, 0xc1, 0x86, 0x32, 0x4e, 0x2b, 0x35, 0x06, 0x38,
		0x36, 0x06, 0x90, 0x7b, 0x6a, 0x7c, 0x02, 0xb0, 0xf9, 0xf6, 0x15, 0x7b, 0x53, 0xc8, 0x67, 0xe4,
		0xb9, 0x16, 0x6c, 0x76, 0x7b, 0x80, 0x4d, 0x46, 0xa5, 0x9b, 0x52, 0x16, 0xcd, 0xe7, 0xa4, 0xe9,
		0x90, 0x40, 0xc5, 0xa4, 0x04, 0x33, 0x22, 0x5e, 0xe2, 0x82, 0xa1, 0xb0, 0xa0, 0x6c, 0x52, 0x3e,
		0xaf, 0x45, 0x34, 0xd7, 0xf8, 0x3f, 0xa1, 0x15, 0x5b, 0x00, 0x47, 0x71, 0x8c, 0xbc, 0x54, 0x6a,
		0x0d, 0x07, 0x2b, 0x04, 0xb3, 0x56, 0x4e, 0xea, 0x1b, 0x42, 0x22, 0x73, 0xf5, 0x48, 0x27, 0x1a,
		0x0b, 0xb2, 0x31, 0x60, 0x53, 0xfa, 0x76, 0x99, 0x19, 0x55, 0xeb, 0xd6, 0x31, 0x59, 0x43, 0x4e,
		0xce, 0xbb, 0x4e, 0x46, 0x6d, 0xae, 0x5a, 0x10, 0x73, 0xa6, 0x72, 0x76, 0x27, 0x09, 0x7a, 0x10,
		0x49, 0xe6, 0x17, 0xd9, 0x1d, 0x36, 0x10, 0x94, 0xfa, 0x68, 0xf0, 0xff, 0x77, 0x98, 0x71, 0x30,
		0x30, 0x5b, 0xea, 0xba, 0x2e, 0xda, 0x04, 0xdf, 0x99, 0x7b, 0x71, 0x4d, 0x6c, 0x6f, 0x2c, 0x29,
		0xa6, 0xad, 0x5c, 0xb4, 0x02, 0x2b, 0x02, 0x70, 0x9b, 0xee, 0xad, 0x9d, 0x67, 0x89, 0x0c, 0xbb,
		0x22, 0x39, 0x23, 0x36, 0xfe}
	if aead, err = NewAEAD_ChaCha20Poly1305(key, noncePrefix); err != nil {
		t.Error(err)
	}
	if err = aead.Open(0, nil, aad, cipherText); err != nil {
		t.Error(err)
	}
}

func Test_KeyGeneratorPoly1305(t *testing.T) {
	var buf [64]byte
	var cipher *ChaCha20Cipher
	var err error

	// Poly1305 Key Generation Using ChaCha20 Test Vectors taken from RFC7539 : http://tools.ietf.org/html/rfc7539

	key := []byte{0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f}
	noncePrefix := []byte{0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7}

	if cipher, err = NewChaCha20Cipher(key, noncePrefix, 0); err != nil {
		t.Error("Key Generator test for Poly1305 : error when calling NewChaCha20Cipher")
	}
	cipher.GetNextKeystream(&buf)
	if !bytes.Equal(buf[:32], []byte{
		0x8a, 0xd5, 0xa0, 0x8b, 0x90, 0x5f, 0x81, 0xcc, 0x81, 0x50, 0x40, 0x27, 0x4a, 0xb2, 0x94, 0x71,
		0xa8, 0x33, 0xb6, 0x37, 0xe3, 0xfd, 0x0d, 0xa5, 0x08, 0xdb, 0xb8, 0xe2, 0xfd, 0xd1, 0xa6, 0x46}) {
		t.Error("Key Generator test for Poly1305 : bad Poly1305 Key Generation test vector")
	}

	/*
	   Test Vector #1:
	   ==============

	   The key:
	   000  00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	   016  00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................

	   The nonce:
	   000  00 00 00 00 00 00 00 00 00 00 00 00              ............

	   Poly1305 one-time key:
	   000  76 b8 e0 ad a0 f1 3d 90 40 5d 6a e5 53 86 bd 28  v.....=.@]j.S..(
	   016  bd d2 19 b8 a0 8d ed 1a a8 36 ef cc 8b 77 0d c7  .........6...w..
	*/
	key = make([]byte, 32)
	noncePrefix = make([]byte, 12)
	if cipher, err = NewChaCha20Cipher(key, noncePrefix, 0); err != nil {
		t.Error("Key Generator test for Poly1305 : error when calling NewChaCha20Cipher")
	}
	cipher.GetNextKeystream(&buf)
	if !bytes.Equal(buf[:32], []byte{
		0x76, 0xb8, 0xe0, 0xad, 0xa0, 0xf1, 0x3d, 0x90, 0x40, 0x5d, 0x6a, 0xe5, 0x53, 0x86, 0xbd, 0x28,
		0xbd, 0xd2, 0x19, 0xb8, 0xa0, 0x8d, 0xed, 0x1a, 0xa8, 0x36, 0xef, 0xcc, 0x8b, 0x77, 0x0d, 0xc7}) {
		t.Error("Key Generator test for Poly1305 : bad Poly1305 Key Generation test vector")
	}

	/*
	   Test Vector #2:
	   ==============

	   The key:
	   000  00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	   016  00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01  ................

	   The nonce:
	   000  00 00 00 00 00 00 00 00 00 00 00 02              ............

	   Poly1305 one-time key:
	   000  ec fa 25 4f 84 5f 64 74 73 d3 cb 14 0d a9 e8 76  ..%O._dts......v
	   016  06 cb 33 06 6c 44 7b 87 bc 26 66 dd e3 fb b7 39  ..3.lD{..&f....9
	*/
	key = []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	noncePrefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	if cipher, err = NewChaCha20Cipher(key, noncePrefix, 0); err != nil {
		t.Error("Key Generator test for Poly1305 : error when calling NewChaCha20Cipher")
	}
	cipher.GetNextKeystream(&buf)
	if !bytes.Equal(buf[:32], []byte{
		0xec, 0xfa, 0x25, 0x4f, 0x84, 0x5f, 0x64, 0x74, 0x73, 0xd3, 0xcb, 0x14, 0x0d, 0xa9, 0xe8, 0x76,
		0x06, 0xcb, 0x33, 0x06, 0x6c, 0x44, 0x7b, 0x87, 0xbc, 0x26, 0x66, 0xdd, 0xe3, 0xfb, 0xb7, 0x39}) {
		t.Error("Key Generator test for Poly1305 : bad Poly1305 Key Generation test vector")
	}

	/*
	   Test Vector #3:
	   ==============

	   The key:
	   000  1c 92 40 a5 eb 55 d3 8a f3 33 88 86 04 f6 b5 f0  ..@..U...3......
	   016  47 39 17 c1 40 2b 80 09 9d ca 5c bc 20 70 75 c0  G9..@+....\. pu.

	   The nonce:
	   000  00 00 00 00 00 00 00 00 00 00 00 02              ............

	   Poly1305 one-time key:
	   000  96 5e 3b c6 f9 ec 7e d9 56 08 08 f4 d2 29 f9 4b  .^;...~.V....).K
	   016  13 7f f2 75 ca 9b 3f cb dd 59 de aa d2 33 10 ae  ...u..?..Y...3..
	*/
	key = []byte{
		0x1c, 0x92, 0x40, 0xa5, 0xeb, 0x55, 0xd3, 0x8a, 0xf3, 0x33, 0x88, 0x86, 0x04, 0xf6, 0xb5, 0xf0,
		0x47, 0x39, 0x17, 0xc1, 0x40, 0x2b, 0x80, 0x09, 0x9d, 0xca, 0x5c, 0xbc, 0x20, 0x70, 0x75, 0xc0}
	noncePrefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}
	if cipher, err = NewChaCha20Cipher(key, noncePrefix, 0); err != nil {
		t.Error("Key Generator test for Poly1305 : error when calling NewChaCha20Cipher")
	}
	cipher.GetNextKeystream(&buf)
	if !bytes.Equal(buf[:32], []byte{
		0x96, 0x5e, 0x3b, 0xc6, 0xf9, 0xec, 0x7e, 0xd9, 0x56, 0x08, 0x08, 0xf4, 0xd2, 0x29, 0xf9, 0x4b,
		0x13, 0x7f, 0xf2, 0x75, 0xca, 0x9b, 0x3f, 0xcb, 0xdd, 0x59, 0xde, 0xaa, 0xd2, 0x33, 0x10, 0xae}) {
		t.Error("Key Generator test for Poly1305 : bad Poly1305 Key Generation test vector")
	}

}
