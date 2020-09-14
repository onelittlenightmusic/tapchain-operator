#!/bin/bash

kubectl apply -f base_v1_generatorclass.yaml
kubectl apply -f base_v1_generatorclass2.yaml
kubectl apply -f base_v1_materialclass.yaml
kubectl apply -f base_v1_material.yaml
# kubectl apply -f base_v1_generator.yaml
kubectl apply -f base_v1_generator2.yaml