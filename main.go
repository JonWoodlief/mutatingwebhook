package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	serviceAccountName = "nginx-serviceaccount"
	nodeLabelKey       = "restricted"
	nodeLabelValue     = "globalcatalog"
)

func main() {
	http.HandleFunc("/mutate", mutateHandler)

	if os.Getenv("TLS") == "true" {
		cert := "/etc/admission-webhook/tls/tls.crt"
		key := "/etc/admission-webhook/tls/tls.key"
		log.Print("Listening on port 443...")
		log.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
	} else {
		log.Print("Listening on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	// Read the admission review request from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse the admission review request

	admissionReview := v1.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		log.Printf("Failed to parse admission review request: %v", err)
		http.Error(w, "Failed to parse admission review request", http.StatusBadRequest)
		return
	}

	// Get the pod object from the admission review request
	pod := corev1.Pod{}
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
		log.Printf("Failed to parse pod object: %v", err)
		http.Error(w, "Failed to parse pod object", http.StatusBadRequest)
		return
	}
	//log name and namespace of the pod
	log.Printf("Initiating review for pod name: %s, namespace: %s", pod.Name, pod.Namespace)

	// Check if the pod is running under the specified service account
	if pod.Spec.ServiceAccountName == serviceAccountName {
		// Check if the pod is scheduled onto a node with the specified label
		if pod.Spec.NodeSelector == nil || pod.Spec.NodeSelector[nodeLabelKey] != nodeLabelValue {
			// If not, add the node selector to the pod
			if pod.Spec.NodeSelector == nil {
				pod.Spec.NodeSelector = make(map[string]string)
			}
			pod.Spec.NodeSelector[nodeLabelKey] = nodeLabelValue

			// Update the admission review response with the mutated pod object
			admissionReview.Response = &v1.AdmissionResponse{
				Allowed: true,
				Patch:   getPatch(&pod),
				PatchType: func() *v1.PatchType {
					pt := v1.PatchTypeJSONPatch
					return &pt
				}(),
			}
		}
	}

	// Write the admission review response to the response writer
	responseBody, err := json.Marshal(admissionReview)
	if err != nil {
		log.Printf("Failed to marshal admission review response: %v", err)
		http.Error(w, "Failed to marshal admission review response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseBody); err != nil {
		log.Printf("Failed to write admission review response: %v", err)
		http.Error(w, "Failed to write admission review response", http.StatusInternalServerError)
		return
	}

	log.Printf("Admission review response: %s", responseBody)
}

func getPatch(pod *corev1.Pod) []byte {
	patch := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/spec/nodeSelector",
			"value": pod.Spec.NodeSelector,
		},
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		log.Printf("Failed to marshal patch: %v", err)
		os.Exit(1)
	}

	return patchBytes
}
