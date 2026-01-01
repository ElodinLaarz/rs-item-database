package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"rs-item-database/internal/db"
	"rs-item-database/internal/ingest"
)

func main() {
	// 1. Init DB
	os.RemoveAll("./test-items.db") // Clean start
	store, err := db.NewStore("./test-items.db")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// 2. Fetch & Ingest
	id := 4151
	fmt.Printf("Fetching item %d...\n", id)
	url := fmt.Sprintf("https://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item=%d", id)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	item, err := ingest.Transform(body)
	if err != nil {
		panic(err)
	}

	if err := store.SaveItem(item); err != nil {
		panic(err)
	}
	fmt.Printf("Saved: %s (Price: %d)\n", item.Name, item.CurrentPrice)

	// 3. Search
	fmt.Println("Searching for 'Abyssal'...")
	results, err := store.SearchItems("Abyssal", 10)
	if err != nil {
		panic(err)
	}

	if len(results) == 0 {
		fmt.Println("No results found!")
	}

	for _, res := range results {
		fmt.Printf("Found: %s - %s\n", res.Name, res.Description)
	}
}
