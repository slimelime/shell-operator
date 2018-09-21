FROM myobplatform/shell-operator:0.1.0-alpha2 as shop

FROM python:3.6
COPY --from=shop /shell-operator /shell-operator
RUN pip install ansible boto boto3
ENTRYPOINT ["/shell-operator"]
