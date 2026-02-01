package services

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/utils"
)

type UrineTestResult struct {
	RawMessage string    `json:"raw_message"`
	SpecimenID string    `json:"specimen_id"`
	DateTime   time.Time `json:"date_time"`

	// Test parameters
	Leukocytes      TestValue `json:"leukocytes"`
	Ketones         TestValue `json:"ketones"`
	Nitrites        TestValue `json:"nitrites"`
	Urobilinogen    TestValue `json:"urobilinogen"`
	Bilirubin       TestValue `json:"bilirubin"`
	Protein         TestValue `json:"protein"`
	Glucose         TestValue `json:"glucose"`
	Blood           TestValue `json:"blood"`
	AscorbicAcid    TestValue `json:"ascorbic_acid"`
	SpecificGravity float64   `json:"specific_gravity"`
	PH              float64   `json:"ph"`
}

type TestValue struct {
	Value       string  `json:"value"`
	Unit        string  `json:"unit,omitempty"`
	Numeric     float64 `json:"numeric,omitempty"`
	IsNumeric   bool    `json:"is_numeric"`
	Qualitative string  `json:"qualitative,omitempty"`
}

func (result *UrineTestResult) Stringify() string {
	b, _ := json.Marshal(result)

	return string(b)
}

func (result *UrineTestResult) Serialize(deviceMessage models.DeviceMessage) models.Serializer {
	seqNum := utils.FindAllInteger(result.SpecimenID)

	res := models.Serializer{
		DeviceID:       deviceMessage.DeviceID,
		Protocol:       string(deviceMessage.Protocol),
		DeviceTypeCode: deviceMessage.DeviceTypeCode,
		SequenceNumber: seqNum,
		PatientID:      result.SpecimenID,
		Timestamp:      result.DateTime,
	}

	// Leukocytes      TestValue `json:"leukocytes"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "leukocytes",
		Unit:          result.Leukocytes.Unit,
		Value:         result.Leukocytes.Value,
		NumericValue:  result.Leukocytes.Numeric,
		Qualitative:   result.Leukocytes.Qualitative,
	})

	// Ketones         TestValue `json:"ketones"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "ketones",
		Unit:          result.Ketones.Unit,
		Value:         result.Ketones.Value,
		NumericValue:  result.Ketones.Numeric,
		Qualitative:   result.Ketones.Qualitative,
	})

	// Nitrites        TestValue `json:"nitrites"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "nitrites",
		Unit:          result.Nitrites.Unit,
		Value:         result.Nitrites.Value,
		NumericValue:  result.Nitrites.Numeric,
		Qualitative:   result.Nitrites.Qualitative,
	})

	// Urobilinogen    TestValue `json:"urobilinogen"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "urobilinogen",
		Unit:          result.Urobilinogen.Unit,
		Value:         result.Urobilinogen.Value,
		NumericValue:  result.Urobilinogen.Numeric,
		Qualitative:   result.Urobilinogen.Qualitative,
	})

	// Bilirubin       TestValue `json:"bilirubin"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "bilirubin",
		Unit:          result.Bilirubin.Unit,
		Value:         result.Bilirubin.Value,
		NumericValue:  result.Bilirubin.Numeric,
		Qualitative:   result.Bilirubin.Qualitative,
	})

	// Protein         TestValue `json:"protein"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "protein",
		Unit:          result.Protein.Unit,
		Value:         result.Protein.Value,
		NumericValue:  result.Protein.Numeric,
		Qualitative:   result.Protein.Qualitative,
	})

	// Glucose         TestValue `json:"glucose"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "glucose",
		Unit:          result.Glucose.Unit,
		Value:         result.Glucose.Value,
		NumericValue:  result.Glucose.Numeric,
		Qualitative:   result.Glucose.Qualitative,
	})

	// Blood           TestValue `json:"blood"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "blood",
		Unit:          result.Blood.Unit,
		Value:         result.Blood.Value,
		NumericValue:  result.Blood.Numeric,
		Qualitative:   result.Blood.Qualitative,
	})

	// AscorbicAcid    TestValue `json:"ascorbic_acid"`
	res.Results = append(res.Results, models.Result{
		ParameterCode: "ascorbic_acid",
		Unit:          result.AscorbicAcid.Unit,
		Value:         result.AscorbicAcid.Value,
		NumericValue:  result.AscorbicAcid.Numeric,
		Qualitative:   result.AscorbicAcid.Qualitative,
	})

	// SpecificGravity
	res.Results = append(res.Results, models.Result{
		ParameterCode: "specific_gravity",
		NumericValue:  result.SpecificGravity,
	})

	// PH
	res.Results = append(res.Results, models.Result{
		ParameterCode: "ph",
		NumericValue:  result.PH,
	})

	return res
}

func parseUrineTestText(text string) (*UrineTestResult, error) {
	lines := strings.Split(text, "\n")
	result := &UrineTestResult{}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Parse specimen ID (first line)
		if result.SpecimenID == "" && strings.HasPrefix(line, "NO.") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				result.SpecimenID = parts[0]
			}
			if len(parts) >= 2 {
				// Parse date
				if date, err := time.Parse("2006-01-02", parts[1]); err == nil {
					result.DateTime = date
				}
			}
			continue
		}

		// Parse time (if on its own line)
		if result.DateTime != (time.Time{}) && strings.Contains(line, ":") {
			if t, err := time.Parse("15:04:05", line); err == nil {
				// Combine date and time
				combined := time.Date(
					result.DateTime.Year(),
					result.DateTime.Month(),
					result.DateTime.Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
					0,
					result.DateTime.Location(),
				)
				result.DateTime = combined
			}
			continue
		}

		// Parse test parameters
		parseTestLine(line, result)
	}
	// add raw data
	result.RawMessage = text

	return result, nil
}

func parseTestLine(line string, result *UrineTestResult) {
	// Regular expressions for different test line formats
	patterns := []*regexp.Regexp{
		// Format: "SG         1.015     "
		regexp.MustCompile(`^(\w{2,4})\s+([\d.*]+)`),
		// Format: "*LEU +3    500 CELL/uL" (with asterisk)
		regexp.MustCompile(`^\*?(\w{2,4})\s+([+\-]?\w+)\s+([\d.]+)\s+(\w+/?\w*)`),
		// Format: "KET -        0 mmol/L"
		regexp.MustCompile(`^(\w{2,4})\s+([+\-]?\w+)\s+([\d.]+)\s+(\w+/?\w*)`),
		// Format: "NIT -   " (no value/unit)
		regexp.MustCompile(`^(\w{2,4})\s+([+\-]?\w+)`),
		// Format: "URO            Normal"
		regexp.MustCompile(`^(\w{2,4})\s+([\w\s]+)$`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		param := strings.ToUpper(strings.TrimSpace(matches[1]))
		value := strings.TrimSpace(matches[2])

		// Extract unit and numeric value if present
		var unit string
		var numericVal float64
		var hasNumeric bool

		if len(matches) > 3 && matches[3] != "" {
			if num, err := strconv.ParseFloat(matches[3], 64); err == nil {
				numericVal = num
				hasNumeric = true
			}
		}

		if len(matches) > 4 {
			unit = strings.TrimSpace(matches[4])
		}

		testValue := TestValue{
			Value:       value,
			Unit:        unit,
			Numeric:     numericVal,
			IsNumeric:   hasNumeric,
			Qualitative: value, // Store the qualitative result
		}

		// Map parameter to struct field
		switch param {
		case "LEU":
			result.Leukocytes = testValue
		case "KET":
			result.Ketones = testValue
		case "NIT":
			result.Nitrites = testValue
		case "URO":
			result.Urobilinogen = testValue
		case "BIL":
			result.Bilirubin = testValue
		case "PRO":
			result.Protein = testValue
		case "GLU":
			result.Glucose = testValue
		case "SG":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				result.SpecificGravity = val
			}
		case "BLD":
			result.Blood = testValue
		case "PH":
			if val, err := strconv.ParseFloat(value, 64); err == nil {
				result.PH = val
			}
		case "VC":
			result.AscorbicAcid = testValue
		}
		break
	}
}
