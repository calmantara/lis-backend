package models

import "time"

type Message interface {
	Stringify() string
	Serialize(DeviceMessage) Serializer
}

type Serializer struct {
	DeviceID       string    `json:"device_id"`
	SequenceNumber int       `json:"sequence_number"`
	Protocol       string    `json:"protocol"`
	PatientID      string    `json:"patient_id"`
	Timestamp      time.Time `json:"timestamp"`
	DeviceTypeCode string    `json:"device_type_code"`
	Results        []Result  `json:"results"`
}

type Result struct {
	ParameterCode  string  `json:"parameter_code"`
	ParameterName  string  `json:"parameter_name"`
	Value          string  `json:"value"`
	NumericValue   float64 `json:"numeric_value"`
	Unit           string  `json:"unit"`
	Qualitative    string  `json:"qualitative"`
	ReferenceRange string  `json:"reference_range"`
	AbnormalFlags  string  `json:"abnormal_flag"`
}
