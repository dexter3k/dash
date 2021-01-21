package avm2

type Global struct {
	Object_Object

	Script *Script
	Domain *Domain
}

func (g *Global) GetDebugName() string {
	return "[global " + g.DebugName + "]"
}
