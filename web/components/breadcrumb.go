package components

import (
	"fmt"
	"mt-hosting-manager/types"
)

type BreadcrumbEntry struct {
	Name   string
	Link   string
	FAIcon string
}

type Breadcrumb struct {
	Entries []*BreadcrumbEntry
}

var HomeBreadcrumb = &BreadcrumbEntry{
	Name:   "Home",
	Link:   "/",
	FAIcon: "home",
}

var ProfileBreadcrumb = &BreadcrumbEntry{
	Name:   "Profile",
	Link:   "/profile",
	FAIcon: "user",
}

var LoginBreadcrumb = &BreadcrumbEntry{
	Name:   "Login",
	Link:   "/login",
	FAIcon: "sign-in",
}

var NodeTypesBreadcrumb = &BreadcrumbEntry{
	Name:   "Node types",
	Link:   "/node_types",
	FAIcon: "server",
}

func NodeTypeBreadcrumb(nt *types.NodeType) *BreadcrumbEntry {
	return &BreadcrumbEntry{
		Name:   nt.ID,
		Link:   fmt.Sprintf("/node_types/%s", nt.ID),
		FAIcon: "server",
	}
}

var NodesBreadcrumb = &BreadcrumbEntry{
	Name:   "Nodes",
	Link:   "/nodes",
	FAIcon: "server",
}

func NodeBreadcrumb(node *types.UserNode) *BreadcrumbEntry {
	return &BreadcrumbEntry{
		Name:   fmt.Sprintf("Node '%s'", node.Alias),
		Link:   fmt.Sprintf("/nodes/%s", node.ID),
		FAIcon: "server",
	}
}

func TransactionBreadcrumb(tx *types.PaymentTransaction) *BreadcrumbEntry {
	return &BreadcrumbEntry{
		Name:   fmt.Sprintf("Payment '%s'", tx.ID),
		Link:   fmt.Sprintf("/nodes/payment/%s", tx.ID),
		FAIcon: "money-bill",
	}
}

var ServersBreadcrumb = &BreadcrumbEntry{
	Name:   "Servers",
	Link:   "/mtserver",
	FAIcon: "list",
}

func ServerBreadcrumb(server *types.MinetestServer) *BreadcrumbEntry {
	return &BreadcrumbEntry{
		Name:   fmt.Sprintf("Server '%s'", server.Name),
		Link:   fmt.Sprintf("/mtservers/%s", server.ID),
		FAIcon: "list",
	}
}
