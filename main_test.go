package main

import "github.com/lpuig/scopelecformvqse/model"

func makeControl() model.Control {
	return model.Control{
		Controls: []model.ControlGroup{
			model.ControlGroup{
				Ref:   "PM_CLIENT",
				Label: "PM CLIENT",
				Items: []model.Item{
					model.Item{
						Ref:     "PM_CLIENT_1",
						Label:   "FERMETURE DU PMI OU PM EXTERIEUR",
						Tooltip: "Porte ou capot non refermé. Non conforme uniquement si certitude que l'ETL est en cause.",
					},
					model.Item{
						Ref:     "PM_CLIENT_2",
						Label:   "POSITION CONFORME AU SI",
						Tooltip: "Vérifier la conformité par rapport à l'OT. Si différent, vérifier obligatoirement la concordance avec IPON. Non conforme si le raccordement terrain est différent d'IPON.",
					},
				},
			},
			model.ControlGroup{
				Ref:   "PB_IMMEUBLE",
				Label: "PB IMMEUBLE",
				Items: []model.Item{
					model.Item{
						Ref:     "PB_IMMEUBLE_1",
						Label:   "ÉPISSURE SOUDÉE",
						Tooltip: "Pas de raccord mécanique!!!Si quadri fibres : les 4 fibres sont soudées avec symétrie des couleurs. Non conforme si les points cités ne sont pas respectés.",
					},
					model.Item{
						Ref:     "PB_IMMEUBLE_2",
						Label:   "FERMETURE DU PBI",
						Tooltip: "Vérifier que le PB est correctement refermé. Non conforme si certitude que l'ETL n'a pas refermé le PB.",
					},
				},
			},
		},
	}
}
