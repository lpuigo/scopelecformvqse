package main

import (
	"fmt"
	"github.com/lpuig/scopelecformvqse/model"
	"github.com/lpuig/scopelecformvqse/recordset"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateDir string = "./template"
	ctrlFile           = `C:\Users\Laurent\Golang\src\github.com\lpuig\scopelecformvqse\test\FTTH_CLIENT.csv`
	ctrlDir            = `C:\Users\Laurent\Golang\src\github.com\lpuig\scopelecformvqse\test`
)

func main() {
	converter := NewConverter()

	files, err := filepath.Glob(filepath.Join(ctrlDir, "*.csv"))
	if err != nil {
		log.Fatalf("could not scan directory: %v", err)
	}

	for _, file := range files {
		err := converter.Convert(file)
		if err != nil {
			log.Printf("could not convert '%s': %v", filepath.Base(file), err)
		}
	}
}

type Converter struct {
	*template.Template
}

func NewConverter() Converter {
	c := Converter{}

	// pattern is the glob pattern used to find all the template files.
	pattern := filepath.Join(templateDir, "*.tmpl")

	// Load the drivers.
	fns := template.FuncMap{"plus1": func(x int) int { return x + 1 }}
	c.Template = template.Must(template.New("abc").Funcs(fns).ParseGlob(pattern))

	return c
}

func (c Converter) Convert(file string) error {
	ctrl, err := loadControlFromCSV(file)
	if err != nil {
		return fmt.Errorf("could not load Control CSV File: %s", err)
	}

	of, err := os.Create(outfile(file))
	if err != nil {
		return fmt.Errorf("could not create result File: %s", err)
	}
	defer of.Close()

	err = c.ExecuteTemplate(of, "formulaire", ctrl)
	if err != nil {
		return fmt.Errorf("template execution: %s", err)
	}
	return nil
}

func outfile(infile string) string {
	return strings.TrimSuffix(infile, filepath.Ext(infile)) + ".json"
}

func loadControlFromCSV(file string) (c model.Control, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	rs := recordset.NewRecordSet()
	err = rs.AddCSVDataFrom(transform.NewReader(f, charmap.Windows1252.NewDecoder()))
	if err != nil {
		return
	}

	cols, err := rs.GetRecordColNumByName(
		"ITEM",
		"CATEGORIE DE CONTROLE",
		"POINT DE CONTROLE",
		"AIDE EN LIGNE",
		"AIDE SEUIL DE NON-CONFORMITE",
		"POIDS",
		"PRIX",
		"POURCENTAGE",
	)
	if err != nil {
		return
	}

	var (
		colSection = cols[0]
		colCtrlGrp = cols[1]
		colItem    = cols[2]
		colAide1   = cols[3]
		colAide2   = cols[4]
		colPoids   = cols[5]
		colPrix    = cols[6]
		colPct     = cols[7]

		ctrlGrp        = model.ControlGroup{}
		curCtrlGrpName = ""
		ctrlGrpName    = ""
		sectionName    = ""
	)

	for i, r := range rs.GetRecords() {
		if r[colCtrlGrp] != curCtrlGrpName {
			if i > 0 {
				c.Controls = append(c.Controls, ctrlGrp)
			}
			curCtrlGrpName = r[colCtrlGrp]
			ctrlGrpName = convertToRef(curCtrlGrpName)
			ctrlGrp = model.ControlGroup{
				Ref:   ctrlGrpName,
				Label: curCtrlGrpName,
			}
			if r[colSection] != sectionName {
				sectionName = r[colSection]
				ctrlGrp.Section = sectionName
			}
		}
		ctrlGrp.Items = append(
			ctrlGrp.Items,
			model.Item{
				Ref:     fmt.Sprintf("%s_%d", ctrlGrpName, len(ctrlGrp.Items)+1),
				Label:   composeItemLabel(r[colItem], r[colPoids], r[colPrix], r[colPct]),
				Tooltip: secureJson(fmt.Sprintf("%s. %s.", r[colAide1], r[colAide2])),
			},
		)
	}
	c.Controls = append(c.Controls, ctrlGrp)
	return
}

func convertToRef(s string) string {
	return strings.Replace(strings.Title(strings.ToLower(s)), " ", "", -1)
}

func secureJson(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}

func composeItemLabel(titre, poids, prix, pct string) string {
	switch {
	case poids != "":
		return fmt.Sprintf("%s (Poids: %s)", titre, poids)
	case prix != "":
		return fmt.Sprintf("%s (%s€)", titre, prix)
	case pct != "":
		return fmt.Sprintf("%s (%s%%)", titre, pct)
	default:
		return titre
	}
}
