let dict = {}
let versionOptions = []
let versionData = false, unitData = false

const interval = 5000

$(document).ready(function() {
    console.log('Getting data')
	//fakeData()
	initWorker()

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

function initWorker() {
	// taken from https://developers.google.com/web/fundamentals/primers/service-workers/
	if ('serviceWorker' in navigator) {
    	navigator.serviceWorker.register('/lab/getData.js').then(
			function(registration) {
      			// Registration was successful
      			console.log('ServiceWorker registration successful with scope: ', registration.scope);
    		}, function(err) {
      			// registration failed :(
      			console.log('ServiceWorker registration failed: ', err);
    		}
		);
	}
}

function fakeData() {
	dict = {
		"17867974393591147666": {
				"version":"3.30",
				"beanID":"720008e5",
				"name":"",
				"state":2
		},"17954946862877046502": {
				"version":"3.31",
				"beanID":"72000886",
				"name":"C2S950",
				"state":0
		},"3818512707708119105": {
				"version":"3.31",
				"beanID":"72001664",
				"name":"C2S900, needs a UPD",
				"state":1
		},"8843070808735976599": {
				"version":"2.30",
				"beanID":"72000bfa",
				"name":"",
				"state":0
		}
	}
	
	versionOptions = ["3.30", "3.31"]
}

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

    let name = ''

    if (dict[key].name) {
        name = dict[key].name
    } else {
        name = dict[key].beanID
    }

    if (!confirm('Do you want to update ' + name + ' to ' + version + '?')) {
        return
    }

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
