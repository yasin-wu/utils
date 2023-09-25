package ldap

type Unit struct {
	Name     string `json:"name"`
	DN       string `json:"dn"`
	ParentDN string `json:"parent_dn"`
}

type Person struct {
	Name   string `json:"name"`
	DN     string `json:"dn"`
	OU     string `json:"ou"`
	OULink string `json:"ou_link"`
}
