package main

import (
	"strings"
	"time"
)

// ScheduleFullTime and JSON unmarshaler
type ScheduleFullTime time.Time

func (s *ScheduleFullTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(fullTimeLayout, value)
	if err != nil {
		return err
	}

	*s = ScheduleFullTime(t)
	return nil
}

// ScheduleTime and JSON unmarshaler
type ScheduleTime time.Time

func (s *ScheduleTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(timeLayout, value)
	if err != nil {
		return err
	}

	*s = ScheduleTime(t)
	return nil
}

type ScheduleItem struct {
	ContentNote      string           `json:"content_note"`
	Description      string           `json:"description"`
	EndDate          ScheduleFullTime `json:"end_date"`
	EndTime          ScheduleTime     `json:"end_time"`
	ID               int              `json:"id"`
	IsFamilyFriendly bool             `json:"is_family_friendly"`
	IsFave           bool             `json:"is_fave"`
	Latlon           []float64        `json:"latlon"`
	Link             string           `json:"link"`
	MapLink          string           `json:"map_link"`
	MayRecord        bool             `json:"may_record"`
	Pronouns         string           `json:"pronouns"`
	Slug             string           `json:"slug"`
	Source           string           `json:"source"`
	Speaker          string           `json:"speaker"`
	StartDate        ScheduleFullTime `json:"start_date"`
	StartTime        ScheduleTime     `json:"start_time"`
	Title            string           `json:"title"`
	Type             string           `json:"type"`
	UserID           int              `json:"user_id"`
	Venue            string           `json:"venue"`
}
