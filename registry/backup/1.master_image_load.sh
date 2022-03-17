image_array=(
"docker.io/istio/proxyv2:1.9.4"
"docker.io/istio/pilot:1.9.4"
"quay.io/kiali/kiali:v1.29"
"ketidevit2/istio-crosscluster-workaround-for-eks:v0.0.1"
"prom/prometheus:v2.21.0"
"jimmidyson/configmap-reload:v0.4.0"
"metallb/controller:v0.8.1"
"metallb/speaker:v0.8.1"
"docker.io/influxdb:1.6.4"
"ketidevit2/openmcp-analytic-engine:v0.0.1"
"ketidevit2/openmcp-apiserver:v0.0.1"
"ketidevit2/openmcp-cluster-manager:v0.0.1"
"ketidevit2/openmcp-configmap-controller:v0.0.1"
"ketidevit2/openmcp-daemonset-controller:v0.0.1"
"ketidevit2/openmcp-deployment-controller:v0.0.1"
"ketidevit2/openmcp-dns-controller:v0.0.1"
"ketidevit2/openmcp-has-controller:v0.0.1"
"ketidevit2/openmcp-ingress-controller:v0.0.1"
"ketidevit2/openmcp-job-controller:v0.0.1"
"ketidevit2/openmcp-loadbalancing-controller:v0.0.1"
"ketidevit2/openmcp-metric-collector:v0.0.1"
"ketidevit2/openmcp-namespace-controller:v0.0.1"
"ketidevit2/openmcp-policy-engine:v0.0.1"
"lkh1434/openmcp-portal:v1.0"
"lkh1434/openmcp-portal-apiserver:v0.0.2"
"lkh1434/openmcp-portal-lstm:v0.0.1"
"postgres:latest"
"ketidevit2/openmcp-pv-controller:v0.0.1"
"ketidevit2/openmcp-pvc-controller:v0.0.1"
"ketidevit2/openmcp-scheduler:v0.0.1"
"ketidevit2/openmcp-secret-controller:v0.0.1"
"ketidevit2/openmcp-service-controller:v0.0.1"
"ketidevit2/openmcp-statefulset-controller:v0.0.1"
"ketidevit2/openmcp-sync-controller:v0.0.1"
)
function pull_and_load() {
  echo "'${1}' image '${2}' pull_and_load start"
  docker pull ${2}
  kind load docker-image ${2} --name ${1}
  echo "--> '${1}'$image '${2}' pull_and_load end"

}
for image_name in "${image_array[@]}"; do
  pull_and_load $1 $image_name &
done

wait
echo "Finished"
