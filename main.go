package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type listOfErrorCodes []string

type errorFileData struct {
	Messages map[string]errorMessage `json:"messages"`
	Types    map[string]string       `json:"types,omitempty"`
	Default  string                  `json:"default"`
}

type errorMessage struct {
	Text  string            `json:"text"`
	Types map[string]string `json:"types,omitempty"`
}

type errorAttribute struct {
	Name      string
	Type      string
	Formatter string
}

type errorTemplateItems []errorTemplateItem

func (s errorTemplateItems) Len() int      { return len(s) }
func (s errorTemplateItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s errorTemplateItems) Less(i, j int) bool {
	return s[i].ErrorCode < s[j].ErrorCode
}

type errorTemplateItem struct {
	ErrorCode      string
	Text           string
	UpperErrorCode string
	Attributes     []errorAttribute
}

func main() {
	type EnvVariables struct {
		YamlFilePath   string
		TargetDir      string
		TargetFilename string
		PackageName    string
	}

	envVars := EnvVariables{
		YamlFilePath:   os.Getenv("ERRORS_YAML_FILE_PATH"),
		TargetDir:      os.Getenv("ERRORS_TARGET_DIR"),
		TargetFilename: os.Getenv("ERRORS_TARGET_FILENAME"),
		PackageName:    os.Getenv("ERRORS_PACKAGE_NAME"),
	}

	targetFullFilename := filepath.Join(envVars.TargetDir, envVars.TargetFilename)

	errorFileDataRaw, err := os.ReadFile(envVars.YamlFilePath)
	if err != nil {
		log.Fatalf("unable to read yaml file. Err: %s", err)
	}

	//nolint:exhaustruct
	errorData := &errorFileData{}

	err = yaml.Unmarshal(errorFileDataRaw, &errorData)
	if err != nil {
		log.Fatalf("unable to unmarshal yaml file. Err: %s", err)
	}

	errorCodes := listOfErrorCodes{}
	for errorCode := range errorData.Messages {
		errorCodes = append(errorCodes, errorCode)
	}

	sort.Strings(errorCodes)

	rawErrorDataJSON, err := json.MarshalIndent(errorData, "", "\t")
	if err != nil {
		log.Fatalf("unable to marshal json data. Err: %s", err)
	}

	items := make(errorTemplateItems, 0)

	for _, errorCode := range errorCodes {
		text := errorData.Messages[errorCode].Text
		globalTypes := errorData.Types
		currentTypes := errorData.Messages[errorCode].Types
		defaultType := errorData.Default

		item := errorTemplateItem{
			ErrorCode:      errorCode,
			Text:           text,
			UpperErrorCode: strings.ToUpper(errorCode[:1]) + errorCode[1:],
			Attributes:     getErrorAttributes(text, currentTypes, globalTypes, defaultType),
		}

		items = append(items, item)
	}

	sort.Sort(items)

	err = os.Remove(targetFullFilename)
	if err != nil {
		log.Printf("unable to remove file %s. Err: %s", targetFullFilename, err)
	}

	//nolint:mnd
	err = os.MkdirAll(envVars.TargetDir, 0o700)
	if err != nil {
		log.Fatalf("unable to create a directory '%s' with permission 0700. Err: %s", envVars.TargetDir, err)
	}

	file, err := os.Create(targetFullFilename)
	if err != nil {
		log.Fatalf("unable to create a file '%s'. Err: %s", targetFullFilename, err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("unable to close the file '%s'. Err: %s", targetFullFilename, err)
		}

		//err = exec.Command("goimports", "-w", targetFullFilename).Start()
		//if err != nil {
		//	log.Printf("unable format file '%s'. Err: %s", targetFullFilename, err)
		//}
	}()

	errorTpl := `package {{ .PackageName }}

import (
	"strings"
)

const (
{{- range .Elements }}
	{{ .UpperErrorCode }} = "{{ .ErrorCode }}"
{{- end }}
)

type DomainError struct {
	Text       string         				` + "`json:\"text\"`" + `
	Code   	   string         				` + "`json:\"code\"`" + `
	Attributes DomainErrorAttributes		` + "`json:\"attributes\"`" + `
	Error      string						` + "`json:\"error\"`" + `
}
	
func (de DomainError) String() string {
	text := de.Text
	for key, value := range de.Attributes {
		text = strings.ReplaceAll(text, fmt.Sprintf("{{ "{{%s}}" }}", key), value)
		text = strings.ReplaceAll(text, fmt.Sprintf("{{ "{{ %s }}" }}", key), value)
	}
	if de.Error != "" {
		text += ": " + de.Error
	}
	
	return text
}

type DomainErrorAttributes map[string]string

func toErrStr(errs []error) string {
	str := ""
	
	for i, err := range errs {
		if err != nil {
			if i > 0 {
				str += ": "
			}
			str += err.Error()
		}
	}

	return str
}

{{- range .Elements }}

// New{{ .UpperErrorCode }} generated from the code "{{ .ErrorCode }}".
func New{{ .UpperErrorCode }}({{ range .Attributes }}{{ .Name }} {{ .Type }}, {{ end }} errs ...error) *DomainError {
	if len(errs) > 0 && errs[0] == nil {
		return nil
	}

	return &DomainError{
		Text: "{{ .Text }}",
		Code: {{ .UpperErrorCode }},
		Attributes: DomainErrorAttributes{
{{- range .Attributes }}
	"{{ .Name }}": fmt.Sprintf("{{ .Formatter }}", {{ .Name }}),
{{- end }}
		},
		Error: toErrStr(errs),
	}
}
{{- end }}

type DomainErrorDeclaration struct {
    ErrorCode   string
    Text       string
    Attributes []string
}

func GetDomainErrors() []DomainErrorDeclaration {
    return []DomainErrorDeclaration{
    {{- range .Elements }}
        {
            ErrorCode:   "{{ .ErrorCode }}",
            Text:       "{{ .Text }}",
            Attributes: []string{ {{- range .Attributes }}"{{.Name}}", {{- end }} },
        },
    {{- end }}
    }
}
`

	errorPackageTemplate := template.Must(template.New("").Parse(errorTpl))

	err = errorPackageTemplate.Execute(
		file,
		struct {
			PackageName string
			Elements    []errorTemplateItem
			RawFileBody template.HTML
		}{
			PackageName: envVars.PackageName,
			Elements:    items,
			RawFileBody: template.HTML(rawErrorDataJSON), //nolint:gosec
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func getErrorAttributes(
	value string,
	currentTypes map[string]string,
	knownTypes map[string]string,
	defaultType string,
) []errorAttribute {
	result := make([]errorAttribute, 0)
	beginnings := strings.Split(value, "{{")

	for _, item := range beginnings {
		values := strings.Split(item, "}}")

		if len(values) > 1 {
			name := strings.TrimSpace(values[0])
			goType := getType(name, currentTypes, knownTypes, defaultType)
			formatter := getFormatter(goType)

			result = append(result, errorAttribute{
				Name:      name,
				Type:      goType,
				Formatter: formatter,
			})
		}
	}

	return result
}

func getType(
	name string,
	currentTypes map[string]string,
	knownTypes map[string]string,
	defaultType string,
) string {
	goType, ok := currentTypes[name]
	if ok {
		return goType
	}

	goType, ok = knownTypes[name]
	if ok {
		return goType
	}

	return defaultType
}

func getFormatter(str string) string {
	switch str {
	case "string", "error":
		return "%s"
	case "int", "int64", "int32":
		return "%d"
	default:
		return "%v"
	}
}
