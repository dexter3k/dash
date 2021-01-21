package avm2

const (
	linkDelimiter = ":"
)

type link struct {
	value     string
	delimiter int
}

func Link(name, space string) link {
	return link{
		value:     space + linkDelimiter + name,
		delimiter: len(space),
	}
}

func (l *link) Name() string {
	return l.value[l.delimiter+len(linkDelimiter):]
}

func (l *link) Space() string {
	return l.value[:l.delimiter]
}

func (l *link) Link() string {
	return l.value
}

func (l *link) IsProtected() bool {
	if l.delimiter == 0 {
		panic("No space")
	}
	return l.value[0] == 'R'
}
