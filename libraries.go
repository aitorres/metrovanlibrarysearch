package main

// Library is one entry in the registry.
type Library struct {
	Name    string
	Adapter Adapter
}

// Libraries is the ordered list of all supported library systems.
// Order is preserved in output.
var Libraries = []Library{
	{Name: "Vancouver Public Library", Adapter: &BiblioCommonsAdapter{Subdomain: "vpl"}},
	{Name: "Burnaby Public Library", Adapter: &BiblioCommonsAdapter{Subdomain: "burnaby"}},
	{Name: "Surrey Libraries", Adapter: &BiblioCommonsAdapter{Subdomain: "surrey"}},
	{Name: "New Westminster Public Library", Adapter: &BiblioCommonsAdapter{Subdomain: "newwestminster"}},
	{Name: "Richmond Public Library", Adapter: &BiblioCommonsAdapter{Subdomain: "yourlibrary"}},
	{Name: "North Vancouver City Library", Adapter: &BiblioCommonsAdapter{Subdomain: "nvcl"}},
	{Name: "West Vancouver Memorial Library", Adapter: &BiblioCommonsAdapter{Subdomain: "westvanlibrary"}},
	{Name: "Coquitlam Public Library", Adapter: &BiblioCommonsAdapter{Subdomain: "cml"}},
	{Name: "Port Moody Public Library", Adapter: &BiblioCommonsAdapter{Subdomain: "portmoody"}},
}
