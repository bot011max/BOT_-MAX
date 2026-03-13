package models

import (
    "time"
    "github.com/google/uuid"
)

type Prescription struct {
    ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    PatientID    uuid.UUID  `json:"patient_id" gorm:"index;not null"`
    DoctorID     uuid.UUID  `json:"doctor_id" gorm:"index;not null"`
    
    Name         string     `json:"name" gorm:"not null"`
    Dosage       string     `json:"dosage"`
    Form         string     `json:"form"`
    Frequency    string     `json:"frequency"`
    Duration     string     `json:"duration"`
    Instructions string     `json:"instructions"`
    
    StartDate    *time.Time `json:"start_date"`
    EndDate      *time.Time `json:"end_date"`
    IsActive     bool       `json:"is_active" gorm:"default:true"`
    
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    
    Patient      Patient    `json:"patient" gorm:"foreignKey:PatientID"`
    Doctor       Doctor     `json:"doctor" gorm:"foreignKey:DoctorID"`
}

type Reminder struct {
    ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    PrescriptionID uuid.UUID  `json:"prescription_id" gorm:"index;not null"`
    PatientID      uuid.UUID  `json:"patient_id" gorm:"index;not null"`
    
    ScheduledTime  time.Time  `json:"scheduled_time"`
    Message        string     `json:"message"`
    Status         string     `json:"status" gorm:"default:'pending'"` // pending, sent, acknowledged
    
    SentAt         *time.Time `json:"sent_at"`
    AcknowledgedAt *time.Time `json:"acknowledged_at"`
    
    CreatedAt      time.Time  `json:"created_at"`
    
    Prescription   Prescription `json:"prescription" gorm:"foreignKey:PrescriptionID"`
    Patient        Patient      `json:"patient" gorm:"foreignKey:PatientID"`
}
