// Use of this source code is governed by a license
// that can be found in the LICENSE file.

package threefish

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
)

func convertVector(i int, v *testVector, t *testing.T) ([]byte, []byte, []byte, []byte) {
	key, err := hex.DecodeString(v.key)
	if err != nil {
		t.Fatalf("Test vector %d: Failed to decode hex key - Cause: %s", i, err)
	}
	tweak, err := hex.DecodeString(v.tweak)
	if err != nil {
		t.Fatalf("Test vector %d: Failed to decode hex tweak - Cause: %s", i, err)
	}
	plaintext, err := hex.DecodeString(v.plaintext)
	if err != nil {
		t.Fatalf("Test vector %d: Failed to decode hex plaintext - Cause: %s", i, err)
	}
	ciphertext, err := hex.DecodeString(v.ciphertext)
	if err != nil {
		t.Fatalf("Test vector %d: Failed to decode hex ciphertext - Cause: %s", i, err)
	}
	return key, tweak, plaintext, ciphertext
}

type testVector struct {
	key, tweak, plaintext, ciphertext string
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

// Test vectors from:
// https://github.com/bcgit/bc-java/blob/master/core/src/test/java/org/bouncycastle/crypto/test/Threefish256Test.java
var testVectors256 = []testVector{
	testVector{
		key:        "0000000000000000000000000000000000000000000000000000000000000000",
		tweak:      "00000000000000000000000000000000",
		plaintext:  "0000000000000000000000000000000000000000000000000000000000000000",
		ciphertext: "84da2a1f8beaee947066ae3e3103f1ad536db1f4a1192495116b9f3ce6133fd8",
	},
	testVector{
		key:        "101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f",
		tweak:      "000102030405060708090a0b0c0d0e0f",
		plaintext:  "FFFEFDFCFBFAF9F8F7F6F5F4F3F2F1F0EFEEEDECEBEAE9E8E7E6E5E4E3E2E1E0",
		ciphertext: "e0d091ff0eea8fdfc98192e62ed80ad59d865d08588df476657056b5955e97df",
	},
}

// Test vectors from:
// https://github.com/bcgit/bc-java/blob/master/core/src/test/java/org/bouncycastle/crypto/test/Threefish512Test.java
var testVectors512 = []testVector{
	testVector{
		key: "0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000",
		tweak: "00000000000000000000000000000000",
		plaintext: "0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000",
		ciphertext: "b1a2bbc6ef6025bc40eb3822161f36e375d1bb0aee3186fbd19e47c5d479947b" +
			"7bc2f8586e35f0cff7e7f03084b0b7b1f1ab3961a580a3e97eb41ea14a6d7bbe",
	},
	testVector{
		key: "101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f" +
			"303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f",
		tweak: "000102030405060708090a0b0c0d0e0f",
		plaintext: "fffefdfcfbfaf9f8f7f6f5f4f3f2f1f0efeeedecebeae9e8e7e6e5e4e3e2e1e0" +
			"dfdedddcdbdad9d8d7d6d5d4d3d2d1d0cfcecdcccbcac9c8c7c6c5c4c3c2c1c0",
		ciphertext: "e304439626d45a2cb401cad8d636249a6338330eb06d45dd8b36b90e97254779" +
			"272a0a8d99463504784420ea18c9a725af11dffea10162348927673d5c1caf3d",
	},
}

// Test vectors from:
// https://github.com/bcgit/bc-java/blob/master/core/src/test/java/org/bouncycastle/crypto/test/Threefish1024Test.java
var testVectors1024 = []testVector{
	testVector{
		key: "0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000",
		tweak: "00000000000000000000000000000000",
		plaintext: "0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000",
		ciphertext: "f05c3d0a3d05b304f785ddc7d1e036015c8aa76e2f217b06c6e1544c0bc1a90d" +
			"f0accb9473c24e0fd54fea68057f43329cb454761d6df5cf7b2e9b3614fbd5a2" +
			"0b2e4760b40603540d82eabc5482c171c832afbe68406bc39500367a592943fa" +
			"9a5b4a43286ca3c4cf46104b443143d560a4b230488311df4feef7e1dfe8391e",
	},
	testVector{
		key: "101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f" +
			"303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f" +
			"505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f" +
			"707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f",
		tweak: "000102030405060708090a0b0c0d0e0f",
		plaintext: "fffefdfcfbfaf9f8f7f6f5f4f3f2f1f0efeeedecebeae9e8e7e6e5e4e3e2e1e0" +
			"dfdedddcdbdad9d8d7d6d5d4d3d2d1d0cfcecdcccbcac9c8c7c6c5c4c3c2c1c0" +
			"bfbebdbcbbbab9b8b7b6b5b4b3b2b1b0afaeadacabaaa9a8a7a6a5a4a3a2a1a0" +
			"9f9e9d9c9b9a999897969594939291908f8e8d8c8b8a89888786858483828180",
		ciphertext: "a6654ddbd73cc3b05dd777105aa849bce49372eaaffc5568d254771bab85531c" +
			"94f780e7ffaae430d5d8af8c70eebbe1760f3b42b737a89cb363490d670314bd" +
			"8aa41ee63c2e1f45fbd477922f8360b388d6125ea6c7af0ad7056d01796e90c8" +
			"3313f4150a5716b30ed5f569288ae974ce2b4347926fce57de44512177dd7cde",
	},
}

func TestVectros256(t *testing.T) {
	for i := range testVectors256 {
		key, tweak, plaintext, ciphertext := convertVector(i, &testVectors256[i], t)

		c, err := New256(key, tweak)
		if err != nil {
			t.Fatal(err)
		}

		dst := make([]byte, blockSize256)

		c.Encrypt(dst, plaintext)
		if !bytesEqual(ciphertext, dst) {
			t.Fatalf("Test vector %d : Encryption failed\nFound:    %v \nExpected: %v", i, dst, ciphertext)
		}

		c.Decrypt(dst, dst)
		if !bytesEqual(plaintext, dst) {
			t.Fatalf("Test vector %d:Decryption failed\nFound: %v \nExpected: %v", i, dst, plaintext)
		}
	}
}

func TestVectros512(t *testing.T) {
	for i := range testVectors512 {
		key, tweak, plaintext, ciphertext := convertVector(i, &testVectors512[i], t)

		c, err := New512(key, tweak)
		if err != nil {
			t.Fatal(err)
		}

		dst := make([]byte, blockSize512)

		c.Encrypt(dst, plaintext)
		if !bytesEqual(ciphertext, dst) {
			t.Fatalf("Test vector %d : Encryption failed\nFound:    %v \nExpected: %v", i, dst, ciphertext)
		}

		c.Decrypt(dst, dst)
		if !bytesEqual(plaintext, dst) {
			t.Fatalf("Test vector %d:Decryption failed\nFound: %v \nExpected: %v", i, dst, plaintext)
		}
	}
}

func TestVectros1024(t *testing.T) {
	for i := range testVectors1024 {
		key, tweak, plaintext, ciphertext := convertVector(i, &testVectors1024[i], t)

		c, err := New1024(key, tweak)
		if err != nil {
			t.Fatal(err)
		}

		dst := make([]byte, blockSize1024)

		c.Encrypt(dst, plaintext)
		if !bytesEqual(ciphertext, dst) {
			t.Fatalf("Test vector %d : Encryption failed\nFound:    %v \nExpected: %v", i, dst, ciphertext)
		}

		c.Decrypt(dst, dst)
		if !bytesEqual(plaintext, dst) {
			t.Fatalf("Test vector %d:Decryption failed\nFound: %v \nExpected: %v", i, dst, plaintext)
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	tweak := make([]byte, TweakSize)

	dst, src := make([]byte, blockSize256), make([]byte, blockSize256)
	dstCp, srcCp := make([]byte, blockSize256), make([]byte, blockSize256)

	_, err := rand.Read(src)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	copy(srcCp, src)
	_, err = rand.Read(dst)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	copy(dstCp, dst)

	// 256 bit key
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	c, err := New256(key, tweak)
	if err != nil {
		t.Fatalf("Failed to create Threefish256 instance: %s", err)
	}
	c.Encrypt(dst, src)
	c.Decrypt(src, dst)
	if !bytesEqual(src, srcCp) {
		t.Fatalf("Threefish256: En/Decryption of dst to src failed\nKey: %v\nSrc:  %v\nDst: %v", key, srcCp, dst)
	}

	dst, src = make([]byte, blockSize512), make([]byte, blockSize512)
	dstCp, srcCp = make([]byte, blockSize512), make([]byte, blockSize512)

	_, err = rand.Read(src)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	copy(srcCp, src)
	_, err = rand.Read(dst)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	copy(dstCp, dst)

	// 512 bit key
	key = make([]byte, 64)
	_, err = rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	c, err = New512(key, tweak)
	if err != nil {
		t.Fatalf("Failed to create Threefish512 instance: %s", err)
	}
	c.Encrypt(dst, src)
	c.Decrypt(src, dst)
	if !bytesEqual(src, srcCp) {
		t.Fatalf("En/Decryption of dst to src failed\nKey: %v\nSrc:  %v\nDst: %v", key, srcCp, dst)
	}

	dst, src = make([]byte, blockSize1024), make([]byte, blockSize1024)
	dstCp, srcCp = make([]byte, blockSize1024), make([]byte, blockSize1024)

	_, err = rand.Read(src)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	copy(srcCp, src)
	_, err = rand.Read(dst)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	copy(dstCp, dst)

	// 1024 bit key
	key = make([]byte, 128)
	_, err = rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to read from random source: %s", err)
	}
	c, err = New1024(key, tweak)
	if err != nil {
		t.Fatalf("Failed to create Threefish1024 instance: %s", err)
	}
	c.Encrypt(dst, src)
	c.Decrypt(src, dst)
	if !bytesEqual(src, srcCp) {
		t.Fatalf("En/Decryption of dst to src failed\nKey: %v\nSrc:  %v\nDst: %v", key, srcCp, dst)
	}
}

func BenchmarkEncrypt256(b *testing.B) {
	c, err := New256(make([]byte, 32), make([]byte, 16))
	if err != nil {
		b.Fatalf("Failed to create threefish-256 instance - Cause: %s", err)
	}
	buf := make([]byte, c.BlockSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(buf, buf)
	}
}

func BenchmarkEncrypt512(b *testing.B) {
	c, err := New512(make([]byte, 64), make([]byte, 16))
	if err != nil {
		b.Fatalf("Failed to create threefish-256 instance - Cause: %s", err)
	}
	buf := make([]byte, c.BlockSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(buf, buf)
	}
}

func BenchmarkEncrypt1024(b *testing.B) {
	c, err := New1024(make([]byte, 128), make([]byte, 16))
	if err != nil {
		b.Fatalf("Failed to create threefish-256 instance - Cause: %s", err)
	}
	buf := make([]byte, c.BlockSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Encrypt(buf, buf)
	}
}

func BenchmarkDecrypt256(b *testing.B) {
	c, err := New256(make([]byte, 32), make([]byte, 16))
	if err != nil {
		b.Fatalf("Failed to create threefish-256 instance - Cause: %s", err)
	}
	buf := make([]byte, c.BlockSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(buf, buf)
	}
}

func BenchmarkDecrypt512(b *testing.B) {
	c, err := New512(make([]byte, 64), make([]byte, 16))
	if err != nil {
		b.Fatalf("Failed to create threefish-256 instance - Cause: %s", err)
	}
	buf := make([]byte, c.BlockSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(buf, buf)
	}
}

func BenchmarkDecrypt1024(b *testing.B) {
	c, err := New1024(make([]byte, 128), make([]byte, 16))
	if err != nil {
		b.Fatalf("Failed to create threefish-256 instance - Cause: %s", err)
	}
	buf := make([]byte, c.BlockSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Decrypt(buf, buf)
	}
}
