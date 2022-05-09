package helper

// LabelsField presents a "f:labels" field subfield of metadataField.
type LabelsField struct {
	Labels map[string]struct{} `json:"f:labels"`
}

// MetadataField presents a "f:metadata" field subfield of v1.FieldsV1.
type MetadataField struct {
	LabelsField `json:"f:metadata"`
}
