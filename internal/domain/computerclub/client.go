package computerclub

type ClientName string

func (clientName *ClientName) String() string {
	return string(*clientName)
}

type Client struct {
	Name        ClientName
	State       uint8
	BusyTableId TableId
}
