package types

type ExecuteResult struct {
	Fields []Field `json:"fields"`
	Rows   []Row   `json:"rows"`
}
