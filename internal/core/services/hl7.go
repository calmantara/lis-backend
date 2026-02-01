package services

import (
	"strings"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/utils"

	"github.com/yehezkel/gohl7"
)

type HL7TestResult struct {
	Patient struct {
		Identifiers []struct {
			ID             string `json:"id"`
			IdentifierType string `json:"identifier_type"`
		} `json:"identifiers"`
		Name struct {
			FamilyName string `json:"family_name"`
			GivenName  string `json:"given_name"`
			MiddleName string `json:"middle_name"`
		} `json:"name"`
		Dob     string `json:"dob"`
		Sex     string `json:"sex"`
		Address struct {
			Street  string `json:"street"`
			City    string `json:"city"`
			State   string `json:"state"`
			Zip     string `json:"zip"`
			Country string `json:"country"`
		} `json:"address"`
	} `json:"patient"`
	Tests []struct {
		Code           string `json:"code"`
		Name           string `json:"name"`
		CodingSystem   string `json:"coding_system"`
		Value          string `json:"value"`
		Units          string `json:"units"`
		ReferenceRange string `json:"reference_range"`
		AbnormalFlags  string `json:"abnormal_flags"`
	} `json:"tests"`
	PrimaryOrderDatetime time.Time `json:"primary_order_datetime"`
}

func (s *HL7TestResult) Stringify() string {
	return ""
}

func (s *HL7TestResult) Serialize(deviceMessage models.DeviceMessage) models.Serializer {
	seqNum := 0
	patientID := ""
	if len(s.Patient.Identifiers) > 0 {
		patientID = s.Patient.Identifiers[0].ID
		seqNum = utils.FindAllInteger(s.Patient.Identifiers[0].ID)
	}

	res := models.Serializer{
		DeviceID:       deviceMessage.DeviceID,
		Protocol:       string(deviceMessage.Protocol),
		DeviceTypeCode: deviceMessage.DeviceTypeCode,
		SequenceNumber: seqNum,
		PatientID:      patientID,
		Timestamp:      s.PrimaryOrderDatetime,
	}

	// transform result
	for _, test := range s.Tests {
		res.Results = append(res.Results, models.Result{
			ParameterCode:  test.Code,
			ParameterName:  test.Name,
			Value:          test.Value,
			Unit:           test.Units,
			ReferenceRange: test.ReferenceRange,
			AbnormalFlags:  test.AbnormalFlags,
		})
	}

	return res
}

func parseHL7Message(deviceMessage models.DeviceMessage) (res *HL7TestResult, err error) {
	msg := strings.TrimSpace(deviceMessage.Message)
	parser, err := gohl7.NewHl7Parser([]byte(msg))
	if err != nil {
		return nil, err
	}

	hl7Payload, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	res = &HL7TestResult{}
	err = utils.ObjectMapper(&hl7Payload, &res)

	return
}
