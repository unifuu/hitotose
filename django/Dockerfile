# Use the official Python image from the Docker Hub
FROM python:3.10.6

# Set the working directory inside the container
WORKDIR /app

# Copy requirements.txt
COPY requirements.txt ./

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY . .

# Expose port 8000 to the host
EXPOSE 8080

# Run the Django development server
# CMD ["python", "manage.py", "runserver", "0.0.0.0:8080"]

# Run the application using Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:8080", "hitotose.wsgi:application"]