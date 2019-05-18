versions = []
versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25']

$(document).ready(function() {
  console.log('page loaded')

  $('#addButton').click(function() {
	let beanID = $('#beanID').val()
    
	if (beanID.length !== 8 || isNaN(beanID)) {
		alert('Please enter a valid bean ID. Must be an 8 digit number.')
		return
	}

	for (key in versions) {
		if (versions[key].beanID === beanID) {
			alert('A unit already exists with that bean ID.')
			return
		}
	}

	addUnitData = { 'header': 'addUnit', 'id': ID(), 'beanID': beanID }
	console.log(addUnitData)

	$.ajax({
		type: 'POST',
		url: 'addUnit',
		data: JSON.stringify(addUnitData),
		success: function() { console.log('addUnit post success') },
		failure: function() { console.log('addUnit post failed') },
		dataType: 'json'
	})

	versions[addUnitData.id] = { 'model': '', 'version': '', 'beanID': addUnitData.beanID }
	buildTable()
  })

  $('#removeButton').click(function() {
	let beanID = $('#beanID').val()
	
	if (beanID.length !== 8 || isNaN(beanID)) {
		alert('Please enter a valid bean ID. Must be an 8 digit number.')
		return
	}

	for (key in versions) {
		if (versions[key].beanID === beanID) {
			let removeUnitData = { 'header': 'removeUnit', 'id': key }
			console.log(removeUnitData)

			$.ajax({
				type: 'POST',
				url: 'removeUnit',
				data: JSON.stringify(removeUnitData),
				success: function() { console.log('removeUnit post success') },
				failure: function() { console.log('removeUnit post failed') },
				dataType: 'json'
			})
	
			delete versions[key]
			buildTable()
			return
		}
	}
  })
})

$(window).on('load', () => {
  console.log('retrieving csv, refresh?')
  $.ajax({
    type: 'GET',
    url: 'versions.csv',
    dataType: 'text',
    success: (data) => { loadDict(data) }
  })	
})

function loadDict(text) {
  console.log('read:\n' + text)

  let lines = text.split(/\n/)

  for (let line of lines) {
    let values = line.split(",")
    
    let key = values[0]
    let currentVersion = values[1]
	let beanID = values[2]
 
    if (key) {
      versions[key] = { 'version': currentVersion, 'beanID': beanID }
    }
  }

  buildTable()
}

function buildTable() {
  $('#dictTable tbody tr').remove()
  console.log('dynamically building table')

  let table = $('#dictTable tbody')
  let count = 0

  for (var key in versions) {
    let tr = $('<tr>')
    let indexTD = $('<td>').html(count++)
    let currentVersionTD = $('<td id="' + key + '">').html(versions[key].version)
    let selectTD = $('<td>')
    let select = $('<select id="' + key + '" class="cs">').appendTo(selectTD) 
	let beanTD = $('<td>').html(versions[key].beanID)   
 
    tr.append(indexTD)
    tr.append(currentVersionTD)
	tr.append(beanTD)
    tr.append(selectTD)
    table.append(tr)


    for (let i = 0; i < versionOptions.length; i++) {
      if (versionOptions[i]) {
        let option = $('<option value="' + i + '">').html(versionOptions[i])
        option.appendTo(select)
      }
    }

    let selectVersionButton = $('<button class="selectButton" id="' + key + '">').html("Change Version")
    selectVersionButton.appendTo(selectTD)
  }
}

$(document).on('click', '#dictTable tbody tr td button.selectButton', function() {
  let clickedID = $(this).attr('id')
  let versionIndex = $('select#' + clickedID + ' option:checked').val()

  // update the dictionary
  versions[clickedID].version = versionOptions[versionIndex]
  // update the displayed current version
  buildTable()
  // post request to backend
  let updateData = { 
	'header': 'updateVersion',
	'version': versionOptions[versionIndex],
	'id': clickedID 
  }
 
  console.log(updateData)
 
  $.ajax({
    type: 'POST',
    url: 'updateVersion',
    data: JSON.stringify(updateData),
    success: function() { console.log('update post success') },
	failure: function() { console.log('update post failed') },
    dataType: 'json'
  })
})

function ID() {
  return '_' + Math.random().toString(36).substr(2, 9)
}
