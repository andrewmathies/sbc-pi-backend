dict = []
versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25', '2.27']

$(document).ready(function() {
    console.log('Getting unit data')

    $.ajax({
        type: 'GET',
        url: '/api/units/',
        dataType: 'json',
        success: (data) => { 
            dict = data
            buildTable()
        }
    })
})

// dropdown listener
$('.dropdown').change(() => {
    let key = event.target.id
    let val = ''

    $('select option:selected').each(function() {
        val += $(this).text() + ' ';
    });

    console.log('key: ' + key + ', val: ' + val)
})

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
                val: curUnit.name
            })
        )

        let dropdown = $('<select id="' + key + '" class="dropdown"/>').appendTo(versionElement)
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