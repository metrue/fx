FROM python:3

RUN pip install flask
ENV FLASK_APP=app.py

COPY . .
EXPOSE 3000
CMD flask run -h 0.0.0.0 -p 3000
