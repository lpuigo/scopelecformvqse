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
	ctrlFile           = `C:\Users\Laurent\Golang\src\github.com\lpuig\scopelecformvqse\test\client_ftth.csv`
)

func outfile(infile string) string {
	return strings.TrimSuffix(infile, filepath.Ext(infile)) + ".json"
}

func main() {
	// pattern is the glob pattern used to find all the template files.
	pattern := filepath.Join(templateDir, "*.tmpl")

	// Load the drivers.
	fns := template.FuncMap{"plus1": func(x int) int { return x + 1 }}
	tmpls := template.Must(template.New("abc").Funcs(fns).ParseGlob(pattern))

	ctrl, err := loadControlFromCSV(ctrlFile)
	if err != nil {
		log.Fatalf("could not load Control CSV File: %s", err)
	}

	of, err := os.Create(outfile(ctrlFile))
	if err != nil {
		log.Fatalf("could not create result File: %s", err)
	}
	defer of.Close()

	err = tmpls.ExecuteTemplate(of, "formulaire", ctrl)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
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
	)
	if err != nil {
		return
	}

	var (
		//colCtrl    = cols[0]
		colCtrlGrp = cols[1]
		colItem    = cols[2]
		colAide1   = cols[3]
		colAide2   = cols[4]

		ctrlGrp        = model.ControlGroup{}
		curCtrlGrpName = ""
		ctrlGrpName    = ""
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
		}
		ctrlGrp.Items = append(
			ctrlGrp.Items,
			model.Item{
				Ref:     fmt.Sprintf("%s_%d", ctrlGrpName, len(ctrlGrp.Items)+1),
				Label:   r[colItem],
				Tooltip: fmt.Sprintf("%s. %s.", r[colAide1], r[colAide2]),
			},
		)
	}
	c.Controls = append(c.Controls, ctrlGrp)
	return
}

func convertToRef(s string) string {
	return strings.Replace(strings.Title(strings.ToLower(s)), " ", "", -1)
}
