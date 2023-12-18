# security-webhook
This is a k8s admission controller webhook for security working with cert-manager

### helm install
```shell
    helm install security-webhook oci://registry-1.docker.io/kubestar/security-webhook -n namespace_to_install
```

### resources annotations setting 
Use below annotation to scape the resource validation
```shell
    annotations:
        security-webhook-check.bypass: "true"
```

### configmap check options 
mount path: /app/configs/default.yml

```shell
checkItems:
  allowPrivilegedContainer: false

```

### CRD ValidatingWebhookConfiguration
```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: security-webhook-validate
  annotations:
    cert-manager.io/inject-ca-from: {{ .Release.Namespace }}/security-webhook-certificate
webhooks:
  - admissionReviewVersions:
      - v1
    clientConfig:
      caBundle: ""
      service:
        name: security-webhook
        namespace: {{ .Release.Namespace }}
        port: 443
        path: /security-validate
    failurePolicy: Fail
    name: validate-privileged-container.{{ .Release.Namespace }}.svc
    namespaceSelector:
      matchExpressions:
        - key: "kubernetes.io/metadata.name"
          operator: NotIn
          values: ["kube-system"]
    rules:
      - apiGroups:
          - "apps"
        apiVersions:
          - v1
        operations:
          - CREATE
          - UPDATE
        resources:
          - deployments
          - statefulsets
          - daemonsets
        scope: '*'
    sideEffects: None
    timeoutSeconds: 3
```