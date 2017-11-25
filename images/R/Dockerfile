FROM r-base

COPY . /usr/local/src/fx
WORKDIR /usr/local/src/fx
RUN R -f packages.R --save

CMD ["Rscript", "app.R"]
EXPOSE 3000
