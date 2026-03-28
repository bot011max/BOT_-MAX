package dashboard

import (
    "encoding/json"
    "time"
)

type DoctorDashboard struct {
    Patients    []PatientSummary `json:"patients"`
    Statistics  Statistics       `json:"statistics"`
    Alerts      []Alert          `json:"alerts"`
}

type PatientSummary struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    LastVisit     time.Time `json:"last_visit"`
    AdherenceRate float64   `json:"adherence_rate"`
    ActiveMeds    int       `json:"active_meds"`
    RiskLevel     string    `json:"risk_level"`
}

type Statistics struct {
    TotalPatients   int     `json:"total_patients"`
    AvgAdherence    float64 `json:"avg_adherence"`
    HighRiskPatients int    `json:"high_risk_patients"`
    CriticalAlerts  int     `json:"critical_alerts"`
}

type Alert struct {
    ID        string    `json:"id"`
    PatientID string    `json:"patient_id"`
    PatientName string  `json:"patient_name"`
    Type      string    `json:"type"` // missed_dose, low_adherence, symptom
    Message   string    `json:"message"`
    CreatedAt time.Time `json:"created_at"`
    Severity  string    `json:"severity"` // low, medium, high, critical
}

func GenerateDashboard(doctorID string) *DoctorDashboard {
    // Симуляция данных дашборда
    return &DoctorDashboard{
        Patients: []PatientSummary{
            {ID: "1", Name: "Иван Петров", AdherenceRate: 94.5, ActiveMeds: 3, RiskLevel: "low"},
            {ID: "2", Name: "Мария Иванова", AdherenceRate: 78.2, ActiveMeds: 5, RiskLevel: "medium"},
        },
        Statistics: Statistics{
            TotalPatients:   42,
            AvgAdherence:    86.4,
            HighRiskPatients: 5,
            CriticalAlerts:  2,
        },
        Alerts: []Alert{
            {ID: "1", PatientName: "Сергей Сидоров", Type: "missed_dose", Severity: "high"},
        },
    }
}

func (d *DoctorDashboard) ToJSON() ([]byte, error) {
    return json.MarshalIndent(d, "", "  ")
}
