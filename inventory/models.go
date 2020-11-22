package inventory

type Inventory struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateInventory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
