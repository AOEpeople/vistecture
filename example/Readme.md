# Vistecture example project

This example can be used as a start for your own documentation project.
The suggested folder structure is:

* projectname: Folder that contains the *.yml files with the service architecture documentation
* templates: Contains the templates that should be used in the generation
* generate-documentation.sh - The main build script - it executes the desired vistecture commands and generates the desired documentation artefacts
* build-with-docker.sh  - Helpful to run the generation local


## Usage

```sh
./build-with-docker.sh
ls -al result
open result/graphfull.png
```
