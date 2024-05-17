package dicom

import (
	"errors"
	"fmt"
	"io"

	"github.com/kulaginds/dicom/lowlevel"
	"github.com/kulaginds/dicom/tag"
	"github.com/kulaginds/dicom/uid"
	"github.com/kulaginds/dicom/vr"
	"github.com/kulaginds/dicom/vr/parse"
)

type fullReader struct {
	lowReader *lowlevel.Reader
}

func NewFullReader(r io.Reader) *fullReader {
	return &fullReader{
		lowReader: lowlevel.NewReader(r),
	}
}

func (r *fullReader) ReadDataset() (*Dataset, error) {
	err := r.lowReader.Header()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	var ds Dataset

	err = r.readMetaInfoElements(&ds)
	if err != nil {
		return nil, fmt.Errorf("read meta info elements: %w", err)
	}

	err = r.readElements(&ds)
	if err != nil {
		return nil, fmt.Errorf("read elements: %w", err)
	}

	return &ds, nil
}

func (r *fullReader) readMetaInfoElements(ds *Dataset) error {
	var (
		transferSyntaxUID string
		err               error
	)

	parentLowReader := r.lowReader

	for {
		var elem Element

		err = r.readElementTag(&elem)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read tag: %w", err)
		}

		err = r.readElementValue(&elem)
		if err != nil {
			return fmt.Errorf("read element value: %w", err)
		}

		ds.Elements = append(ds.Elements, elem)

		if elem.Tag.Equal(tag.FileMetaInformationGroupLength) {
			if elem.VR != vr.UL || elem.VL != 4 {
				return fmt.Errorf("incorrect FileMetaInformationGroupLength: vr=%s, vl=%d", elem.VR, elem.VL)
			}

			metaGroupLength := parse.UL(elem.Value[:elem.VL], r.lowReader.ByteOrder)
			r.lowReader = r.lowReader.CopyWithLimit(metaGroupLength)
		}

		if elem.Tag.Equal(tag.TransferSyntaxUID) {
			transferSyntaxUID = parse.UI(elem.Value)
		}
	}

	r.lowReader = parentLowReader

	if transferSyntaxUID != "" {
		r.lowReader.ByteOrder, r.lowReader.Implicit = uid.ParseTransferSyntaxUID(transferSyntaxUID)
	}

	return nil
}

const (
	itemLengthSize = 4
)

func (r *fullReader) readElements(ds *Dataset) error {
	var err error

	for {
		var elem Element

		err = r.readElementTag(&elem)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("read tag: %w", err)
		}

		if elem.Tag.Equal(tag.ItemDelimitationItem) {
			break
		}

		if elem.VR == vr.SQ {
			err = r.readSequence(&elem)
			if err != nil {
				return fmt.Errorf("read sequence: %w", err)
			}
		} else {
			err = r.readElementValue(&elem)
			if err != nil {
				return fmt.Errorf("read element value: %w", err)
			}
		}

		ds.Elements = append(ds.Elements, elem)
	}

	return nil
}

func (r *fullReader) readElementTag(elem *Element) error {
	var err error

	elem.Tag, err = r.lowReader.Tag()
	if errors.Is(err, io.EOF) {
		return err
	}
	if err != nil {
		return fmt.Errorf("read tag: %w", err)
	}

	if !elem.Tag.Equal(tag.ItemDelimitationItem) {
		elem.VR, err = r.lowReader.VR(elem.Tag)
		if err != nil {
			return fmt.Errorf("read VR: %w", err)
		}
	}

	elem.VL, err = r.lowReader.VL(elem.VR)
	if err != nil {
		return fmt.Errorf("read VL: %w", err)
	}

	return nil
}

func (r *fullReader) readElementValue(elem *Element) error {
	if elem.VL == lowlevel.UndefinedLength {
		return fmt.Errorf("element with undefined length: tag=(%x, %x), vr=%s",
			elem.Tag.GroupNumber, elem.Tag.ElementNumber, elem.VR)
	}

	elem.Value = make([]byte, elem.VL)

	n, err := io.ReadFull(r.lowReader, elem.Value)
	if err != nil {
		return fmt.Errorf("read value: %w", err)
	}

	if uint32(n) != elem.VL {
		return fmt.Errorf("value length mismatch: expected %d, got %d", elem.VL, n)
	}

	return nil
}

func (r *fullReader) readSequence(seqElem *Element) error {
	seqElem.Sequence = new(Sequence)

	parentLowReader := r.lowReader
	if seqElem.VL != lowlevel.UndefinedLength {
		r.lowReader = r.lowReader.CopyWithLimit(seqElem.VL)
	}

	var err error

	for {
		var item SequenceItem

		err = r.readSequenceItem(&item)
		if errors.Is(err, io.EOF) {
			return err
		}
		if err != nil {
			return fmt.Errorf("read item: %w", err)
		}

		if item.Tag.Equal(tag.SequenceDelimitationItem) {
			break
		}

		seqElem.Sequence.Items = append(seqElem.Sequence.Items, item)
	}

	if seqElem.VL != lowlevel.UndefinedLength {
		r.lowReader = parentLowReader
	}

	return nil
}

func (r *fullReader) readSequenceItem(item *SequenceItem) error {
	var err error

	item.Tag, err = r.lowReader.Tag()
	if errors.Is(err, io.EOF) {
		return err
	}
	if err != nil {
		return fmt.Errorf("read item tag: %w", err)
	}

	item.Length, err = r.lowReader.UInt32()
	if err != nil {
		return fmt.Errorf("read item length: %w", err)
	}

	if item.Tag.Equal(tag.SequenceDelimitationItem) {
		return nil
	}
	if !item.Tag.Equal(tag.Item) {
		return fmt.Errorf("unexpected tag instead Item: (%x, %x)", item.Tag.GroupNumber, item.Tag.ElementNumber)
	}

	item.Dataset = new(Dataset)

	parentLowReader := r.lowReader
	if item.Length != lowlevel.UndefinedLength {
		r.lowReader = r.lowReader.CopyWithLimit(item.Length)
	}

	err = r.readElements(item.Dataset)
	if err != nil {
		return fmt.Errorf("read elements: %w", err)
	}

	if item.Length != lowlevel.UndefinedLength {
		r.lowReader = parentLowReader
	}

	return nil
}
