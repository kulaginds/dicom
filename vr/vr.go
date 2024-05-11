package vr

// VR Value Representation
type VR string

const (
	// AE Application Entity
	AE VR = "AE"

	// AS Age String
	AS VR = "AS"

	// AT Attribute Tag
	AT VR = "AT"

	// CS Code String
	CS VR = "CS"

	// DA Date
	DA VR = "DA"

	// DS Decimal String
	DS VR = "DS"

	// DT Date Time
	DT VR = "DT"

	// FL Floating Point Single
	FL VR = "FL"

	// FD Floating Point Double
	FD VR = "FD"

	// IS Integer String
	IS VR = "IS"

	// LO Long String
	LO VR = "LO"

	// LT Long Text
	LT VR = "LT"

	// OB Other Byte
	OB VR = "OB"

	// OD Other Double
	OD VR = "OD"

	// OF Other Float
	OF VR = "OF"

	// OL Other Long
	OL VR = "OL"

	// OV Other 64-bit Very Long
	OV VR = "OV"

	// OW Other Word
	OW VR = "OW"

	// PN Person Name
	PN VR = "PN"

	// SH Short String
	SH VR = "SH"

	// SL Signed Long
	SL VR = "SL"

	// SQ Sequence of Items
	SQ VR = "SQ"

	// SS Signed Short
	SS VR = "SS"

	// ST Short Text
	ST VR = "ST"

	// SV Signed 64-bit Very Long
	SV VR = "SV"

	// TM Time
	TM VR = "TM"

	// UC Unlimited Characters
	UC VR = "UC"

	// UI Unique Identifier (UID)
	UI VR = "UI"

	// UL Unsigned Long
	UL VR = "UL"

	// UN Unknown
	UN VR = "UN"

	// UR Universal Resource Identifier or Universal Resource Locator (URI/URL)
	UR VR = "UR"

	// US Unsigned Short
	US VR = "US"

	// UT Unlimited Text
	UT VR = "UT"

	// UV Unsigned 64-bit Very Long
	UV VR = "UV"
)
