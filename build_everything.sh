#/bin/bash

kind create cluster --config kind-config.yml
kubectl apply -f service-account.yml

echo "Building images..."
cd "http_red container"
docker build . -t http_red:v0.0.1 &

cd ../controller
docker build . -t controller:v0.0.1 &

cd ../resources_monitor
docker build . -t resources_monitor:v0.0.1 &

cd ../tests/example-workload
docker build . -t example_workload:v0.0.1 &

cd ../load-generator
docker build . -t load_generator:v0.0.1 &

for job in `jobs -p`
do
    wait $job
done

cd ../..

echo "Loading images into kind cluster..."
for img in http_red:v0.0.1 controller:v0.0.1 example_workload:v0.0.1 load_generator:v0.0.1 resources_monitor:v0.0.1
do
    kind load docker-image $img &
done

for job in `jobs -p`
do
    wait $job
done

echo "Appling controller manifest..."
kubectl apply -f controller/controller.yml

echo "Installing metrics server"
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
