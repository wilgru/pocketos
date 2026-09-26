package extensions

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

const toolsCollectionName = "_tools"
const defaultAccentColour = "#1055c9"

func Register(app core.App) {
	app.OnCollectionUpdate().BindFunc(func(e *core.CollectionEvent) error {
		if e.Collection.Type != core.CollectionTypeBase || e.Collection.Name == toolsCollectionName {
			return e.Next()
		}

		original, err := e.App.FindCollectionByNameOrId(e.Collection.Id)
		if err != nil {
			return fmt.Errorf("find collection before update: %w", err)
		}
		if original.Name == e.Collection.Name {
			return e.Next()
		}

		tools, err := e.App.FindCollectionByNameOrId(toolsCollectionName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.Next()
			}
			return fmt.Errorf("find _tools collection: %w", err)
		}

		record, err := e.App.FindFirstRecordByFilter(tools, "targetCollectionName = {:name}", map[string]any{
			"name": original.Name,
		})
		if err != nil {
			return fmt.Errorf("find _tools record for %q: %w", original.Name, err)
		}
		record.Set("targetCollectionName", e.Collection.Name)
		if err := e.App.Save(record); err != nil {
			return fmt.Errorf("rename _tools record from %q to %q: %w", original.Name, e.Collection.Name, err)
		}

		return e.Next()
	})

	app.OnCollectionAfterCreateSuccess().BindFunc(func(e *core.CollectionEvent) error {
		if e.Collection.Type != core.CollectionTypeBase || e.Collection.Name == toolsCollectionName {
			return e.Next()
		}

		tools, err := e.App.FindCollectionByNameOrId(toolsCollectionName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.Next()
			}
			return fmt.Errorf("find _tools collection: %w", err)
		}

		record := core.NewRecord(tools)
		record.Set("targetCollectionName", e.Collection.Name)
		record.Set("layout", "detail")
		record.Set("accentColour", defaultAccentColour)
		if err := e.App.Save(record); err != nil {
			return fmt.Errorf("create _tools record for %q: %w", e.Collection.Name, err)
		}

		return e.Next()
	})

	app.OnCollectionAfterDeleteSuccess().BindFunc(func(e *core.CollectionEvent) error {
		if e.Collection.Type != core.CollectionTypeBase || e.Collection.Name == toolsCollectionName {
			return e.Next()
		}

		tools, err := e.App.FindCollectionByNameOrId(toolsCollectionName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.Next()
			}
			return fmt.Errorf("find _tools collection: %w", err)
		}

		record, err := e.App.FindFirstRecordByFilter(tools, "targetCollectionName = {:name}", map[string]any{
			"name": e.Collection.Name,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return e.Next()
			}
			return fmt.Errorf("find _tools record for %q: %w", e.Collection.Name, err)
		}
		if err := e.App.Delete(record); err != nil {
			return fmt.Errorf("delete _tools record for %q: %w", e.Collection.Name, err)
		}

		return e.Next()
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		if _, err := e.App.FindCollectionByNameOrId(toolsCollectionName); err != nil {
			return fmt.Errorf("_tools migration must be applied before serving: %w", err)
		}

		e.Router.GET("/api/healthz/go", func(re *core.RequestEvent) error {
			return re.JSON(http.StatusOK, map[string]string{
				"status": "ok",
				"source": "go",
			})
		})

		return e.Next()
	})
}
