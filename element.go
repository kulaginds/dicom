package dicom

import (
	"github.com/kulaginds/dicom/tag"
	"github.com/kulaginds/dicom/vr"
)

type Element struct {
	Tag      tag.Tag
	VR       vr.VR
	VL       uint32
	Value    []byte
	Sequence *Sequence
}
