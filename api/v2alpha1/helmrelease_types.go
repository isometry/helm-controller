/*
Copyright 2020 The Flux CD contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2alpha1

import (
	"encoding/json"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const HelmReleaseKind = "HelmRelease"
const HelmReleaseFinalizer = "finalizers.fluxcd.io"

// HelmReleaseSpec defines the desired state of HelmRelease.
type HelmReleaseSpec struct {
	// Chart defines the Helm chart name, version and repository.
	// +required
	Chart HelmChartTemplate `json:"chart"`

	// Interval at which to reconcile the Helm release.
	// +required
	Interval metav1.Duration `json:"interval"`

	// Suspend tells the reconciler to suspend reconciliation for this HelmRelease,
	// it does not apply to already started reconciliations. Defaults to false.
	// +optional
	Suspend bool `json:"suspend,omitempty"`

	// ReleaseName used for the Helm release. Defaults to a composition of
	// '[TargetNamespace-]Name'.
	// +optional
	ReleaseName string `json:"releaseName,omitempty"`

	// TargetNamespace to target when performing operations for the HelmRelease.
	// Defaults to the namespace of the HelmRelease.
	// +optional
	TargetNamespace string `json:"targetNamespace,omitempty"`

	// DependsOn may contain a list of HelmReleases that must be ready before this
	// HelmRelease can be reconciled.
	// +optional
	DependsOn []string `json:"dependsOn,omitempty"`

	// Timeout is the time to wait for any individual Kubernetes operation (like Jobs
	// for hooks) during the performance of a Helm action. Defaults to '5m0s'.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// MaxHistory is the number of revisions saved by Helm for this HelmRelease.
	// Use '0' for an unlimited number of revisions; defaults to '10'.
	// +optional
	MaxHistory *int `json:"maxHistory,omitempty"`

	// Install holds the configuration for Helm install actions for this HelmRelease.
	// +optional
	Install *Install `json:"install,omitempty"`

	// Upgrade holds the configuration for Helm upgrade actions for this HelmRelease.
	// +optional
	Upgrade *Upgrade `json:"upgrade,omitempty"`

	// Test holds the configuration for Helm test actions for this HelmRelease.
	// +optional
	Test *Test `json:"test,omitempty"`

	// Rollback holds the configuration for Helm rollback actions for this HelmRelease.
	// +optional
	Rollback *Rollback `json:"rollback,omitempty"`

	// Uninstall holds the configuration for Helm uninstall actions for this HelmRelease.
	// +optional
	Uninstall *Uninstall `json:"uninstall,omitempty"`

	// ValuesFrom holds references to resources containing Helm values for this HelmRelease,
	// and information about how they should be merged.
	ValuesFrom []ValuesReference `json:"valuesFrom,omitempty"`

	// Values holds the values for this Helm release.
	// +optional
	Values *apiextensionsv1.JSON `json:"values,omitempty"`
}

// GetInstall returns the configuration for Helm install actions for the HelmRelease.
func (in HelmReleaseSpec) GetInstall() Install {
	switch in.Install {
	case nil:
		return Install{}
	default:
		return *in.Install
	}
}

// GetUpgrade returns the configuration for Helm upgrade actions for this HelmRelease.
func (in HelmReleaseSpec) GetUpgrade() Upgrade {
	switch in.Upgrade {
	case nil:
		return Upgrade{}
	default:
		return *in.Upgrade
	}
}

// GetTest returns the configuration for Helm test actions for this HelmRelease.
func (in HelmReleaseSpec) GetTest() Test {
	switch in.Test {
	case nil:
		return Test{}
	default:
		return *in.Test
	}
}

// GetRollback returns the configuration for Helm rollback actions for this HelmRelease.
func (in HelmReleaseSpec) GetRollback() Rollback {
	switch in.Rollback {
	case nil:
		return Rollback{}
	default:
		return *in.Rollback
	}
}

// GetUninstall returns the configuration for Helm uninstall actions for this HelmRelease.
func (in HelmReleaseSpec) GetUninstall() Uninstall {
	switch in.Uninstall {
	case nil:
		return Uninstall{}
	default:
		return *in.Uninstall
	}
}

// HelmChartTemplate defines the template from which the controller
// will generate a HelmChart object in the same namespace as the HelmRepository.
type HelmChartTemplate struct {
	// Name of the Helm chart, as made available by the referenced Helm repository.
	// +required
	Name string `json:"name"`

	// Version semver expression, defaults to latest when omitted.
	// +optional
	Version string `json:"version,omitempty"`

	// The name and namespace of the source HelmRepository the chart is available at.
	// +required
	SourceRef CrossNamespaceObjectReference `json:"sourceRef"`

	// Interval at which to check the Helm repository for chart updates.
	// Defaults to 'HelmReleaseSpec.Interval'.
	// +optional
	Interval *metav1.Duration `json:"interval,omitempty"`
}

// GetInterval returns the configured interval for the HelmChart, or the given default.
func (in HelmChartTemplate) GetInterval(defaultInterval metav1.Duration) metav1.Duration {
	switch in.Interval {
	case nil:
		return defaultInterval
	default:
		return *in.Interval
	}
}

// GetNamespace returns the namespace targeted namespace for the HelmChart, or the given default.
func (in HelmChartTemplate) GetNamespace(defaultNamespace string) string {
	switch in.SourceRef.Namespace {
	case "":
		return defaultNamespace
	default:
		return in.SourceRef.Namespace
	}
}

// Install holds the configuration for Helm install actions performed for this HelmRelease.
type Install struct {
	// Timeout is the time to wait for any individual Kubernetes operation (like Jobs
	// for hooks) during the performance of a Helm install action. Defaults to
	// 'HelmReleaseSpec.Timeout'.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// DisableWait disables the waiting for resources to be ready after a
	// Helm install has been performed.
	// +optional
	DisableWait bool `json:"disableWait,omitempty"`

	// DisableHooks prevents hooks from running during the Helm install action.
	// +optional
	DisableHooks bool `json:"disableHooks,omitempty"`

	// DisableOpenAPIValidation prevents the Helm install action from
	// validating rendered templates against the Kubernetes OpenAPI Schema.
	// +optional
	DisableOpenAPIValidation bool `json:"disableOpenAPIValidation,omitempty"`

	// Replace tells the Helm install action to re-use the 'ReleaseName', but
	// only if that name is a deleted release which remains in the history.
	// +optional
	Replace bool `json:"replace,omitempty"`

	// SkipCRDs tells the Helm install action to not install any CRDs. By default,
	// CRDs are installed if not already present.
	// +optional
	SkipCRDs bool `json:"skipCRDs,omitempty"`
}

// GetTimeout returns the configured timeout for the Helm install action,
// or the given default.
func (in Install) GetTimeout(defaultTimeout metav1.Duration) metav1.Duration {
	switch in.Timeout {
	case nil:
		return defaultTimeout
	default:
		return *in.Timeout
	}
}

// Upgrade holds the configuration for Helm upgrade actions for this HelmRelease.
type Upgrade struct {
	// Timeout is the time to wait for any individual Kubernetes operation (like Jobs
	// for hooks) during the performance of a Helm upgrade action. Defaults to
	// 'HelmReleaseSpec.Timeout'.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// MaxRetries is the number of retries that should be attempted on failures before
	// bailing. Defaults to '0', a negative integer equals to unlimited retries.
	// +optional
	MaxRetries int `json:"maxRetries,omitempty"`

	// DisableWait disables the waiting for resources to be ready after a
	// Helm upgrade has been performed.
	// +optional
	DisableWait bool `json:"disableWait,omitempty"`

	// DisableHooks prevents hooks from running during the Helm upgrade action.
	// +optional
	DisableHooks bool `json:"disableHooks,omitempty"`

	// DisableOpenAPIValidation prevents the Helm upgrade action from
	// validating rendered templates against the Kubernetes OpenAPI Schema.
	// +optional
	DisableOpenAPIValidation bool `json:"disableOpenAPIValidation,omitempty"`

	// Force forces resource updates through a replacement strategy.
	// +optional
	Force bool `json:"force,omitempty"`

	// PreserveValues will make Helm reuse the last release's values and merge
	// in overrides from 'Values'. Setting this flag makes the HelmRelease
	// non-declarative.
	// +optional
	PreserveValues bool `json:"preserveValues,omitempty"`

	// CleanupOnFail allows deletion of new resources created during the Helm
	// upgrade action when it fails.
	// +optional
	CleanupOnFail bool `json:"cleanupOnFail,omitempty"`
}

// GetTimeout returns the configured timeout for the Helm upgrade action,
// or the given default.
func (in Upgrade) GetTimeout(defaultTimeout metav1.Duration) metav1.Duration {
	switch in.Timeout {
	case nil:
		return defaultTimeout
	default:
		return *in.Timeout
	}
}

// Test holds the configuration for Helm test actions for this HelmRelease.
type Test struct {
	// Enable enables Helm test actions for this HelmRelease after an
	// Helm install or upgrade action has been performed.
	// +optional
	Enable bool `json:"enable,omitempty"`

	// Timeout is the time to wait for any individual Kubernetes operation
	// during the performance of a Helm test action. Defaults to
	// 'HelmReleaseSpec.Timeout'.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

// GetTimeout returns the configured timeout for the Helm test action,
// or the given default.
func (in Test) GetTimeout(defaultTimeout metav1.Duration) metav1.Duration {
	switch in.Timeout {
	case nil:
		return defaultTimeout
	default:
		return *in.Timeout
	}
}

// Rollback holds the configuration for Helm rollback actions for this HelmRelease.
type Rollback struct {
	// Enable enables Helm rollback actions for this HelmRelease after an
	// Helm install or upgrade action failure.
	// +optional
	Enable bool `json:"enable,omitempty"`

	// Timeout is the time to wait for any individual Kubernetes operation (like Jobs
	// for hooks) during the performance of a Helm rollback action. Defaults to
	// 'HelmReleaseSpec.Timeout'.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// DisableWait disables the waiting for resources to be ready after a
	// Helm rollback has been performed.
	// +optional
	DisableWait bool `json:"disableWait,omitempty"`

	// DisableHooks prevents hooks from running during the Helm rollback action.
	// +optional
	DisableHooks bool `json:"disableHooks,omitempty"`

	// Recreate performs pod restarts for the resource if applicable.
	// +optional
	Recreate bool `json:"recreate,omitempty"`

	// Force forces resource updates through a replacement strategy.
	// +optional
	Force bool `json:"force,omitempty"`

	// CleanupOnFail allows deletion of new resources created during the Helm
	// rollback action when it fails.
	// +optional
	CleanupOnFail bool `json:"cleanupOnFail,omitempty"`
}

// GetTimeout returns the configured timeout for the Helm rollback action,
// or the given default.
func (in Rollback) GetTimeout(defaultTimeout metav1.Duration) metav1.Duration {
	switch in.Timeout {
	case nil:
		return defaultTimeout
	default:
		return *in.Timeout
	}
}

// Uninstall holds the configuration for Helm uninstall actions for this HelmRelease.
type Uninstall struct {
	// Timeout is the time to wait for any individual Kubernetes operation (like Jobs
	// for hooks) during the performance of a Helm uninstall action. Defaults to
	// 'HelmReleaseSpec.Timeout'.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// DisableHooks prevents hooks from running during the Helm rollback action.
	// +optional
	DisableHooks bool `json:"disableHooks,omitempty"`

	// KeepHistory tells Helm to remove all associated resources and mark the release as
	// deleted, but retain the release history.
	// +optional
	KeepHistory bool `json:"keepHistory,omitempty"`
}

// GetTimeout returns the configured timeout for the Helm uninstall action,
// or the given default.
func (in Uninstall) GetTimeout(defaultTimeout metav1.Duration) metav1.Duration {
	switch in.Timeout {
	case nil:
		return defaultTimeout
	default:
		return *in.Timeout
	}
}

// HelmReleaseStatus defines the observed state of HelmRelease
type HelmReleaseStatus struct {
	// ObservedGeneration is the last reconciled generation.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions holds the conditions for the HelmRelease.
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`

	// LastAppliedRevision is the revision of the last successfully applied source.
	// +optional
	LastAppliedRevision string `json:"lastAppliedRevision,omitempty"`

	// LastAttemptedRevision is the revision of the last reconciliation attempt.
	// +optional
	LastAttemptedRevision string `json:"lastAttemptedRevision,omitempty"`

	// LastAttemptedValuesChecksum is the SHA1 checksum of the values of the last reconciliation attempt.
	// +optional
	LastAttemptedValuesChecksum string `json:"lastAttemptedValuesChecksum,omitempty"`

	// LastReleaseRevision is the revision of the last successful Helm release.
	// +optional
	LastReleaseRevision int `json:"lastReleaseRevision,omitempty"`

	// HelmChart is the namespaced name of the HelmChart resource created by
	// the controller for the HelmRelease.
	// +optional
	HelmChart string `json:"helmChart,omitempty"`

	// Failures is the reconciliation failure count. It is reset after a successful
	// reconciliation.
	// +optional
	Failures int64 `json:"failures,omitempty"`
}

// GetHelmChart returns the namespace and name of the HelmChart.
func (in HelmReleaseStatus) GetHelmChart() (string, string) {
	if in.HelmChart == "" {
		return "", ""
	}
	split := strings.Split(in.HelmChart, string(types.Separator))
	return split[0], split[1]
}

// HelmReleaseProgressing resets the conditions of the given HelmRelease to a single
// ReadyCondition with status ConditionUnknown.
func HelmReleaseProgressing(hr HelmRelease) HelmRelease {
	hr.Status.Conditions = []Condition{
		{
			Type:               ReadyCondition,
			Status:             corev1.ConditionUnknown,
			LastTransitionTime: metav1.Now(),
			Reason:             ProgressingReason,
			Message:            "reconciliation in progress",
		},
	}
	return hr
}

// SetHelmReleaseCondition sets the given condition with the given status, reason and message
// on the HelmRelease.
func SetHelmReleaseCondition(hr *HelmRelease, condition string, status corev1.ConditionStatus, reason, message string) {
	hr.Status.Conditions = filterOutCondition(hr.Status.Conditions, condition)
	hr.Status.Conditions = append(hr.Status.Conditions, Condition{
		Type:               condition,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	})
}

// SetHelmReleaseReadiness sets the ReadyCondition, ObservedGeneration, LastAttemptedRevision,
// and LastReleaseRevision, on the HelmRelease.
func SetHelmReleaseReadiness(hr *HelmRelease, status corev1.ConditionStatus, reason, message string, revision string, releaseRevision int, valuesChecksum string) {
	SetHelmReleaseCondition(hr, ReadyCondition, status, reason, message)
	hr.Status.ObservedGeneration = hr.Generation
	hr.Status.LastAttemptedRevision = revision
	hr.Status.LastReleaseRevision = releaseRevision
	hr.Status.LastAttemptedValuesChecksum = valuesChecksum
}

// HelmReleaseNotReady registers a failed release attempt of the given HelmRelease.
func HelmReleaseNotReady(hr HelmRelease, revision string, releaseRevision int, valuesChecksum, reason, message string) HelmRelease {
	SetHelmReleaseReadiness(&hr, corev1.ConditionFalse, reason, message, revision, releaseRevision, valuesChecksum)
	hr.Status.Failures = hr.Status.Failures + 1
	return hr
}

// HelmReleaseReady registers a successful release attempt of the given HelmRelease.
func HelmReleaseReady(hr HelmRelease, revision string, releaseRevision int, valuesChecksum, reason, message string) HelmRelease {
	SetHelmReleaseReadiness(&hr, corev1.ConditionTrue, reason, message, revision, releaseRevision, valuesChecksum)
	hr.Status.LastAppliedRevision = revision
	hr.Status.Failures = 0
	return hr
}

// ShouldUpgrade determines if an Helm upgrade action needs to be performed for the given HelmRelease.
func ShouldUpgrade(hr HelmRelease, revision string, releaseRevision int, valuesChecksum string) bool {
	switch {
	case hr.Status.LastAttemptedRevision != revision:
		return true
	case hr.Status.LastReleaseRevision != releaseRevision:
		return true
	case hr.Generation != hr.Status.ObservedGeneration:
		return true
	case hr.Status.LastAttemptedValuesChecksum != valuesChecksum:
		return true
	case hr.Status.Failures > 0 &&
		(hr.Spec.GetUpgrade().MaxRetries < 0 || hr.Status.Failures < int64(hr.Spec.GetUpgrade().MaxRetries)):
		return true
	default:
		return false
	}
}

// ShouldTest determines if a Helm test actions needs to be performed for the given HelmRelease.
func ShouldTest(hr HelmRelease) bool {
	if hr.Spec.Test.Enable {
		for _, c := range hr.Status.Conditions {
			if c.Status == corev1.ConditionTrue && (c.Type == InstalledCondition || c.Type == UpgradedCondition) {
				return true
			}
		}
	}
	return false
}

// ShouldRollback determines if a Helm rollback action needs to be performed for the given HelmRelease.
func ShouldRollback(hr HelmRelease, releaseRevision int) bool {
	if hr.Spec.GetRollback().Enable {
		if hr.Status.LastReleaseRevision <= releaseRevision {
			return false
		}
		for _, c := range hr.Status.Conditions {
			if c.Type == UpgradedCondition && c.Status == corev1.ConditionFalse {
				return true
			}
		}
	}
	return false
}

// ShouldUninstall determines if a Helm uninstall action needs to be performed for the given HelmRelease.
func ShouldUninstall(hr HelmRelease, releaseRevision int) bool {
	if releaseRevision <= 0 {
		return false
	}
	for _, c := range hr.Status.Conditions {
		if c.Type == InstalledCondition && c.Status == corev1.ConditionFalse {
			return true
		}
	}
	return false
}

const (
	// ReconcileAtAnnotation is the annotation used for triggering a
	// reconciliation outside of the defined schedule.
	ReconcileAtAnnotation string = "fluxcd.io/reconcileAt"

	// SourceIndexKey is the key used for indexing HelmReleases based on
	// their sources.
	SourceIndexKey string = ".metadata.source"
)

// +genclient
// +genclient:Namespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=hr
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description=""
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].message",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// HelmRelease is the Schema for the helmreleases API
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmReleaseSpec   `json:"spec,omitempty"`
	Status HelmReleaseStatus `json:"status,omitempty"`
}

// GetValues unmarshals the raw values to a map[string]interface{}
// and returns the result.
func (in HelmRelease) GetValues() map[string]interface{} {
	var values map[string]interface{}
	if in.Spec.Values != nil {
		_ = json.Unmarshal(in.Spec.Values.Raw, &values)
	}
	return values
}

// GetReleaseName returns the configured release name, or a composition of
// '[TargetNamespace-]Name'.
func (in HelmRelease) GetReleaseName() string {
	if in.Spec.ReleaseName != "" {
		return in.Spec.ReleaseName
	}
	if in.Spec.TargetNamespace != "" {
		return strings.Join([]string{in.Spec.TargetNamespace, in.Name}, "-")
	}
	return in.Name
}

// GetReleaseNamespace returns the configured TargetNamespace, or the namespace
// of the HelmRelease.
func (in HelmRelease) GetReleaseNamespace() string {
	switch {
	case in.Spec.TargetNamespace != "":
		return in.Spec.TargetNamespace
	default:
		return in.Namespace
	}
}

// GetHelmChartName returns the name used by the controller for the HelmChart creation.
func (in HelmRelease) GetHelmChartName() string {
	return strings.Join([]string{in.Namespace, in.Name}, "-")
}

// GetTimeout returns the configured Timeout, or the default of 300s.
func (in HelmRelease) GetTimeout() metav1.Duration {
	switch in.Spec.Timeout {
	case nil:
		return metav1.Duration{Duration: 300 * time.Second}
	default:
		return *in.Spec.Timeout
	}
}

// GetMaxHistory returns the configured MaxHistory, or the default of 10.
func (in HelmRelease) GetMaxHistory() int {
	switch in.Spec.MaxHistory {
	case nil:
		return 10
	default:
		return *in.Spec.MaxHistory
	}
}

// +kubebuilder:object:root=true

// HelmReleaseList contains a list of HelmRelease
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmRelease{}, &HelmReleaseList{})
}

// filterOutCondition returns a new slice of conditions without the
// condition of the given type.
func filterOutCondition(conditions []Condition, condition string) []Condition {
	var newConditions []Condition
	for _, c := range conditions {
		if c.Type == condition {
			continue
		}
		newConditions = append(newConditions, c)
	}
	return newConditions
}
