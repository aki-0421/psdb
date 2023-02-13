package types

type VitessSessionOptions struct {
	IncludedFields  string `json:"includedFields"`
	ClientFoundRows bool   `json:"clientFoundRows"`
}

type VitessSession struct {
	UUID                 string               `json:"SessionUUID"`
	Autocommit           bool                 `json:"autocommit"`
	DDLStrategy          string               `json:"DDLStrategy"`
	EnableSystemSettings bool                 `json:"enableSystemSettings"`
	Options              VitessSessionOptions `json:"options"`
}

type Session struct {
	Signature     string        `json:"signature"`
	VitessSession VitessSession `json:"vitessSession"`
}
