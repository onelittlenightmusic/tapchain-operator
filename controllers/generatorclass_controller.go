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
	// apix "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	basev1 "op/api/v1"
)

var (
	generatorClassOwnerKey = ".metadata.controller"
	generatorClassApiGVStr = basev1.GroupVersion.String()
)

// GeneratorClassReconciler reconciles a GeneratorClass object
type GeneratorClassReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=base.tapchain.org,resources=generatorclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=base.tapchain.org,resources=generatorclasses/status,verbs=get;update;patch

func (r *GeneratorClassReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("generatorclass", req.NamespacedName)

	// getCrdNames := func(generatorClass *basev1.GeneratorClass) (string, string, string, string) {
	// 	name := generatorClass.Name
	// 	group := "base.tapchain.org"
	// 	pluralName := name + "s"
	// 	generatorName := pluralName + "." + group // CRD naming rule
	// 	return name, group, pluralName, generatorName

	// }

	// constructCRDForTable := func(generatorClass *basev1.GeneratorClass) (*apix.CustomResourceDefinition, error) {
	// 	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	// 	name, group, pluralName, generatorName := getCrdNames(generatorClass)

	// 	generator := &basev1.Generator{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Labels:      make(map[string]string),
	// 			Annotations: make(map[string]string),
	// 			Name:        generatorName,
	// 			Namespace:   generatorClass.Namespace,
	// 		},
	// 	}
	// 	generator.Spec = basev1.GeneratorSpec{
	// 		Names: apix.CustomResourceDefinitionNames{
	// 			Plural:   pluralName,
	// 			Singular: name,
	// 			Kind:     strings.Title(name),
	// 		},
	// 		Scope: "Namespaced",
	// 	}
	// 	constructProps := func(columns []basev1.PathClassSpec, props map[string]apix.JSONSchemaProps) map[string]apix.JSONSchemaProps {
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
	// 		return constructProps([]basev1.PathClassSpec{{Name: field, Type: _type}}, props)
	// 	}
	// 	schema := apix.CustomResourceValidation{
	// 		OpenAPIV3Schema: &apix.JSONSchemaProps{
	// 			Type: "object",
	// 			Properties: constructmSimpleProp("spec", "object",
	// 				constructProps(generatorClass.Spec.Columns, nil),
	// 			),
	// 		},
	// 	}
	// 	additionalPrinterColumns := []apix.CustomResourceColumnDefinition{}
	// 	constructPrinterColumn := func(columnSpec basev1.PathClassSpec) apix.CustomResourceColumnDefinition {
	// 		return apix.CustomResourceColumnDefinition{
	// 			Name:     columnSpec.Name,
	// 			Type:     columnSpec.Type,
	// 			JSONPath: ".spec." + columnSpec.Name,
	// 		}
	// 	}
	// 	for _, columnSpec := range generatorClass.Spec.Columns {
	// 		additionalPrinterColumns = append(additionalPrinterColumns, constructPrinterColumn(columnSpec))
	// 	}
	// 	generator.Spec.Versions = append(generator.Spec.Versions,
	// 		apix.CustomResourceDefinitionVersion{
	// 			Name:                     "v1",
	// 			Served:                   true,
	// 			Storage:                  true,
	// 			Schema:                   &schema,
	// 			AdditionalPrinterColumns: additionalPrinterColumns,
	// 		},
	// 	)

	// 	if err := ctrl.SetControllerReference(generatorClass, generator, r.Scheme); err != nil {
	// 		return nil, err
	// 	}

	// 	return generator, nil
	// }

	// var generatorClass basev1.GeneratorClass
	// if err := r.Get(ctx, req.NamespacedName, &generatorClass); err != nil {
	// 	log.Error(err, "unable to fetch Table")
	// 	return ctrl.Result{}, client.IgnoreNotFound(err)
	// }

	// var generator basev1.Generator
	// _, _, _, generatorName := getCrdNames(&generatorClass)
	// if err := r.Get(ctx, client.ObjectKey{Name: generatorName}, &generator); err != nil {
	// 	if errors.IsNotFound(err) {
	// 		var generators basev1.GeneratorList
	// 		// if err := r.List(ctx, &generators, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
	// 		if err := r.List(ctx, &generators, client.MatchingFieldsg{generatorOwnerKey: req.Name}); err != nil {
	// 			log.Error(err, "unable to list CRD")
	// 			return ctrl.Result{}, err
	// 		}

	// 		generator, err := constructCRDForTable(&generatorClass)
	// 		if err != nil {
	// 			log.Error(err, "unable to construct generator from template")
	// 			// don't bother requeuing until we get a change to the spec
	// 			return ctrl.Result{}, nil
	// 		}

	// 		// ...and create it on the cluster
	// 		if err := r.Create(ctx, generator); err != nil {
	// 			log.Error(err, "unable to create CRD for Table", "generator", generator)
	// 			return ctrl.Result{}, err
	// 		}
	// 		log.V(1).Info("created CRD for Table run", "generator", generator)
	// 	}
	// }

	return ctrl.Result{}, nil
}

func (r *GeneratorClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// if err := apix.AddToScheme(mgr.GetScheme()); err != nil {
	// 	return nil
	// }

	// if err := mgr.GetFieldIndexer().IndexField(&basev1.Generator{}, ownerKey, func(rawObj runtime.Object) []string {
	// 	// grab the CRD object, extract the owner...
	// 	generator := rawObj.(*basev1.Generator)
	// 	owner := metav1.GetControllerOf(generator)

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
		For(&basev1.GeneratorClass{}).
		// Owns(&basev1.Generator{}).
		Complete(r)
}
