FROM python:3.6 as test
RUN pip install kubernetes pyyaml voluptuous
WORKDIR /app
COPY ./shelloperator /usr/bin/
ENTRYPOINT ["/usr/bin/shelloperator"]
