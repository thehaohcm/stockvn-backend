Install docker: https://docs.docker.com/engine/install/

Install package:

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

Install k3s:

curl -sfL https://get.k3s.io | sh -

Check status:

sudo k3s kubectl get nodes

Config k3s to run on local:

export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

kubectl get nodes

Install Helm:
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

Install pogres via Helm:

helm repo add bitnami https://charts.bitnami.com/bitnami

helm repo update

helm install stockvn-db bitnami/postgresql -f stockvn-db-service.yaml

helm install my-postgres bitnami/postgresql \
  --set service.type=NodePort \
  --set service.nodePort=5432 \
  --set auth.postgresPassword=mysecurepassword \
  --set primary.service.ports.postgresql=5432


