package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"

	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"

	"github.com/layer5io/meshery-linkerd/internal/config"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var CRDnamesURL map[string]string
var LatestVersion string
var WorkloadPath string
var MeshModelPath string
var AllVersions []string

const Component = "Linkerd"

var Meshmodelmetadata = make(map[string]interface{})

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Cloud Native Network",
	Metadata: Meshmodelmetadata,
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_LINKERD)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			manifests.RemoveHelmTemplatingFromCRD(&manifest)
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}
func init() {
	// Initialize Metadata including logo svgs
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &Meshmodelmetadata)
	wd, _ := os.Getwd()
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("linkerd", "linkerd2")
	if len(AllVersions) == 0 {
		return
	}
	LatestVersion = AllVersions[len(AllVersions)-1]
	DefaultGenerationMethod = adapter.Manifests
	names, err := config.GetFileNames("linkerd", "linkerd2", "charts/linkerd-crds/templates/**")
	if err != nil {
		fmt.Println("dynamic component generation failure: ", err.Error())
		return
	}
	for n := range names {
		if !strings.HasSuffix(n, ".yaml") {
			delete(names, n)
		}
	}
	CRDnamesURL = names
}
