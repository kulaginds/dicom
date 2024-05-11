package dicom

import (
	"os"
	"testing"
)

func TestFullReader(t *testing.T) {
	testCases := []string{
		"testdata/1.dcm",
		"testdata/2.dcm",
		"testdata/3.dcm",
		"testdata/4.dcm",
		"testdata/5.dcm",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			f, err := os.Open(tc)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			reader := NewFullReader(f)
			_, err = reader.ReadDataset()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
