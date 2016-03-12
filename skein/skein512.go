package skein

import (
	"github.com/EncEve/crypto/threefish"
)

func (s *Skein512) BlockSize() int {
	return StateSize512
}

func (s *Skein512) Size() int { return s.hsize }

func (s *Skein512) Reset() {
	s.off = 0
	s.msg = false
	copy(s.hVal[:8], s.initVal[:])
	s.tweak[0] = 0
	s.tweak[1] = messageParam<<56 | firstBlock
}

func (s *Skein512) Write(in []byte) (int, error) {
	n := len(in)

	diff := StateSize512 - s.off
	if n > diff {
		// process buffer.
		copy(s.buf[s.off:], in[:diff])
		s.hashMessage(s.buf[:])
		s.off = 0

		in = in[diff:]
	}
	// process full blocks except for the last
	length := len(in)
	if length > StateSize512 {
		nn := length - (length % StateSize512)
		if nn == length {
			nn -= StateSize512
		}
		s.hashMessage(in[:nn])
		in = in[nn:]
	}
	s.off += copy(s.buf[s.off:], in)

	s.msg = true
	return n, nil
}

func (s *Skein512) Sum(in []byte) []byte {
	s0 := *s // make a copy
	if s0.msg {
		s0.finalize()
	}

	var out [StateSize512]byte
	s0.output(&out, 0)
	return append(in, out[:s0.hsize]...)
}

func (s *Skein512) hashMessage(blocks []byte) {
	var message, msg [8]uint64
	var block [StateSize512]byte
	for i := 0; i < len(blocks); i += StateSize512 {
		// copy the message block into an array for
		// the toWords512 function
		copy(block[:], blocks[i:i+StateSize512])
		toWords512(&msg, &block)
		message = msg

		s.hVal[8] = threefish.C240 ^ s.hVal[0] ^ s.hVal[1] ^ s.hVal[2] ^
			s.hVal[3] ^ s.hVal[4] ^ s.hVal[5] ^ s.hVal[6] ^ s.hVal[7]

		incTweak(&(s.tweak), StateSize512)
		s.tweak[2] = s.tweak[0] ^ s.tweak[1]

		threefish.Encrypt512(&(s.hVal), &(s.tweak), &msg)
		xor512(&(s.hVal), &message, &msg)

		// clear the first block flag
		s.tweak[1] &^= firstBlock
	}
}

// Finalize the hash function with the last message block
func (s *Skein512) finalize() {
	var message, msg [8]uint64
	// flush the buffer
	for i := s.off; i < StateSize512; i++ {
		s.buf[i] = 0
	}

	toWords512(&msg, &(s.buf))
	message = msg

	s.hVal[8] = threefish.C240 ^ s.hVal[0] ^ s.hVal[1] ^ s.hVal[2] ^
		s.hVal[3] ^ s.hVal[4] ^ s.hVal[5] ^ s.hVal[6] ^ s.hVal[7]

	incTweak(&(s.tweak), uint64(s.off))
	s.tweak[1] |= lastBlock // set the last block flag
	s.tweak[2] = s.tweak[0] ^ s.tweak[1]

	threefish.Encrypt512(&(s.hVal), &(s.tweak), &msg)
	xor512(&(s.hVal), &message, &msg)
	s.off = 0
}

// Extract the output from the hash function
func (s *Skein512) output(dst *[StateSize512]byte, ctr uint64) {
	var message, msg [8]uint64
	msg[0], message[0] = ctr, ctr

	s.hVal[8] = threefish.C240 ^ s.hVal[0] ^ s.hVal[1] ^ s.hVal[2] ^
		s.hVal[3] ^ s.hVal[4] ^ s.hVal[5] ^ s.hVal[6] ^ s.hVal[7]

	threefish.Encrypt512(&(s.hVal), &outTweak, &msg)
	xor512(&(s.hVal), &message, &msg)

	for i, v := range s.hVal[:8] {
		j := i * 8
		dst[j+0] = byte(v)
		dst[j+1] = byte(v >> 8)
		dst[j+2] = byte(v >> 16)
		dst[j+3] = byte(v >> 24)
		dst[j+4] = byte(v >> 32)
		dst[j+5] = byte(v >> 40)
		dst[j+6] = byte(v >> 48)
		dst[j+7] = byte(v >> 56)
	}
}

// Add a parameter (secret key, nonce etc.) to the hash function
func (s *Skein512) addParam(ptype uint64, param []byte) {
	s.tweak[0] = 0
	s.tweak[1] = ptype<<56 | firstBlock
	s.Write(param)
	s.finalize()
}

// Add the configuration block to the hash function
func (s *Skein512) addConfig(hashsize int) {
	var c [32]byte
	copy(c[:], schemaId)

	bits := uint64(hashsize * 8)
	c[8] = byte(bits)
	c[9] = byte(bits >> 8)
	c[10] = byte(bits >> 16)
	c[11] = byte(bits >> 24)
	c[12] = byte(bits >> 32)
	c[13] = byte(bits >> 40)
	c[14] = byte(bits >> 48)
	c[15] = byte(bits >> 56)

	s.addParam(configParam, c[:])
}
