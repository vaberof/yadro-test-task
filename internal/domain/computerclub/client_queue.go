package computerclub

type ClientQueue struct {
	queue   []*Client
	maxSize int
}

func NewClientQueue(maxSize int) *ClientQueue {
	return &ClientQueue{
		queue:   make([]*Client, 0, maxSize),
		maxSize: maxSize,
	}
}

func (c *ClientQueue) Push(client *Client) {
	c.queue = append(c.queue, client)
}

func (c *ClientQueue) Pop() *Client {
	var client *Client
	if len(c.queue) > 0 {
		client = c.queue[0]
		c.queue = c.queue[1:]
	}
	return client
}

func (c *ClientQueue) IsEmpty() bool {
	return len(c.queue) == 0
}

func (c *ClientQueue) IsFull() bool {
	return len(c.queue) == c.maxSize
}
