image_array=(
"ketidevit2/cluster-metric-collector:v0.0.1"
)
function pull_and_load() {
  echo "'${1}' image '${2}' pull_and_load start"
  docker pull ${2}
  kind load docker-image ${2} --name $1
  echo "--> '${1}' $image '${2}' pull_and_load end"

}
for image_name in "${image_array[@]}"; do
  pull_and_load $1 $image_name &
done

wait
echo "Finished"
