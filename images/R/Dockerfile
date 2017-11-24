FROM r-base

COPY . /usr/local/src/fx
WORKDIR /usr/local/src/fx
RUN Rscript packages.R

CMD ["Rscript", "app.R"]
EXPOSE 3000
