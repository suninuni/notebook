FROM kong:2.4
USER root

COPY plugins /home/kong/plugins
RUN cd /home/kong/plugins && luarocks make

USER kong
