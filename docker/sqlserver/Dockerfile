# TODO should be mcr.microsoft.com/mssql/server:2022-latest when https://github.com/microsoft/mssql-docker/issues/847 is fixed
FROM mcr.microsoft.com/mssql/server:2022-CU5-ubuntu-20.04

USER mssql

# Create a config directory
RUN mkdir -p /usr/config
WORKDIR /usr/config

# Bundle config source
COPY entrypoint.sh /usr/config/
COPY configure-db.sh /usr/config/
COPY setup.sql /usr/config/

ENTRYPOINT ["./entrypoint.sh"]

# Tail the setup logs to trap the process
CMD ["tail -f /dev/null"]

HEALTHCHECK --interval=15s CMD /opt/mssql-tools/bin/sqlcmd -U sa -P $MSSQL_SA_PASSWORD -Q "select 1" && grep -q "MSSQL CONFIG COMPLETE" ./config.log