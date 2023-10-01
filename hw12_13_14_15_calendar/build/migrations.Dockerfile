FROM gomicro/goose

ARG DBUSER="postgres"
ARG DBPASS="postgres"

ENV DBUSER=$DBUSER
ENV DBPASS=$DBPASS

COPY migrations/*.sql /migrations/
COPY build/migrations.sh /migrations/

CMD ["/migrations/migrations.sh"]