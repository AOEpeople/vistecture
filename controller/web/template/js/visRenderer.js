import visNetworkHelper from "./visNetworkHelper"
import layout from "./layoutFunctions"

import chroma from 'chroma-js'
import vis from "../node_modules/vis-network/dist/vis-network"
/**
 File contains render method to render the graph in given container
 */

export default class visRenderer {}
visRenderer.networkInstance = null
visRenderer.RenderNetwork = function(container, projectData,configurationData) {
    let node;
	const defaultProperties = {
        'hierarchicalSortMethod': "",
        'nodeStyle': "",
        'layout': "",
        'physics': false,
        'physicsStabilization': false,
        'clusterGroups': [],
        'filterGroups': [],
    }
    const config = Object.assign(defaultProperties, configurationData)

    //console.log("renderNetwork",config)

    let edges
    if (config['nodeStyle'] === "detailed") {
        edges = visRenderer.getGroupedEdges(projectData, 300)
    } else {
        edges = visRenderer.getGroupedEdges(projectData, 0)
    }


    let nodes = [];

    for (var app in projectData.applications) {
        let application = projectData.applications[app]

        let node = visRenderer.getBasicNode(application, config['nodeStyle'])
        nodes.push(node);
    }

    //Check if we are missing nodes (this means an application reference an application that is not specified)
    for (let missingAppId in projectData.missingApplications) {
        let missingApp = projectData.missingApplications[missingAppId]
        node = { id: missingApp.id, font: {color: "#ffffff"}, title: "MISSING! " + missingApp.title, label: "MISSING!" + missingApp.name, color: {border: "#ff0000", background: "#ee0000", highlight: {background:"#ff0000", border: "#ffaaaa"}}};
        nodes.push(node);
    }
    //Check if we are missing nodes (this means an application reference an application that is not specified)
    for (let unincludedAppId in projectData.unincludedApplications) {
        let missingApp = projectData.unincludedApplications[unincludedAppId]
		node = {
			id: missingApp.id,
			font: {color: "#999999"},
			title: "OUTSIDE: " + missingApp.title,
			label: "OUTSIDE: " + missingApp.name,
			color: {border: "#ee0000", background: "#666666", highlight: {background: "#aaaaaa", border: "#bbaaaa"}}
		};
		nodes.push(node);
    }

    for (var app in projectData.applications) {
        let application = projectData.applications[app]
        for (let missingAppIndex in application.dependenciesToMissingApplications) {
            let missingApp = application.dependenciesToMissingApplications[missingAppIndex]
            var edge = {color: {color: "#ee0000", highlight: "#ff0000"}, smooth:{enabled: false},arrows:{to: {enabled:true}}, from: application.id, to: missingApp.id}
            edges.push(edge);
        }
        for (let unincludedAppIndex in application.dependenciesToUnincludedApplications) {
            let unincludedApp = application.dependenciesToUnincludedApplications[unincludedAppIndex]
            var edge = {color: {color: "#aaaaaa", highlight: "#eeeeee"}, smooth:{enabled: false},arrows:{to: {enabled:true}}, from: application.id, to: unincludedApp.id}
            edges.push(edge);
        }
    }


    // provide the data in the vis format
	const data = {
		nodes: nodes,
		edges: edges
	};

	let layout = {}
    if (config["layout"] === "hierarchical") {
        layout = {
            hierarchical: {
                enabled: true,
                sortMethod: config['hierarchicalSortMethod']
            }
        }
        if (config['nodeStyle'] === "detailed") {
            layout.hierarchical.levelSeparation= 550
            layout.hierarchical.nodeSpacing= 440
        }
    }


	const options = {
		physics: {
			enabled: config['physics'],
			stabilization: config['physicsStabilization'],
		},
		//physics: false,
		edges: {
			smooth: {
				type: 'continuous'
			},
			arrows: {
				to: {enabled: true, scaleFactor: .5, type: 'arrow'}
			},
			scaling: {
				max: 4,
			}
		},
		nodes: {
			shape: 'box',
			shadow: true
		},
		layout: layout
	};

	//console.log("renderNetworkOptions",options)

    // initialize your network!
    visRenderer.networkInstance = new vis.Network(container, data, options);
    visNetworkHelper.InitNetwork(visRenderer.networkInstance,function(nodeParams) {
        visRenderer.clickEventListener(nodeParams,projectData)
    })
    visNetworkHelper.ClusterByApplicationGroups(visRenderer.networkInstance,config['clusterGroups'],function(groupName) {
        return visRenderer.clusterNodeRenderer(groupName,projectData,configurationData)
    })

}



visRenderer.clickEventListener= function (nodeParams, projectData) {
    if (nodeParams.nodes.length<1) {
        return
    }
    let app = visTectureHelper.FindApp(nodeParams.nodes[0],projectData)
    //console.log("clicked.. ",nodeParams,app)
    if (app !== false) {
        // RENDER COMMONTAB
        let commonTab = ""
        if (typeof app.summary != "undefined") {
            commonTab = app.summary
        }
        let incomingDep = visTectureHelper.GetIncomingDependencies(projectData,app.id)
        commonTab = commonTab + `<p class="pt-1"><strong>Consumers:</strong>There are <strong>#${incomingDep.length}</strong> consumers</p>`
        commonTab = commonTab + `<p class="pt-1"><strong>Group:</strong> ${app.group}</p>`
        if (typeof app.description != "undefined") {
            commonTab = commonTab + `<p><small>${app.description}</small></p>`
        }
        let propContent = ""
        for (var pIndex in app.properties) {
            let property = app.properties[pIndex]
            propContent += `<tr><td>${pIndex}</td><td>${property}</td></tr>`
        }
        if (propContent != "") {
            commonTab = commonTab + `<h6>Properties:</h6> <table class="mt-1 table table-sm small"><tbody>${propContent}</tbody></table> `
        }


        // RENDER SERVICETAB
        let serviceContent = ""
        for (var sIndex in app['provided-services']) {
            let service = app['provided-services'][sIndex]
            let propContent = ""
            for (var pIndex in service.properties) {
                let property = service.properties[pIndex]
                propContent = `<tr><td>${pIndex}</td><td>${property}</td></tr>`
            }

            serviceContent = serviceContent +` <li class="list-group-item">
                            <div class="d-flex w-100 justify-content-between">
                                <h7 class="mb-1">${service.name}</h7>
                                <div class="badge badge-primary text-wrap" style="width: 6rem;">
                                    ${service.type}
                                </div>
                                <small>-</small>
                            </div>
                            <small class="">${service.description}</small>
                            <table class="mt-1 table table-sm small"><tbody>${propContent}</tbody></table>
                        </li>`
        }
        serviceContent = `<ul class="list-group">${serviceContent}</ul>`

        // RENDER DEP TAB
        let depContent = ""
        let renderDepItem = function(dep, extraInfo) {
            let propContent = ""
            for (let pIndex in dep.properties) {
                let property = dep.properties[pIndex]
                propContent = `<tr><td>${pIndex}</td><td>${property}</td></tr>`
            }

            let relationShip = ""
            if (typeof dep.relationship != "undefined") {
                relationShip = dep.relationship
            }
            if (dep.isBrowserBased) {
                relationShip = relationShip + "(browser based)"
            }
            if (relationShip != "") {
                relationShip = `<div class="badge badge-primary text-wrap" style="width: 6rem;">${relationShip}</div>`
            }
            let description = ""
            if (typeof depContent.description != "undefined") {
                description = depContent.description
            }
            return ` <li class="list-group-item">
                            <div class="d-flex w-100 justify-content-between">
                                <h7 class="mb-1">${dep.reference}</h7>
                                ${relationShip}
                            </div>
                            <small class="">${extraInfo} ${description}</small>
                            <table class="mt-1 table table-sm small"><tbody>${propContent}</tbody></table>
                        </li>`
        }
        for (var sIndex in app['dependencies']) {
            let dep = app['dependencies'][sIndex]
            depContent = depContent + renderDepItem(dep,"")
        }
        for (var sIndex in app['provided-services']) {
            let service = app['provided-services'][sIndex]
            for (let dIndex in service['dependencies']) {
                let dep = service['dependencies'][dIndex]
                depContent = depContent + renderDepItem(dep,"From " + service.name)
            }
        }
        depContent = `<ul class="list-group">${depContent}</ul>`

        let title = `[${app.name}]`
        if (typeof  app.title != "undefined") {
            title = title +  ` ${app.title}`
        }
        layout.ShowSideContentModal(title,commonTab,serviceContent,depContent)
    }

}


visRenderer.getBasicNode = function(application, nodeStyle) {
    let colors = visRenderer.getColorsForApplication(application)
	let groupName = "UNDEFINED";
	if (application.group) {
        groupName = application.group
    }
	const node = {
		id: application.id,
		appGroup: groupName,
		font: {color: colors.fontColor},
		title: application.title,
		label: application.name,
		color: {
			border: colors.borderColor,
			background: colors.backgroundColor,
			highlight: {background: colors.highlightColor, border: colors.highLightBorderColor}
		}
	};
	node.title += ' ('
	if (application.group) {
		node.title += 'Group: ' + application.group
	}
	if (application.group && application.team) {
		node.title += ', ';
	}
	if (application.team) {
		node.title += 'Team: ' + application.team
	}
	node.title += ')'

	if (nodeStyle === "detailed") {
        Object.assign(node, { size: 300, image: visRenderer.applicationSvgUrl(application,colors), shape: 'image', borderWidthSelected: 6,shapeProperties: {useImageSize: true, useBorderWithImage: true  }})
        if (application.status === 'planned') {
            node.label = 'planned'
        }
    }
    return node
}


visRenderer.getStandardNodeBgColor = function(application) {
    //if the app has a color defined - use this one:
    if (application.hasOwnProperty('display') && application.display.hasOwnProperty('color') && application.display.color !== "") {
        return application.display.color
    }
    return visRenderer.getColorForGroup(application.group)
}


visRenderer.getColorForGroup = function(groupName) {
    let colorScaleSize = 6
    let colorScaleSizeTotal = 18
    //othrwise choose from a standard set of 6 different colors
    let colorsScale1 = chroma.scale(['#a6f196','#322b84'])
        .mode('lch').colors(colorScaleSize)
    let colorsScale2 = chroma.scale(['#ffd07d','#923069'])
        .mode('lch').colors(colorScaleSize)
    let colorsScale3 = chroma.scale(['#e795ad','#857dad'])
        .mode('lch').colors(colorScaleSize)
    let colors = colorsScale1.concat(colorsScale2.concat(colorsScale3))

    //choose same color for same group - looks nicer:
	let index;
	if (typeof groupName != "undefined") {
        let chr = 0
        for (let i = 0; i < groupName.length; i++) {
            chr   = chr + groupName.charCodeAt(i);
        }
        index = chr % colorScaleSizeTotal
    } else {
        index = Math.floor(Math.random() * colorScaleSizeTotal)
    }
    return colors[index]
}


//getGroupedEdges - returns the vis network edges between nodes - based on the dependencies between applications
// (we are using the grouped dependencies - to only draw one edge between two nodes - even if multiple dependencies to that node exist)
visRenderer.getGroupedEdges = function(projectData, lenght) {
    let edges = [];
    for (let app in projectData.applications) {
        let application = projectData.applications[app]

        let colors = visRenderer.getColorsForApplication(application)

        for (let groupedDepIndex in application.dependenciesGrouped) {
            let groupedDep = application.dependenciesGrouped[groupedDepIndex]
            let isBrowserBased = true
            for (let depIndex in groupedDep.dependencies) {
                let individualDep = groupedDep.dependencies[depIndex]
                if (individualDep.isBrowserBased === false) {
                    isBrowserBased  = false
                }
            }
            let node = {color: {color: colors.borderColor, highlight: colors.highLightBorderColor}, smooth:{enabled: false},arrows:{to: {enabled:true}}, from: groupedDep.sourceApplication.id, to: groupedDep.application.id, value: groupedDep.dependencies.length}
            if (isBrowserBased) {
                node.dashes = true
            }
            if (lenght > 0) {
                node.length= 1000
            }
            edges.push(node)
        }
    }
    return edges
}


visRenderer.applicationSvgUrl = function(application,colors) {
	const iconUrl = application.technology + '.png';
	const icon = '<img src="'+ iconUrl + '" scale="true" >';
	//console.log(colors)
	let tableHeaderColor = colors.backgroundColor; // "#1B4E5E"

    if  (application.category === 'external') {
        tableHeaderColor = "#8e0909"
    }


    let table = `<table border="1" style="width: 100%"><tr><td style="background-color: ${tableHeaderColor}; font-size:35px; color: ${colors.fontColor}">${application.name}</td></tr>`
    if (application.title) {
        table = table + `<tr><td>${application.title}</td></tr>`
    }
    for (let sIndex in application['provided-services']) {
        let service = application['provided-services'][sIndex]
        if (service.status === 'planned') {
            table = table + `<tr><td>${service.type}:${service.name}</td></tr>`
        } else {
            table = table + `<tr><td>${service.type}:${service.name}</td></tr>`
        }
    }

    table = table + `</table>`

    let height = 100
    if (application['provided-services']) {
        height = height+application['provided-services'].length*50
    }
	const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="350px" height="' + height + 'px">' +
		'<rect x="0" y="0" width="100%" height="100%" fill="#efefef" stroke-width="2" stroke="' + colors.borderColor + '" ></rect>' +
		'<foreignObject x="0" y="0" width="100%" height="100%">' +
		'<div xmlns="http://www.w3.org/1999/xhtml" style="font-size:18px; font-family: arial, sans-serif">' +
		table +
		'</div>' +
		'</foreignObject>' +
		'</svg>';
	return "data:image/svg+xml;charset=utf-8,"+ encodeURIComponent(svg);
}


visRenderer.clusterNodeSvgUrl = function(groupName,projectData,borderColor) {

	const tableHeaderColor = "#008482";

	let table = `<table border="1" style="width: 100%"><tr><td style="background-color: ${tableHeaderColor}; font-size:40px; color: #fff; padding: 2px">Group: ${groupName}</td></tr>`

    //table = table + `<tr><td>App Count:${projectData.applications.length}</td></tr>`

    table = table + `</table>`

	const svg = '<svg xmlns="http://www.w3.org/2000/svg" width="500px" height="200px">' +
		'<rect x="0" y="0" width="100%" height="100%" fill="#efefef" stroke-width="2" stroke="' + borderColor + '" ></rect>' +
		'<foreignObject x="0" y="0" width="100%" height="100%">' +
		'<div xmlns="http://www.w3.org/1999/xhtml" style="font-size:25px; font-family: arial, sans-serif">' +
		table +
		'</div>' +
		'</foreignObject>' +
		'</svg>';
	return "data:image/svg+xml;charset=utf-8,"+ encodeURIComponent(svg);
}




visRenderer.clusterNodeRenderer = function(groupName, projectData,configurationData) {
    let groupColor = visRenderer.getColorForGroup(groupName)
    let bgHcl = chroma(groupColor).lch()
    let fontColor = "#333333"
    if (bgHcl[0] < 63) {
        fontColor = "#efefef"
    }
	const node = {
		title: groupName,
		label: groupName,
		id: groupName + 'Cluster',
		borderWidth: 3,
		shape: 'box',
		color: {background: groupColor},
		margin: 12,
		font: {size: 19, color: fontColor}
	};
	if (configurationData['nodeStyle'] == "detailed") {
        Object.assign(node, { size: 300, image: visRenderer.clusterNodeSvgUrl(groupName,projectData), shape: 'image', shapeProperties: {useImageSize: true, useBorderWithImage: true  }})
    }
    return node
}


visRenderer.getColorsForApplication = function(application) {
    let nodeColor = visRenderer.getStandardNodeBgColor(application)
    let highlightColor = chroma(nodeColor).brighten(1).saturate(2).hex()

    let borderColor = chroma(nodeColor).darken(2).saturate(0.5).hex()

    if (application.hasOwnProperty('display') && application.display.hasOwnProperty('bordercolor') && application.display.bordercolor != "") {
        borderColor = application.display.bordercolor
    }
    let highLightBorderColor = chroma(borderColor).brighten(2).saturate(1).hex()


    let bgHcl = chroma(nodeColor).lch()
    let fontColor = "#333333"
    if (bgHcl[0] < 63) {
        fontColor = "#efefef"
    }

    return {
        backgroundColor: nodeColor,
        highlightColor: highlightColor,
        borderColor: borderColor,
        highLightBorderColor: highLightBorderColor,
        fontColor: fontColor
    }

}


