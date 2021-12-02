package iavl

// TreeDelta is the delta for applying on new version tree
type TreeDelta struct {
	NodesDelta         map[string]*NodeJson `json:"nodes_delta"`
	OrphansDelta       map[string]int64   `json:"orphans_delta"`
	CommitOrphansDelta map[string]int64     `json:"commit_orphans_delta"`
}
