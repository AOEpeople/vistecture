# vistecture: Service Architecture Tool

![dependency security scanning](https://github.com/AOEpeople/vistecture/workflows/dependency%20security%20scanning/badge.svg)

A tool for visualizing and analyzing distributed (micro) service oriented architectures.
Just define your applications (microservices) with its dependencies in a simple yaml file.

You can use the online browser:
![Example](doc/onlinebrowser.png)

Or you can use it to render graphviz based images and any kind of text documentations.

![Example](doc/readme-example.png)

## Define your application architecture:

Describe your architecture in YAML in a machine readable format and generate various documentation artefacts out of it.

For an example see the "demoproject" in the example folder.

You can put your definition in one (big) file or split it up in  multiple files in structured directories (which is preferred for structuring bigger definitions).


## Installation Options:
### Go Get

```commandline
go get github.com/AOEpeople/vistecture
```

### Use a Docker container

It has graphviz and vistecture installed and can be used directly:
```commandline
docker pull aoepeople/vistecture
```

Example usages:
```commandline
cd /your/path/with_vistecture_defintions

docker run -v $(pwd):/workspace -p 8080:8080 aoepeople/vistecture vistecture --config=/workspace/projectconfig.yml serve

docker run -v $(pwd):/workspace aoepeople/vistecture  vistecture --config=/workspace/projectconfig.yml analyze

docker run -v $(pwd):/workspace aoepeople/vistecture  sh -c "vistecture --config=/workspace/projectconfig.yml graph --iconPath=templates/icons | dot -Tpng -o /workspace/graph.png"
```


### Or Download Binaries
You can also download a published release from Github:

E.g. for macOS. 

(For Linux use `vistecture-linux` and for Windows `vistecture.exe`)

```commandline
curl -LOk "https://github.com/AOEpeople/vistecture/releases/download/0.2.beta/vistecture"
chmod +x vistecture

# download the templates
curl -LOk "https://github.com/AOEpeople/vistecture/releases/download/0.2.beta/templates.zip"
unzip templates.zip
```


And then discover the command:

```commandline
vistecture help
```


## Vistecture Configuration Format:

Vistecture need a Project Configuration and multiple Application configurations:

### Project Configuration

The project configuration can be used to:

- load the list of applications (key `appDefinitionsPaths`) this can be path to a folder or concrete file that contain the **application configuration**

Here is an example:

```yaml
projectName: "Demoproject"
appDefinitionsPaths:
- external-services
- service-group-1
- service-group-2
appOverrides:
- name: customer-portal
  add-provided-services:
  - name: loyalty
    type: gui
    dependencies:
    - reference: order-workflow
  remove-dependencies:
  - test
subViews:
- name: "Demoproject minimal"
  included-applications:
  - customer-portal
  - single-sign-on
  - order-workflow
```

### Application Configuration

```yaml
name: service1
group: group1
technology: scala
team: team1
display:
  borderColor: "#c922b3"
summary: Short description
properties:
  foo: bar
  my-version: 0.1.latest
description: |
  Use markdown to describe the service.
  * one
  * tow
provided-services:
  - name: someApi
    type: api
    dependencies:
    - reference: service1.someApi
      relationship: partnership
      description: Some description here
  - name: otherApi
    type: api
  - name: eventpublish
    type: exchange
infrastructure-dependencies:
  - type: mysql
dependencies:
  - reference: service2
```

Please also see chapter 'Domain Language / Concepts' for more information

## Usage Options

### Run Browser based view:

```commandline
vistecture --config=pathtodefinitions serve
```

To add a list to other assets you can use the `staticDocumentsFolder` parameter
```commandline
vistecture --config=pathtodefinitions serve --staticDocumentsFolder=/folderwithother_docs
```

### Generate Graphs:



A main feature is generating graphviz compatible graph descriptions that can be used by any of the graphviz layouters like this:

Complete Graph:
```commandline
vistecture --config=pathtodefinitions graph | dot -Tpng -o graph.png
```
(Not all graphviz versions will work!)

Graph for a dedicated project configuration:
```commandline
vistecture --config=pathtodefinitions graph --projectname=nameoftheproject | dot -Tpng -o graph.png
```

Graph for one application and its direct dependencies (including infrastructure dependencies):
```commandline
vistecture --config=pathtodefinitions graph --application=applicationame | dot -Tpng -o graph.png
```

The generation of the graph can add small icons to the applications. Therefore the tool looks in `iconPath` for a .png file matching the defined "technology".

#### Team Graphs
You can also draw the resulting relationships between the teams:
```commandline
vistecture --config=pathtodefinitions teamGraph  --summaryRelation 1 | dot -Tpng -Gbgcolor=white -o teamgraph.png
```

### Generate documentations:
You can also render a documentation - expecting the dot command is executable for the application it will embed svg images:

```commandline
vistecture --config=pathtodefinitions documentation > documentation.html
```
The rendering needs a go html template. The "template" folder comes with a nice example.
You can download the templates to your local filesystem and use or modify them.

E.g.:

```commandline
vistecture --config=pathtodefinitions documentation --templatePath=$GOPATH/github.com/AOEpeople/vistecture/templates/htmldocument.tmpl > documentation.html
```

### Other artefacts:
Check for cyclic dependencies and get a very basic impact analysis:

```commandline
vistecture --config=pathtodefinitions analyze
```

## Concepts and the Domain Language of the Service definition:

This tool defines:

### Project

A project defines which applications to be included for processing at runtime.
In general, a project acts as a repository "overwrite" to minimize configuration effort.

If no project or no "included applications" are configured, all available applications in the repository will be taken.
If multiple projects are available and no one is explicitly mentioned in the command line call, the first found will be taken.

The "included-applications" container references defined applications by name but can overwrite all other root attributes
of an application like title, properties etc.

### Application
An Application is normally something that offers one or more service-components (or interfaces).
Normally an Application is something that is deployed separately - and has a separate build and integration pipeline.

- Supported Categories: external (rendered in red)
- Supported Technologies: go, scala, magento, akeneo, php, anypoint, keycloak (they will get a nice icon)

### Service
An Application offers services (more specific service components - but we use services here).
An application can offer one or more services.
Services are used by other systems or humans.

| Service Properties | Description |
| --- | --- | 
| isPublic      | They can be public or just internal. |
| isOpenHost    | The service is a well designed published API |
| securityLevel | Classification of the API in regard of security (public, internal, confidential, restricted) |
| dependencies  | Array of Dependencies |

### Dependency
An Application or a service can have dependencies.
You can define dependency on application or also on service level (to emphasize that the dependency is only required for a certain service.)
A dependency creates a reference to either an application, or more exactly, to a service. The relation is of a certain relationship type. You can also add a description to explain more details to this dependency.

| Dependency Properties | Description |
| --- | --- | 
| reference      | String in the format `Applicationname.Servicename` (Servicename is optional). |
| relationship   | String - defining the collaboration level between the two bounded contexts / relationship (see below). |
| isSameLevel    | Boolean. Use this to influence graph formatting - to emphasise that the services are semantically on the same level. |
| resilience     | String. Define the implemented resilience pattern. |
| isBrowserBased | If the dependency is established in the browser (and not from the backend.) This results in a dashed line. |

#### Relationship types

| Relationship type | Description |
| --- | --- | 
| partnership         |  Very close collaboration  | 
| customer-supplier   |  A strong dependency exists. The supplier delivers what is required by the customer. A stronger collaboration between the teams of the components needs to exist. | 
| conformist          |  Emphasise a strong dependency that we need the services provided. But there is no chance to influence the interface - so the downstream component is forced to conform to whatever is provided - and needs to make it work.  | 
|                     |  (please note that on a per service level it is possible to define the service as `isOpenHost` - this normally means that the applications consuming this services fall in the relationship "open-host") | 
| acl                 |  Anti corruption layer: If the provided interface is complex or very different from the applications bounded context internal model. The acl emphasizes that the downstream component takes care to isolate his domain with a acl pattern) | 

(See https://www.aoe.com/techradar/methods-and-patterns/strategic-domain-driven-design.html)

### Groups (Business Services)
Applications can be grouped. For example, this can be used to visualize Business Services:

Several service components are typically composed to business services.
For example, an e-commerce shop business service may consist of services from e-commerce application, login application, search application.


## Todos

-  [ ] Better Impact Analysis for application failures - including resilience evaluation
-  [ ] Generate useful artifacts for infrastructure pipeline (e.g. consul acls, service discovery tests...)

## Development

```commandline
go get github.com/AOEpeople/vistecture
cd github.com/AOEpeople/vistecture

// run tests
go test ./...

// build binaries:
make all

// releasing: docker publish
// push changes to github /
// adjust version number in Makefile (!!)  and run
make dockerpublish
```
