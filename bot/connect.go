package bot

// Domain services about connections (ex: twitter) are defined here.

type (
	ConnectConn interface {
		GetIconUrl(id string) (string, error)
		GetName(id string) (string, error)
	}
)
