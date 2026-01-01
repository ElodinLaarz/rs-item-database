package main

import (
	"context"
	"fmt"

	"rs-item-database/internal/db"
	"rs-item-database/internal/ingest"
	"rs-item-database/pb"
)

// App struct
type App struct {
	ctx           context.Context
	store         *db.Store
	ingestService *ingest.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Init DB
	store, err := db.NewStore("./items.db")
	if err != nil {
		// In a real app, maybe show a dialog or log better
		panic(err)
	}
	a.store = store

	// Init Ingest Service
	a.ingestService = ingest.NewService()
}

func (a *App) shutdown(ctx context.Context) {
	if a.store != nil {
		a.store.Close()
	}
	if a.ingestService != nil {
		a.ingestService.Shutdown()
	}
}

// Search items by prefix
func (a *App) Search(query string) []*pb.Item {
	if a.store == nil {
		return nil
	}
	items, err := a.store.SearchItems(query, 50)
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return nil
	}
	return items
}

// IngestItem fetches an item by ID and saves it (Helper for dev)
func (a *App) IngestItem(id int) string {
	if a.ingestService == nil {
		return "Ingest service not initialized"
	}

	item, err := a.ingestService.FetchItem(id)
	if err != nil {
		return fmt.Sprintf("Error fetching: %v", err)
	}

	if err := a.store.SaveItem(item); err != nil {
		return fmt.Sprintf("Error saving: %v", err)
	}

	return fmt.Sprintf("Saved: %s", item.Name)
}
