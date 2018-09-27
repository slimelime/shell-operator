FROM myobplatform/shell-operator:latest as shop

FROM python:3.6
RUN pip install boto boto3 openshift
# need ansible 2.7 to provide k8s_facts
RUN pip install git+https://github.com/ansible/ansible.git@stable-2.7
COPY --from=shop /shell-operator /shell-operator
env KUBE_VERSION="v1.11.3"
run curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl -o /usr/bin/kubectl \
 && chmod +x /usr/bin/kubectl
ENTRYPOINT ["/shell-operator", "-v=1"]
