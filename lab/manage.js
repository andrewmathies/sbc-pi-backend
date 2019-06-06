dict = []
versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25', '2.27']

$(document).ready(function() {
    console.log('Getting unit data')

    $.ajax({
        type: 'GET',
        url: 'api/units/',
        dataType: 'json',
        success: (data) => { 
            dict = data
            buildTable()
        }
    })
})

function buildTable() {
    // remove old table if there is one
    $('#dictTable tbody tr').remove()

    let table = $('#dictTable tbody')

    let options = ''
    for (let i = 0; i < versionOptions.length; i++) {
        if (versionOptions[i]) {
            options += $('<option value="' + i + '">').html(versionOptions[i])
        }
    }

    for (let key in dict) {
        let curUnit = dict[key]

        let row = $('<tr>').appendTo(table)
        let nameElement = $('<td>').appendTo(row)
        let beanElement = $('<td>').appendTo(row)
        let versionElement = $('<td>').appendTo(row)
        let stateElement = $('<td>').appendTo(row)

        nameElement.append(
            $('<input>', {
                type: 'text',
                val: curUnit.name
            })
        )
        beanElement.append($('<div>').html(curUnit.beanID))
        versionElement.append($('<select id="' + key + '" class="cs">').append(options))
        stateElement.append(makeIcon(curUnit.state))
    }
}

function makeIcon(state) {
    switch (state) {
        case 0:
            //idle
            return $('<span class="ui-icon-check">')
        case 1:
            // updating
            return $('<span class="ui-icon-refresh">')
        case 2:
            // failed
            return $('<span class="ui-icon-closethick">')
        default:
            console.log('unexpected state in units response: ' + state)
            return
    }
}