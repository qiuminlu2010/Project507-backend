package service

type AdminMenuItem struct {
	Title    string          `json:"title"`
	Path     string          `json:"path"`
	Icon     string          `json:"icon"`
	Children []AdminMenuItem `json:"children"`
}
