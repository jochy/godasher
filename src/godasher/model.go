package godasher

type Config struct {
	Port           int
	Rotationtime   int
	Theme          string
	Dashboards     []Dashboard
	ExternalConfig map[string]map[string]string
}

type Dashboard struct {
	Title      string
	Components []Component
}

type Component struct {
	Type   string
	View   string
	Title  string
	Width  int
	Height int
	Column int
	Row    int
	Data   map[string]string
}

type RequestContext struct {
	Config          Config
	Dashboard       Dashboard
	DashboardNumber int
}
