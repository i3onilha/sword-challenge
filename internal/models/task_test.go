package models

import (
	"testing"
	"time"
)

func TestTask_Validate(t *testing.T) {
	type fields struct {
		ID           int64
		TechnicianID int64
		Title        string
		Summary      string
		PerformedAt  time.Time
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid task",
			fields: fields{
				Title:       "Fix air conditioning",
				Summary:     "Replaced filters and recharged coolant",
				PerformedAt: time.Date(2024, 3, 20, 14, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "empty title",
			fields: fields{
				Title:       "",
				Summary:     "Some summary",
				PerformedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty summary",
			fields: fields{
				Title:       "Some title",
				Summary:     "",
				PerformedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "title too long",
			fields: fields{
				Title:       string(make([]byte, 256)),
				Summary:     "Valid summary",
				PerformedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "summary too long",
			fields: fields{
				Title:       "Valid title",
				Summary:     string(make([]byte, 2501)),
				PerformedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "performed_at before min date",
			fields: fields{
				Title:       "Valid title",
				Summary:     "Valid summary",
				PerformedAt: time.Date(1899, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "performed_at after max date",
			fields: fields{
				Title:       "Valid title",
				Summary:     "Valid summary",
				PerformedAt: time.Date(2101, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
		{
			name: "title and summary with leading/trailing spaces",
			fields: fields{
				Title:       "   Valid title   ",
				Summary:     "   Valid summary   ",
				PerformedAt: time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Task{
				ID:           tt.fields.ID,
				TechnicianID: tt.fields.TechnicianID,
				Title:        tt.fields.Title,
				Summary:      tt.fields.Summary,
				PerformedAt:  tt.fields.PerformedAt,
				CreatedAt:    tt.fields.CreatedAt,
				UpdatedAt:    tt.fields.UpdatedAt,
			}
			if err := tr.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Task.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Sanitize(t *testing.T) {
	type fields struct {
		ID           int64
		TechnicianID int64
		Title        string
		Summary      string
		PerformedAt  time.Time
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
	tests := []struct {
		name        string
		fields      fields
		wantTitle   string
		wantSummary string
	}{
		{
			name: "sanitize html in title and summary",
			fields: fields{
				Title:   "<script>alert('x')</script>",
				Summary: "<b>bold</b> & <i>italic</i>",
			},
			wantTitle:   "&lt;script&gt;alert(&#39;x&#39;)&lt;/script&gt;",
			wantSummary: "&lt;b&gt;bold&lt;/b&gt; &amp; &lt;i&gt;italic&lt;/i&gt;",
		},
		{
			name: "sanitize with no html",
			fields: fields{
				Title:   "Normal Title",
				Summary: "Normal Summary",
			},
			wantTitle:   "Normal Title",
			wantSummary: "Normal Summary",
		},
		{
			name: "sanitize with special chars",
			fields: fields{
				Title:   "5 > 3 & 2 < 4",
				Summary: "\"quote\" & 'single quote'",
			},
			wantTitle:   "5 &gt; 3 &amp; 2 &lt; 4",
			wantSummary: "&quot;quote&quot; &amp; &#39;single quote&#39;",
		},
		{
			name: "sanitize empty fields",
			fields: fields{
				Title:   "",
				Summary: "",
			},
			wantTitle:   "",
			wantSummary: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Task{
				ID:           tt.fields.ID,
				TechnicianID: tt.fields.TechnicianID,
				Title:        tt.fields.Title,
				Summary:      tt.fields.Summary,
				PerformedAt:  tt.fields.PerformedAt,
				CreatedAt:    tt.fields.CreatedAt,
				UpdatedAt:    tt.fields.UpdatedAt,
			}
			tr.Sanitize()
		})
	}
}
