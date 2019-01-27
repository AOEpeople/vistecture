visTectureHelper = {}

//loadVistectureData - loads the vistecture project data
visTectureHelper.LoadVistectureData = function(project, callback) {


    if (typeof DATAURL == 'undefined') {
        alert('dataurl missing')
        return
    }
    if (!DATAURL) {
        alert('dataurl missing')
        return
    }

    $.getJSON( DATAURL+'?project='+project).done(function(data,statustext,jqXHR) {
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



visTectureHelper.FindApp = function(appId,projectData) {
    for (var i in projectData.applications) {
        let app = projectData.applications[i]
        if (app.id == appId) {
            return app
        }
    }
    return false
}


visTectureHelper.GetIncomingDependencies = function(projectData, appid) {
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


