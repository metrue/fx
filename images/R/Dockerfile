FROM r-base

RUN echo 'install.packages(c("jug", "jsonlite"), repos="http://cran.us.r-project.org", dependencies=TRUE)' > /tmp/packages.R \
    && Rscript /tmp/packages.R

COPY . /usr/local/src/fx
WORKDIR /usr/local/src/fx

CMD ["Rscript", "app.R"]
EXPOSE 3000
