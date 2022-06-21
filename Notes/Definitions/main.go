package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Definition struct {
	Name          string
	Ref           string
	SubDefinition *Definitions
}

const SubDefinitionPre = "\\hspace{2em}"
const RefPre = "def"
const LineBreak = "\\vspace{1em}"

type Definitions struct {
	Def      []*Definition
	SubOrder int
}

func (definitions *Definitions) sort() {
	for i := 0; i < len(definitions.Def)-1; i += 1 {
		for j := len(definitions.Def) - 1; j > i; j -= 1 {
			if definitions.Def[j].Name < definitions.Def[j-1].Name {
				definitions.Def[j], definitions.Def[j-1] = definitions.Def[j-1], definitions.Def[j]
			}
		}
	}
}

func (definitions *Definitions) generateString(subOrder int) string {
	definitions.SubOrder = subOrder
	res := ""
	pre := ""
	for i := 0; i < definitions.SubOrder; i += 1 {
		pre += SubDefinitionPre
	}
	for i := 0; i < len(definitions.Def); i += 1 {
		res += fmt.Sprintf(
			"%s%s, \\pageref{%s:%s}\n\n",
			pre,
			definitions.Def[i].Name,
			RefPre,
			definitions.Def[i].Ref,
		)

		if definitions.Def[i].SubDefinition != nil {
			definitions.Def[i].SubDefinition.sort()
			res += definitions.Def[i].SubDefinition.generateString(definitions.SubOrder + 1)
		}

	}

	return res
}

func main() {
	definitions := &Definitions{
		SubOrder: 0,
		Def: []*Definition{
			{
				"topology",
				"Topology",
				nil,
			},
			{
				"coarser",
				"Comparable",
				&Definitions{
					SubOrder: 1,
					Def: []*Definition{
						{
							"strictly coarser",
							"Comparable",
							nil,
						},
					},
				},
			},
			{
				"convex",
				"Convex",
				nil,
			},
			{
				"ordered square",
				"OrderedSquare",
				nil,
			},
			{
				"subspace topology",
				"SubspaceTopology",
				nil,
			},
			{
				"subspace",
				"SubspaceTopology",
				nil,
			},
			{
				"projection",
				"Projection",
				nil,
			},
			{
				"product topology",
				"ProductTopology",
				nil,
			},
			{
				"ray",
				"Ray",
				&Definitions{
					SubOrder: 1,
					Def: []*Definition{
						{
							"open ray",
							"Ray",
							nil,
						},
						{
							"closed ray",
							"Ray",
							nil,
						},
					},
				},
			},
			{
				"order topology",
				"OrderTopology",
				nil,
			},
			{
				"interval",
				"Interval",
				&Definitions{
					SubOrder: 1,
					Def: []*Definition{
						{
							"open interval",
							"Interval",
							nil,
						},
						{
							"closed interval",
							"Interval",
							nil,
						},
						{
							"half-open interval",
							"Interval",
							nil,
						},
					},
				},
			},
			{
				"topology space",
				"TopologySpace",
				nil,
			},
			{
				"open set",
				"OpenSet",
				nil,
			},
			{
				"open sets",
				"OpenSets",
				nil,
			},
			{
				"discrete topology",
				"DiscreteTopology",
				nil,
			},
			{
				"trivial topology",
				"TrivialTopology",
				nil,
			},
			{
				"finite complement topology",
				"FiniteComplementTopology",
				nil,
			},
			{
				"basis",
				"Basis",
				nil,
			},
			{
				"subbasis",
				"Subbasis",
				nil,
			},
			{
				"K-topology on R",
				"KTopologyOnTheRealLine",
				nil,
			},
			{
				"lower limit topology on R",
				"LowerLimitTopologyOnTheRealLine",
				nil,
			},
			{
				"standard topology on R",
				"StandardTopologyOnTheRealLine",
				nil,
			},
			{
				"topology generated by basis",
				"TopologyGeneratedByBasis",
				nil,
			},
			{
				"finer",
				"Comparable",
				&Definitions{
					SubOrder: 1,
					Def: []*Definition{
						{
							"strictly finer",
							"Comparable",
							nil,
						},
					},
				},
			},
			{
				"larger",
				"Comparable",
				&Definitions{
					SubOrder: 1,
					Def: []*Definition{
						{
							"strictly larger",
							"Comparable",
							nil,
						},
					},
				},
			},
			{
				"smaller",
				"Comparable",
				&Definitions{
					SubOrder: 1,
					Def: []*Definition{
						{
							"strictly smaller",
							"Comparable",
							nil,
						},
					},
				},
			},
		},
	}

	alphabet := map[string]*Definitions{}

	for i := 0; i < len(definitions.Def); i += 1 {
		if _, ok := alphabet[strings.ToUpper(definitions.Def[i].Name[0:1])]; !ok {
			alphabet[strings.ToUpper(definitions.Def[i].Name[0:1])] = &Definitions{
				SubOrder: 0,
				Def:      make([]*Definition, 0, 0),
			}
		}
		alphabet[strings.ToUpper(definitions.Def[i].Name[0:1])].Def = append(alphabet[strings.ToUpper(definitions.Def[i].Name[0:1])].Def, definitions.Def[i])
	}

	file, err := os.Create("./definitions.tex")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write([]byte("\\section*{Definitions}\n\n"))

	file.Write([]byte("\\begin{multicols}{2}\n\n"))

	for k := range alphabet {
		alphabet[k].sort()

		file.Write([]byte(
			fmt.Sprintf("%s\\large{\\textbf{%s}}\n\n", LineBreak, strings.ToUpper(k)),
		))

		file.Write([]byte(alphabet[k].generateString(0)))

	}

	file.Write([]byte("\\end{multicols}"))

	file.Sync()

}
