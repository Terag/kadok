{{/*
Expand the name of the chart.
*/}}
{{- define "kadok.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kadok.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kadok.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kadok.labels" -}}
helm.sh/chart: {{ include "kadok.chart" . }}
{{ include "kadok.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kadok.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kadok.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Pod Annotations
*/}}
{{- define "kadok.annotations" -}}
{{- if .Values.gitlab }}
{{- if .Values.gitlab.env }}
app.gitlab.com/env: {{ .Values.gitlab.env }}
{{- end }}
{{- if .Values.gitlab.app }}
app.gitlab.com/app: {{ .Values.gitlab.app }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create a imagePullSecret value
*/}}
{{- define "kadok.imagePullSecret" }}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.registrySecret.registry (printf "%s:%s" .Values.registrySecret.username .Values.registrySecret.password | b64enc) | b64enc }}
{{- end }}

{{/*
Get ConfigMap hash to ensure the immutability
*/}}
{{- define "kadok.configMapShortHash" }}
{{- if (.Files.Glob (printf "%s/*" .Values.kadok.configs)) }}
{{- (.Files.Glob (printf "%s/*" .Values.kadok.configs)).AsConfig | sha256sum | substr 0 12 }}
{{- end }}
{{- end }}
