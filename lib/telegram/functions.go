package telegram

func (m IncomingMessage) GetCommand() (Command, bool) {
	for _, e := range m.Entities {
		if e.Type == "bot_command" {
			return Command{
				Command: m.Text[e.Offset : e.Offset+e.Length],
				Extra:   m.Text[e.Offset+e.Length:],
			}, true
		}
	}
	return Command{}, false
}

func (l Coordinates) IsZero() bool {
	return l.Latitude == 0 && l.Longitude == 0
}
