FROM r-base

COPY . /usr/local/src/fx
WORKDIR /usr/local/src/fx
RUN R -f packages.R

EXPOSE 3000

CMD ["Rscript", "app.R"]
