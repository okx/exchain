package iavl

import (
	"testing"
)

func TestBatchBaseAsyncCommit(t *testing.T) {
	EnableAsyncCommit = true
	defer t.Cleanup(func() {
		EnableAsyncCommit = false
	})

	// for basic_test.go
	TestBasic(t)
	TestUnit(t)
	TestRemove(t)
	TestIntegration(t)
	TestIterateRange(t)
	TestPersistence(t)
	TestProof(t)
	TestTreeProof(t)

	// for export_test.go
	TestExporter(t)
	TestExporter_Import(t)
	TestExporter_Close(t)

	// for import_test.go
	TestImporter_NegativeVersion(t)
	TestImporter_NotEmpty(t)
	TestImporter_NotEmptyDatabase(t)
	TestImporter_NotEmptyUnsaved(t)
	TestImporter_Add(t)
	TestImporter_Add_Closed(t)
	TestImporter_Close(t)
	TestImporter_Commit(t)
	TestImporter_Commit_Closed(t)
	TestImporter_Commit_Empty(t)

	// for tree_dotgraph_test.go
	TestWriteDOTGraph(t)

	// for tree_fuzz_test.go
	TestMutableTreeFuzz(t)

	// for tree_test.go
	TestVersionedTreeSpecialCase(t)
	TestVersionedTreeErrors(t)
	TestVersionedCheckpointsSpecialCase(t)
	TestVersionedCheckpointsSpecialCase2(t)
	TestVersionedCheckpointsSpecialCase3(t)
	TestVersionedCheckpointsSpecialCase5(t)
	TestVersionedCheckpointsSpecialCase6(t)
	TestVersionedCheckpointsSpecialCase7(t)
	TestNilValueSemantics(t)
	TestCopyValueSemantics(t)
}
