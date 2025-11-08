# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### CONFIG ###
k8s_yaml('./infra/development/k8s/app-configs.yaml')
k8s_yaml('./infra/development/k8s/app-secrets.yaml')
### END OF CONFIG ###

### POSTGRES SERVICE ###
k8s_yaml('./infra/development/k8s/pg-deployment.yaml')
k8s_resource(
    'postgres',
    port_forwards=["5432"],
    labels=["databases"]
)
### END OF POSTGRES ###

### API ###
api_compile_cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/boocol ./cmd/api'
cli_compile_cmd='CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/bin/boocol-cli ./cmd/cli'
local_resource(
    'api-compile',
    api_compile_cmd,
    deps=['./internal', './cmd', './pkg'],
    labels="compiles"
)

local_resource(
    'cli-compile',
    cli_compile_cmd,
    deps=['./internal', './cmd', './pkg'],
    labels="compiles"
)

docker_build_with_restart(
    'boocol/api',
    '.',
    entrypoint=['/app/build/bin/boocol'],
    dockerfile='./infra/development/docker/Dockerfile',
    only=[
    './build/bin/boocol',
    './build/bin/boocol-cli',
    './pkg',
    ],
    live_update=[
    sync('./build', '/app/build'),
    sync('./pkg', '/app/pkg'),
    ],
)

k8s_yaml('./infra/development/k8s/api-deployment.yaml')
k8s_resource(
    'boocol',
    port_forwards=8080,
    resource_deps=['api-compile', 'cli-compile', 'postgres'],
    labels="services"
)
### END OF API ###