package services

import (
	"strings"
)

// CategoryType represents app-level transaction categories
type CategoryType string

const (
	CategoryDining        CategoryType = "Dining"
	CategoryGroceries     CategoryType = "Groceries"
	CategoryTransport     CategoryType = "Transport"
	CategoryShopping      CategoryType = "Shopping"
	CategorySubscriptions CategoryType = "Subscriptions"
	CategoryEntertainment CategoryType = "Entertainment"
	CategoryBills         CategoryType = "Bills"
	CategoryTransfer      CategoryType = "Transfer"
	CategoryMisc          CategoryType = "Misc"
)

// CategorizeTransaction maps Plaid categories to app categories
// Based on the PLAN.md specification for spending categories
func CategorizeTransaction(plaidCategories []string) CategoryType {
	if len(plaidCategories) == 0 {
		return CategoryMisc
	}

	// Get the primary category (first element)
	primary := strings.ToLower(plaidCategories[0])

	// Check subcategories for more specific matching
	var secondary string
	if len(plaidCategories) > 1 {
		secondary = strings.ToLower(plaidCategories[1])
	}

	// Map Plaid categories to app categories
	switch {
	// Transfer - ignore for rules (internal transfers)
	case strings.Contains(primary, "transfer"):
		return CategoryTransfer
	case strings.Contains(primary, "payment"):
		return CategoryTransfer

	// Dining
	case strings.Contains(primary, "food and drink"):
		// Check if it's groceries
		if strings.Contains(secondary, "groceries") || strings.Contains(secondary, "supermarket") {
			return CategoryGroceries
		}
		return CategoryDining
	case strings.Contains(primary, "restaurants"):
		return CategoryDining

	// Groceries
	case strings.Contains(primary, "groceries"):
		return CategoryGroceries
	case strings.Contains(secondary, "supermarket"):
		return CategoryGroceries

	// Transport
	case strings.Contains(primary, "transportation"):
		return CategoryTransport
	case strings.Contains(primary, "travel"):
		return CategoryTransport
	case strings.Contains(secondary, "gas"):
		return CategoryTransport
	case strings.Contains(secondary, "parking"):
		return CategoryTransport
	case strings.Contains(secondary, "public transit"):
		return CategoryTransport
	case strings.Contains(secondary, "ride share"):
		return CategoryTransport

	// Shopping
	case strings.Contains(primary, "shops"):
		return CategoryShopping
	case strings.Contains(primary, "retail"):
		return CategoryShopping
	case strings.Contains(secondary, "clothing"):
		return CategoryShopping
	case strings.Contains(secondary, "electronics"):
		return CategoryShopping

	// Subscriptions
	case strings.Contains(secondary, "subscription"):
		return CategorySubscriptions
	case strings.Contains(primary, "service"):
		// Check if it's a recurring service
		if strings.Contains(secondary, "streaming") ||
			strings.Contains(secondary, "music") ||
			strings.Contains(secondary, "software") {
			return CategorySubscriptions
		}
		return CategoryMisc

	// Entertainment
	case strings.Contains(primary, "recreation"):
		return CategoryEntertainment
	case strings.Contains(primary, "entertainment"):
		return CategoryEntertainment
	case strings.Contains(secondary, "movie"):
		return CategoryEntertainment
	case strings.Contains(secondary, "concert"):
		return CategoryEntertainment
	case strings.Contains(secondary, "sporting"):
		return CategoryEntertainment

	// Bills
	case strings.Contains(primary, "bank fees"):
		return CategoryBills
	case strings.Contains(primary, "interest"):
		return CategoryBills
	case strings.Contains(secondary, "utilities"):
		return CategoryBills
	case strings.Contains(secondary, "internet"):
		return CategoryBills
	case strings.Contains(secondary, "phone"):
		return CategoryBills
	case strings.Contains(secondary, "insurance"):
		return CategoryBills
	case strings.Contains(secondary, "rent"):
		return CategoryBills

	default:
		return CategoryMisc
	}
}

// IsIgnorable checks if a transaction should be ignored for rule processing
// Transfers and very small amounts (under $3) can be ignored per PLAN.md
func IsIgnorable(amount int64, category CategoryType, minAmount int64) bool {
	// Ignore transfers
	if category == CategoryTransfer {
		return true
	}

	// Ignore small transactions if configured
	if minAmount > 0 && amount < minAmount {
		return true
	}

	return false
}
