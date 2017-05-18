# vistecture: Service Architecture Tool

A tool for visualizing and analyzing distributed (micro) service oriented architectures.

![Kiku](templates/example.jpg)

## Define your application architecture:

Describe your architecture in JSON (use .json) or YAML (use .yml). You can do this in two ways
- one json file or
- in multiple json files that are all in one directory. (prefered for structuring bigger definitions)

### Installation Options:
#### Go Get

```
go get github.com/AOEpeople/vistecture
```
#### Use a Docker container

It has graphviz and vistecture installed and can be used directly:
```
docker pull aoepeople/vistecture
```

Example usage with a definition from current folder:
```
docker run -v $(pwd):/workspace vistecture  vistecture --config /workspace analyze

docker run -v $(pwd):/workspace vistecture  sh -c "vistecture --config /workspace/definition graph --iconPath /usr/src/go/src/github.com/AOEpeople/vistecture/templates/icons | dot -Tpng -o /workspace/graph.png"
```


#### Download Binaries
You can also download a published release from github:

E.g. for mac:
(For linux use "vistecture-linux" and for windows "vistecture.exe")

```
curl -LOk "https://github.com/AOEpeople/vistecture/releases/download/0.2.beta/vistecture"
chmod +x vistecture

# dowload the templates
curl -LOk "https://github.com/AOEpeople/vistecture/releases/download/0.2.beta/templates.zip"
unzip templates.zip

```


And then discover the command:

```
vistecture help
```

You can also clone the repository and use golang tools.


### Example project definition ( example.yml ):


```yaml
---
name: Ports and Adapters DDD Architecture
applications:
- name: domain
  group: component-internal-bounded-context
  technology: scala
  display:
    bordercolor: "#c922b3"
  summary: Short description
  properties:
    foo: bar
    my-version: 0.1.latest
  description: |
        Use markdown to describe the service.
        * one
        * tow
  provided-services:
  - name: domain-objects
    type: inbound-port
  - name: repository-interfaces
    type: outbound-port
  - name: eventpublish-interfaces
    type: outbound-port
  infrastructure-dependencies:
  - type: mysql
- name: application
  group: component-internal-bounded-context
  description: Application Use Cases
  provided-services:
  - name: application-services
    type: inbound-port
    description: Main Application API. Also internaly works with eventpublish-interfaces
    dependencies:
    - reference: domain.domain-objects
    - reference: domain.repository-interfaces
      relationship: uses
  - name: eventpublish-interfaces
    type: outbound-port
  - name: domainEventAdapter
    type: "(logic)"
    description: Listen for (some) domain events. Also internaly works with eventpublish-interfaces
    dependencies:
    - reference: domain.eventpublish-interfaces
      relationship: implements

```


## Usage

Currently the main feature is generating graphviz compatible graph descriptions that can be used by any of the graphviz layouters like this:

Complete Graph
```
> vistecture --config=pathtojson graph | dot -Tpng -o graph.png
```

Graph for one Application and its direct dependencies (including infrastructure dependencies):
```
> vistecture --config=pathtojson graph --application=applicationame | dot -Tpng -o graph.png
```
The generation of the graph can add small icons to the applications. Therefore the tool looks in `iconPath` for a .png file matching the defined "technologie".

You can also render a documentation - expecting the dot command is executable for the application it will embedd svg images:

```
> vistecture --config=pathtojson documentation > documentation.html
```
The rendering needs a go html template. The "template" folder comes with a nice example.
You can download the templates to your local filesystem and use or modify them.

E.g.:

```
> vistecture --config=pathtojson documentation --templatePath=$GOPATH/github.com/AOEpeople/vistecture/templates/htmldocument.tmpl > documentation.html
```

Check for cyclic dependencies and get a very basic impact analysis:

```
> vistecture --config=pathtojson analyze
```

## Domain Language / Concepts:

This tool defines:

**Application:**
A Application is normaly something that offers one or more service-components (or interfaces).
Normaly a Application is something that is deployed seperate - and has a seperate build and integration pipeline.

- Supported Categories: external (rendered in red)
- Supported Technologies: go, scala, magento, akeneo, php, anypoint, keycloak (they will get a nice icon)

**Service:**
An Application offers services (more specific service components - but we use services here).
An application can offer one or more services.
Services are used by other systems or humans. 

Service Properties:
- isPublic: They can be public or just internal.
- isOpenHost: The service is a well designed published API
- securityLevel: Classification of the API in regard of security
- dependencies: Array of Dependency

**Dependency:**
A Application or a service can have dependencies. 
You can define dependency on application or also on service level (to emphazize that the dependency is only required for a certain service.)
A depdendency creates a reference to either a application - or more exact to a service. The relation is of a certain relationship type.

Dependency Properties:
- reference: String in the format "Applicationname.Servicename" (Servicename is optional)
- relationship: String - defining the collaboration level between the two bounded contexts / relationship. (See below)
- isSameLevel: Boolean. Use this to influence graph formating - to emphazise that the services are semantically on the same level
- resilience: String. Define the implemented resilience pattern
- isBrowserBased: If the dependency is established in the browser (and not from the backends.) This results in a dashed line.

**supported releationship types:**
- parntership: Very close collaboration
- customer-supplier: (use this where a strong dependency extsis that the supplier delivers whats required by the customer. A stronger collaboration between the teams of the components need to exist.)
- conformist: (use this to emphazise also a strong dependency that we need the services provided. But there is no chance to influence the interface - so the downstream component is forced to be conform to whatever is provided - and need to make it work.)
- acl: (Anti coruption layer: If the provided interface is complex or very different from the applications bounded context internal model. The acl emphazizes that the downstream component takes care to isolate his domain with a acl pattern)

(See https://www.aoe.com/tech-radar/strategic-domain-driven-design.html )

**Business Services:**
Several service components are typically composed to business services.
For example an ecommerce shop business service may consist of services from  ecommerce application, login application, search application.


## Todos

-  [ ] Introduce useful resilience pattern types for the dependencies
-  [ ] Introduce Business Services as Composite of Applications (Service Components)
-  [ ] Better Impact Analysis for application failures - inlcuing resilience evaluation
-  [ ] Generate useful artefacts for infrastructure pipeline (e.g. consul acls, service discovery tests...)

## Development

```
go get github.com/AOEpeople/vistecture
cd github.com/AOEpeople/vistecture

//run tests
go test ./tests/...

//build binaries:
make all

//build docker
docker build --no-cache -t aoepeople/vistecture .



```



## Jenkins Tip

Disable CSP Header in Jenkins to allow inline styles (required for a direct view of the generated documentation as jenkins artefact)
Open Jenkins script console and type:

```
System.setProperty("hudson.model.DirectoryBrowserSupport.CSP", "")
```