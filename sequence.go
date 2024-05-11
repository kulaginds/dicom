package dicom

import (
	"github.com/kulaginds/dicom/tag"
)

type SequenceItem struct {
	Tag     tag.Tag
	Length  uint32
	Dataset *Dataset
}

type Sequence struct {
	Items []SequenceItem
}
