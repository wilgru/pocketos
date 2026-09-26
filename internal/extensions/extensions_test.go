package extensions

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func TestRegisterTracksNewBaseCollection(t *testing.T) {
	app, err := tests.NewTestApp()
	if err != nil {
		t.Fatal(err)
	}
	defer app.Cleanup()

	tools := core.NewBaseCollection(toolsCollectionName)
	tools.System = true
	tools.Fields.Add(&core.TextField{
		Name:        "targetCollectionName",
		Required:    true,
		Presentable: true,
	})
	tools.Fields.Add(&core.SelectField{
		Name:      "layout",
		Required:  true,
		MaxSelect: 1,
		Values:    []string{"detail", "list", "gallery", "table"},
	})
	tools.Fields.Add(&core.TextField{
		Name:     "toolNameOverride",
		Required: false,
	})
	tools.Fields.Add(&core.TextField{
		Name:     "accentColour",
		Required: true,
		Pattern:  "^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$",
	})
	tools.Fields.Add(&core.TextField{
		Name:     "icon",
		Required: false,
	})
	tools.Fields.Add(&core.AutodateField{
		Name:     "created",
		OnCreate: true,
	})
	tools.Fields.Add(&core.AutodateField{
		Name:     "updated",
		OnCreate: true,
		OnUpdate: true,
	})
	tools.Indexes = []string{"CREATE UNIQUE INDEX idx_tools_target_collection_name ON _tools (targetCollectionName)"}
	if err := app.Save(tools); err != nil {
		t.Fatal(err)
	}

	Register(app)

	collection := core.NewBaseCollection("projects")
	if err := app.Save(collection); err != nil {
		t.Fatal(err)
	}

	tools, err = app.FindCollectionByNameOrId(toolsCollectionName)
	if err != nil {
		t.Fatal(err)
	}
	if !tools.System {
		t.Fatal("_tools collection must be a system collection")
	}

	record, err := app.FindFirstRecordByFilter(tools, "targetCollectionName = {:targetCollectionName}", map[string]any{
		"targetCollectionName": collection.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if record.GetString("targetCollectionName") != collection.Name {
		t.Fatalf("got collection name %q, want %q", record.GetString("targetCollectionName"), collection.Name)
	}
	if record.GetString("layout") != "detail" {
		t.Fatalf("got layout %q, want %q", record.GetString("layout"), "detail")
	}
	if record.GetString("accentColour") != defaultAccentColour {
		t.Fatalf("got accentColour %q, want %q", record.GetString("accentColour"), defaultAccentColour)
	}

	collection.Name = "workspaces"
	if err := app.Save(collection); err != nil {
		t.Fatal(err)
	}

	record, err = app.FindFirstRecordByFilter(tools, "targetCollectionName = {:targetCollectionName}", map[string]any{
		"targetCollectionName": collection.Name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if record.GetString("targetCollectionName") != collection.Name {
		t.Fatalf("got renamed collection name %q, want %q", record.GetString("targetCollectionName"), collection.Name)
	}

	if err := app.Delete(collection); err != nil {
		t.Fatal(err)
	}

	if _, err := app.FindFirstRecordByFilter(tools, "targetCollectionName = {:targetCollectionName}", map[string]any{
		"targetCollectionName": collection.Name,
	}); err == nil {
		t.Fatal("expected _tools record to be deleted, but it still exists")
	}
}
