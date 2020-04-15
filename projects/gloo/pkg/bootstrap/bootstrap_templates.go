package bootstrap

import (
	"bytes"
	"os/exec"
	"text/template"

	"github.com/rotisserie/eris"
)

type EnvoyInstance struct {
	Transformations string
}

func (ei *EnvoyInstance) ValidateBootstrap(bootstrapTemplate string) error {
	configYaml := ei.buildBootstrap(bootstrapTemplate)
	validateCmd := exec.Command("/usr/local/bin/envoy", "--mode", "validate", "--config-yaml", configYaml)
	output := bytes.NewBuffer([]byte{})
	validateCmd.Stdout = output
	validateCmd.Stderr = output
	if err := validateCmd.Run(); err != nil {
		return eris.Errorf("%v", string(output.Bytes()), err)
	}
	return nil
}

func (ei *EnvoyInstance) buildBootstrap(bootstrapTemplate string) string {
	var b bytes.Buffer
	parsedTemplate := template.Must(template.New("bootstrap").Parse(bootstrapTemplate))
	if err := parsedTemplate.Execute(&b, ei); err != nil {
		panic(err)
	}
	return b.String()
}

const TransformationBootstrapTemplate = `
node:
  cluster: doesntmatter
  id: imspecial
  metadata:
    role: "gloo-system~gateway-proxy"
static_resources:
  clusters:
  - name: placeholder_cluster
    connect_timeout: 5.000s
  listeners:
  - address:
      socket_address:
        address: 0.0.0.0
        port_value: 8081
    filter_chains:
    - filters:
      - config:
          route_config:
            name: placeholder_route
            virtual_hosts:
            - domains:
              - '*'
              name: placeholder_host
              routes:
              - match:
                  headers:
                  - exact_match: GET
                    name: :method
                  path: /
                route:
                  cluster: placeholder_cluster
              per_filter_config:
                io.solo.transformation:
                  {{.Transformations}}
          stat_prefix: placeholder
        name: envoy.http_connection_manager
    name: placeholder_listener
`
