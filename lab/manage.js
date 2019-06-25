let dict = {}
let versionOptions = []
let versionData = false, unitData = false

const interval = 5000

$(document).ready(function() {
    console.log('Getting data')

	$.ajax({
		type: 'GET',
		url: '/api/versions/',
		contentType: 'application/json',
		success: (resp) => {
			versionOptions = Object.values(resp)
			buildTable()
		},
		failure: (resp) => {
			console.log('request to get versions failed')
			console.log(resp)
		}
	})

    $.ajax({
        type: 'GET',
        url: '/api/units/',
        contentType: 'application/json',
        success: (resp) => { 
            dict = resp
			buildTable()
            poll()
        },
		failure: (resp) => {
			console.log('request to get units failed')
			console.log(resp)
		}
    })
})

function poll() {
    //console.log('polling server')
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
    let id = element.id.split(' ')
    let name = element.value
    let key

    if (id.length == 2)
        key = id[1]
    else {
        console.error('ID didnt have two parts. ' + id)
        return
    }

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
    let id = element.id.split(' ')
    let version = element.value
    let key

    if (id.length == 2)
        key = id[1]
    else {
        console.error('ID didnt have two parts. ' + id)
        return
    }

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

    dict[key].state = 1
    $('#icon-' + key).replaceWith(makeIcon(key))
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

        nameElement.append($('<input type="text" id="textbox ' + key + '" onchange="boxListener(this)" value="' + curUnit.name + '">'))

        let dropdown = $('<select id="select ' + key + '" onchange="dropdownListener(this)"/>').appendTo(versionElement)
        dropdown.value = curUnit.version

        let hasNoOption = true
        versionOptions.forEach(val => {
            hasNoOption = hasNoOption && val !== curUnit.version
            $('<option />', {value: val, text: val, selected: val === curUnit.version}).appendTo(dropdown)
        })

        if (hasNoOption) {
            $('<option />', {text: '', selected: true}).appendTo(dropdown)
        }
        
        beanElement.append($('<div>').html(curUnit.beanID))
        stateElement.append(makeIcon(key))
    }
}

function makeIcon(key) {
    let state = dict[key].state

    switch (state) {
        case 0:
            //idle
            return $('<i id="icon-' + key + '" class="fas fa-check-circle" style="color: #34C53C">')
        case 1:
            // updating
            return $('<i id="icon-' + key + '" class="fas fa-spinner fa-pulse" style="color: #61D7FF">')
        case 2:
            // failed
            return $('<i id="icon-' + key + '" class="fas fa-times-circle" style="color: #FF0104">')
        default:
            console.log('unexpected state in units response: ' + state)
            return
    }
}
