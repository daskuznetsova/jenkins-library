package contrast

type ContrastFindings struct {
	ClassificationName string `json:"classificationName"`
	Total              int    `json:"total"`
	Audited            int    `json:"audited"`
}
