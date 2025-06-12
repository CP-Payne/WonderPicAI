package viewmodel

type PurchaseOption struct {
	Credits   int    // e.g 100, 200, 1000
	Price     int    // e.g., "$5", "€10"
	ActionURL string // For HTMX or link later
}

type PurchaseViewData struct {
	Options []PurchaseOption
}
