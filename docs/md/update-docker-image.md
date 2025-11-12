
# Maintainance and updates of the docker image

## 1. Stop the docker image

`docker-compose down`

## 2. Docker image update

Pull the latest image:
`docker-compose pull`

## 3. Optional remove docker images

List all docker images so far pulled:
`docker images -a`

Remove pulled images with the image id:
`docker imr [image_ID]`

- Probably the best idea is to keep allways two versions. The latest one you pulled and the version   
which you have been running on. That way you can allways easily revert to the last working version.
You can differentiate by the timestamp you pulled your images

## 4. Restart the docker image

Start the latest image DETACHED
`docker-compose up -d`
