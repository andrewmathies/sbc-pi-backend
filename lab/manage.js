versions = []
versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25', '2.27']

$(document).ready(function() {
    console.log('Getting unit data')

    $.ajax({
        type: 'GET',
        url: 'api/units/',
        dataType: 'json',
        success: (data) => { loadDict(data) }
    })	
})

function loadDict(json) {
    dict = JSON.parse(json)
    
    for (key in dict) {
        console.log(dict[key])
    }
}