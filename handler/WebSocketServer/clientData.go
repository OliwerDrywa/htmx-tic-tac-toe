package WebSocketServer

func (c *client) AssignRole() *client {
	if !c.server.existsP1 {
		c.Role = 1
		c.server.existsP1 = true
	} else if !c.server.existsP2 {
		c.Role = 2
		c.server.existsP2 = true
	}

	return c
}

func (c *client) UnassignRole() *client {
	if c.Role == 1 {
		c.Role = 0
		c.server.existsP1 = false
	} else if c.Role == 2 {
		c.Role = 0
		c.server.existsP2 = false
	}

	return c
}

func (c *client) SetName(name string) *client {
	c.Name = name
	return c
}

func (s *Server) GetClients() (cs []*client) {
	for _, c := range s.clients {
		if c.Name != "" {
			cs = append(cs, c)
		}
	}
	return cs
}
