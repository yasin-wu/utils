package ldap

import (
	"bytes"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type Ldap struct {
	addr     string
	username string
	password string
	basedn   string
}

const (
	PersonFilter    = "(|(objectClass=organizationalPerson))"
	UnitFilter      = "(|(objectClass=organizationalUnit))"
	PersonClass     = "organizationalPerson"
	UnitClass       = "organizationalUnit"
	administratorCN = "CN=Administrator,CN=Users,"
)

func New(addr string, username string, password string, basedn string) *Ldap {
	if strings.ToLower(username) == "administrator" {
		username = administratorCN + basedn
	}
	return &Ldap{
		addr:     "ldap://" + addr,
		username: username,
		password: password,
		basedn:   basedn,
	}
}

func (l *Ldap) Search(filter string, pageSize int) ([]*ldap.Entry, error) {
	conn, err := l.conn()
	if err != nil {
		return nil, err
	}
	defer l.close(conn)
	req := ldap.NewSearchRequest(l.basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases,
		0, 0, false,
		filter, []string{"dn", "cn"}, nil,
	)
	if pageSize == 0 {
		pageSize = 1000
	}
	req.Controls = []ldap.Control{
		ldap.NewControlPaging(uint32(pageSize)),
	}
	var entries []*ldap.Entry
	for {
		resp, err := conn.Search(req)
		if err != nil {
			fmt.Println(err)
			break
		}
		entries = append(entries, resp.Entries...)
		if len(resp.Controls) > 0 {
			if c, ok := resp.Controls[0].(*ldap.ControlPaging); ok {
				if len(c.Cookie) == 0 {
					break
				}
				req.Controls[0].(*ldap.ControlPaging).Cookie = c.Cookie
			}
		}
	}
	return entries, nil
}

func (l *Ldap) Add(objectClass string, name ...string) error {
	conn, err := l.conn()
	if err != nil {
		return err
	}
	defer l.close(conn)
	for _, v := range name {
		var buf bytes.Buffer
		buf.WriteString("CN=")
		buf.WriteString(v)
		buf.WriteString(",CN=Users,")
		buf.WriteString(l.basedn)
		dn := ldap.NewAddRequest(buf.String(), nil)
		dn.Attribute("objectClass", []string{objectClass})
		err := conn.Add(dn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Ldap) conn() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(l.addr)
	if err != nil {
		return nil, err
	}
	err = conn.Bind(l.username, l.password)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (l *Ldap) close(conn *ldap.Conn) {
	_ = conn.Close()
}
