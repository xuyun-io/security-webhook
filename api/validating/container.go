package validating

import (
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	PrivilegedSecurityContextForbidden = errors.New("privileged security context forbidden")
	SecurityWebHookByPassAnnotation    = "security-webhook-check.bypass"
)

func skipResource(o metav1.ObjectMeta) bool {
	ann := o.GetAnnotations()
	if len(ann) > 0 && ann[SecurityWebHookByPassAnnotation] == "true" {
		return true
	}
	return false
}

func validateContainer(container corev1.Container) error {
	if container.SecurityContext == nil {
		return nil
	}
	if container.SecurityContext.Privileged == nil {
		return nil
	}
	if *container.SecurityContext.Privileged {
		return PrivilegedSecurityContextForbidden
	}
	return nil
}
