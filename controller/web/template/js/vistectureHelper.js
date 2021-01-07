import $ from 'jquery'

export default class vistectureHelper {}

//loadVistectureData - loads the vistecture project data
vistectureHelper.LoadVistectureData = function(selectedSubView, networkFilterGroups, callback) {


    if (typeof DATAURL == 'undefined') {
        alert('dataurl missing')
        return
    }
    if (!DATAURL) {
        alert('dataurl missing')
        return
    }
    let ajaxUrl = DATAURL
    let params = []
    if (selectedSubView != null) {
        params.push('subview='+selectedSubView)
    }
    if (networkFilterGroups != null) {
        params.push('filterGroups='+networkFilterGroups)
    }

    ajaxUrl = ajaxUrl + '?' + params.join('&')
    $.getJSON( ajaxUrl).done(function(data,statustext,jqXHR) {
        //check if the status is 200(means everything is okay)
        if (jqXHR.status == 200)
        {
            callback(data)
        } else {
            $('#maincontent').html('<div class="mx-auto text-danger m-5" style="max-width: 400px;">Loading error... Server returned invalid statuscode</div>')
            console.log(data,status)
        }
    }).fail(function(data) {
        $('#maincontent').html('<div class="mx-auto text-danger m-5" style="max-width: 400px;">Loading error... check server connection</div>')

    });
}



vistectureHelper.FindApp = function(appId, projectData) {
    for (var i in projectData.applications) {
        let app = projectData.applications[i]
        if (app.id == appId) {
            return app
        }
    }
    return false
}


vistectureHelper.GetIncomingDependencies = function(projectData, appid) {
    let incomingDep = []
    for (var i in projectData.applications) {
        let app = projectData.applications[i]
        for (var j in app.dependenciesGrouped) {
            let appDep = app.dependenciesGrouped[j]
            if (appDep.application.id == appid) {
                incomingDep.push(appDep)
            }
        }
    }
    return incomingDep
}


