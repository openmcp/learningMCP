image_array=(
"ketidevit2/openmcp-analytic-engine:v0.0.1"
"ketidevit2/openmcp-apiserver:v0.0.1"
"ketidevit2/openmcp-configmap-controller:v0.0.1"
"ketidevit2/openmcp-deployment-controller:v0.0.1"
"ketidevit2/openmcp-dns-controller:v0.0.1"
"ketidevit2/openmcp-has-controller:v0.0.1"
"ketidevit2/openmcp-ingress-controller:v0.0.1"
"ketidevit2/openmcp-loadbalancing-controller:v0.0.1"
"ketidevit2/openmcp-metric-collector:v0.0.1"
"ketidevit2/openmcp-policy-engine:v0.0.1"
"ketidevit2/openmcp-scheduler:v0.0.1"
"ketidevit2/openmcp-secret-controller:v0.0.1"
"ketidevit2/openmcp-service-controller:v0.0.1"
"ketidevit2/openmcp-sync-controller:v0.0.1"
"ketidevit2/openmcp-cluster-manager:v0.0.1"
"ketidevit2/cluster-metric-collector:v0.0.1"

)

ketidevit2_image_array=(
#"docker.io/istio/proxyv2:1.9.4"
#"docker.io/istio/pilot:1.9.4"
#"quay.io/kiali/kiali:v1.29"
"istio-crosscluster-workaround-for-eks:v0.0.1"
#"prom/prometheus:v2.21.0"
#"jimmidyson/configmap-reload:v0.4.0"
#"metallb/controller:v0.8.1"
#"metallb/speaker:v0.8.1"
#"docker.io/influxdb:1.6.4"
"openmcp-analytic-engine:v0.0.1"
"openmcp-apiserver:v0.0.1"
"openmcp-cluster-manager:v0.0.1"
"openmcp-configmap-controller:v0.0.1"
"openmcp-daemonset-controller:v0.0.1"
"openmcp-deployment-controller:v0.0.1"
"openmcp-dns-controller:v0.0.1"
"openmcp-has-controller:v0.0.1"
"openmcp-ingress-controller:v0.0.1"
"openmcp-job-controller:v0.0.1"
"openmcp-loadbalancing-controller:v0.0.1"
"openmcp-metric-collector:v0.0.1"
"openmcp-namespace-controller:v0.0.1"
"openmcp-policy-engine:v0.0.1"
#"lkh1434/openmcp-portal:v1.0"
#"lkh1434/openmcp-portal-apiserver:v0.0.2"
#"lkh1434/openmcp-portal-lstm:v0.0.1"
#"postgres:latest"
"openmcp-pv-controller:v0.0.1"
"openmcp-pvc-controller:v0.0.1"
"openmcp-scheduler:v0.0.1"
"openmcp-secret-controller:v0.0.1"
"openmcp-service-controller:v0.0.1"
"openmcp-statefulset-controller:v0.0.1"
"openmcp-sync-controller:v0.0.1"
"cluster-metric-collector:v0.0.1"
"istio-crosscluster-workaround-for-eks:v0.0.1"
)
istio_image_array=(
"proxyv2:1.9.4"
"pilot:1.9.4"
)
function pull_and_push() {
  echo "image '${1}/${2}' pull_and_push start"
  docker pull $1/$2
  docker tag $1/$2 127.0.0.1:5000/openmcp/$2
  docker push 127.0.0.1:5000/openmcp/$2
  echo "--> $image '${1}/${2}' pull_and_push end"

}
for image_name in "${ketidevit2_image_array[@]}"; do
  pull_and_push "ketidevit2" $image_name &
  #docker pull $image_name
  #docker tag $image_name 127.0.0.1:5000/openmcp/$image_name
  #docker push 127.0.0.1:5000/openmcp/$image_name
done
for image_name in "${istio_image_array[@]}"; do
  pull_and_push "istio" $image_name &
  #docker pull $image_name
  #docker tag $image_name 127.0.0.1:5000/openmcp/$image_name
  #docker push 127.0.0.1:5000/openmcp/$image_name
done
wait
echo "Finished"
