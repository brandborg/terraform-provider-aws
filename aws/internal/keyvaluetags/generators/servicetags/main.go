//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

const filename = `service_tags_gen.go`

// Representing types such as []*athena.Tag, []*ec2.Tag, ...
var sliceServiceNames = []string{
	"acm",
	"acmpca",
	"appmesh",
	"athena",
	/* "autoscaling", // includes extra PropagateAtLaunch, skip for now */
	"cloud9",
	"cloudformation",
	"cloudfront",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"codebuild",
	"codedeploy",
	"codepipeline",
	"configservice",
	"databasemigrationservice",
	"datapipeline",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"docdb",
	"dynamodb",
	"ec2",
	"ecr",
	"ecs",
	"efs",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elb",
	"elbv2",
	"emr",
	"firehose",
	"fms",
	"fsx",
	"gamelift",
	"globalaccelerator",
	"iam",
	"inspector",
	"iot",
	"iotanalytics",
	"iotevents",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kms",
	"licensemanager",
	"lightsail",
	"mediastore",
	"neptune",
	"networkmanager",
	"organizations",
	"quicksight",
	"ram",
	"rds",
	"redshift",
	"resourcegroupstaggingapi",
	"route53",
	"route53resolver",
	"s3",
	"sagemaker",
	"secretsmanager",
	"serverlessapplicationrepository",
	"servicecatalog",
	"servicediscovery",
	"sfn",
	"sns",
	"ssm",
	"storagegateway",
	"swf",
	"transfer",
	"waf",
	"wafregional",
	"wafv2",
	"workspaces",
}

var mapServiceNames = []string{
	"accessanalyzer",
	"amplify",
	"apigateway",
	"apigatewayv2",
	"appstream",
	"appsync",
	"backup",
	"batch",
	"cloudwatchlogs",
	"codecommit",
	"codestarnotifications",
	"cognitoidentity",
	"cognitoidentityprovider",
	"dataexchange",
	"dlm",
	"eks",
	"glacier",
	"glue",
	"guardduty",
	"greengrass",
	"kafka",
	"kinesisvideo",
	"imagebuilder",
	"lambda",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mq",
	"opsworks",
	"qldb",
	"pinpoint",
	"resourcegroups",
	"securityhub",
	"sqs",
	"synthetics",
	"worklink",
}

type TemplateData struct {
	MapServiceNames   []string
	SliceServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(mapServiceNames)
	sort.Strings(sliceServiceNames)

	templateData := TemplateData{
		MapServiceNames:   mapServiceNames,
		SliceServiceNames: sliceServiceNames,
	}
	templateFuncMap := template.FuncMap{
		"TagKeyType":        keyvaluetags.ServiceTagKeyType,
		"TagPackage":        keyvaluetags.ServiceTagPackage,
		"TagType":           keyvaluetags.ServiceTagType,
		"TagType2":          keyvaluetags.ServiceTagType2,
		"TagTypeKeyField":   keyvaluetags.ServiceTagTypeKeyField,
		"TagTypeValueField": keyvaluetags.ServiceTagTypeValueField,
		"Title":             strings.Title,
	}

	tmpl, err := template.New("servicetags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/servicetags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"github.com/aws/aws-sdk-go/aws"
{{- range .SliceServiceNames }}
{{- if eq . (. | TagPackage) }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
{{- end }}
)

// map[string]*string handling
{{- range .MapServiceNames }}

// {{ . | Title }}Tags returns {{ . }} service tags.
func (tags KeyValueTags) {{ . | Title }}Tags() map[string]*string {
	return aws.StringMap(tags.Map())
}

// {{ . | Title }}KeyValueTags creates KeyValueTags from {{ . }} service tags.
func {{ . | Title }}KeyValueTags(tags map[string]*string) KeyValueTags {
	return New(tags)
}
{{- end }}

// []*SERVICE.Tag handling
{{- range .SliceServiceNames }}

{{- if . | TagKeyType }}
// {{ . | Title }}TagKeys returns {{ . }} service tag keys.
func (tags KeyValueTags) {{ . | Title }}TagKeys() []*{{ . | TagPackage }}.{{ . | TagKeyType }} {
	result := make([]*{{ . | TagPackage }}.{{ . | TagKeyType }}, 0, len(tags))

	for k := range tags.Map() {
		tagKey := &{{ . | TagPackage }}.{{ . | TagKeyType }}{
			{{ . | TagTypeKeyField }}: aws.String(k),
		}

		result = append(result, tagKey)
	}

	return result
}
{{- end }}

// {{ . | Title }}Tags returns {{ . }} service tags.
func (tags KeyValueTags) {{ . | Title }}Tags() []*{{ . | TagPackage }}.{{ . | TagType }} {
	result := make([]*{{ . | TagPackage }}.{{ . | TagType }}, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &{{ . | TagPackage }}.{{ . | TagType }}{
			{{ . | TagTypeKeyField }}:   aws.String(k),
			{{ . | TagTypeValueField }}: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// {{ . | Title }}KeyValueTags creates KeyValueTags from {{ . }} service tags.
{{- if . | TagType2 }}
// Accepts []*{{ . | TagPackage }}.{{ . | TagType }} and []*{{ . | TagPackage }}.{{ . | TagType2 }}.
func {{ . | Title }}KeyValueTags(tags interface{}) KeyValueTags {
	switch tags := tags.(type) {
	case []*{{ . | TagPackage }}.{{ . | TagType }}:
		m := make(map[string]*string, len(tags))

		for _, tag := range tags {
			m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
		}

		return New(m)
	case []*{{ . | TagPackage }}.{{ . | TagType2 }}:
		m := make(map[string]*string, len(tags))

		for _, tag := range tags {
			m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
		}

		return New(m)
	default:
		return New(nil)
	}
}
{{- else }}
func {{ . | Title }}KeyValueTags(tags []*{{ . | TagPackage }}.{{ . | TagType }}) KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
	}

	return New(m)
}
{{- end }}
{{- end }}
`
