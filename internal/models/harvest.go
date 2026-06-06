package models

// HarvestPlan representa a recomendação gerada pelo motor da AgriAI
type HarvestPlan struct {
	ViabilityScore int    `json:"viability_score" example:"85"`     // 0 a 100
	Recommendation string `json:"recommendation" example:"Colheita recomendada. Condições ideais de temperatura e sem chuvas."`
	Temperature    float64 `json:"temperature"`
	Precipitation  float64 `json:"precipitation"`
}
