package validating

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	admission "k8s.io/api/admission/v1"
	appv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"security-webhook/utils/log"
)

var deploymentKind = v1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
var statefulSetKind = v1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
var cronjobKind = v1.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"}
var jobKind = v1.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}
var daemonSetKind = v1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}

func PrivilegedContainerCheck(c *gin.Context) {
	var requestReview admission.AdmissionReview
	err := c.BindJSON(&requestReview)
	if err != nil {
		log.Logger.Error("Failed to unmarshal Input as admission.AdmissionReview", zap.Error(err))
		return
	}

	responseReview := &admission.AdmissionReview{}
	responseReview.TypeMeta = requestReview.TypeMeta
	responseReview.Request = requestReview.Request

	switch requestReview.Request.Kind {
	case deploymentKind:
		handleDeployment(c, &requestReview, responseReview)
	case statefulSetKind:
		handleStatefulSet(c, &requestReview, responseReview)
	case cronjobKind:
		handleCronJob(c, &requestReview, responseReview)
	case daemonSetKind:
		handleDaemonSet(c, &requestReview, responseReview)
	case jobKind:
		handleJob(c, &requestReview, responseReview)
	default:
		log.Logger.Error("Unsupported group version kind", zap.String("kind", requestReview.Request.Kind.String()))
		responseReview.Response.Allowed = false
		responseReview.Response.Result.Message = fmt.Sprintf("unknown application kind: %s", requestReview.Request.Kind.String())
		c.JSON(http.StatusOK, responseReview)
		return
	}
}

func handleDeployment(ctx *gin.Context, requestReview, responseReview *admission.AdmissionReview) {
	var deployment appv1.Deployment
	if err := json.Unmarshal(requestReview.Request.Object.Raw, &deployment); err != nil {
		log.Logger.Error("Error unmarshaling AdmissionRequest raw into Deployment: %v", zap.Error(err))
		return
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		if container.SecurityContext != nil {
			if container.SecurityContext.Privileged != nil {
				if *container.SecurityContext.Privileged == true {
					responseReview.Response.Allowed = false
					responseReview.Response.Result.Message = fmt.Sprintf("Deployment: \"%s\" container: \"%s\" set as privileged container. Reject it", deployment.Name, container.Name)
					ctx.JSON(http.StatusOK, responseReview)
					return
				}
			}
		}
	}
	// 检查通过
	responseReview.Response.Allowed = true
	ctx.JSON(http.StatusOK, responseReview)
	return
}

func handleStatefulSet(ctx *gin.Context, requestReview, responseReview *admission.AdmissionReview) {
	var statefulSet appv1.StatefulSet
	if err := json.Unmarshal(requestReview.Request.Object.Raw, &statefulSet); err != nil {
		log.Logger.Error("Error unmarshaling AdmissionRequest raw into StatefulSet: %v", zap.Error(err))
		return
	}

	for _, container := range statefulSet.Spec.Template.Spec.Containers {
		if container.SecurityContext != nil {
			if container.SecurityContext.Privileged != nil {
				if *container.SecurityContext.Privileged == true {
					responseReview.Response.Allowed = false
					responseReview.Response.Result.Message = fmt.Sprintf("StatefulSet: \"%s\" container: \"%s\" set as privileged container. Reject it", statefulSet.Name, container.Name)
					ctx.JSON(http.StatusOK, responseReview)
					return
				}
			}
		}
	}
	// 检查通过
	responseReview.Response.Allowed = true
	ctx.JSON(http.StatusOK, responseReview)
	return
}

func handleCronJob(ctx *gin.Context, requestReview, responseReview *admission.AdmissionReview) {
	var cronjob batchv1.CronJob
	if err := json.Unmarshal(requestReview.Request.Object.Raw, &cronjob); err != nil {
		log.Logger.Error("Error unmarshaling AdmissionRequest raw into Cronjob: %v", zap.Error(err))
		return
	}

	for _, container := range cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers {
		if container.SecurityContext != nil {
			if container.SecurityContext.Privileged != nil {
				if *container.SecurityContext.Privileged == true {
					responseReview.Response.Allowed = false
					responseReview.Response.Result.Message = fmt.Sprintf("Cronjob: \"%s\" container: \"%s\" set as privileged container. Reject it", cronjob.Name, container.Name)
					ctx.JSON(http.StatusOK, responseReview)
					return
				}
			}
		}
	}
	// 检查通过
	responseReview.Response.Allowed = true
	ctx.JSON(http.StatusOK, responseReview)
	return
}

func handleJob(ctx *gin.Context, requestReview, responseReview *admission.AdmissionReview) {
	var job batchv1.Job
	if err := json.Unmarshal(requestReview.Request.Object.Raw, &job); err != nil {
		log.Logger.Error("Error unmarshaling AdmissionRequest raw into Job: %v", zap.Error(err))
		return
	}

	for _, container := range job.Spec.Template.Spec.Containers {
		if container.SecurityContext != nil {
			if container.SecurityContext.Privileged != nil {
				if *container.SecurityContext.Privileged == true {
					responseReview.Response.Allowed = false
					responseReview.Response.Result.Message = fmt.Sprintf("Job: \"%s\" container: \"%s\" set as privileged container. Reject it", job.Name, container.Name)
					ctx.JSON(http.StatusOK, responseReview)
					return
				}
			}
		}
	}
	// 检查通过
	responseReview.Response.Allowed = true
	ctx.JSON(http.StatusOK, responseReview)
	return
}

func handleDaemonSet(ctx *gin.Context, requestReview, responseReview *admission.AdmissionReview) {
	var daemonSet appv1.DaemonSet
	if err := json.Unmarshal(requestReview.Request.Object.Raw, &daemonSet); err != nil {
		log.Logger.Error("Error unmarshaling AdmissionRequest raw into DaemonSet: %v", zap.Error(err))
		return
	}

	for _, container := range daemonSet.Spec.Template.Spec.Containers {
		if container.SecurityContext != nil {
			if container.SecurityContext.Privileged != nil {
				if *container.SecurityContext.Privileged == true {
					responseReview.Response.Allowed = false
					responseReview.Response.Result.Message = fmt.Sprintf("DaemonSet: \"%s\" container: \"%s\" set as privileged container. Reject it", daemonSet.Name, container.Name)
					ctx.JSON(http.StatusOK, responseReview)
					return
				}
			}
		}
	}
	// 检查通过
	responseReview.Response.Allowed = true
	ctx.JSON(http.StatusOK, responseReview)
	return
}
