# onedrive-cameraroll-renamer-service

Renames all the files in your OneDrive Camera Roll folder into the format `YYYYMMDD_HHmmss`.

*DISCLAIMER*: This is a tool I built for myself. It is not super extensive or configurable, 
it just gets the thing done for me. See [TODOs](#todos) for open things to improve. 
Feel free to contribute by creating issues or PRs here.

## Background

Phones store fotos and videos with different filenames.
Some store them like `YYYYMMDD_HHmmss.jpg`, some do `IMG_YYYYMMDD_HHmmss.jpg`
and others put the unix timestamp into the name (`IMG_1558292380001_12345.jpg`).

The OneDrive app automatically uploads the fotos and videos to your OneDrive
but keeps the filenames as-is. If the files in the folder do not have a consistent
naming scheme, looking through the files is annoying because you cannot sort them
chronologically.

This website connects to your OneDrive account (as soon as you authorize it) and
renames all files from the CameraRoll folder and moves them to the folder specified by `CAMERA_ROLL_TARGET_FOLDER_ID`
so they are not processed twice.

## Running it yourself

First, register an app as described here https://learn.microsoft.com/en-us/graph/auth-register-app-v2, note down the client ID and client secret.
Then `cp .env.example .env` and set the variables accordingly.

### Running locally

```
env $(cat .env | grep -v "#" | xargs) go run . 
```


## TODOs

* The list of filename patterns is not complete for sure. I just added the cases I encountered with my files
* The files are not really named with 100% consistency. The most important thing for me right now was that all files start with `YYYYMMDD_`
* Errors are only logged. The next time it runs, it will try to process the same files over and over again
* There is no UI to start the authorization flow, instead, the auth URL is printed to stdout
* It's currently not multi-tenant aware. One deployment is only able to work with one onedrive account
