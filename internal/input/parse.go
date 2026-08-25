package input

import (
	"encoding/json"
	"fmt"
	"io"
)

type fileSpec struct {
	Name     string    `json:"name"`
	C        *float64  `json:"c"`
	WMax     *float64  `json:"w_max"`
	PrevWMax *float64  `json:"previous_w_max"`
	RTT      *float64  `json:"rtt_seconds"`
	T        *float64  `json:"t_seconds"`
	Acks     *int64    `json:"acks"`
	Horizon  *float64  `json:"horizon_seconds"`
	Samples  *int      `json:"samples"`
	Sim      *fileSim  `json:"sim"`
	Fair     *fileFair `json:"fair"`
}

type fileSim struct {
	Mode        *string  `json:"mode"`
	Rounds      *int     `json:"rounds"`
	Start       *string  `json:"start"`
	InitialCwnd *float64 `json:"initial_cwnd"`
	Ssthresh    *float64 `json:"ssthresh"`
}

type fileFair struct {
	Capacity *float64 `json:"capacity_segments"`
	Rounds   *int     `json:"rounds"`
	FlowA    *float64 `json:"flow_a_cwnd"`
	FlowB    *float64 `json:"flow_b_cwnd"`
}

func Parse(r io.Reader) (*Spec, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var fs fileSpec
	if err := dec.Decode(&fs); err != nil {
		return nil, fmt.Errorf("spec: json: %w", err)
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("spec: json: trailing data after first object")
		}
		return nil, fmt.Errorf("spec: json: %w", err)
	}
	spec := &Spec{Name: fs.Name}
	applyDefaults(spec, &fs)
	if err := validateRequired(spec, &fs); err != nil {
		return nil, err
	}
	if err := validateRanges(spec, &fs); err != nil {
		return nil, err
	}
	if err := validateTime(spec, &fs); err != nil {
		return nil, err
	}
	if err := validateSim(spec, &fs); err != nil {
		return nil, err
	}
	if err := validateFair(spec, &fs); err != nil {
		return nil, err
	}
	return spec, nil
}
