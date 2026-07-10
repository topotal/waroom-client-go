package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

// loadModel builds the semantic model from the real internal/ogen package
// with the repository's Phase 1 configuration.
func loadModel(t *testing.T) *Model {
	t.Helper()
	cfg, err := loadConfig("../../../.genwaroom.yaml")
	require.NoError(t, err)
	tags, err := loadTags("../../../api/openapi.yaml")
	require.NoError(t, err)
	model, err := load("../../ogen", tags, cfg)
	require.NoError(t, err)
	return model
}

func findOp(t *testing.T, model *Model, id string) *Op {
	t.Helper()
	for _, svc := range model.Services {
		for _, op := range svc.Ops {
			if op.ID == id {
				return op
			}
		}
	}
	t.Fatalf("operation %s not found in model", id)
	return nil
}

func TestLoad_Services(t *testing.T) {
	model := loadModel(t)

	var names []string
	for _, svc := range model.Services {
		names = append(names, svc.Name)
	}
	assert.Equal(t, []string{"IncidentService", "LabelService", "ServiceService"}, names)
}

// One test per structural shape the conventions must recognize.

func TestLoad_ListWithOptions(t *testing.T) {
	op := findOp(t, loadModel(t), "listIncidents")

	assert.Equal(t, "list", op.Success.Kind)
	assert.Equal(t, "IncidentList", op.Success.ListName)
	assert.Equal(t, "IncidentListItem", op.Success.ItemPublic)
	require.NotNil(t, op.Options)
	var fields []string
	for _, f := range op.Options.Fields {
		fields = append(fields, f.Name+" "+f.Type)
	}
	assert.Equal(t, []string{"Page OptInt", "PerPage OptInt", "SortOrder OptListIncidentsSortOrder"}, fields)
	assert.True(t, op.All, "list op with page params must get an All iterator")
}

func TestLoad_ListWithoutParams(t *testing.T) {
	op := findOp(t, loadModel(t), "listServices")

	assert.Equal(t, "list", op.Success.Kind)
	assert.Empty(t, op.PathParams)
	assert.Nil(t, op.Options)
	assert.False(t, op.All, "list op without page params must not get an All iterator")
}

func TestLoad_WholeResponseWithPathParam(t *testing.T) {
	op := findOp(t, loadModel(t), "getIncident")

	assert.Equal(t, "whole", op.Success.Kind)
	assert.Equal(t, "Incident", op.Success.Public)
	assert.Equal(t, "GetIncidentParams", op.ParamsGen)
	assert.Equal(t, []PathParam{{Arg: "uuid", Type: "uuid.UUID", GenField: "UUID"}}, op.PathParams)
}

func TestLoad_RequestUnwrap(t *testing.T) {
	op := findOp(t, loadModel(t), "createServiceLabel")

	assert.Equal(t, "unwrap", op.Success.Kind)
	assert.Equal(t, "CreatedLabel", op.Success.Public)
	assert.Equal(t, &Request{
		Public:      "CreateServiceLabelRequest",
		Gen:         "CreateServiceLabelReq",
		UnwrapField: "Label",
	}, op.Request)
}

func TestLoad_NoContent(t *testing.T) {
	op := findOp(t, loadModel(t), "deleteServiceLabel")

	assert.Equal(t, "none", op.Success.Kind)
	assert.Equal(t, []PathParam{
		{Arg: "name", Type: "string", GenField: "Name"},
		{Arg: "uuid", Type: "uuid.UUID", GenField: "UUID"},
	}, op.PathParams)

	var statuses []int
	for _, e := range op.Errors {
		statuses = append(statuses, e.Status)
	}
	assert.ElementsMatch(t, []int{401, 402, 404}, statuses, "all declared error variants must be mapped")
}

// TestBuildOpView pins the Go fragments the render stage derives from the
// semantic model.
func TestBuildOpView(t *testing.T) {
	model := loadModel(t)

	create, err := buildOpView(findOp(t, model, "createServiceLabel"))
	require.NoError(t, err)
	assert.Equal(t, "ctx context.Context, name string, req CreateServiceLabelRequest", create.ParamsDecl)
	assert.Equal(t, "(*CreatedLabel, error)", create.ResultsDecl)
	assert.Equal(t, "ctx, &gen.CreateServiceLabelReq{Label: req}, gen.CreateServiceLabelParams{Name: name}", create.RawArgs)

	del, err := buildOpView(findOp(t, model, "deleteServiceLabel"))
	require.NoError(t, err)
	assert.Equal(t, "error", del.ResultsDecl)
	assert.Equal(t, "", del.ZeroRet)

	list, err := buildOpView(findOp(t, model, "listServiceLabels"))
	require.NoError(t, err)
	assert.Equal(t, "ctx, name, o", list.FwdArgs)
}

// TestRender_Golden runs the full analyze → YAML → render pipeline
// (including the serialization round-trip) and compares the output with
// the golden files.
func TestRender_Golden(t *testing.T) {
	model := loadModel(t)
	dir := t.TempDir()

	modelPath := filepath.Join(dir, "sdk-model.yaml")
	require.NoError(t, writeModel(modelPath, model))
	restored, err := readModel(modelPath)
	require.NoError(t, err)
	require.NoError(t, render(restored, dir))

	generated, err := filepath.Glob(filepath.Join(dir, "*.go"))
	require.NoError(t, err)
	require.NotEmpty(t, generated)

	goldenDir := filepath.Join("testdata", "golden")
	if *update {
		require.NoError(t, os.RemoveAll(goldenDir))
		require.NoError(t, os.MkdirAll(goldenDir, 0o755))
		for _, path := range generated {
			data, err := os.ReadFile(path)
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(filepath.Join(goldenDir, filepath.Base(path)+".golden"), data, 0o644))
		}
	}

	goldens, err := filepath.Glob(filepath.Join(goldenDir, "*.golden"))
	require.NoError(t, err)
	require.Len(t, goldens, len(generated), "generated file set differs from golden (run with -update)")

	for _, path := range generated {
		got, err := os.ReadFile(path)
		require.NoError(t, err)
		want, err := os.ReadFile(filepath.Join(goldenDir, filepath.Base(path)+".golden"))
		require.NoError(t, err, "missing golden for %s (run with -update)", filepath.Base(path))
		assert.Equal(t, string(want), string(got), "%s differs from golden (run with -update)", filepath.Base(path))
	}
}
