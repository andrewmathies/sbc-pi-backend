dict = []
versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25', '2.27']

$(document).ready(function() {
    console.log('Getting unit data')

    $.ajax({
        type: 'GET',
        url: '/api/units/',
        contentType: 'application/json',
        success: (resp) => { 
            dict = resp
            buildTable()
        }
    })
})

$('.nameBox').on('input propertychange paste', () => {
    console.log(this.id + ' textbox changed')
});

function dropdownListener(element) {
    let key = element.id
    let val = element.value

    console.log('key: ' + key + ', val: ' + val)

    params = {
        version: val,
        beanID: dict[key].beanID,
        name: dict[key].name,
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

        nameElement.append(
            $('<input>', {
                type: 'text',
                id: key,
                class: 'nameBox',
                val: curUnit.name
            })
        )

        let dropdown = $('<select id="' + key + '" onchange="dropdownListener(this)"/>').appendTo(versionElement)
        dropdown.value = curUnit.version

        versionOptions.forEach(val => {
            $('<option />', {value: val, text: val}).appendTo(dropdown)
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