import $ from 'jquery';
window.$ = $;
import '../node_modules/popper.js/dist/popper.js'
import '../node_modules/bootstrap/dist/js/bootstrap.bundle.js'
import visRenderer from "./visRenderer"
import layout from './layoutFunctions'
import vistectureHelper from "./vistectureHelper";

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

var applicationInit = {}

applicationInit.DrawConfiguredGraph = function() {
    let value = $("#select-graph").val()
    let config = layout.GetGraphConfiguration()
    let selectedSubView = $("#select-project").val()
    let networkFilterGroups = $("#networkFilterGroups").val()

    vistectureHelper.LoadVistectureData(selectedSubView,networkFilterGroups,function(projectData) {
        applicationInit.updateProjectDropdown(projectData.availableSubViews, config)
        applicationInit.updateGroups(projectData.applicationsByGroup, projectData.availableGroups, config)
        layout.SetDocumentsMenu(projectData.staticDocumentations)
        visRenderer.RenderNetwork(document.getElementById('maincontent'),projectData, config)
    })
}

applicationInit.updateProjectDropdown = function(availableSubViews) {
    let selected =  $("#select-project").val()
    $("#select-project").find('option').remove()
    $("#select-project").append(new Option("Select project",""));
    for (var i in availableSubViews) {
        let name = availableSubViews[i]
        let selected = false
        if (selected == name) {
            selected = true
        }
        $("#select-project").append(new Option(name,name,selected,selected));
    }
}


applicationInit.updateGroups = function(applicationsByGroup,availableGroups) {
    let selectedClusterOptions =  $("#networkClusterGroups").val()
    $("#networkClusterGroups").find('option').remove()
    applicationInit._addSubGroups("#networkClusterGroups",applicationsByGroup.subGroups,selectedClusterOptions, 0)

    let selectedFilterOptions =  $("#networkFilterGroups").val()
    $("#networkFilterGroups").find('option').remove()
    applicationInit._addSubGroups("#networkFilterGroups",availableGroups.subGroups,selectedFilterOptions, 0)
}



applicationInit._addSubGroups = function(selector, subGroups, selectedOptions, level) {
    if (level > 9) {
        return
    }
    for (var i in subGroups) {
        let subGroup = subGroups[i]
        if (subGroup.groupName == "") {
            continue
        }
        let selected = false
        if ($.inArray( subGroup.qualifiedGroupName, selectedOptions ) != -1) {
            selected = true
        }

        $(selector).append(new Option("  > ".repeat(level)+subGroup.groupName,subGroup.qualifiedGroupName,selected,selected));
        applicationInit._addSubGroups(selector,subGroup.subGroups,selectedOptions,level+1)
    }
}
