package models

// CropProfile representa o perfil climático ideal para uma cultura agrícola
type CropProfile struct {
	ID               int     `json:"id"`
	Name             string  `json:"name" example:"Milho"`
	MinTemp          float64 `json:"min_temp" example:"21.00"`
	MaxTemp          float64 `json:"max_temp" example:"32.00"`
	MinPrecipitation float64 `json:"min_precipitation" example:"500.00"`
	MaxPrecipitation float64 `json:"max_precipitation" example:"800.00"`
	IdealSeason      string  `json:"ideal_season" example:"summer"`
}

// RiskAnalysis representa a análise de risco de uma área
type RiskAnalysis struct {
	OverallRiskScore int      `json:"overall_risk_score" example:"45"`
	RiskFactors      []string `json:"risk_factors" example:"Risco de Geada"`
}

// CropRecommendation representa o nível de compatibilidade de uma cultura com a área
type CropRecommendation struct {
	CropName        string  `json:"crop_name" example:"Milho"`
	MatchPercentage float64 `json:"match_percentage" example:"95.5"`
	Justification   string  `json:"justification" example:"Temperatura e precipitação adequados"`
}
