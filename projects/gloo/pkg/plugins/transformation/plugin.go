package transformation

import (
	"strings"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/proto"
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/go-utils/protoutils"
	"sigs.k8s.io/yaml"
)

const (
	FilterName = "io.solo.transformation"
)

var pluginStage = plugins.AfterStage(plugins.AuthZStage)

type Plugin struct {
	RequireTransformationFilter bool
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.RequireTransformationFilter = false
	return nil
}

// TODO(yuval-k): We need to figure out what\if to do in edge cases where there is cluster weight transform
func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	transformations := in.GetOptions().GetTransformations()
	if transformations == nil {
		return nil
	}

	yml, err := toYaml(transformations)
	if err != nil {
		// should never happen
		return eris.Errorf("Unable to convert transformation to yaml, error: %v", err)
	}

	lines := strings.Split(string(yml), "\n")
	indentedYaml := strings.Join(lines, "\n                  ") // needs to match indentation of fixture template boostrap yaml!

	envoyInstance := bootstrap.EnvoyInstance{Transformations: indentedYaml}
	err = envoyInstance.ValidateBootstrap(bootstrap.TransformationBootstrapTemplate)
	if err != nil {
		return err
	}

	p.RequireTransformationFilter = true
	return pluginutils.SetVhostPerFilterConfig(out, FilterName, transformations)
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	transformations := in.GetOptions().GetTransformations()
	if transformations == nil {
		return nil
	}

	p.RequireTransformationFilter = true
	return pluginutils.SetRoutePerFilterConfig(out, FilterName, transformations)
}

func (p *Plugin) ProcessWeightedDestination(_ plugins.RouteParams, in *v1.WeightedDestination, out *envoyroute.WeightedCluster_ClusterWeight) error {
	transformations := in.GetOptions().GetTransformations()
	if transformations == nil {
		return nil
	}
	p.RequireTransformationFilter = true
	return pluginutils.SetWeightedClusterPerFilterConfig(out, FilterName, transformations)
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	return []plugins.StagedHttpFilter{
		plugins.NewStagedFilter(FilterName, pluginStage),
	}, nil
}

func toYaml(m proto.Message) ([]byte, error) {
	jsn, err := protoutils.MarshalBytes(m)
	if err != nil {
		return nil, err
	}
	return yaml.JSONToYAML(jsn)
}
