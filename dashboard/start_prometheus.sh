echo $PWD
docker run -p 9090:9090 \
  -v ${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml \
  -v ${PWD}/data:/data/prometheus \
  prom/prometheus
