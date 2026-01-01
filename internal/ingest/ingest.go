package ingest

import (
	"encoding/json"
	"strconv"
	"strings"

	"rs-item-database/pb"
)

type RSResponse struct {
	Item RSItem `json:"item"`
}

type RSItem struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Icon        string    `json:"icon"`
	IconLarge   string    `json:"icon_large"`
	Members     string    `json:"members"` // "true"/"false"
	Current     PriceInfo `json:"current"`
	Today       PriceInfo `json:"today"`
}

type PriceInfo struct {
	Trend string      `json:"trend"`
	Price interface{} `json:"price"` // Can be string "75.8k" or number
}

// Transform converts raw JSON bytes from RS API to a Protobuf Item
func Transform(jsonData []byte) (*pb.Item, error) {
	var resp RSResponse
	if err := json.Unmarshal(jsonData, &resp); err != nil {
		return nil, err
	}

	item := resp.Item

	// Convert members string to bool
	isMembers := item.Members == "true"

	// Parse prices
	currentPrice := parsePrice(item.Current.Price)
	todayChange := parsePrice(item.Today.Price)

	pbItem := &pb.Item{
		Id:               item.ID,
		Name:             item.Name,
		Description:      item.Description,
		Type:             item.Type,
		Icon:             item.Icon,
		IconLarge:        item.IconLarge,
		Members:          isMembers,
		CurrentPrice:     currentPrice,
		CurrentTrend:     item.Current.Trend,
		TodayPriceChange: todayChange,
		TodayTrend:       item.Today.Trend,
	}

	return pbItem, nil
}

func parsePrice(p interface{}) int64 {
	s, ok := p.(string)
	if !ok {
		// maybe it's a number?
		if f, ok := p.(float64); ok {
			return int64(f)
		}
		return 0
	}

	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "+", "") // Remove + prefix

	multiplier := 1.0
	if strings.HasSuffix(s, "k") || strings.HasSuffix(s, "K") {
		multiplier = 1000.0
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "m") || strings.HasSuffix(s, "M") {
		multiplier = 1000000.0
		s = s[:len(s)-1]
	} else if strings.HasSuffix(s, "b") || strings.HasSuffix(s, "B") {
		multiplier = 1000000000.0
		s = s[:len(s)-1]
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return int64(val * multiplier)
}
