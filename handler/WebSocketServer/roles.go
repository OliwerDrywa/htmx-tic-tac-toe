package WebSocketServer

var existsP1 bool = false
var existsP2 bool = false

func (c *client) AssignRole() *client {
	if !existsP1 {
		c.Role = 1
		existsP1 = true
	} else if !existsP2 {
		c.Role = 2
		existsP2 = true
	}

	return c
}

func (c *client) UnassignRole() *client {
	if c.Role == 1 {
		c.Role = 0
		existsP1 = false
	} else if c.Role == 2 {
		c.Role = 0
		existsP2 = false
	}

	return c
}
