package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/redhat-developer/tekton-jsongenerator/jsonschema"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	resource "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	triggersv1alpha1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	triggersv1beta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	knative "knative.dev/pkg/apis"
	"os"
	"reflect"
)

func arrayOrStringMapper(i reflect.Type) *jsonschema.Type {
	if (i == reflect.TypeOf(pipelinev1.ParamValue{}) || i == reflect.TypeOf(pipelinev1beta1.ArrayOrString{})) {
		return &jsonschema.Type{
			OneOf: []*jsonschema.Type{
				{
					Type: "string",
				},
				{
					Type: "array",
					Items: &jsonschema.Type{
						Type: "string",
					},
				},
			},
		}
	}
	if (i == reflect.TypeOf(v1.Duration{})) {
		return &jsonschema.Type{
			Type:    "string",
			Pattern: "^[-+]?([0-9]*(\\.[0-9]*)?(ns|us|µs|μs|ms|s|m|h))+$",
		}
	}
	if (i == reflect.TypeOf(v1.Time{})) {
		return &jsonschema.Type{
			Type:   "string",
			Format: "data-time",
		}
	}
	if (i == reflect.TypeOf(triggersv1alpha1.TriggerResourceTemplate{}) || i == reflect.TypeOf(triggersv1beta1.TriggerResourceTemplate{})) {
		return &jsonschema.Type{
			Type:                 "object",
			AdditionalProperties: []byte("true"),
			Properties:           orderedmap.New()}
	}
	if (i == reflect.TypeOf(knative.VolatileTime{})) {
		return &jsonschema.Type{
			Type:   "string",
			Format: "data-time",
		}
	}
	if (i == reflect.TypeOf(pipelinev1.WhenExpression{}) || i == reflect.TypeOf(pipelinev1beta1.WhenExpression{})) {
		properties := orderedmap.New()
		properties.Set("input", &jsonschema.Type{
			Type: "string",
		})
		properties.Set("operator", &jsonschema.Type{
			Type: "string",
		})
		properties.Set("values", &jsonschema.Type{
			Type: "array",
			Items: &jsonschema.Type{
				Type: "string",
			},
		})
		return &jsonschema.Type{
			Type:                 "object",
			AdditionalProperties: []byte("false"),
			Required:             []string{"input", "operator", "values"},
			Properties:           properties}
	}
	if (i == reflect.TypeOf(runtime.RawExtension{})) {
		return &jsonschema.Type{
			Type:                 "object",
			AdditionalProperties: []byte("true"),
			Properties:           orderedmap.New()}
	}
	if (i == reflect.TypeOf(v1.FieldsV1{})) {
		return &jsonschema.Type{
			Type:                 "object",
			AdditionalProperties: []byte("true"),
			Properties:           orderedmap.New()}
	}
	if (i == reflect.TypeOf(apiextensionsv1.JSON{})) {
		return &jsonschema.Type{
			OneOf: []*jsonschema.Type{
				{
					Type: "boolean",
				},
				{
					Type: "integer",
				},
				{
					Type: "number",
				},
				{
					Type: "string",
				},
				{
					Type: "array",
					Items: &jsonschema.Type{
						OneOf: []*jsonschema.Type{
							{
								Type: "string",
							},
							{
								Type:                 "object",
								AdditionalProperties: []byte("true"),
								Properties:           orderedmap.New(),
							},
						},
					},
				},
				{
					Type:                 "object",
					AdditionalProperties: []byte("true"),
					Properties:           orderedmap.New(),
				},
				{
					Type: "null",
				},
			},
		}
	}
	return nil
}

func dump(v interface{}, apiVersion string, kind string) {
	fmt.Printf("Starting generation of %s %s\n", apiVersion, kind)
	filename := fmt.Sprintf("%s_%s.json", apiVersion, kind)
	reflector := jsonschema.Reflector{
		TypeMapper: arrayOrStringMapper,
	}
	reflect := reflector.Reflect(v)
	JSON, _ := reflect.MarshalJSON()
	file, _ := os.Create(filename)
	defer file.Close()
	var out bytes.Buffer
	json.Indent(&out, JSON, "", "  ")
	out.WriteTo(file)
	index, _ := os.OpenFile("index.properties", os.O_WRONLY|os.O_APPEND, 0)
	index.WriteString(filename)
	index.WriteString("\n")
}

func main() {
	os.Create("index.properties")
	os.Mkdir("tekton.dev", os.ModePerm)
	os.Mkdir("triggers.tekton.dev", os.ModePerm)
	dump(&resource.PipelineResource{}, "tekton.dev/v1alpha1", "PipelineResource")
	dump(&resource.PipelineResourceList{}, "tekton.dev/v1alpha1", "PipelineResourceList")

	dump(&pipelinev1.Pipeline{}, "tekton.dev/v1", "Pipeline")
	dump(&pipelinev1.PipelineList{}, "tekton.dev/v1", "PipelineList")
	dump(&pipelinev1.PipelineRun{}, "tekton.dev/v1", "PipelineRun")
	dump(&pipelinev1.PipelineRunList{}, "tekton.dev/v1", "PipelineRunList")
	dump(&pipelinev1.Task{}, "tekton.dev/v1", "Task")
	dump(&pipelinev1.TaskList{}, "tekton.dev/v1", "TaskList")
	dump(&pipelinev1.TaskRun{}, "tekton.dev/v1", "TaskRun")
	dump(&pipelinev1.TaskRunList{}, "tekton.dev/v1", "TaskRunList")

	dump(&pipelinev1beta1.Pipeline{}, "tekton.dev/v1beta1", "Pipeline")
	dump(&pipelinev1beta1.PipelineList{}, "tekton.dev/v1beta1", "PipelineList")
	dump(&pipelinev1beta1.PipelineRun{}, "tekton.dev/v1beta1", "PipelineRun")
	dump(&pipelinev1beta1.PipelineRunList{}, "tekton.dev/v1beta1", "PipelineRunList")
	dump(&pipelinev1beta1.Task{}, "tekton.dev/v1beta1", "Task")
	dump(&pipelinev1beta1.TaskList{}, "tekton.dev/v1beta1", "TaskList")
	dump(&pipelinev1beta1.TaskRun{}, "tekton.dev/v1beta1", "TaskRun")
	dump(&pipelinev1beta1.TaskRunList{}, "tekton.dev/v1beta1", "TaskRunList")
	dump(&pipelinev1beta1.ClusterTask{}, "tekton.dev/v1beta1", "ClusterTask")
	dump(&pipelinev1beta1.ClusterTaskList{}, "tekton.dev/v1beta1", "ClusterTaskList")

	dump(&triggersv1alpha1.EventListener{}, "triggers.tekton.dev/v1alpha1", "EventListener")
	dump(&triggersv1alpha1.EventListenerList{}, "triggers.tekton.dev/v1alpha1", "EventListenerList")
	dump(&triggersv1alpha1.Trigger{}, "triggers.tekton.dev/v1alpha1", "Trigger")
	dump(&triggersv1alpha1.Trigger{}, "triggers.tekton.dev/v1alpha1", "TriggerList")
	dump(&triggersv1alpha1.TriggerTemplate{}, "triggers.tekton.dev/v1alpha1", "TriggerTemplate")
	dump(&triggersv1alpha1.TriggerTemplateList{}, "triggers.tekton.dev/v1alpha1", "TriggerTemplateList")
	dump(&triggersv1alpha1.TriggerBinding{}, "triggers.tekton.dev/v1alpha1", "TriggerBinding")
	dump(&triggersv1alpha1.TriggerBindingList{}, "triggers.tekton.dev/v1alpha1", "TriggerBindingList")
	dump(&triggersv1alpha1.ClusterTriggerBinding{}, "triggers.tekton.dev/v1alpha1", "ClusterTriggerBinding")
	dump(&triggersv1alpha1.ClusterTriggerBindingList{}, "triggers.tekton.dev/v1alpha1", "ClusterTriggerBindingList")
	dump(&triggersv1alpha1.ClusterInterceptor{}, "triggers.tekton.dev/v1alpha1", "ClusterInterceptor")
	dump(&triggersv1alpha1.ClusterInterceptorList{}, "triggers.tekton.dev/v1alpha1", "ClusterInterceptorList")

	dump(&triggersv1beta1.EventListener{}, "triggers.tekton.dev/v1beta1", "EventListener")
	dump(&triggersv1beta1.EventListenerList{}, "triggers.tekton.dev/v1beta1", "EventListenerList")
	dump(&triggersv1beta1.Trigger{}, "triggers.tekton.dev/v1beta1", "Trigger")
	dump(&triggersv1beta1.Trigger{}, "triggers.tekton.dev/v1beta1", "TriggerList")
	dump(&triggersv1beta1.TriggerTemplate{}, "triggers.tekton.dev/v1beta1", "TriggerTemplate")
	dump(&triggersv1beta1.TriggerTemplateList{}, "triggers.tekton.dev/v1beta1", "TriggerTemplateList")
	dump(&triggersv1beta1.TriggerBinding{}, "triggers.tekton.dev/v1beta1", "TriggerBinding")
	dump(&triggersv1beta1.TriggerBindingList{}, "triggers.tekton.dev/v1beta1", "TriggerBindingList")
	dump(&triggersv1beta1.ClusterTriggerBinding{}, "triggers.tekton.dev/v1beta1", "ClusterTriggerBinding")
	dump(&triggersv1beta1.ClusterTriggerBindingList{}, "triggers.tekton.dev/v1beta1", "ClusterTriggerBindingList")
}
