// Copyright 2020 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const codeTemplate string = `/* Code generated by codegen/main.go. DO NOT EDIT. */

{{- if ne .SubNamespace "" }}
namespace Nakama.{{.SubNamespace}}
{{- else }}
namespace Nakama
{{- end }}
{
    using System;
    using System.Collections.Generic;
    using System.Runtime.Serialization;
    using System.Text;
    using System.Threading.Tasks;
    using TinyJson;

    /// <summary>
    /// An exception generated for <c>HttpResponse</c> objects don't return a success status.
    /// </summary>
    public sealed class ApiResponseException : Exception
    {
        public long StatusCode { get; }

        public int GrpcStatusCode { get; }

        public ApiResponseException(long statusCode, string content, int grpcCode) : base(content)
        {
            StatusCode = statusCode;
            GrpcStatusCode = grpcCode;
        }

        public ApiResponseException(string message, Exception e) : base(message, e)
        {
            StatusCode = -1L;
            GrpcStatusCode = -1;
        }

        public ApiResponseException(string content) : this(-1L, content, -1)
        {
        }

        public override string ToString()
        {
            return $"ApiResponseException(StatusCode={StatusCode}, Message='{Message}', GrpcStatusCode={GrpcStatusCode})";
        }
    }

    {{- range $defname, $definition := .Definitions }}
    {{- $classname := $defname | title }}

    /// <summary>
    /// {{ $definition.Description | stripNewlines }}
    /// </summary>
    public interface I{{ $classname }}
    {
        {{- range $propname, $property := $definition.Properties }}
        {{- $fieldname := $propname | pascalCase }}

        /// <summary>
        /// {{ $property.Description | stripNewlines }}
        /// </summary>
        {{- if eq $property.Type "integer"}}
        int {{ $fieldname }} { get; }
        {{- else if eq $property.Type "number" }}
        double {{ $fieldname }} { get; }
        {{- else if eq $property.Type "boolean" }}
        bool {{ $fieldname }} { get; }
        {{- else if eq $property.Type "string"}}
        string {{ $fieldname }} { get; }
        {{- else if eq $property.Type "array"}}
            {{- if eq $property.Items.Type "string"}}
        List<string> {{ $fieldname }} { get; }
            {{- else if eq $property.Items.Type "integer"}}
        List<int> {{ $fieldname }} { get; }
            {{- else if eq $property.Items.Type "number"}}
        List<double> {{ $fieldname }} { get; }
            {{- else if eq $property.Items.Type "boolean"}}
        List<bool> {{ $fieldname }} { get; }
            {{- else}}
        IEnumerable<I{{ $property.Items.Ref | cleanRef }}> {{ $fieldname }} { get; }
            {{- end }}
        {{- else if eq $property.Type "object"}}
            {{- if eq $property.AdditionalProperties.Type "string"}}
        IDictionary<string, string> {{$fieldname}} { get; }
            {{- else if eq $property.AdditionalProperties.Type "integer"}}
        IDictionary<string, int> {{$fieldname}} { get; }
            {{- else if eq $property.AdditionalProperties.Type "number"}}
        IDictionary<string, double> {{$fieldname}} { get; }
            {{- else if eq $property.AdditionalProperties.Type "boolean"}}
        IDictionary<string, bool> {{$fieldname}} { get; }
            {{- else}}
        IDictionary<string, {{$property.AdditionalProperties | cleanRef}}> {{$fieldname}} { get; }
            {{- end}}
        {{- else }}
        I{{ $property.Ref | cleanRef }} {{ $fieldname }} { get; }
        {{- end }}
        {{- end }}
    }

    /// <inheritdoc />
    internal class {{ $classname }} : I{{ $classname }}
    {
        {{- range $propname, $property := $definition.Properties }}
        {{- $fieldname := $propname | pascalCase }}
        {{- $attrDataName := $propname | snakeCase }}

        /// <inheritdoc />
        {{- if eq $property.Type "integer" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public int {{ $fieldname }} { get; set; }
        {{- else if eq $property.Type "number" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public double {{ $fieldname }} { get; set; }
        {{- else if eq $property.Type "boolean" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public bool {{ $fieldname }} { get; set; }
        {{- else if eq $property.Type "string" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public string {{ $fieldname }} { get; set; }
        {{- else if eq $property.Type "array" }}
            {{- if eq $property.Items.Type "string" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public List<string> {{ $fieldname }} { get; set; }
            {{- else if eq $property.Items.Type "integer" }}
        [DataMember(Name="{{ $propname }}"), Preserve]
        public List<int> {{ $fieldname }} { get; set; }
            {{- else if eq $property.Items.Type "number" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public List<double> {{ $fieldname }} { get; set; }
            {{- else if eq $property.Items.Type "boolean" }}
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public List<bool> {{ $fieldname }} { get; set; }
            {{- else}}
        public IEnumerable<I{{ $property.Items.Ref | cleanRef }}> {{ $fieldname }} => _{{ $propname | camelCase }} ?? new List<{{ $property.Items.Ref | cleanRef }}>(0);
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public List<{{ $property.Items.Ref | cleanRef }}> _{{ $propname | camelCase }} { get; set; }
            {{- end }}
        {{- else if eq $property.Type "object"}}
            {{- if eq $property.AdditionalProperties.Type "string"}}
        public IDictionary<string, string> {{ $fieldname }} => _{{ $propname | camelCase }} ?? new Dictionary<string, string>();
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public Dictionary<string, string> _{{ $propname | camelCase }} { get; set; }
            {{- else if eq $property.Items.Type "integer"}}
        public IDictionary<string, int> {{ $fieldname }} => _{{ $propname | camelCase }} ?? new Dictionary<string, int>();
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
           {{- else if eq $property.Items.Type "number"}}
        public IDictionary<string, double> {{ $fieldname }} => _{{ $propname | camelCase }} ?? new Dictionary<string, double>();
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public Dictionary<string, int> _{{ $propname | camelCase }} { get; set; }
            {{- else if eq $property.Items.Type "boolean"}}
        public IDictionary<string, bool> {{ $fieldname }} => _{{ $propname | camelCase }} ?? new Dictionary<string, bool>();
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public Dictionary<string, bool> _{{ $propname | camelCase }} { get; set; }
            {{- else}}
        public IDictionary<string, {{$property.AdditionalProperties | cleanRef}}> {{ $fieldname }}  => _{{ $propname | camelCase }} ?? new Dictionary<string, {{$property.AdditionalProperties | cleanRef}}>();
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public Dictionary<string, {{$property.AdditionalProperties | cleanRef}}> _{{ $propname | camelCase }} { get; set; }
            {{- end}}
        {{- else }}
        public I{{ $property.Ref | cleanRef }} {{ $fieldname }} => _{{ $propname | camelCase }};
        [DataMember(Name="{{ $attrDataName }}"), Preserve]
        public {{ $property.Ref | cleanRef }} _{{ $propname | camelCase }} { get; set; }
        {{- end }}
        {{- end }}

        public override string ToString()
        {
            var output = "";
            {{- range $fieldname, $property := $definition.Properties }}
            {{- if eq $property.Type "array" }}
            output = string.Concat(output, "{{ $fieldname | pascalCase }}: [", string.Join(", ", {{ $fieldname | pascalCase }}), "], ");
            {{- else if eq $property.Type "object" }}

            var mapString = "";
            foreach (var kvp in {{ $fieldname | pascalCase }})
            {
                mapString = string.Concat(mapString, "{" + kvp.Key + "=" + kvp.Value + "}");
            }
            output = string.Concat(output, "{{ $fieldname | pascalCase }}: [" + mapString + "]");
            {{- else }}
            output = string.Concat(output, "{{ $fieldname | pascalCase }}: ", {{ $fieldname | pascalCase }}, ", ");
            {{- end }}
            {{- end }}
            return output;
        }
    }
    {{- end }}

    /// <summary>
    /// The low level client for the Nakama API.
    /// </summary>
    internal class ApiClient
    {
        public readonly IHttpAdapter HttpAdapter;
        public int Timeout { get; set; }

        private readonly Uri _baseUri;

        public ApiClient(Uri baseUri, IHttpAdapter httpAdapter, int timeout = 10)
        {
            _baseUri = baseUri;
            HttpAdapter = httpAdapter;
            Timeout = timeout;
        }

        {{- range $url, $path := .Paths }}
        {{- range $method, $operation := $path}}

        /// <summary>
        /// {{ $operation.Summary | stripNewlines }}
        /// </summary>
        {{- if $operation.Responses.Ok.Schema.Ref }}
        public async Task<I{{ $operation.Responses.Ok.Schema.Ref | cleanRef }}> {{ $operation.OperationId | stripOperationPrefix | pascalCase }}Async(
        {{- else }}
        public async Task {{ $operation.OperationId | stripOperationPrefix |pascalCase }}Async(
        {{- end}}

        {{- $isPreviousParam := false}}

        {{- if $operation.Security }}
        {{- with (index $operation.Security 0) }}
            {{- range $key, $value := . }}
                {{- if eq $key "BasicAuth" }}
            string basicAuthUsername,
            string basicAuthPassword
            {{- $isPreviousParam = true}}
                {{- else if eq $key "HttpKeyAuth" }}
            {{- $isPreviousParam = true}}
            string bearerToken
                {{- end }}
            {{- end }}
        {{- end }}
        {{- else }}
           {{- $isPreviousParam = true}}
            string bearerToken
        {{- end }}


        {{- range $parameter := $operation.Parameters }}

        {{- if eq $isPreviousParam true}},{{- end}}
        {{- if eq $parameter.In "path" }}
            {{ $parameter.Type }}{{- if not $parameter.Required }}?{{- end }} {{ $parameter.Name }}
        {{- else if eq $parameter.In "body" }}
            {{- if eq $parameter.Schema.Type "string" }}
            string{{- if not $parameter.Required }}?{{- end }} {{ $parameter.Name }}
            {{- else }}
            {{ $parameter.Schema.Ref | cleanRef }}{{- if not $parameter.Required }}?{{- end }} {{ $parameter.Name }}
            {{- end }}
        {{- else if eq $parameter.Type "array"}}
            IEnumerable<{{ $parameter.Items.Type }}> {{ $parameter.Name | camelCase }}
        {{- else if eq $parameter.Type "object"}}
            {{- if eq $parameter.AdditionalProperties.Type "string"}}
        IDictionary<string, string> {{ $parameter.Name }}
            {{- else if eq $parameter.Items.Type "integer"}}
        IDictionary<string, int> {{ $parameter.Name }}
            {{- else if eq $parameter.Items.Type "boolean"}}
        IDictionary<string, int> {{ $parameter.Name }}
            {{- else}}
        IDictionary<string, {{ $parameter.Items.Type }}> {{ $parameter.Name }}
            {{- end}}
        {{- else if eq $parameter.Type "integer" }}
            int? {{ $parameter.Name }}
        {{- else if eq $parameter.Type "boolean" }}
            bool? {{ $parameter.Name }}
        {{- else if eq $parameter.Type "string" }}
            string {{ $parameter.Name }}
        {{- else }}
            {{ $parameter.Type }} {{ $parameter.Name }}
        {{- end }}
        {{- $isPreviousParam = true}}
        {{- end }})
        {
            {{- range $parameter := $operation.Parameters }}
            {{- if $parameter.Required }}
            if ({{ $parameter.Name | camelCase}} == null)
            {
                throw new ArgumentException("'{{ $parameter.Name | camelCase }}' is required but was null.");
            }
            {{- end }}
            {{- end }}

            var urlpath = "{{- $url }}";


            {{- range $parameter := $operation.Parameters }}
            {{- $snakecase := $parameter.Name | snakeCase }}
            {{- if eq $parameter.In "path" }}
            urlpath = urlpath.Replace("{{- print "{" $parameter.Name "}"}}", Uri.EscapeDataString({{- $parameter.Name }}));
            {{- end }}
            {{- end }}

            var queryParams = "";
            {{- range $parameter := $operation.Parameters }}
            {{- $snakecase := $parameter.Name | snakeCase }}
            {{- if eq $parameter.In "query"}}
                {{- if eq $parameter.Type "integer" }}
            if ({{ $parameter.Name }} != null) {
                queryParams = string.Concat(queryParams, "{{- $snakecase }}=", {{ $parameter.Name }}, "&");
            }
                {{- else if eq $parameter.Type "string" }}
            if ({{ $parameter.Name }} != null) {
                queryParams = string.Concat(queryParams, "{{- $snakecase }}=", Uri.EscapeDataString({{ $parameter.Name }}), "&");
            }
                {{- else if eq $parameter.Type "boolean" }}
            if ({{ $parameter.Name }} != null) {
                queryParams = string.Concat(queryParams, "{{- $snakecase }}=", {{ $parameter.Name }}.ToString().ToLower(), "&");
            }
                {{- else if eq $parameter.Type "array" }}
            foreach (var elem in {{ $parameter.Name | camelCase }} ?? new {{ $parameter.Items.Type }}[0])
            {
                queryParams = string.Concat(queryParams, "{{- $snakecase }}=", elem, "&");
            }
                {{- else }}
            {{ $parameter }} // ERROR
                {{- end }}
            {{- end }}
            {{- end }}

            var uri = new UriBuilder(_baseUri)
            {
                Path = urlpath,
                Query = queryParams
            }.Uri;

            var method = "{{- $method | uppercase }}";
            var headers = new Dictionary<string, string>();

            {{- if $operation.Security }}
            {{- with (index $operation.Security 0) }}
                {{- range $key, $value := . }}
                    {{- if eq $key "BasicAuth" }}
            var credentials = Encoding.UTF8.GetBytes(basicAuthUsername + ":" + basicAuthPassword);
            var header = string.Concat("Basic ", Convert.ToBase64String(credentials));
            headers.Add("Authorization", header);

                    {{- else if eq $key "HttpKeyAuth" }}
            if (!string.IsNullOrEmpty(bearerToken))
            {
                var header = string.Concat("Bearer ", bearerToken);
                headers.Add("Authorization", header);
            }
                    {{- end }}
                {{- end }}
            {{- end }}
            {{- else }}
            var header = string.Concat("Bearer ", bearerToken);
            headers.Add("Authorization", header);
            {{- end }}

            byte[] content = null;
            {{- range $parameter := $operation.Parameters }}
            {{- if eq $parameter.In "body" }}
            var jsonBody = {{ $parameter.Name }}.ToJson();
            content = Encoding.UTF8.GetBytes(jsonBody);
            {{- end }}
            {{- end }}

            {{- if $operation.Responses.Ok.Schema.Ref }}
            var contents = await HttpAdapter.SendAsync(method, uri, headers, content, Timeout);
            return contents.FromJson<{{ $operation.Responses.Ok.Schema.Ref | cleanRef }}>();
            {{- else }}
            await HttpAdapter.SendAsync(method, uri, headers, content, Timeout);
            {{- end}}
        }
        {{- end }}
        {{- end }}
    }
}
`

func convertRefToClassName(input string) (className string) {
	cleanRef := strings.TrimPrefix(input, "#/definitions/")
	className = strings.Title(cleanRef)
	return
}

func snakeCaseToCamelCase(input string) (camelCase string) {
	isToUpper := false
	for k, v := range input {
		if k == 0 {
			camelCase = strings.ToLower(string(input[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}

	}
	return
}

func snakeCaseToPascalCase(input string) (output string) {
	isToUpper := false
	for k, v := range input {
		if k == 0 {
			output = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				output += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					output += string(v)
				}
			}
		}
	}
	return
}

func isSnakeCase(input string) (output bool) {

	output = true

	for _, v := range input {
		vString := string(v)
		if vString != "_" && strings.ToUpper(vString) == vString {
			output = false
		}
	}

	return
}

func camelCaseToSnakeCase(input string) (output string) {
	output = ""

	if isSnakeCase(input) {
		output = input
		return
	}

	for _, v := range input {
		vString := string(v)
		if vString == strings.ToUpper(vString) {
			output += "_" + strings.ToLower(vString)
		} else {
			output += vString
		}
	}

	return
}

func stripNewlines(input string) (output string) {
	output = strings.Replace(input, "\n", " ", -1)
	return
}

func stripOperationPrefix(input string) string {
	return strings.Replace(input, "Nakama_", "", 1)
}

func main() {
	// Argument flags
	var output = flag.String("output", "", "The output for generated code.")
	flag.Parse()

	inputs := flag.Args()
	if len(inputs) < 1 {
		fmt.Printf("No input file found: %s\n\n", inputs)
		fmt.Println("openapi-gen [flags] inputs...")
		flag.PrintDefaults()
		return
	}

	inputFile := inputs[0]
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", err)
		return
	}

	var subnamespace (string) = ""

	if len(inputs) > 1 {
		if len(inputs[1]) <= 0 {
			fmt.Println("Empty Sub-Namespace provided.")
			return
		}

		subnamespace = inputs[1]
	}

	var schema struct {
		SubNamespace string
		Paths        map[string]map[string]struct {
			Summary     string
			OperationId string
			Responses   struct {
				Ok struct {
					Schema struct {
						Ref string `json:"$ref"`
					}
				} `json:"200"`
			}
			Parameters []struct {
				Name     string
				In       string
				Required bool
				Type     string   // used with primitives
				Items    struct { // used with type "array"
					Type string
				}
				Schema struct { // used with http body
					Type string
					Ref  string `json:"$ref"`
				}
				Format string // used with type "boolean"
			}
			Security []map[string][]struct {
			}
		}
		Definitions map[string]struct {
			Properties map[string]struct {
				Type  string
				Ref   string   `json:"$ref"` // used with object
				Items struct { // used with type "array"
					Type string
					Ref  string `json:"$ref"`
				}
				AdditionalProperties struct {
					Type string // used with type "map"
				}
				Format      string // used with type "boolean"
				Description string
			}
			Description string
		}
	}

	schema.SubNamespace = subnamespace

	if err := json.Unmarshal(content, &schema); err != nil {
		fmt.Printf("Unable to decode input file %s : %s\n", inputFile, err)
		return
	}

	fmap := template.FuncMap{
		"camelCase":            snakeCaseToCamelCase,
		"cleanRef":             convertRefToClassName,
		"pascalCase":           snakeCaseToPascalCase,
		"stripNewlines":        stripNewlines,
		"title":                strings.Title,
		"uppercase":            strings.ToUpper,
		"snakeCase":            camelCaseToSnakeCase,
		"stripOperationPrefix": stripOperationPrefix,
	}

	tmpl, err := template.New(inputFile).Funcs(fmap).Parse(codeTemplate)
	if err != nil {
		fmt.Printf("Template parse error: %s\n", err)
		return
	}

	if len(*output) < 1 {
		tmpl.Execute(os.Stdout, schema)
		return
	}

	f, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Unable to create file: %s\n", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	tmpl.Execute(writer, schema)
	writer.Flush()
}
