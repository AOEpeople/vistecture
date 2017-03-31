# appdependency: Service Architecture Tool

A tool for visualizing and analyzing distributed (micro) service oriented architectures.

![Kiku](templates/example.jpg)

## Define your application architecture:

Describe your architecture in JSON. You can do this in two ways
- one json file or
- in multiple json files that are all in one directory. (prefered for structuring bigger definitions)

### Installation:

Download a published release from github:

E.g. for mac:
(For linux use "appdependency-linux" and for windows "appdependency.exe")

```
curl -LOk "https://github.com/danielpoe/appdependency/releases/download/0.2.alpha/appdependency"
chmod +x appdependency

curl -LOk "https://github.com/danielpoe/appdependency/releases/download/0.2.alpha/templates.zip"

```


And then discover the command:

```
appdependency help
```

You can also clone the repository and use golang tools.


### Example:

```
{
  "name": "my application suite",
  "components": [
    {
      "name": "Name of Component",
      "group": "Optional a Group",
      "technology": "scala",
      "category": "Optional a category"
      "description": "Some short description / New line",
      "summary": "Optional shorter summary",
      "display":{
        "bordercolor": "#3971ad"
      },
      "provided-services": [
        {
            "name": "auth",
            "type": "api",
            "description": "A description of the service",
            "isPublic": true,
            "dependencies": [
              {
                "reference": "othercomponent"
              }
            ]
        }
      ],
      "infrastructure-dependencies": [
        {
          "type": "redis"
        }
      ],
      "dependencies": [
          {
            "reference": "keycloak.login",
            "relationship": "serviceapi",
            "isSameLevel": false,
            "isBrowserBased": true,
            "resilience": true
          }
      ]
    }

    ...

```


## Usage

Currently the main feature is generating graphviz compatible graph descriptions that can be used by any of the graphviz layouters like this:

```
> appdependency --config=pathtojson graph | dot -Tpng -o graph.png
```

You can also render a documentation - expecting the dot command is executable for the application it will embedd svg images:

```
> appdependency --config=pathtojson documentation > documentation.html
```

## Domain Language / Concepts:

This tool defines:

**Component:**
A component is normaly something that offers one or more services (or interfaces).
Normaly a component is something that is deployed seperate - and has a seperate build and integration pipeline.

- Supported Categories: external (rendered in red)
- Supported Technologies: go, scala, magento, akeneo, php, anypoint, keycloak

**Service:**
A component offers services (one or more). Services can have a type. Services are used by other systems or humans. The can be public or just internal.

**Dependency:**
A component or a service can have dependencies. You can add dependency on service level to emphazize that the dependency is only required for a certain service.
(This is used for impact analyses).
A depdendency creates a reference to either a component - or more exact to a service. The relation is of a certain relationship type.

**supported releationship types:**
- customer-supplier (use this where a strong dependency extsis that the supplier delivers whats required by the customer. A stronger collaboration between the teams of the components need to exist.)
- conformist (use this to emphazise also a strong dependency that we need the services provided. But there is no chance to influence the interface - so the downstream component is forced to be conform to whatever is provided - and need to make it work.)
- serviceapi (The used api is designed for integration. Its nice and offers multiple services and the format is published (documented). Most modern Rest API should fall under this section. (see open host / published language). This is the default
- acl (Anti coruption layer: If the provided interface is complex or very different from the components internal model. The acl emphazizes that the downstream component takes care to isolate his domain with a acl pattern)


## Todos


-  [X] Graph for single component including infrastructure
-  [ ] Impact Analysis for component failures
-  [ ] Create complete documentation
-  [ ] Generate useful artefacts for infrastructure pipeline (e.g. consul acls, service discovery tests...)
