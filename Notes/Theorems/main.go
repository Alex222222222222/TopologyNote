package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Definition struct {
	Name          string
	Ref           string
	SubDefinition *Definitions
}

const SubDefinitionPre = "\\hspace{1em}"
const RefPre = "theorem"
const LineBreak = "\\vspace{1em}"

type Definitions struct {
	Def      []*Definition
	SubOrder int
}

func (definitions *Definitions) sort() {
	for i := 0; i < len(definitions.Def)-1; i += 1 {
		for j := len(definitions.Def) - 1; j > i; j -= 1 {
			if strings.ToLower(definitions.Def[j].Ref) < strings.ToLower(definitions.Def[j-1].Ref) {
				definitions.Def[j], definitions.Def[j-1] = definitions.Def[j-1], definitions.Def[j]
			} else if strings.ToLower(definitions.Def[j].Ref) == strings.ToLower(definitions.Def[j-1].Ref) {
				if strings.ToLower(definitions.Def[j].Name) < strings.ToLower(definitions.Def[j-1].Name) {
					definitions.Def[j], definitions.Def[j-1] = definitions.Def[j-1], definitions.Def[j]
				}
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
				"Comparison of the box and product topologies",
				"ComparisonOfBoxProductTopology",
				nil,
			},
			{
				"Uniform limit theorem",
				"UniformLimitTheorem",
				nil,
			},
			{
				"The sequence lemma",
				"TheSequenceLemma",
				nil,
			},
			{
				"Rules for constructing continuous functions",
				"RulesForConstructingContinuousFunctions",
				nil,
			},
			{
				"The pasting lemma",
				"ThePastingLemma",
				nil,
			},
			{
				"Maps into products",
				"MapsIntoProducts",
				nil,
			},
		},
	}

	alphabet := map[string]*Definitions{}

	for i := 0; i < len(definitions.Def); i += 1 {
		if _, ok := alphabet[strings.ToUpper(definitions.Def[i].Ref[0:1])]; !ok {
			alphabet[strings.ToUpper(definitions.Def[i].Ref[0:1])] = &Definitions{
				SubOrder: 0,
				Def:      make([]*Definition, 0, 0),
			}
		}
		alphabet[strings.ToUpper(definitions.Def[i].Ref[0:1])].Def = append(alphabet[strings.ToUpper(definitions.Def[i].Ref[0:1])].Def, definitions.Def[i])
	}

	file, err := os.Create("./theorems.tex")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write([]byte("\\section*{Theorems}\n\n"))

	// file.Write([]byte("\\begin{multicols}{2}\n\n"))

	keys := make([]string, 0, len(alphabet))
	for k := range alphabet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		alphabet[k].sort()

		file.Write([]byte(
			fmt.Sprintf("%s\\noindent\\large{\\textbf{%s}}\n\n", LineBreak, strings.ToUpper(k)),
		))

		file.Write([]byte(alphabet[k].generateString(0)))

	}

	// file.Write([]byte("\\end{multicols}"))

	file.Sync()

}
