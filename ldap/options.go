package ldap

const (
	ScopeBaseObject   = "base"
	ScopeSingleLevel  = "single"
	ScopeWholeSubtree = "sub"
)

var scopeMap = map[string]int{
	ScopeBaseObject:   0,
	ScopeSingleLevel:  1,
	ScopeWholeSubtree: 2,
}
