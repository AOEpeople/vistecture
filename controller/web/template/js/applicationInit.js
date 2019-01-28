


// Dom Ready: Shorthand for $( document ).ready()
$(function() {

    layout.SetGraphPresets($( "#select-graphpreset" ).val())
    applicationInit.DrawConfiguredGraph()

    $( "#select-graphpreset" ).change(function() {
        layout.SetGraphPresets($( "#select-graphpreset" ).val())
        applicationInit.DrawConfiguredGraph()
    });
    $( "#select-project" ).change(applicationInit.DrawConfiguredGraph);
    $( "#updateGraphConfiguration" ).click(applicationInit.DrawConfiguredGraph);
    $('#networkConfigureForm').change(applicationInit.DrawConfiguredGraph)
});

applicationInit = {}

applicationInit.DrawConfiguredGraph = function() {
    let value = $("#select-graph").val()
    let config = layout.GetGraphConfiguration()
    let selectedProject = $("#select-project").val()

    visTectureHelper.LoadVistectureData(selectedProject,function(projectData) {
        applicationInit.updateProjectDropdown(projectData.availableProjectNames, config)
        applicationInit.updateGroups(projectData.applicationsByGroup, config)
        layout.SetDocumentsMenu(projectData.staticDocumentations)
        visRenderer.RenderNetwork(document.getElementById('maincontent'),projectData, config)
    })
}

applicationInit.updateProjectDropdown = function(availableProjectNames) {
    let selected =  $("#select-project").val()
    $("#select-project").find('option').remove()
    $("#select-project").append(new Option("Select project",""));
    for (var i in availableProjectNames) {
        let name = availableProjectNames[i]
        let selected = false
        if (selected == name) {
            selected = true
        }
        $("#select-project").append(new Option(name,name,selected,selected));
    }
}


applicationInit.updateGroups = function(applicationsByGroup) {
    let selectedOptions =  $("#networkClusterGroups").val()
    $("#networkClusterGroups").find('option').remove()
    applicationInit._addSubGroups(applicationsByGroup.subGroups,selectedOptions,0)
}


applicationInit._addSubGroups = function(subGroups, selectedOptions,level) {
    if (level > 9) {
        return
    }
    for (var i in subGroups) {
        let subGroup = subGroups[i]
        let selected = false
        if ($.inArray( subGroup.groupName, selectedOptions ) != -1) {
            selected = true
        }

        $("#networkClusterGroups").append(new Option("   ".repeat(level)+subGroup.groupName,subGroup.groupName,selected,selected));
        applicationInit._addSubGroups(subGroups.subGroups,level+1)
    }
}