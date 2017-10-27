FROM python:3

RUN pip install flask
ENV FLASK_APP=app.py

COPY . .
CMD flask run -h 0.0.0.0 -p 3000
