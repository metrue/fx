package config

// Host host entity
type Host struct {
	Host        string
	User        string
	Password    string
	Enabled     bool
	Provisioned bool
}

// NewHost new a host
func NewHost(addr, user, password string) Host {
	return Host{
		Host:        addr,
		User:        user,
		Password:    password,
		Enabled:     false,
		Provisioned: false,
	}
}

// Valid if host is valid
func (h Host) Valid() bool {
	// TODO stronger check
	return h.Host != ""
}

// IsLocal if host is localhost
func (h Host) IsLocal() bool {
	if !h.Valid() {
		return false
	}
	return h.Host == "127.0.0.1" || h.Host == "localhost"
}

// IsRemote is host is remote
func (h Host) IsRemote() bool {
	return !h.IsLocal()
}
