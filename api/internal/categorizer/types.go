package categorizer

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

// AllCategories returns a list of all available categories
func AllCategories() []CategoryType {
	return []CategoryType{
		CategoryDining,
		CategoryGroceries,
		CategoryTransport,
		CategoryShopping,
		CategorySubscriptions,
		CategoryEntertainment,
		CategoryBills,
		CategoryTransfer,
		CategoryMisc,
	}
}

// String returns the string representation of the category
func (c CategoryType) String() string {
	return string(c)
}
