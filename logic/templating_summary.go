package logic

import (
	"strings"
	"text/template"
	"time"
)

type SummaryMetadata struct {
	Title     string
	CreatedAt time.Time
}

func TemplatingSummary(summaryLines []string, summaryMetadata SummaryMetadata) (string, error) {
	data := map[string]interface{}{
		"title":      summaryMetadata.Title,
		"created_at": summaryMetadata.CreatedAt.Format("2006-01-02 15:04:05"),
		"content":    strings.Join(summaryLines, "\n"),
	}

	tmpl, err := template.New("").Parse(`---
title: "{{ .title }}"
created_at: "{{ .created_at }}"
---

{{ .content }}`)

	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
