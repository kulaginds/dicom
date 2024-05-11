package lowlevel

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kulaginds/dicom/tag"
	"github.com/kulaginds/dicom/vr"
)

type Reader struct {
	r         io.Reader
	buf       [128]byte
	ByteOrder binary.ByteOrder
	Implicit  bool
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:         r,
		ByteOrder: binary.LittleEndian,
		Implicit:  false,
	}
}

func (r *Reader) CopyWithLimit(limit uint32) *Reader {
	return &Reader{
		r:         io.LimitReader(r, int64(limit)),
		buf:       r.buf,
		ByteOrder: r.ByteOrder,
		Implicit:  r.Implicit,
	}
}

const magicWord = "DICM"

func (r *Reader) Header() error {
	// read offset
	n, err := r.r.Read(r.buf[:])
	if err != nil {
		return fmt.Errorf("offset: %w", err)
	}

	if n != len(r.buf) {
		return ErrIncorrectHeader
	}

	// read magic word
	n, err = r.r.Read(r.buf[0:len(magicWord)])
	if err != nil {
		return fmt.Errorf("magic word: %w", err)
	}

	if n != len(magicWord) {
		return ErrIncorrectHeader
	}

	return nil
}

const tagLength = 4

func (r *Reader) Tag() (tag.Tag, error) {
	n, err := r.r.Read(r.buf[:tagLength])
	if err != nil {
		return tag.Tag{}, fmt.Errorf("tag: %w", err)
	}

	if n != tagLength {
		return tag.Tag{}, ErrIncorrectTag
	}

	return tag.Tag{
		GroupNumber:   r.ByteOrder.Uint16(r.buf[0 : tagLength/2]),
		ElementNumber: r.ByteOrder.Uint16(r.buf[tagLength/2 : tagLength]),
	}, nil
}

const (
	vrLength         = 2
	vrReservedLength = 2
)

func (r *Reader) VR(t tag.Tag) (vr.VR, error) {
	if r.Implicit {
		rep, find := vr.Tag2VR[t]
		if !find {
			return vr.UN, nil
		}

		return rep, nil
	}

	n, err := r.r.Read(r.buf[:vrLength])
	if err != nil {
		return "", fmt.Errorf("vr: %w", err)
	}

	if n != vrLength {
		return "", ErrIncorrectValueRepresentation
	}

	valueRep := vr.VR(r.buf[:vrLength])

	if _, is16BitVl := vrWith16bitVl[valueRep]; !is16BitVl {
		n, err = r.r.Read(r.buf[:vrReservedLength]) // ignore two reserved bytes (0000H)
		if err != nil {
			return "", fmt.Errorf("skip reserved bytes: %w", err)
		}

		if n != vrReservedLength {
			return "", ErrIncorrectValueRepresentation
		}
	}

	return valueRep, nil
}

var vrWith16bitVl = map[vr.VR]struct{}{
	vr.AE: {}, vr.AS: {}, vr.AT: {}, vr.CS: {},
	vr.DA: {}, vr.DS: {}, vr.DT: {}, vr.FL: {},
	vr.FD: {}, vr.IS: {}, vr.LO: {}, vr.LT: {},
	vr.PN: {}, vr.SH: {}, vr.SL: {}, vr.SS: {},
	vr.ST: {}, vr.TM: {}, vr.UI: {}, vr.UL: {},
	vr.US: {},
}

const (
	smallVLLength        = 2
	bigVLLength          = 4
	smallUndefinedLength = 0xFFFF
)

func (r *Reader) VL(valueRep vr.VR) (uint32, error) {
	_, is16BitVl := vrWith16bitVl[valueRep]
	if !r.Implicit && is16BitVl {
		n, err := r.r.Read(r.buf[:smallVLLength])
		if err != nil {
			return 0, fmt.Errorf("vl16: %w", err)
		}

		if n != smallVLLength {
			return 0, ErrIncorrectValueLength
		}

		vl := uint32(r.ByteOrder.Uint16(r.buf[:smallVLLength]))
		// Rectify Undefined Length VL
		if vl == smallUndefinedLength {
			vl = UndefinedLength
		}

		return vl, nil
	}

	n, err := r.r.Read(r.buf[:bigVLLength])
	if err != nil {
		return 0, fmt.Errorf("vl32: %w", err)
	}

	if n != bigVLLength {
		return 0, ErrIncorrectValueLength
	}

	return r.ByteOrder.Uint32(r.buf[:bigVLLength]), nil
}

func (r *Reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *Reader) UInt32() (uint32, error) {
	const uint32Size = 4

	n, err := r.r.Read(r.buf[:uint32Size])
	if err != nil {
		return 0, fmt.Errorf("uint32: %w", err)
	}

	if n != uint32Size {
		return 0, ErrIncorrectUInt32
	}

	return r.ByteOrder.Uint32(r.buf[:uint32Size]), nil
}

func (r *Reader) Skip(n int64) error {
	k, err := io.CopyBuffer(io.Discard, io.LimitReader(r.r, n), r.buf[:])
	if err != nil {
		return fmt.Errorf("skip: %w", err)
	}

	if k != n {
		return fmt.Errorf("skip: expected %d, got %d", n, k)
	}

	return nil
}
