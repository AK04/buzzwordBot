package common

import (
	"buzzwordBot/models"

	"golang.org/x/exp/slices"
)

func RemoveDuplicateItems(items []models.Item, itemIds []string) []models.Item {
	var newItems []models.Item
	for _, item := range items {
		if !slices.Contains(itemIds, item.Id) {
			newItems = append(newItems, item)
		}
	}
	return newItems
}
