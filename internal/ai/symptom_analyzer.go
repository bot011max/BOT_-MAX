package ai

import (
    "strings"
)

type Diagnosis struct {
    Name        string  `json:"name"`
    Probability float64 `json:"probability"`
    Description string  `json:"description"`
    Advice      string  `json:"advice"`
}

type SymptomAnalyzer struct {
    symptomsDB map[string][]string
}

func NewSymptomAnalyzer() *SymptomAnalyzer {
    return &SymptomAnalyzer{
        symptomsDB: map[string][]string{
            "ОРВИ":     {"кашель", "насморк", "температура", "головная боль"},
            "Грипп":    {"высокая температура", "ломота", "кашель", "слабость"},
            "Ангина":   {"боль в горле", "температура", "затрудненное глотание"},
            "Аллергия": {"чихание", "зуд", "слезотечение", "насморк"},
        },
    }
}

func (a *SymptomAnalyzer) Analyze(symptoms []string) []Diagnosis {
    results := []Diagnosis{}
    
    for disease, diseaseSymptoms := range a.symptomsDB {
        matches := 0
        for _, s := range symptoms {
            for _, ds := range diseaseSymptoms {
                if strings.Contains(strings.ToLower(s), strings.ToLower(ds)) {
                    matches++
                }
            }
        }
        
        if matches > 0 {
            probability := float64(matches) / float64(len(diseaseSymptoms))
            results = append(results, Diagnosis{
                Name:        disease,
                Probability: probability,
                Description: getDescription(disease),
                Advice:      getAdvice(disease),
            })
        }
    }
    
    return results
}

func getDescription(disease string) string {
    descriptions := map[string]string{
        "ОРВИ":     "Острая респираторная вирусная инфекция",
        "Грипп":    "Острое инфекционное заболевание дыхательных путей",
        "Ангина":   "Воспаление небных миндалин",
        "Аллергия": "Реакция иммунной системы на аллергены",
    }
    return descriptions[disease]
}

func getAdvice(disease string) string {
    advices := map[string]string{
        "ОРВИ":     "Рекомендуется: постельный режим, обильное питье, витамин C",
        "Грипп":    "Рекомендуется: обратиться к врачу, противовирусные препараты",
        "Ангина":   "Рекомендуется: полоскание горла, антибиотики по назначению врача",
        "Аллергия": "Рекомендуется: антигистаминные препараты, избегать аллергенов",
    }
    return advices[disease]
}
