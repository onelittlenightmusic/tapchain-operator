/*


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

package controllers

import (
	"context"
	// "strings"

	"github.com/go-logr/logr"
	apix "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	basev1 "op/api/v1"
)

var (
	materialOwnerKey = ".metadata.controller"
	materialApiGVStr = apix.SchemeGroupVersion.String()
)

// MaterialClassReconciler reconciles a MaterialClass object
type MaterialClassReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=base.tapchain.org,resources=materialclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=base.tapchain.org,resources=materialclasses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions/status,verbs=get;update;patch

func (r *MaterialClassReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("materialclass", req.NamespacedName)

	// getCrdNames := func(materialClass *basev1.MaterialClass) (string, string, string, string) {
	// 	name := materialClass.Name
	// 	group := "user.k8sasdb.org"
	// 	pluralName := name + "s"
	// 	crdName := pluralName + "." + group // CRD naming rule
	// 	return name, group, pluralName, crdName

	// }

	// constructCRDForMaterial := func(materialClass *basev1.MaterialClass) (*apix.CustomResourceDefinition, error) {
	// 	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	// 	name, group, pluralName, crdName := getCrdNames(materialClass)

	// 	crd := &apix.CustomResourceDefinition{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Labels:      make(map[string]string),
	// 			Annotations: make(map[string]string),
	// 			Name:        crdName,
	// 			Namespace:   materialClass.Namespace,
	// 		},
	// 	}
	// 	crd.Spec = apix.CustomResourceDefinitionSpec{
	// 		Group: group,
	// 		Names: apix.CustomResourceDefinitionNames{
	// 			Plural:   pluralName,
	// 			Singular: name,
	// 			Kind:     strings.Title(name),
	// 		},
	// 		Scope: "Namespaced",
	// 	}
	// 	constructProps := func(columns []basev1.FieldSpec, props map[string]apix.JSONSchemaProps) map[string]apix.JSONSchemaProps {
	// 		rtn := map[string]apix.JSONSchemaProps{}
	// 		for _, columnSpec := range columns {
	// 			rtn[columnSpec.Name] = apix.JSONSchemaProps{
	// 				Type:       columnSpec.Type,
	// 				Properties: props,
	// 			}
	// 		}
	// 		return rtn
	// 	}
	// 	constructmSimpleProp := func(field string, _type string, props map[string]apix.JSONSchemaProps) map[string]apix.JSONSchemaProps {
	// 		return constructProps([]basev1.FieldSpec{{Name: field, Type: _type}}, props)
	// 	}
	// 	schema := apix.CustomResourceValidation{
	// 		OpenAPIV3Schema: &apix.JSONSchemaProps{
	// 			Type: "object",
	// 			Properties: constructmSimpleProp("spec", "object",
	// 				constructProps(materialClass.Spec.Fields, nil),
	// 			),
	// 		},
	// 	}
	// 	additionalPrinterColumns := []apix.CustomResourceColumnDefinition{}
	// 	constructPrinterColumn := func(columnSpec basev1.FieldSpec) apix.CustomResourceColumnDefinition {
	// 		return apix.CustomResourceColumnDefinition{
	// 			Name:     columnSpec.Name,
	// 			Type:     columnSpec.Type,
	// 			JSONPath: ".spec." + columnSpec.Name,
	// 		}
	// 	}
	// 	for _, columnSpec := range materialClass.Spec.Fields {
	// 		additionalPrinterColumns = append(additionalPrinterColumns, constructPrinterColumn(columnSpec))
	// 	}
	// 	crd.Spec.Versions = append(crd.Spec.Versions,
	// 		apix.CustomResourceDefinitionVersion{
	// 			Name:                     "v1",
	// 			Served:                   true,
	// 			Storage:                  true,
	// 			Schema:                   &schema,
	// 			AdditionalPrinterColumns: additionalPrinterColumns,
	// 		},
	// 	)

	// 	if err := ctrl.SetControllerReference(materialClass, crd, r.Scheme); err != nil {
	// 		return nil, err
	// 	}

	// 	return crd, nil
	// }

	// var materialClass basev1.MaterialClass
	// if err := r.Get(ctx, req.NamespacedName, &materialClass); err != nil {
	// 	log.Error(err, "unable to fetch Material")
	// 	return ctrl.Result{}, client.IgnoreNotFound(err)
	// }

	// var crd apix.CustomResourceDefinition
	// _, _, _, crdName := getCrdNames(&materialClass)
	// if err := r.Get(ctx, client.ObjectKey{Name: crdName}, &crd); err != nil {
	// 	if errors.IsNotFound(err) {
	// 		var crds apix.CustomResourceDefinitionList
	// 		// if err := r.List(ctx, &crds, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
	// 		if err := r.List(ctx, &crds, client.MatchingFields{materialOwnerKey: req.Name}); err != nil {
	// 			log.Error(err, "unable to list CRD")
	// 			return ctrl.Result{}, err
	// 		}

	// 		crd, err := constructCRDForMaterial(&materialClass)
	// 		if err != nil {
	// 			log.Error(err, "unable to construct crd from template")
	// 			// don't bother requeuing until we get a change to the spec
	// 			return ctrl.Result{}, nil
	// 		}

	// 		// ...and create it on the cluster
	// 		if err := r.Create(ctx, crd); err != nil {
	// 			log.Error(err, "unable to create CRD for Material", "crd", crd)
	// 			return ctrl.Result{}, err
	// 		}
	// 		log.V(1).Info("created CRD for Material run", "crd", crd)
	// 	}
	// }

	return ctrl.Result{}, nil
}

func (r *MaterialClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// if err := apix.AddToScheme(mgr.GetScheme()); err != nil {
	// 	return nil
	// }

	// if err := mgr.GetFieldIndexer().IndexField(&apix.CustomResourceDefinition{}, materialOwnerKey, func(rawObj runtime.Object) []string {
	// 	// grab the CRD object, extract the owner...
	// 	crd := rawObj.(*apix.CustomResourceDefinition)
	// 	owner := metav1.GetControllerOf(crd)

	// 	if owner == nil {
	// 		return nil
	// 	}
	// 	if owner.APIVersion != apiGVStr || owner.Kind != "Material" {
	// 		return nil
	// 	}

	// 	return []string{owner.Name}
	// }); err != nil {
	// 	return err
	// }
	return ctrl.NewControllerManagedBy(mgr).
		For(&basev1.MaterialClass{}).
		// Owns(&apix.CustomResourceDefinition{}).
		Complete(r)
}
