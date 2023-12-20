package contrast

type ContrastFindings struct {
	ClassificationName string `json:"classificationName"`
	Total              int    `json:"total"`
	Audited            int    `json:"audited"`
}

type ContrastAudit struct {
	ToolName       string `json:"toolName"`
	ApplicationURL string `json:"applicationURL"`
}
