// Code generated by protoc-gen-go-form. DO NOT EDIT.
// Source: label.proto

package pb

import (
	url "net/url"
	strconv "strconv"

	urlenc "github.com/erda-project/erda-infra/pkg/urlenc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the "github.com/erda-project/erda-infra/pkg/urlenc" package it is being compiled against.
var _ urlenc.URLValuesUnmarshaler = (*PipelineLabel)(nil)
var _ urlenc.URLValuesUnmarshaler = (*PipelineLabelBatchInsertRequest)(nil)
var _ urlenc.URLValuesUnmarshaler = (*PipelineLabelBatchInsertResponse)(nil)
var _ urlenc.URLValuesUnmarshaler = (*PipelineLabelListRequest)(nil)
var _ urlenc.URLValuesUnmarshaler = (*PipelineLabelListResponse)(nil)

// PipelineLabel implement urlenc.URLValuesUnmarshaler.
func (m *PipelineLabel) UnmarshalURLValues(prefix string, values url.Values) error {
	for key, vals := range values {
		if len(vals) > 0 {
			switch prefix + key {
			case "ID":
				val, err := strconv.ParseUint(vals[0], 10, 64)
				if err != nil {
					return err
				}
				m.ID = val
			case "type":
				m.Type = vals[0]
			case "targetID":
				val, err := strconv.ParseUint(vals[0], 10, 64)
				if err != nil {
					return err
				}
				m.TargetID = val
			case "pipelineSource":
				m.PipelineSource = vals[0]
			case "pipelineYmlName":
				m.PipelineYmlName = vals[0]
			case "key":
				m.Key = vals[0]
			case "value":
				m.Value = vals[0]
			case "timeCreated":
				if m.TimeCreated == nil {
					m.TimeCreated = &timestamppb.Timestamp{}
				}
			case "timeCreated.seconds":
				if m.TimeCreated == nil {
					m.TimeCreated = &timestamppb.Timestamp{}
				}
				val, err := strconv.ParseInt(vals[0], 10, 64)
				if err != nil {
					return err
				}
				m.TimeCreated.Seconds = val
			case "timeCreated.nanos":
				if m.TimeCreated == nil {
					m.TimeCreated = &timestamppb.Timestamp{}
				}
				val, err := strconv.ParseInt(vals[0], 10, 32)
				if err != nil {
					return err
				}
				m.TimeCreated.Nanos = int32(val)
			case "timeUpdated":
				if m.TimeUpdated == nil {
					m.TimeUpdated = &timestamppb.Timestamp{}
				}
			case "timeUpdated.seconds":
				if m.TimeUpdated == nil {
					m.TimeUpdated = &timestamppb.Timestamp{}
				}
				val, err := strconv.ParseInt(vals[0], 10, 64)
				if err != nil {
					return err
				}
				m.TimeUpdated.Seconds = val
			case "timeUpdated.nanos":
				if m.TimeUpdated == nil {
					m.TimeUpdated = &timestamppb.Timestamp{}
				}
				val, err := strconv.ParseInt(vals[0], 10, 32)
				if err != nil {
					return err
				}
				m.TimeUpdated.Nanos = int32(val)
			}
		}
	}
	return nil
}

// PipelineLabelBatchInsertRequest implement urlenc.URLValuesUnmarshaler.
func (m *PipelineLabelBatchInsertRequest) UnmarshalURLValues(prefix string, values url.Values) error {
	return nil
}

// PipelineLabelBatchInsertResponse implement urlenc.URLValuesUnmarshaler.
func (m *PipelineLabelBatchInsertResponse) UnmarshalURLValues(prefix string, values url.Values) error {
	return nil
}

// PipelineLabelListRequest implement urlenc.URLValuesUnmarshaler.
func (m *PipelineLabelListRequest) UnmarshalURLValues(prefix string, values url.Values) error {
	for key, vals := range values {
		if len(vals) > 0 {
			switch prefix + key {
			case "pipelineSource":
				m.PipelineSource = vals[0]
			case "pipelineYmlName":
				m.PipelineYmlName = vals[0]
			case "targetIDs":
				list := make([]uint64, 0, len(vals))
				for _, text := range vals {
					val, err := strconv.ParseUint(text, 10, 64)
					if err != nil {
						return err
					}
					list = append(list, val)
				}
				m.TargetIDs = list
			case "matchKeys":
				m.MatchKeys = vals
			case "unMatchKeys":
				m.UnMatchKeys = vals
			}
		}
	}
	return nil
}

// PipelineLabelListResponse implement urlenc.URLValuesUnmarshaler.
func (m *PipelineLabelListResponse) UnmarshalURLValues(prefix string, values url.Values) error {
	for key, vals := range values {
		if len(vals) > 0 {
			switch prefix + key {
			case "total":
				val, err := strconv.ParseInt(vals[0], 10, 64)
				if err != nil {
					return err
				}
				m.Total = val
			}
		}
	}
	return nil
}