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

	// "encoding/json"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	basev1 "op/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GeneratorReconciler reconciles a Generator object
type GeneratorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type MaterialInfo struct {
	Class *basev1.MaterialClass
	Name  string
}

type MaterialInfoList map[string]*MaterialInfo

type GeneratorError struct {
	errorString string
}

func (g *GeneratorError) Error() string {
	return g.errorString
}

var (
	generatorController = ".metadata.controller"
	generatorApiGVStr   = basev1.GroupVersion.String()
)

const (
	materialClass = "Material"
)

// +kubebuilder:rbac:groups=base.tapchain.org,resources=generators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=base.tapchain.org,resources=generators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=base.tapchain.org,resources=materials,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=base.tapchain.org,resources=materials/status,verbs=get;update;patch

func (r *GeneratorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("generator", req.NamespacedName)

	// get generator if it exists
	var generator basev1.Generator
	if err := r.Get(ctx, req.NamespacedName, &generator); err != nil {
		log.Error(err, "unable to fetch generator")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// get generatorClass
	generatorClass, err := r.getGeneratorClass(ctx, &generator)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("GeneratorClass found")

	// get MaterialClass for outputs of generator
	var materialInfoList MaterialInfoList
	if err := r.getIOMaterialInfoList(ctx, generator.Namespace, generatorClass.Spec.Outputs, generator.Spec.Outputs, &materialInfoList); err != nil {
		log.Error(err, "unable to get MaterialInfo")
		return ctrl.Result{}, err
	}

	// create Material for output name
	for pathName, materialInfo := range materialInfoList {
		if err := r.createPath(ctx, generator, pathName, materialInfo); err != nil {
			log.Error(err, "unable to create Path")
			return ctrl.Result{}, err
		}
	}

	// get MaterialClass for inputs of generator
	if err := r.getIOMaterialInfoList(ctx, generator.Namespace, generatorClass.Spec.Inputs, generator.Spec.Inputs, &materialInfoList); err != nil {
		log.Error(err, "unable to get MaterialInfo")
		return ctrl.Result{}, err
	}

	// create Material for input name
	for pathName, materialInfo := range materialInfoList {
		if err := r.createPath(ctx, generator, pathName, materialInfo); err != nil {
			log.Error(err, "unable to create Path")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *GeneratorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&basev1.Material{}, generatorController, func(rawObj runtime.Object) []string {
		material := rawObj.(*basev1.Material)
		owner := metav1.GetControllerOf(material)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != generatorApiGVStr || owner.Kind != materialClass {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&basev1.Generator{}).
		Owns(&basev1.Material{}).
		Complete(r)
}

func (r *GeneratorReconciler) getGeneratorClass(ctx context.Context, generator *basev1.Generator) (basev1.GeneratorClass, error) {
	var generatorClass basev1.GeneratorClass
	if err := r.Get(ctx, client.ObjectKey{Name: generator.Spec.GeneratorClassName, Namespace: generator.Namespace}, &generatorClass); err != nil {
		return basev1.GeneratorClass{}, err
	}
	return generatorClass, nil
}

func (r *GeneratorReconciler) getIOMaterialInfoList(ctx context.Context, namespace string, pathClass []basev1.PathClassSpec, ioSpec basev1.IOSpec, materialInfoList *MaterialInfoList) error {
	log := r.Log
	rtn := make(MaterialInfoList)

	for _, outputPathClass := range pathClass {
		var materialClass basev1.MaterialClass
		if err := r.Get(ctx, client.ObjectKey{Name: outputPathClass.MaterialClassName, Namespace: namespace}, &materialClass); err != nil {
			log.Error(err, "unable to find MaterialClass", "materialClass", outputPathClass.MaterialClassName)
			return err
		}
		log.Info("found MaterialClass" + outputPathClass.MaterialClassName + " in Namespace " + namespace)

		materialIndex := outputPathClass.Name

		outputMaterialName, ok := ioSpec[materialIndex]
		if !ok {
			err := &GeneratorError{errorString: "unable to find PathClass defined GeneratorClass in Generator spec.outputs"}
			log.Error(err, "test")
			return err
		}

		log.Info("found path index" + outputMaterialName)
		rtn[materialIndex] = &MaterialInfo{Class: &materialClass, Name: outputMaterialName}
	}
	*materialInfoList = rtn
	// jsonString, _ := json.Marshal(rtn)
	// log.Info(string(jsonString))
	return nil
}

func getMaterialName(materialInfo *MaterialInfo) (string, string) {
	name := materialInfo.Name
	group := "base.tapchain.org"
	return name, group
}

func (r *GeneratorReconciler) createPath(ctx context.Context, generator basev1.Generator, pathName string, materialInfo *MaterialInfo) error {
	log := r.Log.WithValues("generator", client.ObjectKey{Name: generator.Name})
	var material basev1.Material

	materialName, _ := getMaterialName(materialInfo)
	if err := r.Get(ctx, client.ObjectKey{Name: materialName, Namespace: generator.Namespace}, &material); err != nil {
		if apierrors.IsNotFound(err) {
			// var materials basev1.MaterialList
			// if err := r.List(ctx, &materials, client.MatchingFields{materialOwnerKey: req.Name}); err != nil {
			// 	log.Error(err, "unable to list materials")
			// 	return ctrl.Result{}, err
			// }

			basicLabels := makeLabels(pathName, materialInfo.Name, materialInfo.Class.Name)

			material, err := r.constructMaterial(&generator, materialInfo, basicLabels)
			if err != nil {
				log.Error(err, "unable to construct material from template")
				return &GeneratorError{errorString: "unable to construct material from template"}
				// don't bother requeuing until we get a change to the spec
			}

			// ...and create it on the cluster
			if err := r.Create(ctx, material); err != nil {
				log.Error(err, "unable to create material object", "material", material)
				return &GeneratorError{errorString: "unable to create material object"}
			}
			log.V(1).Info("created material", "material", material)
			return nil
		}
		log.Error(err, "unable to get material")
		return &GeneratorError{errorString: "unable to get material"}
	}
	return nil
}

func (r *GeneratorReconciler) constructMaterial(generator *basev1.Generator, materialInfo *MaterialInfo, labels map[string]string) (*basev1.Material, error) {
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	name, _ := getMaterialName(materialInfo)

	material := &basev1.Material{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      labels,
			Annotations: make(map[string]string),
			Name:        name,
			Namespace:   generator.Namespace,
		},
	}
	material.Spec = basev1.MaterialSpec{
		MaterialClassName: materialInfo.Class.Name,
		Metrics:           []basev1.Metric{},
	}

	if err := ctrl.SetControllerReference(generator, material, r.Scheme); err != nil {
		return nil, err
	}

	return material, nil
}

func makeLabels(pathName string, materialName string, materialClassName string) map[string]string {
	return map[string]string{
		"pathName":          pathName,
		"materialName":      materialName,
		"materialClassName": materialClassName,
	}
}
