Install docker: https://docs.docker.com/engine/install/
For instance (ubuntu OS):

# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

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



