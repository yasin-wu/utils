package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type LDAP struct {
	addr     string
	network  string
	username string
	password string
	baseDN   string
	dc       string
	pageSize uint32
}

const (
	UnitClass   = "organizationalUnit"
	PersonClass = "organizationalPerson"
)

func New(addr, username, password, baseDN string) *LDAP {
	baseDN = toUpper(baseDN)
	dc := dc(baseDN)
	if strings.ToLower(username) == "administrator" {
		username = "CN=Administrator,CN=Users," + dc
	}
	return &LDAP{
		addr:     addr,
		network:  "tcp",
		username: username,
		password: password,
		baseDN:   baseDN,
		dc:       dc,
		pageSize: 5000,
	}
}

func (l *LDAP) SearchGroup() ([]*GroupResult, error) {
	entries, err := l.search(UnitClass)
	if err != nil {
		return nil, err
	}
	var result []*GroupResult
	for _, v := range entries {
		if !strings.HasPrefix(v.DN, "OU=") {
			continue
		}
		name, parentDN := l.handleOU(v.DN)
		if l.isRoot(v.DN) {
			parentDN = ""
			name = l.rootName(v.DN)
		}
		result = append(result, &GroupResult{
			Name:     name,
			DN:       v.DN,
			ParentDN: parentDN,
		})
	}
	return result, nil
}

func (l *LDAP) SearchPerson() ([]*PersonResult, error) {
	entries, err := l.search(PersonClass)
	if err != nil {
		return nil, err
	}
	var result []*PersonResult
	for _, v := range entries {
		if !strings.HasPrefix(v.DN, "CN=") {
			continue
		}
		name, ou, ouLink := l.handleCN(v.DN)
		result = append(result, &PersonResult{
			Name:   name,
			DN:     v.DN,
			OU:     ou,
			OULink: ouLink,
		})
	}
	return result, nil
}

func (l *LDAP) search(objectClass string) ([]*ldap.Entry, error) {
	conn, err := l.conn()
	if err != nil {
		return nil, err
	}
	defer l.close(conn)
	attributes := []string{"DN", "CN"}
	filter := fmt.Sprintf("(&(objectClass=%s))", objectClass)
	req := ldap.NewSearchRequest(l.baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false, filter, attributes, nil)
	resp, err := conn.SearchWithPaging(req, l.pageSize)
	if err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

func (l *LDAP) conn() (*ldap.Conn, error) {
	conn, err := ldap.Dial(l.network, l.addr)
	if err != nil {
		return nil, err
	}
	err = conn.Bind(l.username, l.password)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (l *LDAP) close(conn *ldap.Conn) {
	_ = conn.Close()
}

func (l *LDAP) isRoot(dn string) bool {
	dns := strings.Split(strings.ReplaceAll(dn, ","+l.dc, ""), ",")
	return len(dns) == 1
}

func (l *LDAP) rootName(dn string) string {
	var tmp []string
	for _, v := range strings.Split(dn, ",") {
		if arr := strings.Split(v, "="); len(arr) == 2 {
			tmp = append(tmp, arr[1])
		}
	}
	return strings.Join(tmp, ",")
}

func (l *LDAP) handleOU(dn string) (string, string) {
	ou := ""
	var parentOU []string
	for i, v := range strings.Split(dn, ",") {
		if i == 0 {
			if ous := strings.Split(v, "="); len(ous) == 2 {
				ou = ous[1]
			}
		} else {
			parentOU = append(parentOU, v)
		}

	}
	return ou, strings.Join(parentOU, ",")
}

func (l *LDAP) handleCN(dn string) (string, string, string) {
	cn := ""
	ou := ""
	var ouLinks []string
	for i, v := range strings.Split(dn, ",") {
		if i == 0 {
			if cns := strings.Split(v, "="); len(cns) == 2 {
				cn = cns[1]
			}
		} else if i == 1 {
			if ous := strings.Split(v, "="); len(ous) == 2 {
				ou = ous[1]
			}
			ouLinks = append(ouLinks, v)
		} else {
			ouLinks = append(ouLinks, v)
		}
	}
	return cn, ou, strings.Join(ouLinks, ",")
}

func toUpper(baseDN string) string {
	baseDN = strings.ReplaceAll(baseDN, "ou=", "OU=")
	baseDN = strings.ReplaceAll(baseDN, "dc=", "DC=")
	return baseDN
}

func dc(baseDN string) string {
	var dcs []string
	for _, v := range strings.Split(baseDN, ",") {
		if strings.Contains(v, "DC") {
			dcs = append(dcs, v)
		}
	}
	return strings.Join(dcs, ",")
}
