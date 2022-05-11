export POD=$3
export NS=$2
export PROFILER=$1
kubectl exec -n "$NS" "$POD" -c istio-proxy -- curl -X POST -s "http://localhost:15000/${PROFILER}profiler?enable=y"
sleep 15
kubectl exec -n "$NS" "$POD" -c istio-proxy -- curl -X POST -s "http://localhost:15000/${PROFILER}profiler?enable=n"
rm -rf ./tmp/envoy
kubectl cp -n "$NS" "$POD":/var/lib/istio/data ./tmp/envoy -c istio-proxy
kubectl cp -n "$NS" "$POD":/lib/x86_64-linux-gnu ./tmp/envoy/lib -c istio-proxy
kubectl cp -n "$NS" "$POD":/usr/local/bin/envoy ./tmp/envoy/lib/envoy -c istio-proxy
