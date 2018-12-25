1. Prepare your Dockerfile within this [directory](https://github.com/metrue/fx/tree/master/assets/dockerfiles/fx). you can simply test it:

```
docker build -t <foo-bar> .
docker run -p 3000:3000 foo-bar
```
if everything works as you expected, fx will support it without any extra effort.

2.  then you have to update assetsMap in the  [file](https://github.com/metrue/fx/blob/05a0372672b813f9afa855b0c849692fa2fe0ee5/image/image.go#L23), and add your language to support.

3. you can test your work with
```
make build
./build/fx up <your-function-source>
```
