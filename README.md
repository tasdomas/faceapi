faceapi
=======

A facedetector API server using openCV.

The server can be started using:
    go run server/server.go -haar=data/haarcascade_frontalface_alt.xml

To test it, simply perform:
    curl -i -F image=@data/lena.jpg localhost:8080

