FROM myobplatform/shell-operator:latest as shop

FROM python:3.6
COPY --from=shop /shell-operator /shell-operator
RUN pip install ansible boto boto3
COPY shell-conf.yaml /app/shell-config.yaml
ENTRYPOINT ["/shell-operator", "-v=1"]
