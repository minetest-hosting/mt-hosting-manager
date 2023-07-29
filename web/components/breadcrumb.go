package components

type BreadcrumbEntry struct {
	Name   string
	Link   string
	Active bool
}

type Breadcrumb struct {
	Entries []*BreadcrumbEntry
}
