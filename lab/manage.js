let dict = []
let versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25', '2.27']
const interval = 30000

$(document).ready(function() {
    console.log('Getting unit data')

    $.ajax({
        type: 'GET',
        url: '/api/units/',
        contentType: 'application/json',
        success: (resp) => { 
            dict = resp
            buildTable()
            poll()
        }
    })
})

function poll() {
    console.log('polling server')
    setTimeout(() => {
        $.ajax({
            type: 'GET',
            url: '/api/units/',
            contentType: 'application/json',
            success: (resp) => { 
                dict = resp
                buildTable()
                poll()
            }
        })
    }, interval)
}

function boxListener(element) {
    let key = element.id
    let name = element.value

    console.log('textbox change, key: ' + key + ', name: ' + name)

    params = {
        version: dict[key].version,
        beanID: dict[key].beanID,
        name: name,
        state: dict[key].state
    }

    $.ajax({
        type: 'PUT',
        url: '/api/units/' + key,
        contentType: 'application/json',
        data: JSON.stringify(params),
        success: (resp) => {
            console.log('update request successful')
            console.log(resp)
        }
    })
}

function dropdownListener(element) {
    let key = element.id
    let version = element.value

    console.log('select change, key: ' + key + ', version: ' + version)

    params = {
        version: version,
        beanID: dict[key].beanID,
        name: dict[key].name,
        state: 1 // 1 is State.Updating
    }

    $.ajax({
        type: 'PUT',
        url: '/api/units/' + key,
        contentType: 'application/json',
        data: JSON.stringify(params),
        success: (resp) => {
            console.log('update request successful')
            console.log(resp)
        }
    })
}

function buildTable() {
    // remove old table if there is one
    $('#dictTable tbody tr').remove()

    let table = $('#dictTable tbody')

    for (let key in dict) {
        let curUnit = dict[key]

        let row = $('<tr>').appendTo(table)
        let nameElement = $('<td>').appendTo(row)
        let beanElement = $('<td>').appendTo(row)
        let versionElement = $('<td>').appendTo(row)
        let stateElement = $('<td>').appendTo(row)

        nameElement.append($('<input type="text" id="' + key + '" onchange="boxListener(this)" value="' + curUnit.name + '">'))

        let dropdown = $('<select id="' + key + '" onchange="dropdownListener(this)"/>').appendTo(versionElement)
        dropdown.value = curUnit.version

        versionOptions.forEach(val => {
            $('<option />', {value: val, text: val, selected: val === curUnit.version}).appendTo(dropdown)
        })
        
        beanElement.append($('<div>').html(curUnit.beanID))
        stateElement.append(makeIcon(curUnit.state))
    }
}

function makeIcon(state) {
    switch (state) {
        case 0:
            //idle
            return $('<i class="fas fa-check-circle" style="color: #34C53C">')
        case 1:
            // updating
            return $('<i class="fas fa-spinner fa-pulse" style="color: #61D7FF">')
        case 2:
            // failed
            return $('<i class="fas fa-times-circle" style="color: #FF0104">')
        default:
            console.log('unexpected state in units response: ' + state)
            return
    }
}