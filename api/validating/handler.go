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
	"security-webhook/configs"
	"security-webhook/utils/log"
)

var deploymentKind = v1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
var statefulSetKind = v1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
var cronjobKind = v1.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"}
var jobKind = v1.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"}
var daemonSetKind = v1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"}

func SecurityValidate(c *gin.Context) {
	log.Logger.Debug("receive request")
	var reviewContext admission.AdmissionReview
	err := c.BindJSON(&reviewContext)
	if err != nil {
		log.Logger.Error("Failed to unmarshal Input as admission.AdmissionReview", zap.Error(err))
		return
	}

	log.Logger.Info("review data: " + fmt.Sprintf("%s/%s", reviewContext.Request.Namespace, reviewContext.Request.Name))

	switch reviewContext.Request.Kind {
	case deploymentKind:
		handleDeployment(&reviewContext)
	case statefulSetKind:
		handleStatefulSet(&reviewContext)
	case daemonSetKind:
		handleDaemonSet(&reviewContext)
	case cronjobKind:
		handleCronJob(&reviewContext)
	case jobKind:
		handleJob(&reviewContext)
	default:
		reviewContext.Response.Allowed = false
		c.JSON(http.StatusOK, reviewContext)
		return
	}
}

func handleDeployment(reviewContext *admission.AdmissionReview) {
	if configs.GlobalConfig.CheckItems.ForbiddenPrivilegedContainer {
		var deployment appv1.Deployment
		if err := json.Unmarshal(reviewContext.Request.Object.Raw, &deployment); err != nil {
			log.Logger.Error("failed to parse deployment object raw ", zap.String("data", string(reviewContext.Request.Object.Raw)))
			reviewContext.Response.Allowed = true
			return
		}

		if skipResource(deployment.ObjectMeta) {
			reviewContext.Response.Allowed = true
			return
		}

		for _, container := range deployment.Spec.Template.Spec.Containers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}

		for _, container := range deployment.Spec.Template.Spec.InitContainers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}
	}

	reviewContext.Response.Allowed = true
	return
}

func handleStatefulSet(reviewContext *admission.AdmissionReview) {
	if configs.GlobalConfig.CheckItems.ForbiddenPrivilegedContainer {
		var deployment appv1.StatefulSet
		if err := json.Unmarshal(reviewContext.Request.Object.Raw, &deployment); err != nil {
			log.Logger.Error("failed to parse statefulSet object raw ", zap.String("data", string(reviewContext.Request.Object.Raw)))
			reviewContext.Response.Allowed = true
			return
		}

		if skipResource(deployment.ObjectMeta) {
			reviewContext.Response.Allowed = true
			return
		}

		for _, container := range deployment.Spec.Template.Spec.Containers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}

		for _, container := range deployment.Spec.Template.Spec.InitContainers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}
	}

	reviewContext.Response.Allowed = true
	return
}

func handleDaemonSet(reviewContext *admission.AdmissionReview) {
	if configs.GlobalConfig.CheckItems.ForbiddenPrivilegedContainer {
		var deployment appv1.DaemonSet
		if err := json.Unmarshal(reviewContext.Request.Object.Raw, &deployment); err != nil {
			log.Logger.Error("failed to parse daemonSet object raw ", zap.String("data", string(reviewContext.Request.Object.Raw)))
			reviewContext.Response.Allowed = true
			return
		}

		if skipResource(deployment.ObjectMeta) {
			reviewContext.Response.Allowed = true
			return
		}

		for _, container := range deployment.Spec.Template.Spec.Containers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}

		for _, container := range deployment.Spec.Template.Spec.InitContainers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}
	}

	reviewContext.Response.Allowed = true
	return
}

func handleCronJob(reviewContext *admission.AdmissionReview) {
	if configs.GlobalConfig.CheckItems.ForbiddenPrivilegedContainer {
		var cronjob batchv1.CronJob
		if err := json.Unmarshal(reviewContext.Request.Object.Raw, &cronjob); err != nil {
			log.Logger.Error("failed to parse cronjob object raw ", zap.String("data", string(reviewContext.Request.Object.Raw)))
			reviewContext.Response.Allowed = true
			return
		}

		if skipResource(cronjob.ObjectMeta) {
			reviewContext.Response.Allowed = true
			return
		}

		for _, container := range cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}

		for _, container := range cronjob.Spec.JobTemplate.Spec.Template.Spec.InitContainers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}
	}

	reviewContext.Response.Allowed = true
	return
}

func handleJob(reviewContext *admission.AdmissionReview) {
	if configs.GlobalConfig.CheckItems.ForbiddenPrivilegedContainer {
		var job batchv1.Job
		if err := json.Unmarshal(reviewContext.Request.Object.Raw, &job); err != nil {
			log.Logger.Error("failed to parse job object raw ", zap.String("data", string(reviewContext.Request.Object.Raw)))
			reviewContext.Response.Allowed = true
			return
		}

		if skipResource(job.ObjectMeta) {
			reviewContext.Response.Allowed = true
			return
		}

		for _, container := range job.Spec.Template.Spec.Containers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}

		for _, container := range job.Spec.Template.Spec.InitContainers {
			if err := validateContainer(container); err != nil {
				reviewContext.Response.Allowed = false
				reviewContext.Response.Result.Message = err.Error()
				return
			}
		}
	}

	reviewContext.Response.Allowed = true
	return
}
