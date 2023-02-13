package types

type Field struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Table        string `json:"table"`
	OrgTable     string `json:"orgTable"`
	Database     string `json:"database"`
	OrgName      string `json:"orgName"`
	ColumnLength int64  `json:"columnLength"`
	Charset      int    `json:"charset"`
	Flags        int    `json:"flags"`
}
