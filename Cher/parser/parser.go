package parser

type ParsedCommands struct {
	ParsedCommands map[string]bool
}

type Parser struct {
	Commands []string
}

func NewParser() *Parser {
	p := new(Parser)

	p.Commands = make([]string, 0)

	p.Commands = append(p.Commands, "init")
	p.Commands = append(p.Commands, "help")
	p.Commands = append(p.Commands, "add")
	p.Commands = append(p.Commands, "new")
	p.Commands = append(p.Commands, "remove")
	p.Commands = append(p.Commands, "edit")

	return p
}

func (p *Parser) Parse(args []string) (*ParsedCommands, error) {
	ParsedCommands := new(ParsedCommands)
	ParsedCommands.ParsedCommands = make(map[string]bool)
	for j := 0; j < len(p.Commands); j++ {
		ParsedCommands.ParsedCommands[p.Commands[j]] = false
	}

	for i := 0; i < len(args); i++ {
		for j := 0; j < len(p.Commands); j++ {
			if args[i] == p.Commands[j] {
				ParsedCommands.ParsedCommands[p.Commands[j]] = true
			}
		}
	}

	return ParsedCommands, nil
}
