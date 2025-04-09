{{- define "app.resources" -}}
resources:
{{- if .Values.resources.limits.memory }}
  limits:
    memory: "{{ .Values.resources.limits.memory }}"
{{- end }}
  requests:
    memory: "{{ .Values.resources.requests.memory }}"
{{- end }}

{{- define "app.service" -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ .Values.service.name }}
  namespace: apps
spec:
  type: ClusterIP
  selector:
    {{- range $k, $v := .Values.labels }}
    {{ $k }}: {{ $v }}
    {{- end }}
  ports:
  {{- range .Values.service.port }}
  - port: {{ . }}
    targetPort: {{ . }}
  {{- end }}
{{- end }}
