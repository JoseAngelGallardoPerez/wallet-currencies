print("Wallet Currencies")

load("ext://restart_process", "docker_build_with_restart")

cfg = read_yaml(
    "tilt.yaml",
    default = read_yaml("tilt.yaml.sample")
)

local_resource(
    "currencies-build-binary",
    "make fast_build",
    deps = ["./cmd", "./internal"],
)
local_resource(
    "currencies-generate-protpbuf",
    "make gen-protobuf",
    deps = ["./rpc/currencies/currencies.proto"],
)

docker_build(
    "velmie/wallet-currencies-db-migration",
    ".",
    dockerfile = "Dockerfile.migrations",
    only = "migrations",
)
k8s_resource(
    "wallet-currencies-db-migration",
    trigger_mode = TRIGGER_MODE_MANUAL,
    resource_deps = ["wallet-currencies-db-init"],
)

wallet_currencies_options = dict(
    entrypoint = "/app/service_currencies",
    dockerfile = "Dockerfile.prebuild",
    port_forwards = [],
    helm_set = [],
)

if cfg['debug']:
    wallet_currencies_options["entrypoint"] = "$GOPATH/bin/dlv --continue --listen :%s --accept-multiclient --api-version=2 --headless=true exec /app/service_currencies" % cfg["debug_port"]
    wallet_currencies_options["dockerfile"] = "Dockerfile.debug"
    wallet_currencies_options["port_forwards"] = cfg["debug_port"]
    wallet_currencies_options["helm_set"] = ["containerLivenessProbe.enabled=false", "containerPorts[0].containerPort=%s" % cfg["debug_port"]]

docker_build_with_restart(
    "velmie/wallet-currencies",
    ".",
    dockerfile = wallet_currencies_options["dockerfile"],
    entrypoint = wallet_currencies_options["entrypoint"],
    only = [
        "./build",
        "zoneinfo.zip",
    ],
    live_update = [
        sync("./build", "/app/"),
    ],
)
k8s_resource(
    "wallet-currencies",
    resource_deps = ["wallet-currencies-db-migration"],
    port_forwards = wallet_currencies_options["port_forwards"],
)

yaml = helm(
    "./helm/wallet-currencies",
    # The release name, equivalent to helm --name
    name = "wallet-currencies",
    # The values file to substitute into the chart.
    values = ["./helm/values-dev.yaml"],
    set = wallet_currencies_options["helm_set"],
)

k8s_yaml(yaml)
