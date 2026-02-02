package parsers

import (
	"fmt"
	"strings"
	"time"
)

// Result types returned by the parser
type Result struct {
	Patient              Patient `json:"patient"`
	Tests                []Test  `json:"tests"`
	PrimaryOrderDatetime string  `json:"primary_order_datetime,omitempty"`
}

type Patient struct {
	Identifiers []Identifier `json:"identifiers"`
	Name        Name         `json:"name"`
	DOB         string       `json:"dob,omitempty"`
	Sex         string       `json:"sex,omitempty"`
	Address     Address      `json:"address"`
}

type Identifier struct {
	ID             string `json:"id"`
	IdentifierType string `json:"identifier_type,omitempty"`
}

type Name struct {
	Family string `json:"family_name,omitempty"`
	Given  string `json:"given_name,omitempty"`
	Middle string `json:"middle_name,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

type Address struct {
	Street  string `json:"street,omitempty"`
	Other   string `json:"other,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}

type Test struct {
	Code                string `json:"code,omitempty"`
	Name                string `json:"name,omitempty"`
	CodingSystem        string `json:"coding_system,omitempty"`
	Value               string `json:"value,omitempty"`
	Units               string `json:"units,omitempty"`
	ReferenceRange      string `json:"reference_range,omitempty"`
	AbnormalFlags       string `json:"abnormal_flags,omitempty"`
	ObservationDatetime string `json:"observation_datetime,omitempty"`
}

// parseHL7Timestamp converts HL7 TS (e.g., 20250101120000 or 20250101 or 202501011200+0500)
// into ISO 8601 string (RFC3339 for timestamps or YYYY-MM-DD for date-only).
// Returns empty string if input empty or unparseable.
func parseHL7Timestamp(ts string) string {
	if ts == "" {
		return ""
	}

	// strip fractional seconds if present
	if dot := strings.Index(ts, "."); dot != -1 {
		ts = ts[:dot]
	}

	// find timezone (+/-) if present (only if not the leading character)
	tzPos := strings.IndexAny(ts, "+-")
	var tzPart string
	var tsMain string
	if tzPos > 0 {
		tsMain = ts[:tzPos]
		tzPart = ts[tzPos:]
	} else {
		tsMain = ts
		tzPart = ""
	}

	var layout string
	switch len(tsMain) {
	case 14:
		layout = "20060102150405"
	case 12:
		layout = "200601021504"
	case 10:
		layout = "2006010215"
	case 8:
		layout = "20060102"
	default:
		// unsupported granularities - return empty
		return ""
	}

	// If timezone present, append timezone layout
	if tzPart != "" {
		// HL7 uses +/-HHMM; Go layout for that is -0700
		layout = layout + "-0700"
		parsed, err := time.Parse(layout, tsMain+tzPart)
		if err != nil {
			return ""
		}
		// If only date (8 chars) produce date-only string
		if len(tsMain) == 8 {
			return parsed.Format("2006-01-02")
		}
		return parsed.Format(time.RFC3339)
	}

	// No timezone
	parsed, err := time.Parse(layout, tsMain)
	if err != nil {
		return ""
	}
	if len(tsMain) == 8 {
		return parsed.Format("2006-01-02")
	}
	// Make it RFC3339 (naive, local set to UTC for consistent output)
	return parsed.Format(time.RFC3339)
}

// parseHL7Message parses a single HL7 message (string) and extracts patient, tests, and primary order datetime.
// It reads delimiters from MSH and supports component and repetition separators.
func ParseHL7Message(msg string) (Result, error) {
	var res Result

	// Normalize line endings and split segments
	lines := []string{}
	for _, ln := range strings.Split(msg, "\n") {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "MSH") {
		return res, fmt.Errorf("message does not start with MSH")
	}

	msh := lines[0]
	// field separator is the 4th byte in MSH (index 3)
	if len(msh) < 4 {
		return res, fmt.Errorf("invalid MSH segment")
	}
	fieldSep := string(msh[3])

	fields := strings.Split(msh, fieldSep)
	encChars := "^~\\&"
	if len(fields) > 1 && fields[1] != "" {
		encChars = fields[1]
	}
	// encoding chars: component, repetition, escape, subcomponent
	compSep := "^"
	repSep := "~"
	subCompSep := "&"
	if len(encChars) > 0 {
		compSep = string(encChars[0])
	}
	if len(encChars) > 1 {
		repSep = string(encChars[1])
	}
	if len(encChars) > 3 {
		subCompSep = string(encChars[3])
	}

	// helper to split by fieldSep
	splitFields := func(segment string) []string {
		return strings.Split(segment, fieldSep)
	}
	// helper to split components
	splitComp := func(field string) []string {
		if field == "" {
			return []string{}
		}
		return strings.Split(field, compSep)
	}

	var primaryOrderDatetime string
	var patient Patient
	var tests []Test

	for _, seg := range lines {
		if len(seg) < 3 {
			continue
		}
		tag := seg[:3]
		switch tag {
		case "PID":
			pidFields := splitFields(seg)
			safe := func(i int) string {
				if i < len(pidFields) {
					return pidFields[i]
				}
				return ""
			}
			// PID-3 Patient Identifier List
			pid3 := safe(3)
			identifiers := []Identifier{}
			if pid3 != "" {
				for _, rep := range strings.Split(pid3, repSep) {
					comps := splitComp(rep)
					id := ""
					idType := ""
					if len(comps) > 0 {
						id = comps[0]
					}
					if len(comps) > 4 {
						idType = comps[4]
					}
					if id != "" {
						identifiers = append(identifiers, Identifier{ID: id, IdentifierType: idType})
					}
				}
			}
			// PID-5 Name
			nameRaw := safe(5)
			nc := splitComp(nameRaw)
			name := Name{
				Family: "",
				Given:  "",
				Middle: "",
				Suffix: "",
			}
			if len(nc) > 0 {
				name.Family = nc[0]
			}
			if len(nc) > 1 {
				name.Given = nc[1]
			}
			if len(nc) > 2 {
				name.Middle = nc[2]
			}
			if len(nc) > 3 {
				name.Suffix = nc[3]
			}
			// PID-7 DOB
			dobRaw := safe(7)
			dob := parseHL7Timestamp(dobRaw)
			// PID-8 Sex
			sex := safe(8)
			// PID-11 Address
			addrRaw := safe(11)
			ac := splitComp(addrRaw)
			address := Address{}
			if len(ac) > 0 {
				address.Street = ac[0]
			}
			if len(ac) > 1 {
				address.Other = ac[1]
			}
			if len(ac) > 2 {
				address.City = ac[2]
			}
			if len(ac) > 3 {
				address.State = ac[3]
			}
			if len(ac) > 4 {
				address.Zip = ac[4]
			}
			if len(ac) > 5 {
				address.Country = ac[5]
			}

			patient = Patient{
				Identifiers: identifiers,
				Name:        name,
				DOB:         dob,
				Sex:         sex,
				Address:     address,
			}

		case "OBR":
			obrFields := splitFields(seg)
			safe := func(i int) string {
				if i < len(obrFields) {
					return obrFields[i]
				}
				return ""
			}
			// OBR-7 Observation Date/Time (index 7)
			if dt := safe(7); dt != "" {
				if parsed := parseHL7Timestamp(dt); parsed != "" {
					primaryOrderDatetime = parsed
				}
			}
			// fallback OBR-14 (specimen collected)
			if primaryOrderDatetime == "" {
				if dt := safe(14); dt != "" {
					if parsed := parseHL7Timestamp(dt); parsed != "" {
						primaryOrderDatetime = parsed
					}
				}
			}

		case "OBX":
			obxFields := splitFields(seg)
			safe := func(i int) string {
				if i < len(obxFields) {
					return obxFields[i]
				}
				return ""
			}
			obsIDRaw := safe(3)
			obsIdComps := splitComp(obsIDRaw)
			code := ""
			name := ""
			codingSystem := ""
			if len(obsIdComps) > 0 {
				code = obsIdComps[0]
			}
			if len(obsIdComps) > 1 {
				name = obsIdComps[1]
			}
			if len(obsIdComps) > 2 {
				codingSystem = obsIdComps[2]
			}
			value := safe(5)
			unitsRaw := safe(6)
			unitsComps := splitComp(unitsRaw)
			units := unitsRaw
			if len(unitsComps) > 0 && unitsComps[0] != "" {
				units = unitsComps[0]
			}
			refRange := safe(7)
			flags := safe(8)
			obsDtRaw := safe(14)
			obsDt := ""
			if obsDtRaw != "" {
				obsDt = parseHL7Timestamp(obsDtRaw)
			}

			tests = append(tests, Test{
				Code:                code,
				Name:                name,
				CodingSystem:        codingSystem,
				Value:               value,
				Units:               units,
				ReferenceRange:      refRange,
				AbnormalFlags:       flags,
				ObservationDatetime: obsDt,
			})
		default:
			// ignore other segments
			_ = subCompSep // silence unused
		}
	}

	res = Result{
		Patient:              patient,
		Tests:                tests,
		PrimaryOrderDatetime: primaryOrderDatetime,
	}
	return res, nil
}
